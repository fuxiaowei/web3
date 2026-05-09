import "dotenv/config";

import hardhatEthersPlugin from "@nomicfoundation/hardhat-ethers";
import hardhatIgnitionPlugin from "@nomicfoundation/hardhat-ignition";
import hardhatMochaPlugin from "@nomicfoundation/hardhat-mocha";
import hardhatViemPlugin from "@nomicfoundation/hardhat-viem";
import { configVariable, defineConfig } from "hardhat/config";

export default defineConfig({
  plugins: [
    hardhatEthersPlugin,
    hardhatViemPlugin,
    hardhatMochaPlugin,
    hardhatIgnitionPlugin,
  ],
  solidity: {
    version: "0.8.28",
    npmFilesToBuild: [
      "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol",
      "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol",
    ],
  },
  networks: {
    sepolia: {
      type: "http",
      chainType: "l1",
      chainId: 11155111,
      url: configVariable("SEPOLIA_RPC_URL"),
      accounts: [configVariable("SEPOLIA_PRIVATE_KEY")],
    },
  },
});