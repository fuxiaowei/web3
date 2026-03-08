# 第4期：Ethereum 基础与 Go-Eth Client 入门

本讲义将带你从以太坊核心概念出发，一步步学习如何使用 Go 语言与 `go-ethereum` 提供的客户端进行区块、交易和账户相关的开发实践。

## 1. 学习目标

通过本部分学习，学员将能够：

- 理解以太坊的账户模型、交易与区块基本结构  
- 说清楚 Gas 与手续费的构成，理解 EIP-1559 之后的费用模型  
- 区分主网 / 测试网 / 本地私有链以及 RPC 节点的角色  
- 安装并使用 `go-ethereum` 中的 `ethclient` 连接节点  
- 编写简单程序查询区块、交易与账户余额  
- 实现 ETH 转账交易的构造、签名与发送  

## 2. Ethereum 基础架构

### 2.1 账户模型

- 外部账户（EOA，Externally Owned Account）  
  - 由私钥控制  
  - 常见字段：地址、余额、nonce  
- 合约账户（Contract Account）  
  - 没有私钥，由合约代码控制  
  - 由部署交易创建，代码存储在状态中  

你可以结合 EVM 账户状态图来帮助理解账户之间的关系。

### 2.2 交易结构

- 核心字段（按直观顺序讲解）：
  - `nonce`：账户已发出交易数  
  - `to`：接收方地址（部署合约时为空）  
  - `value`：转账金额（ETH）  
  - `data`：调用合约时的编码数据（函数选择器 + 参数）  
  - `gasLimit`：愿意为本次执行支付的最大 Gas 数量  
  - `gasPrice` / `maxFeePerGas` / `maxPriorityFeePerGas`：Gas 单价相关字段  

你可以在区块浏览器中选取一笔真实交易，对照这些字段来加深理解。

### 2.3 ETH 交易的生命周期

一笔以太坊交易从创建到最终确认，会经历以下几个阶段：

#### 阶段 1：创建与签名（本地）
- **构造交易对象**：填充 `nonce`、`to`、`value`、`data`、`gasLimit`、`gasPrice` 等字段
- **使用私钥签名**：使用发送方账户的私钥对交易进行 ECDSA 签名
- **生成交易哈希**：签名后的交易会生成唯一的交易哈希（`txHash`），用于后续查询和追踪

#### 阶段 2：广播到网络（Pending）
- **发送到节点**：通过 RPC 调用（如 `eth_sendRawTransaction`）将签名后的交易发送到以太坊节点
- **进入交易池（Mempool）**：节点验证交易有效性（签名、nonce、余额等）后，将交易放入本地交易池
- **网络传播**：节点将交易广播给其他对等节点，交易在网络中传播
- **状态：Pending**：此时交易尚未被打包进区块，可以通过 `eth_getTransactionByHash` 查询到交易，但 `eth_getTransactionReceipt` 返回 `null`

#### 阶段 3：打包进区块（Mined）
- **验证者选择**：验证者（Validator）从交易池中选择交易打包进新区块
- **执行交易**：在区块中执行交易，更新账户状态（余额、nonce 等）
- **生成交易回执**：执行完成后生成交易回执（Receipt），包含执行结果、Gas 消耗、事件日志等
- **状态：已打包**：此时可以通过 `eth_getTransactionReceipt` 获取交易回执，`Status` 字段为 1（成功）或 0（失败）

#### 阶段 4：确认（Confirmed）
- **区块确认**：交易所在区块被后续区块引用，每增加一个后续区块，确认数 +1
- **最终性**：在 PoS 机制下，经过一定数量的确认后，交易达到最终性（Finality），几乎不可能被回滚
- **状态：已确认**：交易成为链上历史的一部分，可以查询任意历史区块的状态

#### 生命周期状态查询示例

```go
// 查询交易状态
tx, isPending, err := client.TransactionByHash(ctx, txHash)
if err != nil {
    // 交易不存在或查询失败
}

if isPending {
    // 交易还在 pending 状态，等待被打包
} else {
    // 交易已打包，可以查询回执
    receipt, err := client.TransactionReceipt(ctx, txHash)
    if err == nil {
        // 获取确认数：当前区块号 - 交易所在区块号
        currentBlock, _ := client.BlockNumber(ctx)
        confirmations := currentBlock - receipt.BlockNumber.Uint64()
    }
}
```

#### 常见问题与注意事项

- **交易可能失败的情况**：
  - Gas 不足：`gasLimit` 设置过小，执行中途耗尽 Gas
  - 余额不足：账户余额不足以支付 `value + gasFee`
  - 合约执行失败：合约代码执行过程中 revert
  - Nonce 错误：nonce 不连续或重复使用

- **交易可能卡住的情况**：
  - Gas 价格过低：交易池中 Gas 价格更高的交易优先被打包
  - 网络拥堵：大量交易排队等待打包
  - 节点同步问题：本地节点未同步到最新状态

- **最佳实践**：
  - 发送交易后定期轮询查询回执，设置合理的超时时间
  - 根据网络情况动态调整 Gas 价格
  - 确保 nonce 连续递增，避免交易顺序错乱

你可以在 `examples/03-tx-ops/main.go` 中看到如何查询交易状态和回执的完整示例。

### 2.4 区块结构

- 区块头（Header）：
  - 区块号、时间戳、父区块哈希  
  - 状态根（stateRoot）、交易根（txRoot）、日志根（receiptsRoot）  
  - Gas 上限与已用 Gas、难度等  
- 区块体（Body）：
  - 交易列表  
  - 自以太坊转向 PoS 后，可简化对工作量证明相关字段的讲解重点  

可以重点记住："区块是交易的有序集合 + 状态快照索引"。

### 2.5 Gas 机制与 EIP-1559 简介

- Gas 的含义：消耗计算与存储资源的度量单位  
- 费用计算（简化版）：
  - 旧模型：`fee = gasUsed * gasPrice`  
  - 新模型（EIP-1559）：`fee ≈ gasUsed * (baseFee + priorityFee)`  
- `baseFee` 由协议动态调整，`priorityFee` 用于激励打包交易的验证者  

你可以通过下面这些例子来建立直观感觉：
- 一笔普通转账大约消耗多少 Gas  
- 调用复杂 DeFi 合约时 Gas 为什么会更高  

### 2.6 网络类型

- 主网（Mainnet）：真实价值  
- 公共测试网（如 Goerli / Sepolia）：用于测试，无真实价值  
- 本地私有链 / 开发链：
  - `geth --dev`、`anvil`、`ganache` 等  
  - 优点：可控环境、区块快速出块、免费 Gas  

动手实践时，你可以尝试：
- 启动一个本地开发链  
- 在本地链上发一笔交易，并在控制台中观察相关日志  

## 3. Go-Eth Client 安装与连接

### 3.1 安装 `go-ethereum`

在示例工程中引入模块依赖：

```bash
go get github.com/ethereum/go-ethereum
```

核心包：
- `github.com/ethereum/go-ethereum/ethclient`：高层客户端封装  
- `github.com/ethereum/go-ethereum/common`：地址、哈希类型  
- `github.com/ethereum/go-ethereum/core/types`：区块与交易类型  

### 3.2 使用 `ethclient.Client` 连接节点

基础模式：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client, err := ethclient.DialContext(ctx, rpcURL)
if err != nil {
    // 处理错误
}
defer client.Close()
```

你可以在 `examples/01-connect-node/main.go` 中看到对应示例代码，主要展示：
- 读取环境变量 `ETH_RPC_URL`  
- 使用 `DialContext` 连接（支持超时控制）
- 连接成功后获取链 ID（`ChainID`）、最新区块号并打印区块信息  

### 3.3 本地节点 vs 远程节点

- 本地节点：
  - 优点：控制完全、隐私更好、可做归档查询  
  - 缺点：运维成本高，需要同步数据  
- 远程 RPC 服务（Infura 等）：
  - 优点：开箱即用、提供高可用与多地区出口  
  - 缺点：受限于请求速率与接口类型  

可以在课程中给出配置示例：
- `ETH_RPC_URL=http://127.0.0.1:8545`  
- `ETH_RPC_URL=https://sepolia.infura.io/v3/<project-id>`  

### 3.4 简单连接池与多节点策略（概念级）

- 多个 `ethclient.Client` 连接不同节点  
- 读操作可做简单负载均衡（轮询 / 随机）  
- 写操作（发交易）建议选定主节点，以减少状态不一致问题  
- 遇到节点不可用时自动切换并打告警日志  

本课程基础篇可仅讲概念，进阶篇或项目中再给具体实现示例。  
对应示例代码（简单读写分离与多节点轮询）：`examples/10-multi-node-pool/main.go`

## 4. 区块操作（对应示例代码：`examples/02-block-ops/main.go`）

### 4.1 获取最新区块

- 使用 `HeaderByNumber(ctx, nil)` 获取最新区块头  
- 使用 `BlockByNumber(ctx, nil)` 获取完整区块  

解析要点：
- 高度、时间戳、矿工、Gas 使用情况  
- 交易数量  

### 4.2 获取指定区块与区块列表

#### 4.2.1 使用大整数表示区块号

在 `go-ethereum` 中，区块号使用 `*big.Int` 类型表示，这是因为区块号可能非常大（当前主网已超过 1800 万）。

```go
import "math/big"

// 创建指定区块号
blockNumber := big.NewInt(123456)

// 从 uint64 创建
num := uint64(18500000)
blockNumber := big.NewInt(0).SetUint64(num)

// nil 表示最新区块
latestBlock, err := client.BlockByNumber(ctx, nil)
```

#### 4.2.2 查询指定区块

基本用法：

```go
blockNumber := big.NewInt(123456)
block, err := client.BlockByNumber(ctx, blockNumber)
if err != nil {
    log.Fatalf("failed to get block: %v", err)
}
```

在示例代码 `examples/02-block-ops/main.go` 中，你可以通过命令行参数 `--number` 查询指定区块：

```bash
go run main.go -number 123456
```

#### 4.2.3 批量查询区块范围

在实际项目中，经常需要循环获取一段区块范围（例如 `[start, end]`）。这需要注意以下几个关键点：

**1. 控制请求频率（避免触发限流）**

RPC 服务提供商会设置请求频率限制，过快的请求可能导致：
- 429 Too Many Requests 错误
- 临时封禁 IP
- 服务质量下降

使用 `time.Ticker` 实现频率控制：

```go
rateLimit := 200 * time.Millisecond // 每 200ms 一次请求
ticker := time.NewTicker(rateLimit)
defer ticker.Stop()

for num := start; num <= end; num++ {
    <-ticker.C // 等待下一个时间窗口
    
    blockNumber := big.NewInt(0).SetUint64(num)
    block, err := client.BlockByNumber(ctx, blockNumber)
    // ... 处理区块
}
```

**最佳实践建议：**
- Infura/Alchemy 等公共服务：建议 100-500ms 间隔
- 本地节点：可以更快，建议 50-200ms
- 根据实际 RPC 服务的限制动态调整

**2. 错误重试机制**

网络请求可能因为各种原因失败（网络波动、节点临时不可用等）。实现指数退避重试：

```go
func fetchBlockWithRetry(ctx context.Context, client *ethclient.Client, 
    blockNumber *big.Int, maxRetries int) (*types.Block, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
        block, err := client.BlockByNumber(reqCtx, blockNumber)
        cancel()
        
        if err == nil {
            return block, nil
        }
        
        lastErr = err
        if i < maxRetries-1 {
            // 指数退避：500ms, 1000ms, 2000ms...
            backoff := time.Duration(1<<uint(i)) * 500 * time.Millisecond
            time.Sleep(backoff)
        }
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

**重试策略选择：**
- 临时错误（网络超时、连接重置）：重试 2-3 次
- 永久错误（区块不存在、参数错误）：不重试，直接跳过
- 限流错误（429）：延长等待时间后重试

**3. 超时控制**

为每个请求设置合理的超时时间，避免长时间阻塞：

```go
// 为整个批量操作设置总体超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

// 为单个请求设置超时
reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
block, err := client.BlockByNumber(reqCtx, blockNumber)
cancel()
```

**超时时间建议：**
- 单个区块查询：5-10 秒
- 批量查询总体超时：根据区块数量动态计算（例如：区块数 × 请求间隔 × 1.5）

**4. 错误处理与跳过策略**

在批量查询中，某些区块可能查询失败，需要决定是终止还是跳过：

```go
successCount := 0
skipCount := 0

for num := start; num <= end; num++ {
    block, err := fetchBlockWithRetry(ctx, client, blockNumber, 2)
    if err != nil {
        log.Printf("[ERROR] Block %d: %v (skipping)", num, err)
        skipCount++
        continue // 跳过当前区块，继续下一个
    }
    
    successCount++
    // 处理区块...
}

// 最后输出统计信息
fmt.Printf("Success: %d, Skipped: %d\n", successCount, skipCount)
```

**5. 上下文取消支持**

支持通过上下文取消来优雅地停止批量查询：

```go
for num := start; num <= end; num++ {
    select {
    case <-ctx.Done():
        log.Printf("Context cancelled, stopping at block %d", num)
        return
    default:
    }
    
    // 执行查询...
}
```

#### 4.2.4 完整示例

在 `examples/02-block-ops/main.go` 中提供了完整的批量查询示例：

```bash
# 批量查询区块范围 [100, 105]
go run main.go -range-start 100 -range-end 105

# 自定义请求间隔（毫秒）
go run main.go -range-start 100 -range-end 105 -rate-limit 500
```

该示例演示了：
- 频率控制（`-rate-limit` 参数）
- 错误重试（每个区块最多重试 2 次）
- 超时控制（总体超时 30 秒，单次请求 10 秒）
- 错误统计和跳过策略

### 4.3 区块信息解析

区块对象（`*types.Block`）包含丰富的信息，理解这些字段对于构建区块浏览器、数据分析工具等应用至关重要。

#### 4.3.1 基本信息

**区块号（Number）**
```go
blockNumber := block.Number().Uint64()
fmt.Printf("Block Number: %d\n", blockNumber)
```

**区块哈希（Hash）**
```go
hash := block.Hash().Hex()
fmt.Printf("Block Hash: %s\n", hash)
// 输出：0x1234...5678
```
区块哈希是区块头的 Merkle 根，用于唯一标识一个区块。

**父区块哈希（ParentHash）**
```go
parentHash := block.ParentHash().Hex()
fmt.Printf("Parent Hash: %s\n", parentHash)
```
用于构建区块链的链式结构，每个区块都指向其父区块。

#### 4.3.2 时间信息

**时间戳（Timestamp）**
```go
timestamp := block.Time()
blockTime := time.Unix(int64(timestamp), 0)

// RFC3339 格式（ISO 8601）
fmt.Printf("Time (UTC): %s\n", blockTime.Format(time.RFC3339))
// 输出：2024-01-15T10:30:45Z

// 本地时间
fmt.Printf("Time (Local): %s\n", blockTime.Local().Format("2006-01-02 15:04:05 MST"))
// 输出：2024-01-15 18:30:45 CST
```

**注意事项：**
- 区块时间戳以秒为单位（Unix 时间戳）
- 在 PoS 机制下，出块时间相对固定（约 12 秒）
- 可以通过时间戳计算区块间隔，分析网络状况

#### 4.3.3 Gas 信息

**Gas 使用情况**
```go
gasUsed := block.GasUsed()
gasLimit := block.GasLimit()

// 计算使用率
usagePercent := float64(gasUsed) / float64(gasLimit) * 100
fmt.Printf("Gas Used: %d / %d (%.2f%%)\n", gasUsed, gasLimit, usagePercent)
```

**字段说明：**
- `GasLimit`：区块允许的最大 Gas 总量，由网络共识决定
- `GasUsed`：区块中所有交易实际消耗的 Gas 总和
- 使用率可以反映网络拥堵程度：接近 100% 表示网络繁忙

**实际应用：**
- 监控网络拥堵：Gas 使用率持续高于 90% 表示网络繁忙
- Gas 费用预测：使用率高的时期，Gas 费用通常更高

#### 4.3.4 交易信息

**交易数量**
```go
txCount := len(block.Transactions())
fmt.Printf("Transaction Count: %d\n", txCount)
```

**交易哈希**
```go
if txCount > 0 {
    // 第一笔交易
    firstTx := block.Transactions()[0]
    fmt.Printf("First Tx Hash: %s\n", firstTx.Hash().Hex())
    
    // 最后一笔交易
    lastTx := block.Transactions()[txCount-1]
    fmt.Printf("Last Tx Hash: %s\n", lastTx.Hash().Hex())
}
```

**注意事项：**
- 某些区块可能没有交易（空区块）
- 交易在区块中按顺序排列
- 可以通过交易数量分析链上活动活跃度

#### 4.3.5 Merkle 树根信息

以太坊使用 Merkle 树来高效验证区块数据的完整性：

**状态根（State Root）**
```go
stateRoot := block.Root().Hex()
fmt.Printf("State Root: %s\n", stateRoot)
```
代表区块执行后所有账户状态的 Merkle 根。

**交易根（Transaction Root）**
```go
txRoot := block.TxHash().Hex()
fmt.Printf("Transaction Root: %s\n", txRoot)
```
代表区块中所有交易的 Merkle 根。

**回执根（Receipt Root）**
```go
receiptRoot := block.ReceiptHash().Hex()
fmt.Printf("Receipt Root: %s\n", receiptRoot)
```
代表所有交易回执的 Merkle 根。

**实际应用：**
- 轻节点可以通过 Merkle 证明验证特定交易是否在区块中
- 状态根用于验证账户状态的快照

#### 4.3.6 其他重要字段

**难度（Difficulty）**
```go
difficulty := block.Difficulty()
fmt.Printf("Difficulty: %s\n", difficulty.String())
```
- PoW 机制下用于调整挖矿难度
- PoS 机制下基本固定或不再使用

**验证者地址（Coinbase / Beneficiary）**
```go
if block.Coinbase() != nil {
    coinbase := block.Coinbase().Hex()
    fmt.Printf("Coinbase: %s\n", coinbase)
}
```
- PoW 中称为矿工地址（Miner）
- PoS 中称为验证者地址（Validator）
- 接收区块奖励和交易费用的地址

#### 4.3.7 完整解析示例

在 `examples/02-block-ops/main.go` 的 `printBlockInfo` 函数中，提供了完整的区块信息解析示例，包括：
- 所有基本信息
- 时间戳的多格式显示
- Gas 使用率计算
- 交易列表信息
- Merkle 根信息

**对比区块浏览器**

建议你选择一个区块（例如主网区块 18500000），分别：
1. 在 Etherscan 等区块浏览器查看
2. 使用示例代码查询相同区块
3. 对比两者显示的信息，建立直观认知

这样可以更好地理解：
- 区块浏览器是如何展示这些数据的
- 哪些字段在哪些场景下更重要
- 如何组织这些信息以便用户理解

## 5. 交易与账户操作（对应示例代码：`examples/03-tx-ops/main.go` 与 `examples/04-account-balance/main.go`）

本章节涵盖交易的完整生命周期：从查询已有交易到构造并发送新交易。

### 5.1 获取交易信息

- 使用 `TransactionByHash(ctx, txHash)` 获取交易  
- 注意：返回值包含交易本身和一个布尔值 `isPending`，用于判断交易是否还在等待打包
- 交易对象中不包含执行结果，需要结合回执判断交易是否成功  

解析字段：
- `To()` / `Value()` / `Gas()` / `GasPrice()` / `Data()` / `Nonce()` 等
- 示例代码通过 `common.HexToHash` 将十六进制字符串转换为交易哈希  

### 5.2 交易回执解析

- 使用 `TransactionReceipt(ctx, txHash)` 获取回执  
- 重点字段：
  - `Status`：1=成功，0=失败  
  - `GasUsed`：实际消耗的 Gas  
  - `Logs`：事件日志列表（为后续事件解析做铺垫）  

你可以通过一笔成功 / 失败交易进行对比，观察 `Status` 与 `GasUsed` 的差异。

### 5.3 交易状态查询

- 已打包（mined）：`TransactionReceipt` 返回非空  
- pending：`TransactionByHash` 能查到，`TransactionReceipt` 暂无结果  
- 建议实现一个简单轮询函数：
  - 每隔几秒查询回执  
  - 超过一定时间仍未打包给出提示  

### 5.4 账户余额查询

- 使用 `BalanceAt(ctx, address, blockNumber)` 查询余额  
  - `blockNumber = nil` 表示最新状态  
  - 返回值为 Wei（`*big.Int`），需要换算为 ETH：`eth = wei / 1e18`  
- 示例代码实现了 `weiToEth` 函数，使用 `big.Float` 进行精确的数值转换
- 支持通过 `--block` 参数查询历史区块时的余额

你可以练习：
- 使用 `--address` 参数查询指定地址在主网/测试网的余额  
- 使用 `--block` 参数查询历史区块时的余额状态
- 在本地链上发起一次转账，查看前后余额的变化  

### 5.5 发送交易（发起转账）

在 `examples/03-tx-ops/main.go` 中，除了查询交易功能外，还提供了发送 ETH 转账交易的完整示例。

#### 5.5.1 发送交易的基本步骤

发送一笔 ETH 转账交易需要完成以下步骤：

1. **准备私钥和地址**
   - 从环境变量 `SENDER_PRIVATE_KEY` 读取发送方私钥（十六进制格式，可带或不带 `0x` 前缀）
   - 使用 `crypto.HexToECDSA()` 解析私钥
   - 通过 `crypto.PubkeyToAddress()` 获取发送方地址

2. **获取链信息**
   - 使用 `client.ChainID()` 获取链 ID（用于签名）
   - 使用 `client.PendingNonceAt()` 获取发送方账户的 nonce（必须连续递增）

3. **设置 Gas 费用（EIP-1559 动态费用）**
   - 使用 `client.SuggestGasTipCap()` 获取建议的优先费用（tip）
   - 通过 `client.HeaderByNumber()` 获取当前 base fee
   - 计算 fee cap：`feeCap = baseFee * 2 + tipCap`（简单策略）

4. **构造交易**
   - 使用 `types.DynamicFeeTx` 构造 EIP-1559 动态费用交易
   - 设置 `To`（接收方地址）、`Value`（转账金额，Wei 单位）、`Gas`（普通转账固定为 21000）

5. **签名和发送**
   - 使用 `types.NewLondonSigner()` 创建签名器
   - 使用 `types.SignTx()` 对交易进行签名
   - 使用 `client.SendTransaction()` 广播交易到网络

#### 5.5.2 余额检查

在发送交易前，示例代码会检查账户余额是否足够：
- 总费用 = `value + gasFeeCap * gasLimit`
- 如果余额不足，程序会直接报错退出

#### 5.5.3 使用示例

**发送交易模式：**
```bash
# 设置环境变量
export ETH_RPC_URL="http://127.0.0.1:8545"  # 本地节点或测试网
export SENDER_PRIVATE_KEY="your_private_key_hex"  # 发送方私钥（不含 0x 前缀也可）

# 发送 0.1 ETH 到指定地址
go run main.go --send --to 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb --amount 0.1
```

**查询交易状态：**
```bash
# 发送交易后会输出交易哈希，可以用它查询状态
go run main.go --tx 0x1234...5678
```

#### 5.5.4 安全注意事项

⚠️ **重要安全提示：**

1. **私钥管理**
   - 永远不要在代码仓库中提交私钥
   - 使用环境变量或专门的密钥管理服务
   - 示例代码仅用于测试网或本地开发链

2. **测试环境**
   - 建议在本地开发链（如 `anvil`、`ganache`）或公共测试网（如 Sepolia）上测试
   - 不要在主网使用包含真实资产的私钥

3. **Gas 费用**
   - 示例代码使用简单的 Gas 费用策略，生产环境可能需要更复杂的费用估算
   - 可以通过 `client.EstimateGas()` 更准确地估算 Gas Limit

4. **Nonce 管理**
   - 确保 nonce 连续递增，避免交易顺序错乱
   - 如果多笔交易并发发送，需要妥善管理 nonce

#### 5.5.5 交易生命周期回顾

发送交易后，交易会经历以下状态：

1. **Pending**：交易已广播到网络，等待被打包
   - 此时 `TransactionByHash` 可以查询到交易
   - `TransactionReceipt` 返回 `null`

2. **已打包**：交易被打包进区块
   - `TransactionReceipt` 可以查询到回执
   - `Status` 字段为 1（成功）或 0（失败）

3. **已确认**：交易所在区块被后续区块引用
   - 确认数 = 当前区块号 - 交易所在区块号
   - 确认数越多，交易越不可能被回滚

你可以结合 5.3 节的交易状态查询功能，实现一个简单的轮询机制来跟踪交易状态。

## 6. 小结与过渡

在基础篇结束时，你已经具备通过 Go 代码读取链上**静态数据**的能力（区块、交易、账户）。  
接下来在进阶篇中，你将进一步学习：
- **事件订阅**：实时监听新区块、pending 交易与合约日志  
- **合约交互**：通过 ABI 进行读写与事件解析  

在此基础上，你最终会用这些能力搭建一个迷你区块浏览器 / 监听服务，把链上的变化变成自己系统可以消费的数据。
