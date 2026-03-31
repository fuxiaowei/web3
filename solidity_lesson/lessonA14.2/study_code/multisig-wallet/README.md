# 多签钱包项目 (Multi-Signature Wallet)

一个基于 Hardhat 3 开发的多签钱包智能合约项目，支持透明代理模式的合约升级。

## 项目简介

多签钱包是一个需要多个所有者确认才能执行交易的智能合约钱包。本项目实现了完整的多签钱包功能，包括所有者管理、交易提案、确认机制和执行交易等核心功能。

## 功能特性

- ✅ 所有者管理：添加/删除所有者，修改确认阈值
- ✅ 交易提案：创建和查询交易提案
- ✅ 确认机制：确认和撤销确认交易提案
- ✅ 执行交易：在达到阈值后执行交易
- ✅ 接收 ETH：支持接收以太币
- ✅ 安全防护：重入攻击防护、权限控制、输入验证
- ✅ **合约升级**：支持透明代理模式的合约升级

## 项目结构

```
multisig-wallet/
├── contracts/
│   ├── MultiSigWallet.sol              # 原始不可升级版本
│   ├── MultiSigWalletUpgradeable.sol   # 可升级版本（V1）
│   └── MultiSigWalletV2.sol            # 升级版本（V2，演示用）
├── test/
│   └── MultiSigWallet.test.ts          # 测试文件
├── scripts/
│   ├── deploy.ts                       # 传统部署脚本
│   ├── deployWithProxy.ts              # 使用透明代理部署（推荐）
│   ├── deployWithIgnition.ts           # 结合 Ignition 部署
│   ├── upgrade.ts                      # 升级合约脚本
│   └── prepareUpgrade.ts               # 准备升级（验证）
├── ignition/
│   └── modules/
│       └── DeployMultiSigWallet.ts     # Ignition 部署模块
├── hardhat.config.ts                   # Hardhat 配置
├── package.json                        # 项目依赖
├── tsconfig.json                       # TypeScript 配置
├── README.md                           # 项目说明
└── README_UPGRADE.md                   # 升级指南
```

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 编译合约

```bash
npm run compile
```

### 3. 运行测试

```bash
npm run test
```

### 4. 部署合约（使用透明代理）

```bash
# 推荐方式：使用代理部署（支持升级）
npm run deploy:proxy

# 或使用 Ignition（需要额外配置）
npm run deploy:ignition
```

**📖 详细部署指南请参阅 [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)**

## 部署方式

### 方式一：使用透明代理部署（推荐）

```bash
npm run deploy:proxy
# 或
npx hardhat run scripts/deployWithProxy.ts --network localhost
```

这会部署：
- 实现合约（Implementation）
- 代理合约（Proxy）
- ProxyAdmin

### 方式二：使用 Ignition 部署

```bash
npm run deploy:ignition
# 或
npx hardhat ignition deploy ignition/modules/DeployMultiSigWallet.ts --network localhost
```

## 合约升级

### 准备升级（验证）

```bash
PROXY_ADDRESS=0x... npm run prepare:upgrade
```

### 执行升级

```bash
PROXY_ADDRESS=0x... npm run upgrade
```

**📖 详细部署和升级指南请参阅：**
- [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) - **完整部署和升级流程（推荐阅读）**
- [README_UPGRADE.md](./README_UPGRADE.md) - 升级技术细节

## 使用说明

### 部署多签钱包

部署时需要指定所有者和确认阈值：

```typescript
const owners = [owner1.address, owner2.address, owner3.address];
const numConfirmationsRequired = 2; // 需要 2 个确认
const wallet = await MultiSigWallet.deploy(owners, numConfirmationsRequired);
```

### 创建交易提案

```typescript
await wallet.submitTransaction(
  recipient.address,      // 目标地址
  ethers.parseEther("1"), // 转账金额
  "0x"                    // 调用数据
);
```

### 确认提案

```typescript
await wallet.connect(owner1).confirmTransaction(0);
await wallet.connect(owner2).confirmTransaction(0);
```

### 执行交易

当确认数达到阈值后，可以执行交易：

```typescript
await wallet.connect(owner1).executeTransaction(0);
```

### 所有者管理

```typescript
// 添加所有者
await wallet.addOwner(newOwner.address);

// 删除所有者
await wallet.removeOwner(owner.address);

// 修改阈值
await wallet.changeThreshold(3);
```

## 测试

项目包含完整的测试套件，涵盖以下场景：

- ✅ 部署测试：验证合约初始化
- ✅ 所有者管理测试：添加/删除所有者，修改阈值
- ✅ 提案功能测试：创建和查询提案
- ✅ 确认机制测试：确认和撤销确认
- ✅ 执行交易测试：执行 ETH 转账
- ✅ 接收 ETH 测试：接收以太币

运行所有测试：

```bash
npx hardhat test
```

## 技术栈 

- **Solidity**: 0.8.24
- **Hardhat**: 3.0.0
- **Hardhat Ignition**: 3.0.0
- **OpenZeppelin Contracts**: 5.0.0
- **OpenZeppelin Upgrades**: 3.0.0
- **TypeScript**: Latest
- **Ethers.js**: v6 (via Hardhat Toolbox)

## 安全考虑

本项目实现了以下安全措施：

1. **重入攻击防护**：使用 Checks-Effects-Interactions 模式
2. **权限控制**：onlyOwner 修饰符限制访问
3. **输入验证**：验证所有输入参数
4. **状态检查**：防止重复执行和重复确认
5. **使用最新 Solidity 版本**：利用内置安全检查
6. **代理模式安全**：使用 OpenZeppelin 标准代理实现

## 升级相关

### 存储布局规则

升级时必须遵循存储布局规则：
- ✅ 可以添加新的状态变量（在末尾）
- ❌ 不能删除状态变量
- ❌ 不能改变变量类型
- ❌ 不能改变变量顺序

详细说明请参阅 [README_UPGRADE.md](./README_UPGRADE.md)

## 许可证

MIT License

## 参考资源

- [OpenZeppelin Upgrades 文档](https://docs.openzeppelin.com/upgrades-plugins/1.x/)
- [Hardhat Ignition 文档](https://hardhat.org/ignition/docs)
- [代理模式详解](https://docs.openzeppelin.com/upgrades-plugins/1.x/proxies)

## 贡献

欢迎提交 Issue 和 Pull Request！
