import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

// 部署 MockOracle（默认价格按 8 位小数表示：3000 * 1e8）。
const mockOracleModule = buildModule("MockOracleModule", (m) => {
  const initialPrice = m.getParameter("initialPrice", 300000000000n);

  const mockOracle = m.contract("MockOracle", [initialPrice]);

  return { mockOracle };
});

export default mockOracleModule;
