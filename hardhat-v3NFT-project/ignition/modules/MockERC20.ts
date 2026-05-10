import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

// 部署 MockERC20（默认 6 位精度，初始供应量 1,000,000 枚）。
const mockERC20Module = buildModule("MockERC20Module", (m) => {
  const name = m.getParameter("name", "Mock USD Coin");
  const symbol = m.getParameter("symbol", "mUSDC");
  const decimals = m.getParameter("decimals", 6);
  const initialSupply = m.getParameter("initialSupply", 1000000000000n);

  const mockERC20 = m.contract("MockERC20", [
    name,
    symbol,
    decimals,
    initialSupply,
  ]);

  return { mockERC20 };
});

export default mockERC20Module;
