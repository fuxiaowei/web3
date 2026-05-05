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

export default defineConfig({
  plugins: [hardhatEthersPlugin],
  solidity: "0.8.28",
  networks: {
    sepolia: {
      type: "http",
      chainType: "l1",
      chainId: 11155111,
      url: sepoliaRpcUrlOrConfigVariable(),
      accounts: [sepoliaPrivateKeyOrConfigVariable()],
    },
  },
});