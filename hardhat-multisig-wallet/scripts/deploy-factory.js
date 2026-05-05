/**
 * 部署 CrowdfundingFactory（无构造参数）。
 * 用法: npx hardhat run scripts/deploy-factory.js --network <sepolia|ganache|...>
 */
import { network } from "hardhat";

const { ethers } = await network.connect();

const factory = await ethers.deployContract("CrowdfundingFactory");
await factory.waitForDeployment();

const address = await factory.getAddress();
console.log("CrowdfundingFactory deployed to:", address);
console.log(`前端可设置 NEXT_PUBLIC_FACTORY_ADDRESS=${address}`);
