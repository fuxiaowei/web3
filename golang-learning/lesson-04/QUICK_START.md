# 第4期快速上手指南：Go-Eth Client 与以太坊基础

## 1. 环境准备

### ✅ 必备工具

- Go 1.22+  
- Git  
- 一个可用的以太坊节点 RPC 地址（满足其一即可）：  
  - 本地 `geth` 节点（推荐开发用 `--dev` 或测试网模式）  
  - 远程节点服务（如 Infura / Alchemy / 自建 RPC 网关）  

### ✅ 可选工具

- VS Code / GoLand + Go 插件  
- curl 或 httpie（命令行测试 HTTP 接口）  
- Postman / Bruno / Insomnia（如你计划为项目加 HTTP API）  

## 2. 准备示例代码

本期示例将使用 `go-ethereum` 提供的 `ethclient` 包进行开发。

> 代码示例目录结构预留为：
>
> - `examples/basics`：基础区块 / 交易 / 账户查询  
> - `examples/advanced`：事件订阅、合约交互  
> - `examples/project`：迷你区块浏览器项目  

进入第4期示例目录：

```bash
cd /Users/zhouxing/Documents/rcc/new-golang/lesson-04/examples
```

在后续补充代码时，建议：

- 为每个子目录（如 `basics/`、`advanced/`、`project/`）单独创建 `go.mod`，保持和第3期 Gin 示例相同的**独立工程**风格  
- 或者在 `examples/` 下创建统一的 `go.mod`，通过子包组织示例代码  

## 3. 配置以太坊节点 RPC

建议通过环境变量配置 RPC 地址，便于在本地链、测试网与主网之间切换：

```bash
export ETH_RPC_URL="http://127.0.0.1:8545"   # 本地 geth / anvil / ganache
# 或者
export ETH_RPC_URL="https://sepolia.infura.io/v3/your-key"
```

在示例代码中读取该环境变量：

```go
rpcURL := os.Getenv("ETH_RPC_URL")
if rpcURL == "" {
    log.Fatal("ETH_RPC_URL is not set")
}
```

## 4. 运行基础连接与区块查询示例（示意）

> 下面是你在 `examples/basics` 中可以实现的典型运行方式，具体文件可按课程进度逐步补充。

### 4.1 连接节点并获取链 ID

```bash
cd examples/basics
go run 01-connect-node.go
```

预期行为：
- 输出当前连接的链 ID（如 1=主网，11155111=Sepolia 等）  
- 输出最新区块高度  

### 4.2 查询最新区块与指定区块

```bash
go run 02-block-ops.go
```

预期行为：
- 打印最近一个区块的高度、哈希、时间戳、交易数量  
- 根据给定区块号（或从命令行参数读取）打印指定区块的关键信息  

## 5. 运行交易与账户示例（示意）

### 5.1 按哈希查询交易与回执

```bash
go run 03-tx-ops.go \
  --tx 0x你的交易哈希
```

预期行为：
- 打印交易的 from/to、value、gas、input data 等字段  
- 打印回执中的 status、gasUsed、logs 数量等信息  

### 5.2 查询账户余额

```bash
go run 04-account-balance.go \
  --address 0x你的地址
```

预期行为：
- 打印该地址的 ETH 余额（以 ETH 与 Wei 两种单位展示）  
- 如后续补充 ERC-20 功能，可增加 `--token` 参数查询代币余额  

## 6. 运行事件订阅与合约交互（示意）

### 6.1 订阅新区块头

```bash
cd ../advanced
go run 01-subscribe-blocks.go
```

预期行为：
- 当有新区块产生时，在终端打印区块高度与哈希  
- 在网络断开 / 节点重启时能优雅处理错误，必要时自动重连  

### 6.2 订阅日志事件（如 ERC-20 Transfer）

```bash
go run 03-subscribe-logs.go \
  --contract 0xERC20合约地址
```

预期行为：
- 打印匹配到的 `Transfer` 事件中的 from、to、value 信息  

### 6.3 调用合约方法与发送交易

```bash
go run 05-contract-interact.go \
  --contract 0xERC20合约地址 \
  --method balanceOf \
  --address 0x你的地址
```

预期行为：
- 通过 ABI 编解码调用智能合约方法，并打印返回值  

## 7. 项目实战：迷你区块浏览器（示意）

在 `examples/project` 中，你可以逐步实现：

```bash
cd ../project
go run main.go
```

预期功能（可按阶段逐步实现）：

- 提供 HTTP API 或 CLI 命令  
- 支持按区块号/哈希查询区块  
- 按哈希查询交易与回执  
- 按地址查看最近交易与余额  
- 后台协程持续监听区块 / ERC-20 转账事件并缓存最近记录  

## 8. 常见问题

### Q: 连接节点时提示超时或 403？
A: 检查 `ETH_RPC_URL` 是否正确、网络是否连通，若使用远程服务（Infura 等）需确认 API Key 是否配置正确、项目是否启用了对应网络。

### Q: 区块高度/交易内容为空？
A: 确认当前网络是否有最新区块（本地链是否正在出块），对于测试网/主网可以在区块浏览器对比高度是否一致。

### Q: 日志订阅收不到事件？
A: 检查 filter 是否正确（合约地址、topics）、节点是否支持 `eth_subscribe`，以及合约事件是否在当前网络上真正发生。

## 9. 下一步

- 按顺序阅读并实现 `courseware/basics.md` 与 `courseware/advanced.md` 中的代码示例  
- 在 `examples/` 中补全对应示例文件，并跑通基础功能  
- 按 `courseware/project.md` 搭建自己的迷你区块浏览器或链上监听服务  

---

**完成以上步骤，你就可以顺利开始第4期 Go-Eth Client 与以太坊基础的学习与实战了！** 🚀


