import "dotenv/config";

import hardhatEthersPlugin from "@nomicfoundation/hardhat-ethers";
import { configVariable, defineConfig } from "hardhat/config";

/**
 * Sepolia 部署账户：必须是 32 字节 ECDSA 私钥（64 个十六进制字符，可选 0x 前缀）。
 * 填成「钱包地址」或长度不对的十六进制会在底层报 invalid private key ... got object。
 */
function sepoliaPrivateKeyOrConfigVariable() {
  const raw = process.env.SEPOLIA_PRIVATE_KEY?.trim();
  if (raw === undefined || raw === "") {
    return configVariable("SEPOLIA_PRIVATE_KEY");
  }
  const body = raw.startsWith("0x") || raw.startsWith("0X") ? raw.slice(2) : raw;
  if (!/^[0-9a-fA-F]{64}$/.test(body)) {
    throw new Error(
      "SEPOLIA_PRIVATE_KEY 无效：需要 32 字节私钥（恰好 64 个十六进制字符，可选 0x 前缀）。不要填账户地址、助记词或截断的密钥。",
    );
  }
  return `0x${body.toLowerCase()}`;
}

/**
 * Infura 主网端点是 mainnet.infura.io；Sepolia 必须是 sepolia.infura.io（或 URL 中含 sepolia 的其它服务商）。
 */
function sepoliaRpcUrlOrConfigVariable() {
  const raw = process.env.SEPOLIA_RPC_URL?.trim();
  if (raw === undefined || raw === "") {
    return configVariable("SEPOLIA_RPC_URL");
  }
  let parsed;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error("SEPOLIA_RPC_URL 不是合法 URL（需以 https:// 开头）。");
  }
  if (parsed.protocol !== "https:" && parsed.protocol !== "http:") {
    throw new Error("SEPOLIA_RPC_URL 仅支持 http(s) 协议。");
  }
  const host = parsed.hostname.toLowerCase();
  if (host === "mainnet.infura.io" || host === "ethereum.publicnode.com") {
    throw new Error(
      "SEPOLIA_RPC_URL 指向了以太坊主网或其它非 Sepolia 节点。请改为 Sepolia，例如：https://sepolia.infura.io/v3/<你的密钥>",
    );
  }
  if (!/sepolia/i.test(raw)) {
    throw new Error(
      "SEPOLIA_RPC_URL 应指向 Sepolia 测试网（完整 URL 中通常包含 sepolia，例如 https://sepolia.infura.io/v3/... 或 eth-sepolia.g.alchemy.com）。",
    );
  }
  return raw;
}

/**
 * 本地 Ganache：默认 http://127.0.0.1:7545（Ganache 桌面常见端口）；CLI 常用 8545 时可设 GANACHE_RPC_URL。
 */
function ganacheRpcUrlOrDefault() {
  const raw = process.env.GANACHE_RPC_URL?.trim();
  if (raw === undefined || raw === "") {
    return "http://127.0.0.1:7545";
  }
  let parsed;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error("GANACHE_RPC_URL 不是合法 URL（需以 http:// 或 https:// 开头）。");
  }
  if (parsed.protocol !== "https:" && parsed.protocol !== "http:") {
    throw new Error("GANACHE_RPC_URL 仅支持 http(s) 协议。");
  }
  return raw;
}

function ganacheChainIdOrDefault() {
  const raw = process.env.GANACHE_CHAIN_ID?.trim();
  if (raw === undefined || raw === "") {
    return 1337;
  }
  const n = Number(raw);
  if (!Number.isInteger(n) || n <= 0) {
    throw new Error("GANACHE_CHAIN_ID 无效：需为正整数（例如 1337 或 5777，与 Ganache 中 Network Chain ID 一致）。");
  }
  return n;
}

/**
 * 未设置 GANACHE_PRIVATE_KEY 时使用经典确定性助记词「candy maple ... sweet treat」下 HD 路径默认的第一个账户，
 * 与 ganache-cli --wallet.deterministic / 常见 Ganache 固定钱包一致；随机生成账户时请在 .env 填写 Ganache 显示的私钥。
 */
const GANACHE_DEFAULT_DETERMINISTIC_FIRST_PRIVATE_KEY =
  "0x82d52fc43999a9ee0da164d048c797902dd1893cedfb7001b5a11bc89c1fc9e3";

function ganacheAccounts() {
  const raw = process.env.GANACHE_PRIVATE_KEY?.trim();
  if (raw === undefined || raw === "") {
    return [GANACHE_DEFAULT_DETERMINISTIC_FIRST_PRIVATE_KEY];
  }
  const body = raw.startsWith("0x") || raw.startsWith("0X") ? raw.slice(2) : raw;
  if (!/^[0-9a-fA-F]{64}$/.test(body)) {
    throw new Error(
      "GANACHE_PRIVATE_KEY 无效：需要 32 字节私钥（恰好 64 个十六进制字符，可选 0x 前缀）。不要填账户地址。",
    );
  }
  return [`0x${body.toLowerCase()}`];
}

export default defineConfig({
  plugins: [hardhatEthersPlugin],
  // solc 0.8.x 默认目标常为 cancun；Ganache 等本地链 EVM 较旧，会报 invalid opcode。paris 与多数 Ganache / 测试网兼容。
  solidity: {
    version: "0.8.28",
    settings: {
      evmVersion: "paris",
      optimizer: { enabled: true, runs: 200 },
    },
  },
  networks: {
    sepolia: {
      type: "http",
      chainType: "l1",
      chainId: 11155111,
      url: sepoliaRpcUrlOrConfigVariable(),
      accounts: [sepoliaPrivateKeyOrConfigVariable()],
    },
    ganache: {
      type: "http",
      chainType: "l1",
      chainId: ganacheChainIdOrDefault(),
      url: ganacheRpcUrlOrDefault(),
      accounts: ganacheAccounts(),
    },
  },
});