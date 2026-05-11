import { createPublicClient, createWalletClient, formatUnits, http, parseEther, zeroAddress } from "viem";
import { privateKeyToAccount } from "viem/accounts";
import hre from "hardhat";

/**
 * 新手流程脚本（Ganache 友好版）
 *
 * 目标：把拍卖主流程串起来，避免手工一步步敲命令容易漏。
 * 默认执行顺序：
 * 1) set-oracle
 * 2) mint-nft（若已存在则跳过）
 * 3) approve-nft
 * 4) start
 * 5) bid-eth
 * 6) end（可选自动快进时间）
 * 7) status（打印最终状态）
 *
 * 用法：
 * - 执行全部：npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- all
 * - 执行单步：npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- start
 */

type Step =
  | "all"
  | "set-oracle"
  | "mint-nft"
  | "approve-nft"
  | "start"
  | "bid-eth"
  | "end"
  | "status";

// Hardhat CLI 会在 argv 中包含 "run"，优先读取环境变量可避免被误判为步骤名。
const STEP = ((process.env.STEP || process.argv[2] || "all").trim()) as Step;
const RPC_URL = process.env.RPC_URL || "http://127.0.0.1:7545";
const AUCTION_ADDRESS = (process.env.AUCTION_ADDRESS || "") as `0x${string}`;
const NFT_ADDRESS = (process.env.NFT_ADDRESS || "") as `0x${string}`;
const PAYMENT_TOKEN = (process.env.PAYMENT_TOKEN || "") as `0x${string}`;
const ORACLE_ADDRESS = (process.env.ORACLE_ADDRESS || "") as `0x${string}`;

const NFT_ID = BigInt(process.env.NFT_ID || "1");
const STARTING_PRICE_USD = BigInt(process.env.STARTING_PRICE_USD || "1000");
const DURATION_SEC = BigInt(process.env.DURATION_SEC || "300");
const BID_AMOUNT_ETH = process.env.BID_AMOUNT_ETH || "0.1";
const AUCTION_ID_ENV = process.env.AUCTION_ID ? BigInt(process.env.AUCTION_ID) : undefined;
const AUTO_TIME_TRAVEL = (process.env.AUTO_TIME_TRAVEL || "true").toLowerCase() === "true";
const END_WAIT_SEC = Number(process.env.END_WAIT_SEC || "0");

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value || value.trim().length === 0) {
    throw new Error(`缺少环境变量: ${name}`);
  }
  return value.trim();
}

function requireAddress(name: string, value: string): `0x${string}` {
  if (!value || value === zeroAddress) {
    throw new Error(`${name} 未配置或是零地址`);
  }
  return value as `0x${string}`;
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function getAuctionABI() {
  const artifact = await hre.artifacts.readArtifact("MetaNFTAuction");
  return artifact.abi;
}

async function getNftABI() {
  const artifact = await hre.artifacts.readArtifact("MetaNFT");
  return artifact.abi;
}

async function waitTx(
  label: string,
  publicClient: any,
  walletClient: any,
  params: {
    address: `0x${string}`;
    abi: any;
    functionName: string;
    args?: readonly unknown[];
    value?: bigint;
  }
) {
  console.log(`\n[${label}] 发送交易中...`);
  const hash = await walletClient.writeContract(params);
  console.log(`[${label}] tx hash:`, hash);
  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  console.log(`[${label}] 已确认，区块:`, receipt.blockNumber.toString());
}

async function resolveAuctionId(publicClient: any, auctionAbi: any): Promise<bigint> {
  if (AUCTION_ID_ENV !== undefined) {
    return AUCTION_ID_ENV;
  }
  const nextAuctionId = (await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: auctionAbi,
    functionName: "auctionId",
  })) as bigint;
  if (nextAuctionId === 0n) {
    throw new Error("当前还没有拍卖，无法自动推导 AUCTION_ID");
  }
  return nextAuctionId - 1n;
}

async function run() {
  requireAddress("AUCTION_ADDRESS", AUCTION_ADDRESS);
  requireAddress("NFT_ADDRESS", NFT_ADDRESS);
  requireAddress("PAYMENT_TOKEN", PAYMENT_TOKEN);
  requireAddress("ORACLE_ADDRESS", ORACLE_ADDRESS);

  const adminAccount = privateKeyToAccount(requireEnv("ADMIN_PRIVATE_KEY") as `0x${string}`);
  const sellerAccount = privateKeyToAccount(requireEnv("SELLER_PRIVATE_KEY") as `0x${string}`);
  const bidderAccount = privateKeyToAccount(requireEnv("BIDDER_PRIVATE_KEY") as `0x${string}`);
  const enderPrivateKey = (process.env.ENDER_PRIVATE_KEY || process.env.ADMIN_PRIVATE_KEY) as `0x${string}`;
  const enderAccount = privateKeyToAccount(enderPrivateKey);

  const sellerAddressFromEnv = process.env.SELLER_ADDRESS as `0x${string}` | undefined;
  const sellerAddress = sellerAddressFromEnv || sellerAccount.address;

  const transport = http(RPC_URL);
  const publicClient = createPublicClient({ transport } as any) as any;
  const adminWalletClient = createWalletClient({ account: adminAccount, transport } as any) as any;
  const sellerWalletClient = createWalletClient({ account: sellerAccount, transport } as any) as any;
  const bidderWalletClient = createWalletClient({ account: bidderAccount, transport } as any) as any;
  const enderWalletClient = createWalletClient({ account: enderAccount, transport } as any) as any;

  const auctionAbi = (await getAuctionABI()) as any;
  const nftAbi = (await getNftABI()) as any;

  const chainId = await publicClient.getChainId();
  console.log("=== 新手拍卖流程脚本 ===");
  console.log("STEP:", STEP);
  console.log("链 ID:", chainId);
  console.log("AUCTION:", AUCTION_ADDRESS);
  console.log("NFT:", NFT_ADDRESS);
  console.log("SELLER:", sellerAddress);
  console.log("BIDDER:", bidderAccount.address);

  // 仅提示，不强制：如果你手填了 SELLER_ADDRESS，但与 SELLER_PRIVATE_KEY 导出的地址不同，容易混淆。
  if (sellerAddressFromEnv && sellerAddressFromEnv.toLowerCase() !== sellerAccount.address.toLowerCase()) {
    console.log(
      "[警告] SELLER_ADDRESS 与 SELLER_PRIVATE_KEY 对应地址不一致。请确认 NFT 真正在 SELLER_ADDRESS 名下。"
    );
  }

  const runSetOracle = async () => {
    await waitTx("set-oracle", publicClient, adminWalletClient, {
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "setTokenOracle",
      args: [zeroAddress, ORACLE_ADDRESS],
    });
  };

  const runMintNft = async () => {
    // 先检查 token 是否已经存在。
    let owner: `0x${string}` | undefined;
    try {
      owner = (await publicClient.readContract({
        address: NFT_ADDRESS,
        abi: nftAbi,
        functionName: "ownerOf",
        args: [NFT_ID],
      })) as `0x${string}`;
    } catch {
      owner = undefined;
    }

    if (!owner) {
      // MetaNFT 的 mint 无权限限制，任意地址都可调用。这里让卖家自己 mint，更直观。
      await waitTx("mint-nft", publicClient, sellerWalletClient, {
        address: NFT_ADDRESS,
        abi: nftAbi,
        functionName: "mint",
        args: [sellerAddress, NFT_ID],
      });
      return;
    }

    if (owner.toLowerCase() !== sellerAddress.toLowerCase()) {
      throw new Error(`NFT #${NFT_ID} 已存在，但所有者是 ${owner}，不是 SELLER_ADDRESS ${sellerAddress}`);
    }
    console.log(`[mint-nft] NFT #${NFT_ID} 已存在且归属卖家，跳过`);
  };

  const runApproveNft = async () => {
    // 先查是否已经授权给拍卖合约，已授权则跳过。
    const approved = (await publicClient.readContract({
      address: NFT_ADDRESS,
      abi: nftAbi,
      functionName: "getApproved",
      args: [NFT_ID],
    })) as `0x${string}`;
    if (approved.toLowerCase() === AUCTION_ADDRESS.toLowerCase()) {
      console.log(`[approve-nft] NFT #${NFT_ID} 已授权给拍卖合约，跳过`);
      return;
    }

    await waitTx("approve-nft", publicClient, sellerWalletClient, {
      address: NFT_ADDRESS,
      abi: nftAbi,
      functionName: "approve",
      args: [AUCTION_ADDRESS, NFT_ID],
    });
  };

  const runStart = async () => {
    await waitTx("start", publicClient, adminWalletClient, {
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "start",
      args: [sellerAddress, NFT_ID, NFT_ADDRESS, STARTING_PRICE_USD, DURATION_SEC, PAYMENT_TOKEN],
    });
  };

  const runBidEth = async () => {
    const targetAuctionId = await resolveAuctionId(publicClient, auctionAbi);
    const amountWei = parseEther(BID_AMOUNT_ETH);
    await waitTx("bid-eth", publicClient, bidderWalletClient, {
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "bid",
      args: [targetAuctionId, amountWei],
      value: amountWei,
    });
  };

  const runEnd = async () => {
    const targetAuctionId = await resolveAuctionId(publicClient, auctionAbi);
    const ended = (await publicClient.readContract({
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "isEnded",
      args: [targetAuctionId],
    })) as boolean;

    // 新手最常见卡点：刚 start + bid 就 end，会报 not ended。
    // 在 Ganache 本地链上可自动快进时间，省去手工操作。
    if (!ended && AUTO_TIME_TRAVEL) {
      const jumpSec = Number(DURATION_SEC) + 5;
      console.log(`[end] 拍卖未结束，自动快进 ${jumpSec} 秒...`);
      await publicClient.request({ method: "evm_increaseTime", params: [jumpSec] });
      await publicClient.request({ method: "evm_mine", params: [] });
    } else if (!ended && END_WAIT_SEC > 0) {
      console.log(`[end] 拍卖未结束，等待 ${END_WAIT_SEC} 秒后继续...`);
      await sleep(END_WAIT_SEC * 1000);
    }

    await waitTx("end", publicClient, enderWalletClient, {
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "end",
      args: [targetAuctionId],
    });
  };

  const runStatus = async () => {
    const version = await publicClient.readContract({
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "getVersion",
    });
    const nextAuctionId = (await publicClient.readContract({
      address: AUCTION_ADDRESS,
      abi: auctionAbi,
      functionName: "auctionId",
    })) as bigint;
    console.log("\n=== 状态检查 ===");
    console.log("合约版本:", version);
    console.log("下一个 auctionId:", nextAuctionId.toString());
    if (nextAuctionId > 0n) {
      const targetAuctionId = AUCTION_ID_ENV ?? nextAuctionId - 1n;
      const ended = (await publicClient.readContract({
        address: AUCTION_ADDRESS,
        abi: auctionAbi,
        functionName: "isEnded",
        args: [targetAuctionId],
      })) as boolean;
      const auctionData = (await publicClient.readContract({
        address: AUCTION_ADDRESS,
        abi: auctionAbi,
        functionName: "auctions",
        args: [targetAuctionId],
      })) as any[];
      console.log("拍卖ID:", targetAuctionId.toString());
      console.log("是否结束:", ended);
      console.log("最高出价者:", auctionData[4]);
      console.log("最高出价(美元):", formatUnits(auctionData[9], 8));
    }
  };

  if (STEP === "all") {
    await runSetOracle();
    await runMintNft();
    await runApproveNft();
    await runStart();
    await runBidEth();
    await runEnd();
    await runStatus();
    console.log("\n全部步骤执行完成。");
    return;
  }

  if (STEP === "set-oracle") return runSetOracle();
  if (STEP === "mint-nft") return runMintNft();
  if (STEP === "approve-nft") return runApproveNft();
  if (STEP === "start") return runStart();
  if (STEP === "bid-eth") return runBidEth();
  if (STEP === "end") return runEnd();
  if (STEP === "status") return runStatus();

  throw new Error(`不支持的 STEP: ${STEP}`);
}

run()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
