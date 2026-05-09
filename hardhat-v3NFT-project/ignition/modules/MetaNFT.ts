import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
// 引入 Ignition 的模块构建函数，用来定义一个可复用的部署单元。

// 定义名为 MetaNFTModule 的部署模块。
const metaNFTModule = buildModule("MetaNFTModule", (m) => {
  // 部署名为 MetaNFT 的合约（合约名需与 artifacts 中编译产物一致）。
  const metaNFT = m.contract("MetaNFT");

  // 将部署出来的合约实例导出，供其他模块或脚本继续使用。
  return { metaNFT };
});

// 导出该模块，供 `npx hardhat ignition deploy` 指定使用。
export default metaNFTModule;