# Hardhat 3 项目示例

这是一个完整的Hardhat 3项目示例，展示了从环境搭建到合约部署的完整流程。

## 项目结构

```
.
├── contracts/              # Solidity合约目录
│   ├── Counter.sol        # 计数器合约示例
│   ├── Token.sol          # 代币合约示例
│   └── Vault.sol          # 保险库合约示例
├── scripts/                # 部署和交互脚本
│   └── deploy.ts          # 传统脚本部署示例
├── test/                   # 测试文件目录
│   └── Counter.test.ts    # Counter合约测试
├── ignition/               # Ignition部署模块（Hardhat 3新增）
│   └── modules/           # 部署模块定义
│       ├── Counter.ts     # Counter部署模块
│       └── VaultSystem.ts # 多合约系统部署模块
├── hardhat.config.ts      # Hardhat配置文件
├── package.json           # npm配置文件
└── .env.example           # 环境变量示例文件
```

## 环境要求

- **Node.js**: >= 22.0.0（硬性要求）
- **npm**: >= 10.0.0（随Node.js自动安装）

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量

复制`.env.example`为`.env`并填入实际值：

```bash
cp .env.example .env
```

编辑`.env`文件，填入：
- RPC节点URL（从Infura或Alchemy获取）
- 私钥（测试账户私钥）
- API密钥（Etherscan、Coinmarketcap等）

### 3. 编译合约

```bash
npm run compile
```

或使用Hardhat命令：

```bash
npx hardhat compile
```

### 4. 运行测试

```bash
npm run test
```

运行测试并生成Gas报告：

```bash
npm run test:gas
```

### 5. 启动本地节点

```bash
npm run node
```

这会启动一个本地Hardhat网络，提供HTTP和WebSocket JSON-RPC服务。

### 6. 部署合约

**使用传统脚本部署**：

```bash
npm run deploy
```

部署到Sepolia测试网：

```bash
npm run deploy:sepolia
```

**使用Ignition部署**：

```bash
npm run ignition:deploy
```

部署到Sepolia测试网：

```bash
npm run ignition:deploy:sepolia
```

查看部署状态：

```bash
npm run ignition:status
```

## 合约说明

### Counter.sol

一个简单的计数器合约，演示了：
- 状态变量的定义和使用
- 函数的定义和调用
- 事件的触发

### Token.sol

一个简单的代币合约，实现了基础的转账功能。

### Vault.sol

一个保险库合约，依赖于Token合约，演示了：
- 合约之间的依赖关系
- 构造函数参数的使用

## 测试说明

### Counter.test.ts

包含以下测试用例：
- 设置数字功能测试
- 事件触发测试
- 递增功能测试
- 多次递增测试

## 部署方式

### 传统脚本部署

使用`scripts/deploy.ts`进行部署，适合简单的部署场景。

**优点**：
- 简单直接
- 控制灵活
- 适合一次性部署

**缺点**：
- 需要手动管理状态
- 错误恢复困难
- 依赖关系需要手动处理

### Ignition部署

使用`ignition/modules/`下的模块进行部署，适合复杂的系统部署。

**优点**：
- 自动状态管理
- 错误恢复支持
- 依赖自动处理
- 可重复部署

**缺点**：
- 学习曲线稍陡
- 适合复杂系统

## 常用命令

### 编译相关

```bash
# 编译所有合约
npx hardhat compile

# 清理缓存
npx hardhat clean
```

### 测试相关

```bash
# 运行所有测试
npx hardhat test

# 运行特定测试文件
npx hardhat test test/Counter.test.ts

# 显示Gas报告
REPORT_GAS=true npx hardhat test
```

### 部署相关

```bash
# 传统脚本部署
npx hardhat run scripts/deploy.ts
npx hardhat run scripts/deploy.ts --network sepolia

# Ignition部署
npx hardhat ignition deploy ignition/modules/Counter.ts
npx hardhat ignition deploy ignition/modules/Counter.ts --network sepolia

# 查看部署状态
npx hardhat ignition status chain-11155111
```

### 网络相关

```bash
# 启动本地节点
npx hardhat node

# 打开Hardhat控制台
npx hardhat console
```

## 网络配置

### 本地网络

- **类型**: `edr-simulated`
- **链ID**: 31337
- **RPC URL**: 自动提供
- **账户**: 10个测试账户，每个10000 ETH

### Sepolia测试网

- **类型**: `http`
- **链ID**: 11155111
- **RPC URL**: 从环境变量`SEPOLIA_URL`读取
- **账户**: 从环境变量`PRIVATE_KEY`读取

### 主网

- **类型**: `http`
- **链ID**: 1
- **RPC URL**: 从环境变量`MAINNET_URL`读取
- **账户**: 从环境变量`PRIVATE_KEY`读取
- **警告**: 主网使用真实资金，部署前务必充分测试

## 安全注意事项

1. **私钥安全**：
   - 永远不要提交私钥到Git
   - 使用环境变量管理私钥
   - 为不同环境使用不同的私钥

2. **测试网优先**：
   - 先在测试网上充分测试
   - 确认无误后再部署到主网

3. **备份重要信息**：
   - 备份私钥和助记词
   - 记录部署地址和交易哈希

4. **验证合约**：
   - 部署后立即验证合约
   - 确保源码透明

## 故障排除

### Node.js版本问题

如果遇到版本不兼容问题：

```bash
# 检查Node.js版本
node --version

# 使用nvm切换版本
nvm install 22
nvm use 22
```

### 编译错误

如果编译失败：

```bash
# 清理缓存
npx hardhat clean

# 重新编译
npx hardhat compile
```

### 网络连接问题

如果无法连接到网络：

1. 检查RPC URL是否正确
2. 检查网络配置中的type字段
3. 确认环境变量已正确设置

### 部署失败

如果部署失败：

1. 检查账户余额是否足够
2. 检查Gas限制设置
3. 查看详细错误信息：`--show-stack-traces`

## 学习资源

- [Hardhat官方文档](https://hardhat.org/docs)
- [Hardhat 3新特性](https://hardhat.org/hardhat-runner/docs/guides/whats-new)
- [Ignition部署指南](https://hardhat.org/ignition/docs)
- [Ethers.js文档](https://docs.ethers.org/)

## 许可证

MIT

