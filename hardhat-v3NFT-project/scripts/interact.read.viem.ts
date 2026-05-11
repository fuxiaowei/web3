import { createPublicClient, http, formatEther, formatUnits, zeroAddress } from "viem";
import hre from "hardhat";

const AUCTION_ADDRESS = (process.env.AUCTION_ADDRESS || "") as `0x${string}`;
const RPC_URL = process.env.RPC_URL || "http://127.0.0.1:8545";
const QUERY_AUCTION_ID = process.env.QUERY_AUCTION_ID ? BigInt(process.env.QUERY_AUCTION_ID) : undefined;

async function getAuctionABI() {
  const artifact = await hre.artifacts.readArtifact("MetaNFTAuction");
  return artifact.abi;
}

function ensureAddress(name: string, value: string) {
  if (!value || value === zeroAddress) {
    throw new Error(`${name} 未设置或是零地址，请在 .env 中配置`);
  }
}

async function main() {
  ensureAddress("AUCTION_ADDRESS", AUCTION_ADDRESS);

  const AUCTION_ABI = (await getAuctionABI()) as any;
  const publicClient = createPublicClient({
    transport: http(RPC_URL),
  } as any) as any;

  console.log("=== MetaNFTAuction 只读查询脚本 (viem) ===\n");
  console.log("连接地址:", AUCTION_ADDRESS);
  console.log("网络 ID:", await publicClient.getChainId(), "\n");

  const version = await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: AUCTION_ABI,
    functionName: "getVersion",
  });
  console.log("1. 合约版本:", version);

  const auctionId = (await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: AUCTION_ABI,
    functionName: "auctionId",
  })) as bigint;
  console.log("2. 当前拍卖ID(下一个ID):", auctionId.toString());

  if (auctionId === 0n) {
    console.log("\n当前暂无拍卖。");
    return;
  }

  const targetAuctionId = QUERY_AUCTION_ID ?? auctionId - 1n;
  if (targetAuctionId < 0n || targetAuctionId >= auctionId) {
    throw new Error(`QUERY_AUCTION_ID 超出范围，可选范围: 0 ~ ${auctionId - 1n}`);
  }

  const auctionData = (await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: AUCTION_ABI,
    functionName: "auctions",
    args: [targetAuctionId],
  })) as any[];

  console.log(`\n3. 拍卖 #${targetAuctionId} 详情:`);
  console.log("   - NFT地址:", auctionData[0]);
  console.log("   - NFT ID:", auctionData[1].toString());
  console.log("   - 卖家:", auctionData[2]);
  console.log("   - 开始时间:", new Date(Number(auctionData[3]) * 1000).toISOString());
  console.log("   - 最高出价者:", auctionData[4]);
  console.log("   - 起拍价(美元):", formatUnits(auctionData[5], 8));
  console.log("   - 持续时间:", auctionData[6].toString(), "秒");
  console.log("   - 支付代币:", auctionData[7]);
  console.log("   - 最高出价(raw):", auctionData[8].toString());
  console.log("   - 最高出价(美元):", formatUnits(auctionData[9], 8));
  console.log("   - 最高出价代币:", auctionData[10]);
  if (auctionData[10] === zeroAddress) {
    console.log("   - 最高出价(按ETH展示):", formatEther(auctionData[8]), "ETH");
  }

  const ended = (await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: AUCTION_ABI,
    functionName: "isEnded",
    args: [targetAuctionId],
  })) as boolean;
  console.log(`\n4. 拍卖 #${targetAuctionId} 是否已结束:`, ended);

  const ethOracle = (await publicClient.readContract({
    address: AUCTION_ADDRESS,
    abi: AUCTION_ABI,
    functionName: "tokenToOracle",
    args: [zeroAddress],
  })) as `0x${string}`;
  console.log("5. ETH Oracle 地址:", ethOracle);

  if (ethOracle !== zeroAddress) {
    const ethPrice = (await publicClient.readContract({
      address: AUCTION_ADDRESS,
      abi: AUCTION_ABI,
      functionName: "getPriceInDollar",
      args: [zeroAddress],
    })) as bigint;
    console.log("6. ETH 价格(美元):", formatUnits(ethPrice, 8));
  } else {
    console.log("6. ETH 价格(美元): 未配置预言机，跳过查询");
  }

  console.log("\n只读查询完成。");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
