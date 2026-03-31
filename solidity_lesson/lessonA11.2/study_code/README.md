# Hardhat 3 单元测试示例代码

本目录包含Hardhat 3单元测试的完整示例代码，展示了各种测试技巧和最佳实践。

## 文件说明

### 合约文件

- **Counter.sol**: 简单的计数器合约，用于演示基础测试

### 测试文件

- **Counter.test.ts**: Counter合约的完整测试示例
  - 部署测试
  - 功能测试
  - 事件测试
  - 错误处理测试
  - 状态隔离测试

- **Token.test.ts**: 代币合约测试示例（演示更多测试场景）
  - 转账测试
  - 余额变化测试
  - 错误处理测试

- **TimeLock.test.ts**: 时间锁合约测试示例
  - 时间操作演示
  - 时间旅行测试

- **ErrorHandling.test.ts**: 错误处理测试示例
  - require回退测试
  - 自定义错误测试
  - Panic错误测试

## 运行测试

### 运行所有测试

```bash
npx hardhat test
```

### 运行特定测试文件

```bash
npx hardhat test test/Counter.test.ts
```

### 运行匹配的测试

```bash
npx hardhat test --grep "Deployment"
```

### 生成Gas报告

```bash
npx hardhat test --gas-stats
```

### 生成覆盖率报告

```bash
npx hardhat test --coverage
```

### 生成HTML覆盖率报告

```bash
npx hardhat test --coverage --report html
```

### 详细输出

```bash
npx hardhat test --verbosity 3
```

### 显示堆栈跟踪

```bash
npx hardhat test --show-stack-traces
```

## 测试要点

### 1. Hardhat 3新API

- 使用`await network.connect()`获取ethers对象
- 使用`networkHelpers.loadFixture()`加载Fixture
- 使用`ethers.deployContract()`部署合约

### 2. Fixture使用

- 定义Fixture函数进行初始化
- 使用`loadFixture`实现快照恢复
- 每个测试都有干净的初始状态

### 3. 断言方法

- 基础断言：`equal`, `above`, `below`等
- 事件断言：`.to.emit()`和`.withArgs()`
- 回退断言：`.to.be.revertedWith()`等
- 余额变化断言：`.to.changeEtherBalance()`

### 4. 时间操作

- 使用`time.increase()`增加时间
- 使用`time.increaseTo()`跳转到特定时间
- 使用`time.setNextBlockTimestamp()`设置时间戳

### 5. 区块操作

- 使用`mine()`挖掘区块
- 使用`mineUpTo()`挖掘到特定区块号

## 最佳实践

1. **描述性命名**：使用清晰的测试名称
2. **按功能分组**：使用嵌套describe组织测试
3. **保持独立**：每个测试应该能够独立运行
4. **全面覆盖**：测试正常流程、异常情况和边界条件
5. **使用Fixture**：使用loadFixture而不是beforeEach
6. **验证状态**：确保状态变化正确
7. **测试事件**：验证事件触发和参数
8. **错误处理**：测试所有回退路径

## 常见问题

### 忘记使用await

```typescript
// 错误
expect(counter.inc()).to.emit(counter, "Increment");

// 正确
await expect(counter.inc()).to.emit(counter, "Increment");
```

### 类型不匹配

```typescript
// 错误：number类型
expect(await counter.x()).to.equal(0);

// 正确：bigint类型
expect(await counter.x()).to.equal(0n);
```

### 事件名称错误

```typescript
// 错误：大小写不匹配
.to.emit(counter, "increment");

// 正确
.to.emit(counter, "Increment");
```

### 错误消息不匹配

```typescript
// 错误消息必须完全匹配
await expect(counter.incBy(0))
  .to.be.revertedWith("incBy: increment should be positive");
```

## 学习资源

- [Hardhat测试文档](https://hardhat.org/docs/testing)
- [Chai断言库](https://www.chaijs.com/)
- [Mocha测试框架](https://mochajs.org/)
- [Hardhat Network Helpers](https://hardhat.org/hardhat-network-helpers/docs/overview)

## 许可证

MIT

