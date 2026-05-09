import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
// 引入 Ignition 的模块构建函数。

// 定义部署 MetaNFTAuctionV2 实现合约的模块。
const metaNFTAuctionModule = buildModule("MetaNFTAuctionV2", (m) => {
  // 部署 MetaNFTAuctionV2 合约本体（通常作为升级后的实现逻辑）。
  const metaNFTAuction = m.contract("MetaNFTAuctionV2");

  // 将部署结果导出给其他模块使用。
  return { metaNFTAuction };
});

// 默认导出该模块。
export default metaNFTAuctionModule;