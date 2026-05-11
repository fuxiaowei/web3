import {
  createPublicClient,
  createWalletClient,
  http,
  numberToHex,
  parseEther,
  zeroAddress,
} from "viem";
import { privateKeyToAccount } from "viem/accounts";
import hre from "hardhat";

/** OpenZeppelin ERC1967 `_ADMIN_SLOT`（与链上 TransparentUpgradeableProxy 一致） */
const ERC1967_PROXY_ADMIN_SLOT =
  "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103" as const;

/**
 * MetaNFTAuction 代理存储：`Initializable` 的 uint8/bool 与首个字段 `address admin` 会被 Solidity **打包在同一 slot 0**，
 * 因此业务 admin 不在 slot 1（此前误判会导致读到 0）。若未来布局改为独占槽，可回退读 slot 1。
 */
const META_NFT_AUCTION_SLOT0 = numberToHex(0, { size: 32 });
const META_NFT_AUCTION_ADMIN_FALLBACK_SLOT = numberToHex(1, { size: 32 });

type Action = "set-oracle" | "start" | "bid-eth" | "bid-erc20" | "end";

const AUCTION_ADDRESS = (process.env.AUCTION_ADDRESS || "") as `0x${string}`;
const RPC_URL = process.env.RPC_URL || "http://127.0.0.1:8545";
// Hardhat CLI 会在 argv 中包含 "run"，优先读取环境变量可避免被误判为动作名。
const ACTION = ((process.env.ACTION || process.argv[2] || "").trim()) as Action;

async function getAuctionABI() {
  const artifact = await hre.artifacts.readArtifact("MetaNFTAuction");
  return artifact.abi;
}

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) {
    throw new Error(`缺少环境变量: ${name}`);
  }
  return value;
}

function requireAddressEnv(name: string): `0x${string}` {
  const value = requireEnv(name) as `0x${string}`;
  if (value === zeroAddress) {
    throw new Error(`${name} 不能是零地址`);
  }
  return value;
}

function parseBigIntEnv(name: string): bigint {
  const value = requireEnv(name);
  try {
    return BigInt(value);
  } catch {
    throw new Error(`${name} 不是合法整数: ${value}`);
  }
}

function getEnvByKeys(keys: string[]): string | undefined {
  for (const key of keys) {
    const value = process.env[key];
    if (value && value.trim().length > 0) {
      return value.trim();
    }
  }
  return undefined;
}

function getActionPrivateKey(action: Action): `0x${string}` {
  const keys: string[] = [
    `ACTION_${action.toUpperCase().replace(/-/g, "_")}_PRIVATE_KEY`,
  ];
  // 与 Ignition 一致：account(0) 常为 Proxy 的 admin，不能对代理做业务写入；业务管理员为 account(1)。
  // 本仓库 hardhat.config.js 里 ganache 的 accounts[1] = SEPOLIA_PRIVATE_KEY2，与 getAccount(1) 对齐；优先用它可“只跑脚本”。
  if (action === "set-oracle" || action === "start") {
    keys.push("SEPOLIA_PRIVATE_KEY2", "APP_ADMIN_PRIVATE_KEY");
  }
  if (action.startsWith("bid-")) {
    keys.push("BIDDER_PRIVATE_KEY");
  }
  if (action === "end") {
    keys.push("ENDER_PRIVATE_KEY");
  }
  keys.push("ADMIN_PRIVATE_KEY", "PRIVATE_KEY");

  const value = getEnvByKeys(keys);

  if (!value || value === "0x") {
    throw new Error(
      `未找到 ${action} 对应私钥。set-oracle/start 请配置 SEPOLIA_PRIVATE_KEY2（本仓库 ganache 与 Ignition getAccount(1) 一致）或 APP_ADMIN_PRIVATE_KEY / PRIVATE_KEY`
    );
  }
  return value as `0x${string}`;
}

function addressFromStorageWord(word: `0x${string}`): `0x${string}` {
  const hex = word.toLowerCase();
  if (hex.length < 42) {
    return zeroAddress;
  }
  return `0x${hex.slice(-40)}` as `0x${string}`;
}

/** 从 slot0 解析与 OZ Initializable 打包在一起的 `admin`（当前合约实测：hex[20:60]）。 */
function parseAppAdminFromInitializableSlot0(raw: `0x${string}` | undefined): `0x${string}` {
  if (!raw || raw === "0x") {
    return zeroAddress;
  }
  let hex = raw.toLowerCase().replace(/^0x/, "");
  if (hex.length < 64) {
    hex = hex.padStart(64, "0");
  }
  if (/^0+$/.test(hex)) {
    return zeroAddress;
  }
  const addrPart = hex.slice(20, 60);
  if (/^0+$/.test(addrPart)) {
    return zeroAddress;
  }
  return `0x${addrPart}` as `0x${string}`;
}

/** 透明代理的 admin 不能通过代理调用实现合约（含 estimateGas），否则会 revert。 */
async function assertNotTransparentProxyAdmin(
  publicClient: { getStorageAt: (args: unknown) => Promise<`0x${string}` | undefined> },
  proxyAddress: `0x${string}`,
  signerAddress: `0x${string}`
) {
  const raw = await publicClient.getStorageAt({
    address: proxyAddress,
    slot: ERC1967_PROXY_ADMIN_SLOT,
  });
  if (!raw || raw === "0x" + "0".repeat(64)) {
    return;
  }
  const proxyAdmin = addressFromStorageWord(raw);
  if (proxyAdmin === zeroAddress) {
    return;
  }
  if (proxyAdmin.toLowerCase() === signerAddress.toLowerCase()) {
    throw new Error(
      "当前账户是 TransparentUpgradeableProxy 的链上 admin（常为部署时 account(0)），不能对代理执行业务写入（setTokenOracle/start/bid/end 等）。\n" +
        "请设置 APP_ADMIN_PRIVATE_KEY 为 initialize 时传入的业务管理员私钥（Ignition 中 metaNFTAuctionProxyModule 里通常为 getAccount(1)，\n" +
        "且必须与**部署该 AUCTION_ADDRESS 时** Hardhat 同一网络的第 2 个账户一致）。\n" +
        "或设置 BIDDER_PRIVATE_KEY / ENDER_PRIVATE_KEY 等非 admin 账户；勿对写操作使用与 Proxy admin 相同的 PRIVATE_KEY。"
    );
  }
}

/** 从代理合约存储读取 `MetaNFTAuction.admin`（优先 slot0 打包布局，其次 slot1 单字对齐）。 */
async function readMetaNFTAuctionAppAdmin(
  publicClient: { getStorageAt: (args: unknown) => Promise<`0x${string}` | undefined> },
  proxyAddress: `0x${string}`
): Promise<`0x${string}`> {
  const raw0 = await publicClient.getStorageAt({
    address: proxyAddress,
    slot: META_NFT_AUCTION_SLOT0,
  });
  const fromSlot0 = parseAppAdminFromInitializableSlot0(raw0);
  if (fromSlot0 !== zeroAddress) {
    return fromSlot0;
  }
  const raw1 = await publicClient.getStorageAt({
    address: proxyAddress,
    slot: META_NFT_AUCTION_ADMIN_FALLBACK_SLOT,
  });
  if (!raw1) {
    return zeroAddress;
  }
  return addressFromStorageWord(raw1);
}

async function assertSignerIsMetaNFTAuctionAppAdmin(
  publicClient: { getStorageAt: (args: unknown) => Promise<`0x${string}` | undefined> },
  proxyAddress: `0x${string}`,
  signerAddress: `0x${string}`
) {
  const onChainAdmin = await readMetaNFTAuctionAppAdmin(publicClient, proxyAddress);
  if (onChainAdmin === zeroAddress) {
    throw new Error(
      "从代理存储读取的 MetaNFTAuction.admin 为零地址，合约可能未初始化或存储布局与脚本不一致。"
    );
  }
  if (onChainAdmin.toLowerCase() === signerAddress.toLowerCase()) {
    return;
  }

  let hardhatAccount1 = "";
  try {
    const signers = await hre.ethers.getSigners();
    const s1 = signers[1];
    if (s1) {
      hardhatAccount1 = await s1.getAddress();
    }
  } catch {
    /* 非 Hardhat 环境或未配置 ethers 时忽略 */
  }

  const repoHint =
    "本仓库 `hardhat.config.js` 里 `ganache` 网络的 accounts[1] 来自环境变量 **SEPOLIA_PRIVATE_KEY2**（Ignition 模块里 `getAccount(1)` 与之对应）。\n" +
    "请把 **SEPOLIA_PRIVATE_KEY2** 设为 APP_ADMIN_PRIVATE_KEY，或直接使用 `set APP_ADMIN_PRIVATE_KEY=%SEPOLIA_PRIVATE_KEY2%`（PowerShell 用 `$env:APP_ADMIN_PRIVATE_KEY=$env:SEPOLIA_PRIVATE_KEY2`）。\n";

  const hint =
    hardhatAccount1 !== ""
      ? `\n当前 \`npx hardhat run ... --network <name>\` 下 ethers.getSigners()[1] = ${hardhatAccount1}。\n` +
        "若该地址与「链上 admin」一致，请使用**该 signer 对应私钥**作为 APP_ADMIN_PRIVATE_KEY。\n"
      : "\n请在 Hardhat 网络配置中找到与「链上 admin」一致的账户私钥，并设为 APP_ADMIN_PRIVATE_KEY。\n";

  throw new Error(
    `onlyAdmin：链上业务管理员地址 MetaNFTAuction.admin = ${onChainAdmin}\n` +
      `当前 APP_ADMIN_PRIVATE_KEY 对应 signer = ${signerAddress}\n` +
      `二者不一致，故 setTokenOracle/start 会 revert。\n` +
      repoHint +
      hint +
      "若链已重置而 AUCTION_ADDRESS 仍是旧部署，地址也会对不上；请重新部署或恢复同一链状态。"
  );
}

function requireActionEnv(action: Action, suffix: string, fallback?: string): string {
  const key = `ACTION_${action.toUpperCase().replace(/-/g, "_")}_${suffix}`;
  const value = getEnvByKeys([key, fallback || ""].filter(Boolean));
  if (!value) {
    throw new Error(`缺少环境变量: ${key}${fallback ? ` (或 ${fallback})` : ""}`);
  }
  return value;
}

function requireActionAddressEnv(action: Action, suffix: string, fallback?: string): `0x${string}` {
  const value = requireActionEnv(action, suffix, fallback) as `0x${string}`;
  if (value === zeroAddress) {
    throw new Error(`ACTION_${action.toUpperCase().replace(/-/g, "_")}_${suffix} 不能是零地址`);
  }
  return value;
}

function parseActionBigIntEnv(action: Action, suffix: string, fallback?: string): bigint {
  const value = requireActionEnv(action, suffix, fallback);
  try {
    return BigInt(value);
  } catch {
    throw new Error(`环境变量值不是合法整数: ${value}`);
  }
}

async function waitTx(
  publicClient: any,
  walletClient: any,
  params: {
    address: `0x${string}`;
    abi: any;
    functionName: string;
    args?: readonly unknown[];
    value?: bigint;
  },
  label: string
) {
  console.log(`${label} -> 发送交易中...`);
  const hash = await walletClient.writeContract(params);
  console.log(`${label} -> 交易哈希:`, hash);
  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  console.log(`${label} -> 已确认，区块:`, receipt.blockNumber.toString());
}

async function main() {
  if (!AUCTION_ADDRESS || AUCTION_ADDRESS === zeroAddress) {
    throw new Error("AUCTION_ADDRESS 未设置或是零地址");
  }
  if (!ACTION) {
    throw new Error("请设置 ACTION（或命令行参数）: set-oracle | start | bid-eth | bid-erc20 | end");
  }

  const privateKey = getActionPrivateKey(ACTION);
  const account = privateKeyToAccount(privateKey);
  const AUCTION_ABI = (await getAuctionABI()) as any;
  const transport = http(RPC_URL);

  const publicClient = createPublicClient({ transport } as any) as any;
  const walletClient = createWalletClient({ account, transport } as any) as any;

  console.log("=== MetaNFTAuction 写操作脚本 (viem) ===");
  console.log("ACTION:", ACTION);
  console.log("合约地址:", AUCTION_ADDRESS);
  console.log("执行账户:", account.address);
  console.log("网络 ID:", await publicClient.getChainId(), "\n");

  await assertNotTransparentProxyAdmin(publicClient, AUCTION_ADDRESS, account.address);

  if (ACTION === "set-oracle" || ACTION === "start") {
    const appAdmin = await readMetaNFTAuctionAppAdmin(publicClient, AUCTION_ADDRESS);
    console.log("链上 MetaNFTAuction.admin (storage):", appAdmin);
    // 本仓库 hardhat.config.js：ganache 的 accounts[0..] 对应 SEPOLIA_PRIVATE_KEY、SEPOLIA_PRIVATE_KEY2…
    // Ignition `getAccount(1)` = 第二个私钥，即 SEPOLIA_PRIVATE_KEY2。
    const pk2 = getEnvByKeys(["SEPOLIA_PRIVATE_KEY2"]);
    if (pk2 && pk2.startsWith("0x") && pk2.length > 2) {
      try {
        const fromPk2 = privateKeyToAccount(pk2 as `0x${string}`).address;
        console.log("环境变量 SEPOLIA_PRIVATE_KEY2 对应地址:", fromPk2);
      } catch {
        /* ignore */
      }
    }
    try {
      const signers = await hre.ethers.getSigners();
      const a1 = signers[1] ? await signers[1]!.getAddress() : "";
      if (a1) {
        console.log("Hardhat ethers.getSigners()[1]:", a1, "\n");
      }
    } catch {
      console.log("");
    }
    await assertSignerIsMetaNFTAuctionAppAdmin(publicClient, AUCTION_ADDRESS, account.address);
  }

  if (ACTION === "set-oracle") {
    const token = (getEnvByKeys([
      "ACTION_SET_ORACLE_TOKEN_ADDRESS",
      "TOKEN_ADDRESS",
    ]) || zeroAddress) as `0x${string}`;
    const oracle = requireActionAddressEnv(ACTION, "ORACLE_ADDRESS", "ORACLE_ADDRESS");
    await waitTx(
      publicClient,
      walletClient,
      {
        address: AUCTION_ADDRESS,
        abi: AUCTION_ABI,
        functionName: "setTokenOracle",
        args: [token, oracle],
      },
      "setTokenOracle"
    );
    return;
  }

  if (ACTION === "start") {
    const seller = requireActionAddressEnv(ACTION, "SELLER_ADDRESS", "SELLER_ADDRESS");
    const nft = requireActionAddressEnv(ACTION, "NFT_ADDRESS", "NFT_ADDRESS");
    const nftId = parseActionBigIntEnv(ACTION, "NFT_ID", "NFT_ID");
    const startingPriceUsd = parseActionBigIntEnv(ACTION, "STARTING_PRICE_USD", "STARTING_PRICE_USD");
    const durationSec = parseActionBigIntEnv(ACTION, "DURATION_SEC", "DURATION_SEC");
    const paymentToken = requireActionAddressEnv(ACTION, "PAYMENT_TOKEN", "PAYMENT_TOKEN");

    await waitTx(
      publicClient,
      walletClient,
      {
        address: AUCTION_ADDRESS,
        abi: AUCTION_ABI,
        functionName: "start",
        args: [seller, nftId, nft, startingPriceUsd, durationSec, paymentToken],
      },
      "start"
    );
    return;
  }

  if (ACTION === "bid-eth") {
    const auctionId = parseActionBigIntEnv(ACTION, "AUCTION_ID", "AUCTION_ID");
    const bidAmountEth = requireActionEnv(ACTION, "BID_AMOUNT_ETH", "BID_AMOUNT_ETH");
    const amountWei = parseEther(bidAmountEth);

    await waitTx(
      publicClient,
      walletClient,
      {
        address: AUCTION_ADDRESS,
        abi: AUCTION_ABI,
        functionName: "bid",
        args: [auctionId, amountWei],
        value: amountWei,
      },
      "bid(ETH)"
    );
    return;
  }

  if (ACTION === "bid-erc20") {
    const auctionId = parseActionBigIntEnv(ACTION, "AUCTION_ID", "AUCTION_ID");
    const amount = parseActionBigIntEnv(ACTION, "BID_AMOUNT", "BID_AMOUNT");
    await waitTx(
      publicClient,
      walletClient,
      {
        address: AUCTION_ADDRESS,
        abi: AUCTION_ABI,
        functionName: "bid",
        args: [auctionId, amount],
      },
      "bid(ERC20)"
    );
    return;
  }

  if (ACTION === "end") {
    const auctionId = parseActionBigIntEnv(ACTION, "AUCTION_ID", "AUCTION_ID");
    await waitTx(
      publicClient,
      walletClient,
      {
        address: AUCTION_ADDRESS,
        abi: AUCTION_ABI,
        functionName: "end",
        args: [auctionId],
      },
      "end"
    );
    return;
  }

  throw new Error(`不支持的 ACTION: ${ACTION}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
