import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
// 引入 Ignition 的模块构建函数，用于声明部署流程。

// 定义一个名为 MetaNFTAuction 的部署模块。
const metaNFTAuctionModule = buildModule("MetaNFTAuction", (m) => {
  // 部署 MetaNFTAuction 合约。
  const metaNFTAuction = m.contract("MetaNFTAuction");

  // 返回部署后的合约句柄，后续模块可以复用。
  return { metaNFTAuction };
});

// 导出模块，供 Ignition 命令直接调用。
export default metaNFTAuctionModule;