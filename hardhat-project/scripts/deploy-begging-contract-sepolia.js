/**
 * 将 BeggingContract 部署到 Ethereum Sepolia 测试网。
 *
 * 前置条件：
 *   1. 配置变量 SEPOLIA_RPC_URL（Infura / Alchemy 等提供的 Sepolia HTTPS RPC）
 *   2. 配置变量 SEPOLIA_PRIVATE_KEY（部署账户私钥，需带少量 Sepolia ETH 支付 gas）
 *
 * 设置方式（任选其一）：
 *   推荐：复制 env.example 为项目根目录下的 .env，填写 SEPOLIA_RPC_URL 与 SEPOLIA_PRIVATE_KEY
 *   或：npx hardhat keystore set SEPOLIA_RPC_URL / SEPOLIA_PRIVATE_KEY
 *   或：在 PowerShell 中 $env:SEPOLIA_RPC_URL="..." ; $env:SEPOLIA_PRIVATE_KEY="0x..."
 *
 * 执行：
 *   npx hardhat run scripts/deploy-begging-contract-sepolia.js --network sepolia
 *
 * 若报错请求 mainnet.infura.io：说明 .env 里误填了主网 RPC，应改为含 sepolia 的 URL。
 * 若 ETIMEDOUT：多为网络访问 Infura 受限，可换 Alchemy / 其它 Sepolia RPC，或检查代理与防火墙。
 */

import hre from "hardhat";

const SEPOLIA_CHAIN_ID = 11155111n;

async function main() {
  const conn = await hre.network.connect();
  const ethers = conn.ethers;
  if (!ethers) throw new Error("conn.ethers is undefined");

  try {
    const net = await ethers.provider.getNetwork();
    if (net.chainId !== SEPOLIA_CHAIN_ID) {
      throw new Error(
        `当前 chainId 为 ${net.chainId}，预期 Sepolia ${SEPOLIA_CHAIN_ID}。请使用: npx hardhat run scripts/deploy-begging-contract-sepolia.js --network sepolia`,
      );
    }

    const [deployer] = await ethers.getSigners();
    console.log("部署账户:", deployer.address);

    const BeggingContract = await ethers.getContractFactory("BeggingContract");
    const contract = await BeggingContract.deploy();
    await contract.waitForDeployment();

    const address = await contract.getAddress();
    console.log("BeggingContract 已部署到 Sepolia，合约地址:", address);
    console.log("浏览器: https://sepolia.etherscan.io/address/" + address);
  } finally {
    await conn.close();
  }
}

main().catch((e) => {
  console.error(e);
  process.exitCode = 1;
});
