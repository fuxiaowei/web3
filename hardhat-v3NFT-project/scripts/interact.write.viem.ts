import { createPublicClient, createWalletClient, http, parseEther, zeroAddress } from "viem";
import { privateKeyToAccount } from "viem/accounts";
import hre from "hardhat";

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
  const value = getEnvByKeys([
    `ACTION_${action.toUpperCase().replace(/-/g, "_")}_PRIVATE_KEY`,
    action.startsWith("bid-") ? "BIDDER_PRIVATE_KEY" : "",
    action === "end" ? "ENDER_PRIVATE_KEY" : "",
    "ADMIN_PRIVATE_KEY",
    "PRIVATE_KEY",
  ].filter(Boolean));

  if (!value || value === "0x") {
    throw new Error(`未找到 ${action} 对应私钥。请配置 PRIVATE_KEY 或动作专用私钥`);
  }
  return value as `0x${string}`;
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
