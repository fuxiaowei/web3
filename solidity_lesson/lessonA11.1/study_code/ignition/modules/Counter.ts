import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

/**
 * @title Counter部署模块
 * @dev 演示如何使用Hardhat Ignition进行声明式部署
 * @notice Ignition提供了自动状态管理、错误恢复和依赖处理
 */
export default buildModule("CounterModule", (m) => {
  // 部署Counter合约
  const counter = m.contract("Counter");
  
  // 部署后调用incBy方法，初始化值为10
  // 这展示了如何在部署后自动执行合约方法
  m.call(counter, "incBy", [10n]);
  
  // 返回部署的合约引用
  // 这些引用可以在其他模块中使用
  return { counter };
});

