# 第4期目录结构说明：区块链基础与 Go-Eth Client

```
lesson-04/
├── README.md                  # 课程简介与学习指南
├── COURSE_OVERVIEW.md         # 课程总览（目标、内容、收获）
├── QUICK_START.md             # 快速上手指南
├── DIRECTORY_STRUCTURE.md     # 本文档
├── courseware/                # 课件资料
│   ├── basics.md              # Ethereum 基础与 Go-Eth Client 入门课件大纲
│   ├── advanced.md            # 事件订阅、合约交互与工程实践课件大纲
│   └── project.md             # 区块浏览器项目实战课件大纲
└── examples/                  # 代码示例（可按需要拆分为独立工程）
    ├── basics/                # 基础示例：区块/交易/账户查询
    ├── advanced/              # 进阶示例：事件订阅、合约交互
    └── project/               # 综合项目：迷你区块浏览器
```

## 目录使用说明

- `courseware/`：用于录制课程与备课的详细大纲，内容覆盖从 Ethereum 基础到 Go-Eth Client、事件订阅与项目实战。  
- `examples/`：配套代码示例目录，建议延续第3期的做法，将每个子目录设计为**独立的 Go 工程**，便于单独运行与演示。  

### `courseware/` 设计

- `basics.md`：侧重概念与基础 API 使用，如账户/交易/区块结构、Gas 机制、网络类型、`ethclient.Client` 的基本用法。  
- `advanced.md`：侧重事件订阅、合约交互、ABI 编解码、断线重连与错误处理等工程实践。  
- `project.md`：围绕“迷你区块浏览器 / ERC-20 监听服务”给出项目结构、模块拆分与实现步骤。  

### `examples/` 设计建议

> 可根据实际教学节奏逐步补充代码，下面是推荐的组织方式。

#### 1. 基础示例（`examples/basics`）

- `01-connect-node.go`：连接节点、获取链 ID 与最新区块号  
- `02-block-ops.go`：查询最新区块、指定区块、区块区间  
- `03-tx-ops.go`：按哈希查询交易与回执，解析关键字段  
- `04-account-balance.go`：查询账户 ETH 余额与（可选）ERC-20 余额  

#### 2. 进阶示例（`examples/advanced`）

- `01-subscribe-blocks.go`：订阅新区块头  
- `02-subscribe-pending.go`：订阅 pending 交易  
- `03-subscribe-logs.go`：订阅合约日志事件（如 ERC-20 Transfer）  
- `04-reconnect-strategy.go`：订阅断线后的重连与恢复示例  
- `05-contract-interact.go`：通过 ABI 调用合约方法、发送交易、解析事件  

#### 3. 综合项目（`examples/project`）

- `main.go`：区块浏览器 / 监听服务入口  
- 其他文件可按功能拆分为：
  - `cmd/`：命令行入口（可选）  
  - `internal/client/`：以太坊客户端封装  
  - `internal/service/`：区块 / 交易 / 地址查询、事件监听逻辑  
  - `internal/api/`：HTTP API 层（如采用 Gin / 标准库）  
  - `internal/store/`：缓存或持久化实现（内存 / 文件 / DB）  

## 设计理念

### 独立示例，便于教学

参考第3期 Gin 的实践经验，推荐：

- 为 `examples/basics`、`examples/advanced`、`examples/project` 下的示例**分别创建独立的 `go.mod`**  
- 每个示例可以直接 `cd` 进入目录，然后使用 `go run` 运行  
- 这样可以避免多个 `main` 函数导致的冲突，也便于针对单个示例进行讲解和调试  

### 聚焦“链上读写能力”

本期课程的代码结构以“链上数据读取 + 事件监听 + 合约交互”为核心，目录划分也围绕这三类能力展开，便于学员从基础查询逐步过渡到事件驱动与项目实战。  

---

**按照本结构组织第4期课程材料，可以保持与前几期一致的风格，同时突出区块链领域的特点与实践重点。** 🚀


