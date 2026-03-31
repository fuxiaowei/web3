import hre from "hardhat";

/**
 * @title 部署脚本
 * @dev 演示如何使用Hardhat 3的新API部署合约
 * @notice 这是传统的脚本部署方式，适合简单的部署场景
 */
async function main() {
  // 获取签名者（部署账户）
  const [deployer] = await hre.ethers.getSigners();
  
  // 获取账户余额
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  
  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(balance), "ETH");
  
  // 检查余额是否足够
  if (balance === 0n) {
    throw new Error("Account has no balance. Please fund the account first.");
  }
  
  // 部署合约（Hardhat 3新API）
  console.log("\nDeploying Counter contract...");
  const counter = await hre.ethers.deployContract("Counter");
  
  // 等待部署确认
  await counter.waitForDeployment();
  
  // 获取合约地址（Hardhat 3新API）
  const address = await counter.getAddress();
  console.log("Counter deployed to:", address);
  
  // 验证部署
  const code = await hre.ethers.provider.getCode(address);
  if (code === "0x") {
    throw new Error("Contract deployment failed - no code at address");
  }
  
  console.log("Contract deployment verified!");
  
  // 可选：与合约交互
  console.log("\nInteracting with contract...");
  const initialValue = await counter.number();
  console.log("Initial number value:", initialValue.toString());
  
  // 调用increment函数
  const tx = await counter.increment();
  await tx.wait();
  
  const newValue = await counter.number();
  console.log("Number after increment:", newValue.toString());
  
  console.log("\nDeployment and interaction completed successfully!");
}

// 执行主函数
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

