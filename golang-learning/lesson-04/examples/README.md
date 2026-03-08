# 课程目录

## 第一阶段: Ethereum 基础与 Go-Eth Client 入门

### 1.1 Ethereum 基础架构

**状态:** ✅  
**学时:** 20分钟  
**教学内容:**
- 账户模型：外部账户（EOA）与合约账户
- 交易结构：nonce、gasPrice / maxFeePerGas、gasLimit、to、value、data 等字段
- 区块结构：header / body、区块链状态、叔块
- Gas 机制：Gas 消耗、Gas 价格、EIP-1559 费用模型
- 智能合约简介：ABI、EVM、部署与调用流程
- 网络类型：主网、测试网（Goerli / Sepolia 等）、本地开发链（geth / anvil / ganache）

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_694e9bcce4b0694ca16111c5?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [Ethereum 基础与 Go-Eth Client 入门(课件).md](./courseware/basics.md#2-ethereum-基础架构)

---

### 1.2 Go-Eth Client 安装与连接

**状态:** ✅  
**学时:** 15分钟  
**教学内容:**
- 安装 `go-ethereum` 库与基础依赖
- 连接本地节点（`geth` / `erigon` 等）
- 连接远程节点（如 Infura / Alchemy / 自建 RPC）
- Client 对象与上下文管理
- 简单连接池与多节点容错思路

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_69649fb7e4b0694ca16b2d3a?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [Ethereum 基础与 Go-Eth Client 入门(课件).md](./courseware/basics.md#3-go-ethereum-客户端安装与连接)

**代码示例:**
- [01-connect-node](./examples/01-connect-node/)
- [10-multi-node-pool](./examples/10-multi-node-pool/)

---

### 1.3 区块操作

**状态:** ✅  
**学时:** 20分钟  
**教学内容:**
- 获取最新区块（By Number / By Hash）
- 获取指定区块、区块区间列表
- 遍历区块中的交易与日志
- 区块信息解析：时间戳、矿工、难度、Gas 使用情况
- 实战：编写小工具持续拉取最新区块

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_69649fb7e4b0694ca16b2d3a?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [Ethereum 基础与 Go-Eth Client 入门(课件).md](./courseware/basics.md#4-区块操作)

**代码示例:**
- [02-block-ops](./examples/02-block-ops/)

---

### 1.4 交易操作

**状态:** ✅  
**学时:** 20分钟  
**教学内容:**
- 按哈希查询交易详情
- 解析交易回执（Receipt）：状态、Gas 使用、日志
- 交易状态查询：pending / mined / failed
- 查询账户余额（ETH / 代币）
- 实战：实现简单的交易追踪脚本

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6964a009e4b0694ca16b2d93?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [Ethereum 基础与 Go-Eth Client 入门(课件).md](./courseware/basics.md#5-交易操作)

**代码示例:**
- [03-tx-ops](./examples/03-tx-ops/)
- [04-account-balance](./examples/04-account-balance/)

---

## 第二阶段: 事件订阅与智能合约交互

### 2.1 事件订阅

**状态:** ✅  
**学时:** 25分钟  
**教学内容:**
- 订阅新区块头（`SubscribeNewHead`）
- 订阅链上日志事件（基于 filter）
- 订阅待处理交易（pending tx）
- 事件订阅的通道模型与错误处理

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6964a01de4b0694ca16b2d9e?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [事件订阅、智能合约交互与工程实践(课件).md](./courseware/advanced.md#2-事件订阅基础)

**代码示例:**
- [05-subscribe-blocks](./examples/05-subscribe-blocks/)
- [06-subscribe-logs](./examples/06-subscribe-logs/)

---

### 2.2 断线重连机制

**状态:** ✅  
**学时:** 20分钟  
**教学内容:**
- 常见失败场景分析
- 基本重连策略设计
- 指数退避算法
- 与手动扫描结合的健壮性设计

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6964a01de4b0694ca16b2d9e?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [事件订阅、智能合约交互与工程实践(课件).md](./courseware/advanced.md#3-断线重连与健壮性设计)

**代码示例:**
- [07-reconnect-strategy](./examples/07-reconnect-strategy/)

---

### 2.3 智能合约交互

**状态:** ✅  
**学时:** 30分钟  
**教学内容:**
- ABI 编码解码：加载 ABI、构造调用数据
- 读取合约方法（`call`）与链下解析返回值
- 发送合约交易（`transact`）：构造、签名、发送
- 解析合约事件日志（topics + data）
- 实战：与 ERC-20 合约交互，查询余额与转账

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6964a034e4b0694ca16b2dbf?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [事件订阅、智能合约交互与工程实践(课件).md](./courseware/advanced.md#4-智能合约交互)

**代码示例:**
- [08-contract-interact](./examples/08-contract-interact/)
- [11-erc20-foundry](./examples/11-erc20-foundry/)

---

## 第三阶段: 综合实战项目

### 3.1 迷你区块浏览器项目实战

**状态:** ✅  
**学时:** 40分钟  
**教学内容:**
- 需求分析：支持区块/交易/地址/事件基础查询
- 数据获取策略：实时 RPC 拉取 vs 本地缓存
- 使用 Go-Eth Client 构建后端服务
- 监听 ERC-20 转账事件并记录
- 提供简单 HTTP API 或命令行界面

**文档资料:**
- 课件: [区块浏览器项目实战(课件).md](./courseware/project.md)

**代码示例:**
- [09-project](./examples/09-project/)

---

## 📁 目录结构

> **重要说明**：每个示例现在都是**独立的工程**，都有自己的 `go.mod`，可以单独运行！

```
lesson-04/
├── README.md                  # 本文档
├── COURSE_OVERVIEW.md         # 课程总览
├── QUICK_START.md             # 快速上手指南
├── DIRECTORY_STRUCTURE.md     # 目录结构说明
├── courseware/                # 课件目录
│   ├── basics.md              # Ethereum 基础与 Go-Eth Client 入门课件
│   ├── advanced.md            # 事件订阅、智能合约交互与工程实践课件
│   └── project.md             # 区块浏览器项目实战课件
└── examples/                  # 代码示例（每个都是独立工程）
    ├── 01-connect-node/       # 连接节点（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 02-block-ops/          # 区块操作（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 03-tx-ops/             # 交易操作（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 04-account-balance/    # 账户余额查询（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 05-subscribe-blocks/   # 订阅新区块（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 06-subscribe-logs/     # 订阅日志事件（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 07-reconnect-strategy/ # 断线重连策略（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 08-contract-interact/  # 智能合约交互（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 09-project/            # 综合项目（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    ├── 10-multi-node-pool/    # 多节点连接池（独立工程）
    │   ├── main.go
    │   ├── go.mod
    │   └── go.sum
    └── 11-erc20-foundry/      # ERC-20 Foundry 示例（独立工程）
        ├── src/
        ├── script/
        ├── foundry.toml
        └── go.mod
```

### 为什么采用独立工程？

**之前的问题**：多个示例文件共享一个 `go.mod`，由于都是 `package main` 且都有 `main()` 函数，会导致冲突无法运行。

**现在的优势**：
- ✅ 每个示例都可以**独立运行**：`cd 01-connect-node && go run main.go`
- ✅ 避免依赖冲突
- ✅ 便于教学演示
- ✅ 便于学员理解和实践

## 🚀 快速开始

### 方式一：运行单个示例（推荐）

每个示例都是独立的工程，可以直接运行：

```bash
# 1. 连接节点
cd lesson-04/examples/01-connect-node
go run main.go

# 2. 区块操作
cd lesson-04/examples/02-block-ops
go run main.go

# 3. 交易操作
cd lesson-04/examples/03-tx-ops
go run main.go

# 4. 账户余额查询
cd lesson-04/examples/04-account-balance
go run main.go

# 5. 订阅新区块
cd lesson-04/examples/05-subscribe-blocks
go run main.go

# 6. 订阅日志事件
cd lesson-04/examples/06-subscribe-logs
go run main.go

# 7. 断线重连策略
cd lesson-04/examples/07-reconnect-strategy
go run main.go

# 8. 智能合约交互
cd lesson-04/examples/08-contract-interact
go run main.go

# 9. 多节点连接池
cd lesson-04/examples/10-multi-node-pool
go run main.go
```

### 方式二：编译后运行

```bash
cd lesson-04/examples/01-connect-node
go build
./01-connect-node
```

### 运行综合项目

```bash
cd lesson-04/examples/09-project
go run main.go
```

### 环境配置

运行前需要设置以太坊节点 RPC 地址：

```bash
# 本地节点
export ETH_RPC_URL="http://127.0.0.1:8545"

# 远程节点（测试网）
export ETH_RPC_URL="https://sepolia.infura.io/v3/your-api-key"

# WebSocket 连接（用于事件订阅）
export ETH_WS_URL="wss://sepolia.infura.io/v3/your-api-key"
```

## 📚 示例代码与课件对应关系

| 示例文件 | 对应课件章节 | 主要知识点 | 运行命令 |
|---------|------------|-----------|---------|
| **基础示例** | | | |
| `01-connect-node` | 第二部分：Go-Eth Client 安装与连接 | 连接节点、获取链ID、最新区块号 | `cd examples/01-connect-node && go run main.go` |
| `10-multi-node-pool` | 第二部分：Go-Eth Client 安装与连接 | 多节点连接池、读写分离、容错 | `cd examples/10-multi-node-pool && go run main.go` |
| `02-block-ops` | 第三部分：区块操作 | 区块查询、区块信息解析 | `cd examples/02-block-ops && go run main.go` |
| `03-tx-ops` | 第四部分：交易操作 | 交易查询、交易回执解析 | `cd examples/03-tx-ops && go run main.go` |
| `04-account-balance` | 第四部分：交易操作 | 账户余额查询（ETH/代币） | `cd examples/04-account-balance && go run main.go` |
| **进阶示例** | | | |
| `05-subscribe-blocks` | 第一部分：事件订阅 | 订阅新区块头 | `cd examples/05-subscribe-blocks && go run main.go` |
| `06-subscribe-logs` | 第一部分：事件订阅 | 订阅日志事件 | `cd examples/06-subscribe-logs && go run main.go` |
| `07-reconnect-strategy` | 第二部分：断线重连机制 | 重连策略、错误处理 | `cd examples/07-reconnect-strategy && go run main.go` |
| `08-contract-interact` | 第三部分：智能合约交互 | ABI编解码、合约调用、事件解析 | `cd examples/08-contract-interact && go run main.go` |
| `11-erc20-foundry` | 第三部分：智能合约交互 | ERC-20 合约交互实战 | `cd examples/11-erc20-foundry && 参考 README.md` |
| **项目实战** | | | |
| `09-project` | 迷你区块浏览器项目实战 | 区块浏览器后端服务 | `cd examples/09-project && go run main.go` |

## 🎯 学习目标

完成本课程后，你应该能够：

1. ✅ 理解以太坊的核心概念：账户、交易、区块、Gas、智能合约与网络类型
2. ✅ 使用 Go-Eth Client 连接本地/远程节点，查询区块与交易数据
3. ✅ 实现交易、区块、账户余额等常用 RPC 查询逻辑
4. ✅ 订阅新区块、pending 交易与合约日志事件，并处理断线重连
5. ✅ 通过 ABI 与 Go 代码进行智能合约读写与事件解析
6. ✅ 基于 Go-Eth Client 实现一个简单的区块浏览器或链上监控服务

## 💡 实践建议

### 课后练习

1. **基础练习：**
   - 实现一个区块监控工具，持续打印最新区块信息
   - 实现一个交易追踪工具，监控指定地址的所有交易

2. **进阶练习：**
   - 实现一个 ERC-20 转账监听服务
   - 实现一个多链数据采集服务
   - 实现一个 Gas 价格监控工具

3. **综合练习：**
   - 构建完整的区块浏览器后端（区块/交易/地址查询）
   - 实现 DeFi 协议事件监听（如 Uniswap Swap 事件）
   - 添加数据缓存层，提高查询性能

### 扩展阅读

- [Go-Ethereum 官方文档](https://geth.ethereum.org/)
- [以太坊开发者文档](https://ethereum.org/developers/)
- [Web3.py 文档](https://web3py.readthedocs.io/)（Python 参考）
- [Ethers.js 文档](https://docs.ethers.org/)（JavaScript 参考）
- [EIP-1559 提案](https://eips.ethereum.org/EIPS/eip-1559)

## 📞 常见问题

### Q: 连接节点时提示超时或 403？
A: 检查 `ETH_RPC_URL` 是否正确、网络是否连通，若使用远程服务（Infura 等）需确认 API Key 是否配置正确、项目是否启用了对应网络。

### Q: 区块高度/交易内容为空？
A: 确认当前网络是否有最新区块（本地链是否正在出块），对于测试网/主网可以在区块浏览器对比高度是否一致。

### Q: 日志订阅收不到事件？
A: 检查 filter 是否正确（合约地址、topics）、节点是否支持 `eth_subscribe`，以及合约事件是否在当前网络上真正发生。

### Q: 订阅功能需要 WebSocket 连接吗？
A: 是的，订阅功能（如 `SubscribeNewHead`、`SubscribeFilterLogs`）通常需要 WebSocket 连接（`ws://` 或 `wss://`），普通的 HTTP RPC 连接不支持订阅。

### Q: 如何选择本地节点还是远程节点？
A: 开发测试推荐使用本地节点（如 `anvil`、`ganache`），生产环境建议使用专业的 RPC 服务商（如 Infura、Alchemy）或自建节点。

## 🔗 相关资源

- [课件文档](./courseware/)
- [代码示例](./examples/)
- [Go-Ethereum GitHub 仓库](https://github.com/ethereum/go-ethereum)
- [以太坊官方文档](https://ethereum.org/developers/)

## 📅 课程计划

- **第1周**：Ethereum 基础与 Go-Eth Client 入门学习
- **第2周**：事件订阅与智能合约交互学习
- **第3周**：项目实践和复习

## 👨‍💻 作者

本课程由Metaland团队制作，如有问题欢迎反馈。

## 📄 许可证

MIT License

---

**祝学习愉快！** 🎉
