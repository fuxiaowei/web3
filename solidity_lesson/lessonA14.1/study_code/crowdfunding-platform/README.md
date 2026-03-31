# 众筹平台 - Hardhat 3 实战项目

一个基于 Hardhat 3 开发的去中心化众筹平台，完整展示了智能合约开发、测试和部署的全流程。本项目采用**工厂模式**和**状态机模式**，确保了系统的可扩展性与资金安全性。

## 项目概述

这是一个生产级的众筹平台实战项目，旨在解决传统众筹透明度低、手续费高的问题。通过智能合约，我们将募资逻辑、资金托管和退款机制彻底去中心化。

### 核心特性
- 多项目管理：通过工厂合约，任何用户都可以一键发布自己的众筹项目。
- 状态机驱动：严格管理项目的五个生命周期阶段（准备、活跃、成功、失败、关闭）。
- 资金安全：遵循“检查-影响-交互（CEI）”设计模式，彻底杜绝重入攻击。
- 透明度：所有筹款进度、参与者记录和资金流向均在链上可查。

---

## 快速开始

### 环境要求
- Node.js >= 22.0.0
- npm 或 yarn

### 安装与编译
```bash
# 1. 克隆并进入目录
# 2. 安装依赖
npm install

# 3. 编译合约
npm run compile
```

---

## 核心功能模块详解

### 1. 众筹工厂合约 (CrowdfundingFactory.sol)
工厂合约是整个平台的管理中心，负责创建和索引所有的众筹项目。
- 功能：
    - `createCampaign(goal, duration)`：部署一个新的众筹项目合约实例。
    - `getCampaigns()`：返回平台上所有已创建项目的地址列表。
    - `getUserCampaigns(user)`：查询特定用户创建的所有项目，方便前端仪表盘展示。
- 设计价值：实现了合约的“按需生产”，且每个项目合约地址独立，安全性相互隔离。

### 2. 众筹项目合约 (CrowdfundingCampaign.sol)
这是处理具体募资业务逻辑的核心合约。
- 核心逻辑：
    - `contribute()`：用户通过此函数发送 ETH 参与众筹。如果达到目标金额，状态自动切换为 `Success`。
    - `finalize()`：截止时间到达后，如果未达标，可调用此函数将状态切换为 `Failed`。
    - `withdraw()`：项目成功后，创建者提取资金。包含严格的重入防护。
    - `refund()`：项目失败后，贡献者自助领取退款。
- 内部状态机：
    - `Preparing` -> `Active` -> `Success`/`Failed` -> `Closed`。

---

## 运行示例与操作流程

### 1. 本地模拟环境运行
您可以在本地 Hardhat 网络中完整模拟从部署到募资成功的全过程。

```bash
# 启动本地节点（在一个终端窗口）
npm run node

# 执行部署脚本（在另一个终端窗口）
# 该脚本会自动部署工厂合约，并通过工厂创建一个初始项目
npm run deploy:local
```

### 2. 交互操作示例（使用 Hardhat Console）
```bash
# 进入交互式控制台
npx hardhat console --network localhost

# 获取合约实例并参与众筹
const factory = await ethers.getContractAt("CrowdfundingFactory", "部署的工厂地址");
const campaigns = await factory.getCampaigns();
const campaign = await ethers.getContractAt("CrowdfundingCampaign", campaigns[0]);

# 用户注资 1 ETH
await campaign.contribute({ value: ethers.parseEther("1.0") });
```

### 3. 自动化测试
本项目包含近 100% 覆盖率的测试用例，涵盖了所有边缘情况。
```bash
# 运行标准测试
npm test

# 运行测试并生成 Gas 报告
npm run test:gas
```

---

## 项目结构

```
crowdfunding-platform/
├── contracts/              # 智能合约源码
│   ├── CrowdfundingCampaign.sol    # 业务逻辑合约
│   └── CrowdfundingFactory.sol     # 工厂管理合约
├── test/                   # 基于 Chai 的完整测试集
│   ├── CrowdfundingCampaign.test.ts
│   └── CrowdfundingFactory.test.ts
├── ignition/              # Hardhat Ignition 声明式部署模块
│   └── modules/
│       └── Crowdfunding.ts # 部署逻辑定义
├── hardhat.config.ts       # Hardhat 工程配置
└── package.json            # 脚本与依赖管理
```

---

## 安全最佳实践
- CEI 模式：所有资金变动函数均采用“先修改状态，后转账”的原则，防止重入攻击。
- Immutable 优化：关键参数（如目标金额、截止日期）使用 `immutable` 存储，既安全又节省 Gas。
- 输入验证：构造函数和公开函数均包含详尽的 `require` 断言，防止无效参数输入。

---

## 学习资源
- [Hardhat 官方文档](https://hardhat.org/docs)
- [Solidity 状态机模式指南](https://docs.soliditylang.org/en/v0.8.24/common-patterns.html#state-machine)

---

## 许可证
MIT License
