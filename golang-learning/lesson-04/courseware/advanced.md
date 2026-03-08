# 第4期：事件订阅、智能合约交互与工程实践

本讲义对应课程的进阶篇与项目实战部分，结合 `examples/05-subscribe-blocks` 至 `examples/09-project` 目录下的示例，帮助你把前面学到的基础能力应用到更完整的工程场景中。

## 1. 学习目标

通过本部分学习，学员将能够：

- 使用 Go-Eth Client 订阅新区块、pending 交易与日志事件  
- 设计基本的断线检测与重连机制，避免订阅意外中断  
- 通过 ABI 编解码与合约进行读写交互  
- 解析合约事件日志（尤其是 ERC-20 Transfer）  
- 实现一个简单的链上监听服务 / 区块浏览器后端  

## 2. 事件订阅基础

### 2.1 事件源与订阅接口

以 `ethclient` 为例，可以通过以下方式订阅事件：

- 订阅新区块头：`SubscribeNewHead`  
- 订阅日志（合约事件）：`SubscribeFilterLogs`  

对于 pending 交易：
- 部分节点支持 `eth_subscribe` 类型的 `newPendingTransactions`  
- 实际可用性取决于所连接节点与 RPC 提供方  

### 2.2 订阅新区块头

你可以参考示例文件：`examples/05-subscribe-blocks/main.go`

关键步骤：
- 创建 `context.Context` 与 `chan *types.Header`  
- 调用 `client.SubscribeNewHead(ctx, headers)`  
- 使用 `select` 监听：
  - `case h := <-headers`：收到新区块头  
  - `case err := <-sub.Err()`：订阅错误  
  - `case <-ctx.Done()`：主动退出  
  - `case sig := <-sigCh`：捕获系统信号（SIGINT/SIGTERM）实现优雅退出

理解重点：
- **重要**：订阅功能通常需要 WebSocket 连接（`ws://` 或 `wss://`），示例代码使用 `ETH_WS_URL` 环境变量，并提供了 `ETH_RPC_URL` 的回退机制
- 为什么使用 channel 和 goroutine  
- 如何优雅退出（通过信号量捕获 Ctrl+C，或 context 取消）
- 示例代码展示了如何处理 `nil` header 的情况  

### 2.3 订阅日志事件

你可以参考示例文件：`examples/06-subscribe-logs/main.go`

使用 `ethereum.FilterQuery` 构建过滤条件：

- `Addresses`：合约地址列表（示例代码通过 `--contract` 命令行参数传入）  
- `Topics`：事件 topic 过滤（如 ERC-20 Transfer 的第一个 topic，示例中暂未使用）  

调用：
- `client.SubscribeFilterLogs(ctx, query, logsCh)`  

在日志中可获取：
- `log.Address`：合约地址  
- `log.Topics`：事件 topic 列表  
- `log.Data`：编码后的事件参数  
- `log.BlockNumber`、`log.TxHash`、`log.Index` 等元数据

**注意**：与区块订阅类似，日志订阅也需要 WebSocket 连接，示例代码同样使用 `ETH_WS_URL` 环境变量。  

### 2.4 订阅 pending 交易（可选内容）

根据节点支持情况，你可以尝试：

- 使用原始 WebSocket 客户端订阅 `newPendingTransactions`  
- 了解不同托管 RPC 提供商对该接口的支持限制  

这部分属于知识扩展内容，在掌握前面章节后再动手实现即可。

## 3. 断线重连与健壮性设计

### 3.1 常见失败场景

- 网络闪断、节点重启  
- RPC 提供方主动关闭空闲连接  
- 本地应用重启或 panic  

### 3.2 基本重连策略

你可以参考示例文件：`examples/07-reconnect-strategy/main.go`

设计要点：

- 使用一个封装函数（`runWithReconnect`）管理订阅生命周期  
- 当 `sub.Err()` 返回错误时：
  - 记录错误日志  
  - 关闭旧连接（`client.Close()`）  
  - 等待一小段时间（指数退避，示例代码中实现了 `sleepWithBackoff`，最大等待 60 秒）  
  - 重新建立连接与订阅  
- 示例代码使用 `goto RECONNECT` 跳出内层循环，在外层循环中重新尝试连接
- 在重新订阅时，可记录上一次成功处理的区块高度，必要时补拉中间缺失的数据（仅对日志订阅 / 区块扫描场景适用）
- 支持通过系统信号（Ctrl+C）优雅退出，退出时会取消 context，停止重连循环  

### 3.3 与手动扫描结合

订阅适合实时性高的场景，但仍需配合：

- 启动时的“回放”扫描（从某个高度到最新高度）  
- 出现长时间断线后的“补数据”扫描  

可以在项目中抽象出：
- 扫描模块：按区块范围批量拉取数据  
- 订阅模块：从最新高度往前持续监听  

## 4. 智能合约交互基础

### 4.1 ABI 的角色

- ABI（Application Binary Interface）：定义合约对外暴露的函数与事件接口  
- 对于函数调用：
  - 将函数名与参数编码为 `data` 字段  
  - 将返回值从字节流解码为 Go 类型  
- 对于事件：
  - 将 `topics` 与 `data` 解码为结构化类型  

### 4.2 加载 ABI

两种常见方式：

1. 直接加载 JSON 文件内容：  
   - 使用 `abi.JSON(strings.NewReader(abiJSON))`  
2. 使用 `abigen` 生成 Go 绑定（项目中可作为后续扩展）  

建议你先掌握"通用 ABI"方式，再根据需要了解代码生成方式。

#### 4.2.1 从 Foundry 编译结果中提取 ABI

使用 Foundry 编译合约后，会在 `out/` 目录下生成包含 ABI 的 JSON 文件。提取 ABI 有以下几种方法：

**方法一：使用 jq 提取（推荐）**

```bash
# 查看 ABI
jq '.abi' out/MyERC20.sol/MyERC20.json

# 保存 ABI 到单独文件
jq '.abi' out/MyERC20.sol/MyERC20.json > MyERC20.abi.json
```

**方法二：手动提取**

直接打开 `out/MyERC20.sol/MyERC20.json` 文件，复制其中的 `abi` 字段内容。

**生成的 JSON 文件结构：**

编译后的 JSON 文件包含以下内容：
- `abi`：合约的 ABI（应用二进制接口），包含函数、事件、错误等定义
- `bytecode`：完整字节码（包含构造函数参数）
- `deployedBytecode`：部署后的字节码（不包含构造函数）
- `metadata`：编译元数据，包括编译器版本、源码信息等

在 Go 代码中使用提取的 ABI：

```go
// 读取 ABI JSON 文件
abiBytes, err := os.ReadFile("MyERC20.abi.json")
if err != nil {
    log.Fatal(err)
}

// 解析 ABI
parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
if err != nil {
    log.Fatal(err)
}

// 现在可以使用 parsedABI 进行函数调用和事件解析
```

### 4.3 部署自己的 ERC20 合约用于测试

在实际学习和测试智能合约交互时，使用自己部署的合约比使用主网上的合约更加灵活和安全。本课程提供了 `examples/11-erc20-foundry` 项目，使用 Foundry 工具链部署一个基于 OpenZeppelin 的 ERC20 代币合约。

#### 4.3.1 项目结构

`11-erc20-foundry` 项目包含：

- `src/MyERC20.sol`：ERC20 代币合约，基于 OpenZeppelin 的 `ERC20` 实现
- `script/DeployERC20.s.sol`：Foundry 部署脚本
- `foundry.toml`：Foundry 配置文件
- `lib/`：依赖库（forge-std、openzeppelin-contracts）

#### 4.3.2 合约说明

`MyERC20.sol` 是一个简单的 ERC20 代币合约，支持：

- 自定义代币名称和符号
- 设置初始供应量
- 将初始代币铸造给指定地址
- 可铸造新代币的功能（`mint` 函数）

**构造函数参数：**
- `name`：代币名称（例如：MyToken）
- `symbol`：代币符号（例如：MTK）
- `initialSupply`：初始供应量（wei 单位，18 位小数）
- `recipient`：初始代币接收地址

#### 4.3.3 部署步骤

**1. 启动本地开发链**

首先启动一个本地 Anvil 节点（Foundry 自带的本地开发链）：

```bash
anvil
```

Anvil 会启动一个本地以太坊节点，默认监听 `http://127.0.0.1:8545`，并自动创建 10 个测试账户，每个账户有 10000 ETH。

**2. 编译合约**

在 `11-erc20-foundry` 目录下运行：

```bash
cd examples/11-erc20-foundry
forge build
```

**3. 部署到本地链**

使用 Foundry 的 `forge script` 命令部署合约：

```bash
# 使用 Anvil 默认的第一个账户私钥部署
forge script script/DeployERC20.s.sol:DeployERC20Script \
  --rpc-url http://localhost:8545 \
  --broadcast \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

**重要输出：**
- 部署脚本会输出合约地址，例如：`ERC20 Token deployed at: 0x5FbDB2315678afecb367f032d93F642f64180aa3`
- 记录这个地址，后续与合约交互时需要用到

**4. 验证部署**

部署成功后，可以使用 `cast` 命令验证合约：

```bash
# 查询代币名称
cast call <CONTRACT_ADDRESS> "name()(string)" --rpc-url http://localhost:8545

# 查询代币符号
cast call <CONTRACT_ADDRESS> "symbol()(string)" --rpc-url http://localhost:8545

# 查询总供应量
cast call <CONTRACT_ADDRESS> "totalSupply()(uint256)" --rpc-url http://localhost:8545

# 查询部署者余额（使用 Anvil 默认的第一个账户地址）
cast call <CONTRACT_ADDRESS> "balanceOf(address)(uint256)" 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --rpc-url http://localhost:8545
```

**5. 提取 ABI**

部署完成后，需要提取合约的 ABI 用于后续的 Go 代码交互。编译后，Foundry 会在 `out/MyERC20.sol/MyERC20.json` 文件中生成完整的编译结果。

**使用 jq 提取（推荐）：**

```bash
# 在 11-erc20-foundry 目录下
cd examples/11-erc20-foundry

# 提取 ABI 并保存到文件
jq '.abi' out/MyERC20.sol/MyERC20.json > MyERC20.abi.json

# 验证 ABI 文件
cat MyERC20.abi.json
```

**手动提取：**

如果系统没有安装 `jq` 或 `cast`，可以：
1. 打开 `out/MyERC20.sol/MyERC20.json` 文件
2. 找到 `"abi"` 字段（通常在文件开头）
3. 复制整个 `abi` 数组内容（从 `[` 开始到对应的 `]` 结束）
4. 保存为 `MyERC20.abi.json` 文件

**验证 ABI 文件：**

提取的 ABI 文件应该是一个 JSON 数组，包含合约的函数、事件和错误定义。可以使用以下命令验证：

```bash
# 检查文件格式
jq '.' MyERC20.abi.json

# 查看函数列表
jq '.[] | select(.type == "function") | .name' MyERC20.abi.json

# 查看事件列表
jq '.[] | select(.type == "event") | .name' MyERC20.abi.json
```

**在 Go 代码中使用 ABI：**

提取 ABI 后，可以在 Go 代码中加载并使用：

```go
package main

import (
    "os"
    "strings"
    "log"
    
    "github.com/ethereum/go-ethereum/accounts/abi"
)

func main() {
    // 读取 ABI 文件
    abiBytes, err := os.ReadFile("MyERC20.abi.json")
    if err != nil {
        log.Fatal("Failed to read ABI file:", err)
    }
    
    // 解析 ABI
    parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
    if err != nil {
        log.Fatal("Failed to parse ABI:", err)
    }
    
    // 现在可以使用 parsedABI 进行函数调用和事件解析
    // 例如：parsedABI.Pack("balanceOf", address)
    // 例如：parsedABI.Unpack("Transfer", logData)
}
```

**注意事项：**

- ABI 文件是纯 JSON 格式，不包含字节码信息
- 同一个合约的 ABI 在不同编译环境下应该是一致的（只要合约代码相同）
- 可以将 ABI 文件提交到代码仓库，方便团队共享
- 在生产环境中，建议将 ABI 文件放在项目的 `abi/` 或 `contracts/` 目录下统一管理

#### 4.3.4 与部署的合约交互

部署完成后，你可以使用 `examples/08-contract-interact/main.go` 与部署的合约进行交互：

**1. 查询代币余额**

```bash
export ETH_RPC_URL="http://127.0.0.1:8545"
go run examples/08-contract-interact/main.go --mode balance \
  --contract <DEPLOYED_CONTRACT_ADDRESS> \
  --address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

**2. 发送 ERC-20 转账交易**

```bash
export ETH_RPC_URL="http://127.0.0.1:8545"
export SENDER_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
go run examples/08-contract-interact/main.go --mode transfer \
  --contract <DEPLOYED_CONTRACT_ADDRESS> \
  --to 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 \
  --amount 100
```

**3. 解析转账事件**

转账成功后，可以使用交易哈希解析 Transfer 事件：

```bash
go run examples/08-contract-interact/main.go --mode parse-event \
  --tx <TRANSACTION_HASH>
```

#### 4.3.6 为什么使用自己的合约？

使用自己部署的 ERC20 合约进行测试有以下优势：

1. **完全控制**：可以自由铸造代币，无需担心余额不足
2. **安全性**：在本地链上测试，不会影响真实资产
3. **学习价值**：理解从合约编写、编译、部署到交互的完整流程
4. **灵活性**：可以修改合约代码，测试不同的场景

#### 4.3.7 部署到测试网（可选）

如果需要部署到公共测试网（如 Sepolia），可以：

```bash
# 设置环境变量
export SEPOLIA_RPC_URL="https://sepolia.infura.io/v3/<your-key>"
export PRIVATE_KEY="your_private_key_hex"
export ETHERSCAN_API_KEY="your_etherscan_api_key"

# 部署到 Sepolia
forge script script/DeployERC20.s.sol:DeployERC20Script \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast \
  --verify \
  --etherscan-api-key $ETHERSCAN_API_KEY \
  --private-key $PRIVATE_KEY
```

**注意事项：**
- 测试网需要真实的 ETH 作为 Gas 费用
- 部署后可以在 Etherscan 等区块浏览器查看合约
- 建议仅在测试网使用，不要在主网部署测试合约

更多详细信息请参考 `examples/11-erc20-foundry/README.md`。

### 4.4 调用只读方法（`call`）

你可以参考示例文件：`examples/08-contract-interact/main.go`

步骤：

- 获取合约 ABI（参考 4.2.1 和 4.3.3 中的提取方法，或使用示例代码中直接定义的 `erc20ABIJSON` 字符串）  
- 使用 `abi.JSON(strings.NewReader(abiJSON))` 解析 ABI  
- 使用 `parsedABI.Pack("functionName", args...)` 编码函数调用数据  
- 构造 `ethereum.CallMsg`，指定 `To`（合约地址）和 `Data`（编码后的函数调用数据）  
- 使用 `client.CallContract(ctx, msg, blockNumber)` 执行（`blockNumber = nil` 表示最新状态）  
- 使用 `parsedABI.UnpackIntoInterface(&output, "functionName", result)` 将返回值解码为 Go 类型  

示例方法：
- `balanceOf(address)`：查询 ERC-20 代币余额  
- 示例代码通过 `--contract` 和 `--address` 命令行参数传入合约地址和查询地址

**实践建议：**
- 使用 `11-erc20-foundry` 项目部署的合约地址进行测试
- 在本地 Anvil 链上测试，避免使用真实资产

### 4.5 发送交易（`transact`）

发送写入型交易需要：

- 账户私钥（本地管理或通过外部签名服务）  
- 构造裸交易：
  - `types.NewTransaction` / `NewTx`（EIP-1559）  
  - 填充 `nonce`、`gasLimit`、`gasFeeCap`、`gasTipCap` 等字段  
- 使用 `types.Signer` 与私钥进行签名  
- 使用 `client.SendTransaction(ctx, signedTx)` 广播  

课程中可以简化为：

- 使用测试网 / 本地链  
- 给出最小可运行的 ERC-20 转账示例（不一定完整覆盖 Gas 估算与费率优化）  

**实践建议：**
- 使用 `11-erc20-foundry` 部署的合约进行转账测试
- 在本地 Anvil 链上测试，Gas 费用为 0，可以自由测试
- 参考 `examples/08-contract-interact/main.go` 中的 `handleTransfer` 函数了解完整实现

### 4.6 解析合约事件日志

合约事件是以太坊中重要的数据发布机制，理解如何从日志（Log）中解析事件参数是进行链上数据监听和分析的基础。

#### 4.6.1 事件的存储结构

以太坊中的事件日志存储在交易的回执（Receipt）中，每个日志（`types.Log`）包含以下字段：

- **Address**：发出事件的合约地址
- **Topics**：主题数组，最多 4 个 topic（每个 topic 是 32 字节的哈希值）
  - `Topics[0]`：事件签名的 keccak256 哈希（用于识别事件类型）
  - `Topics[1..3]`：indexed 参数的编码值（最多 3 个 indexed 参数）
- **Data**：非 indexed 参数的编码数据（ABI 编码格式）
- **BlockNumber**、**TxHash**、**Index** 等元数据

#### 4.6.2 indexed 参数 vs 非 indexed 参数

在 Solidity 中定义事件时，参数可以标记为 `indexed`：

```solidity
event Transfer(
    address indexed from,  // indexed: 可以在 Topics 中查询和过滤
    address indexed to,    // indexed: 可以在 Topics 中查询和过滤
    uint256 value          // 非 indexed: 存储在 Data 中
);
```

**indexed 参数的特点：**
- 存储在 `Topics` 数组中，可以用于过滤和搜索
- 限制：最多 3 个 indexed 参数（因为 Topics 数组最多 4 个元素，Topics[0] 用于事件签名）
- 对于复杂类型（如数组、结构体），indexed 参数存储的是类型哈希，而非实际值

**非 indexed 参数的特点：**
- 存储在 `Data` 字段中，使用 ABI 编码
- 没有数量限制，但会占用更多的 Gas
- 不能用于事件过滤，但可以完整保存复杂数据结构

#### 4.6.3 识别事件类型

事件类型通过 `Topics[0]` 识别，它是事件签名的 keccak256 哈希值：

```go
// 事件签名格式：EventName(type1,type2,...)
eventSig := "Transfer(address,address,uint256)"
eventSigHash := crypto.Keccak256Hash([]byte(eventSig))

// 在日志中匹配
if vLog.Topics[0] == eventSigHash {
    // 这是 Transfer 事件
}
```

在实际应用中，可以通过 ABI 自动匹配：

```go
// 遍历 ABI 中定义的所有事件
for name, event := range parsedABI.Events {
    eventSigHash := crypto.Keccak256Hash([]byte(event.Sig))
    if eventSigHash == vLog.Topics[0] {
        eventName = name
        break
    }
}
```

#### 4.6.4 解析 indexed 参数

indexed 参数从 `Topics[1..N]` 中解析。以 ERC-20 `Transfer` 事件为例：

```go
// Transfer(address indexed from, address indexed to, uint256 value)
// Topics[0]: 事件签名哈希
// Topics[1]: from 地址（indexed）
// Topics[2]: to 地址（indexed）

if len(vLog.Topics) >= 3 {
    // 解析 address 类型的 indexed 参数
    // address 在 topic 中是 32 字节，前 12 字节为 0，后 20 字节是地址
    fromAddr := common.BytesToAddress(vLog.Topics[1].Bytes())
    toAddr := common.BytesToAddress(vLog.Topics[2].Bytes())
    
    fmt.Printf("From: %s\n", fromAddr.Hex())
    fmt.Printf("To: %s\n", toAddr.Hex())
}
```

**不同类型 indexed 参数的解析：**

```go
topic := vLog.Topics[index]

switch inputType {
case "address":
    // address: 32 字节，前 12 字节为 0，后 20 字节是地址
    addr := common.BytesToAddress(topic.Bytes())
    
case "uint256", "int256":
    // 整数类型：32 字节直接转换为 big.Int
    value := new(big.Int).SetBytes(topic.Bytes())
    
case "bool":
    // bool: 检查最后一个字节（topic[31]）
    isTrue := topic[31] != 0
    
case "bytes32":
    // bytes32: 直接使用 32 字节数据
    bytesValue := topic.Bytes()
    
default:
    // 其他类型：显示原始十六进制
    hexValue := topic.Hex()
}
```

#### 4.6.5 解析非 indexed 参数

非 indexed 参数存储在 `Data` 字段中，使用 ABI 编码。使用 `abi` 包进行解码：

**方法 1：使用 `Unpack` 返回切片**

```go
// 解码所有参数（包括 indexed 和非 indexed）
values, err := parsedABI.Unpack(eventName, vLog.Data)
if err != nil {
    log.Printf("failed to unpack: %v", err)
    return
}

// values 是一个 []interface{}，只包含非 indexed 参数
// 注意：Unpack 只解码 Data 字段，不包含 Topics 中的 indexed 参数
if len(values) > 0 {
    value := values[0].(*big.Int) // Transfer 事件中，value 是 uint256
    fmt.Printf("Value: %s\n", value.String())
}
```

**方法 2：使用 `UnpackIntoInterface` 解码到结构体**

```go
// 定义接收结构体（只包含非 indexed 参数）
var event struct {
    Value *big.Int // Transfer 事件中的 value 参数
}

// 解码 Data 部分
err := parsedABI.UnpackIntoInterface(&event, "Transfer", vLog.Data)
if err != nil {
    log.Printf("failed to unpack: %v", err)
    return
}

fmt.Printf("Value: %s\n", event.Value.String())
```

**重要提示：**
- `Unpack` 和 `UnpackIntoInterface` 只解码 `Data` 字段
- `Data` 中**只包含非 indexed 参数**
- indexed 参数需要从 `Topics` 中单独解析

#### 4.6.6 完整解析示例

以 ERC-20 `Transfer(address indexed from, address indexed to, uint256 value)` 事件为例，完整解析流程：

```go
// 1. 检查日志有效性
if len(vLog.Topics) == 0 {
    return // 无效日志
}

// 2. 识别事件类型
eventSigHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
if vLog.Topics[0] != eventSigHash {
    return // 不是 Transfer 事件
}

// 3. 解析 indexed 参数（从 Topics 中）
if len(vLog.Topics) >= 3 {
    from := common.BytesToAddress(vLog.Topics[1].Bytes())
    to := common.BytesToAddress(vLog.Topics[2].Bytes())
    fmt.Printf("From: %s, To: %s\n", from.Hex(), to.Hex())
}

// 4. 解析非 indexed 参数（从 Data 中）
var event struct {
    Value *big.Int
}
if err := parsedABI.UnpackIntoInterface(&event, "Transfer", vLog.Data); err != nil {
    log.Printf("failed to unpack data: %v", err)
    return
}
fmt.Printf("Value: %s\n", event.Value.String())
```

#### 4.6.7 实际应用示例

在 `examples/08-contract-interact/main.go` 的 `handleParseEvent` 函数中提供了详细的事件解析示例，展示了：

- 如何通过 ABI 自动识别事件类型
- 如何解析 indexed 参数（支持 address、uint256、bool 等类型）
- 如何解析非 indexed 参数（使用 `Unpack` 方法）
- 完整的错误处理和格式化输出

**实践建议：**
- 使用 `11-erc20-foundry` 部署的合约进行转账，然后使用 `--mode parse-event` 解析 Transfer 事件
- 这样可以完整地体验从合约部署、转账到事件解析的全流程

在 `examples/06-subscribe-logs/main.go` 中展示了如何实时订阅和解析事件：

- 订阅指定合约地址的事件日志
- 实时解析并输出事件参数

在 `examples/09-project/main.go` 的项目示例中，展示了如何在生产环境中使用事件解析：

- 定义事件结构体用于存储解析结果
- 将解析后的事件存入内存或数据库
- 通过 HTTP API 提供事件查询接口

#### 4.6.8 常见问题和注意事项

1. **Topics 数量限制**
   - Topics 数组最多 4 个元素（Topics[0..3]）
   - 因此最多只能有 3 个 indexed 参数
   - 超过 3 个 indexed 参数会导致编译错误

2. **复杂类型的 indexed 参数**
   - 数组、结构体等复杂类型标记为 indexed 时，存储的是类型哈希（keccak256），而非实际值
   - 如果需要完整数据，应该使用非 indexed 参数

3. **Data 字段为空**
   - 如果事件所有参数都是 indexed，Data 字段可能为空
   - 在解析前应检查 `len(vLog.Data) > 0`

4. **事件签名哈希计算**
   - 事件签名必须完全匹配，包括参数类型
   - `Transfer(address,address,uint256)` 和 `Transfer(address indexed,address indexed,uint256)` 的哈希是不同的
   - 使用 ABI 定义可以自动计算正确的签名哈希  

## 5. 实战项目：迷你区块浏览器 / ERC-20 监听服务

### 5.1 需求与功能列表

在 `courseware/project.md` 中有更细化的说明，这里给出高层结构：

- 提供基础查询能力：
  - 按区块号/哈希查询区块  
  - 按哈希查询交易与回执  
  - 按地址查询最近交易与余额（可先实现简单版本）  
- 后台持续监听：
  - 新区块  
  - 指定 ERC-20 合约的 `Transfer` 事件  
- 数据存储：
  - 初始版本可只保存在内存中（如环形缓冲区）  
  - 后续版本可以接入 SQLite / PostgreSQL 等  

### 5.2 推荐目录结构（示例）

你可以参考 `lesson-04/examples/09-project`：

- `main.go`：程序入口，负责：
  - 读取配置（RPC 地址、监听合约列表等）  
  - 初始化客户端与服务  
  - 启动 HTTP 服务或命令行界面  
- `internal/client/`：封装 Go-Eth Client 初始化与基础调用  
- `internal/service/`：实现区块 / 交易 / 地址 / 事件相关业务逻辑  
- `internal/api/`：对外暴露接口（RESTful API 或简单 HTTP handler）  
- `internal/store/`：数据存储层（内存 / DB）  

### 5.3 最小可行版本（MVP）建议

第一版项目可以只实现：

- 后台 goroutine 监听新区块与指定 `Transfer` 事件  
- 将最近 N 条事件存入内存（slice 或 ring buffer）  
- 暴露一个 HTTP 接口：`GET /events` 返回最近事件列表  

在此基础上逐步扩展：

- 查询历史区块与交易  
- 为地址增加过滤参数（只看某个地址相关的事件）  
- 将事件写入数据库，并支持分页查询  

## 6. 工程实践与最佳实践提示

### 6.1 上下文与超时

- 为每个 RPC 调用设置合理超时（如 5~10 秒）  
- 对长时间订阅使用无超时的根 context，但在程序退出时要能取消  

### 6.2 日志与监控

- 为关键事件（连接失败、重连、数据解码错误）记录日志  
- 对重连次数、失败率等做基础监控（可先打印日志，生产环境对接监控系统）  

### 6.3 错误处理与重试

- 对瞬时网络错误进行重试（带退避）  
- 对明显配置错误（如无效 RPC URL）直接失败并提示  

### 6.4 安全性提示

- 切勿在代码仓库中提交明文私钥  
- 如需演示签名交易，应使用测试网或本地链  
- 生产环境应考虑使用专门的签名服务或 HSM  

## 7. 实践练习示例

你可以尝试完成以下练习：

1. **练习 1：实现简易块高监控**  
   - 每隔 10 秒打印当前最新区块号  
   - 当块高长时间不变时给出告警日志  

2. **练习 2：实现指定地址交易监控**  
   - 订阅新区块  
   - 对每个区块中的交易，筛选 `from` 或 `to` 为指定地址的交易并打印  

3. **练习 3：扩展 ERC-20 监听**  
   - 在日志订阅基础上，支持监听多个 ERC-20 合约地址  
   - 增加一个简单的内存索引：按地址维护最近收到/发出的转账列表  

