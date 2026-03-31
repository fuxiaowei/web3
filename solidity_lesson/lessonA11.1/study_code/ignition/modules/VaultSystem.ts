import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

/**
 * @title VaultSystem部署模块
 * @dev 演示多合约系统的部署，展示Ignition的依赖管理能力
 * @notice Vault合约依赖于Token合约，Ignition会自动处理依赖关系
 */
export default buildModule("VaultSystem", (m) => {
  // 先部署Token合约
  // Ignition会自动识别这是Vault的依赖
  const token = m.contract("Token");
  
  // 然后部署Vault合约，传入Token地址作为构造函数参数
  // Ignition会确保Token先部署，然后使用Token的地址部署Vault
  const vault = m.contract("Vault", {
    args: [token],
  });
  
  // 可选：给部署者转账一些Token
  // 这展示了如何在部署后执行额外的操作
  const deployer = m.getAccount(0);
  m.call(token, "transfer", [deployer, 1000n]);
  
  // 返回部署的合约引用
  return { token, vault };
});

