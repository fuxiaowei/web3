# NFT 拍卖市场测试报告

## 1. 项目与环境信息
- 项目名称：`hardhat-v3nft-project`
- 测试时间：`YYYY-MM-DD HH:mm`
- Node.js 版本：`vXX.X.X`
- Hardhat 版本：`v3.x`
- Solidity 版本：`0.8.19`
- 测试网络：`Hardhat Network / Ganache(1337)`

## 2. 测试范围
本次测试覆盖以下模块：

### 2.1 合约模块
- `MetaNFTAuction.sol`（主拍卖逻辑）
- `MetaNFTAuctionV2.sol`（升级逻辑）
- `MetaNFT.sol`（NFT）
- `MockERC20.sol`（支付代币）
- `MockOracle.sol`（价格预言机 Mock）

### 2.2 功能范围
- 初始化与权限控制
- 拍卖创建（start）
- 出价（ETH / ERC20）
- 结算（end）
- 价格换算（Chainlink 接口）
- 代理升级与升级后回归

## 3. 测试用例结果
### 3.1 主流程测试（`test/MetaNFTAuction.core.viem.ts`）
- [x] 获取版本号
- [x] 获取价格
- [x] 重复初始化失败
- [x] 非管理员创建拍卖失败
- [x] auctionId 递增
- [x] 拍卖结束后禁止出价
- [x] 低价出价失败
- [x] 多轮竞价状态正确
- [x] 未到期禁止结算
- [x] 无出价禁止结算
- [x] ETH 结算（NFT 与资金流转）
- [x] ERC20 结算（扣款、退款、卖家收款）

### 3.2 升级测试（`test/MetaNFTAuctionV2.viem.ts`）
- [x] 升级成功（状态保留 + 新版本可读）
- [x] 非管理员升级失败
- [x] 升级后管理员功能可正常调用（修改预言机并生效）

## 4. 测试命令与结果
- 执行命令：`npx hardhat test`
- 执行结果：`15 passing`
- 失败用例：`0`

## 5. 覆盖率结果
- 覆盖率文件：`coverage/lcov.info`
- 行覆盖率（Line Coverage）：`92.38%`（97/105）
- 说明：`MetaNFTAuction` 覆盖较高，`MetaNFT` 和 `MockOracle` 仍有少量函数边界待补

## 6. 已知限制与后续优化
- `MetaNFT.mintNext / burn` 需增加独立测试
- `MockOracle.setPrice / getPrice` 需增加独立测试
- 可增加更多边界场景（allowance 不足、amount/value 不一致、不存在拍卖ID等）

## 7. 结论
核心功能（NFT 拍卖、价格换算、透明代理升级）已实现并通过测试，满足课程/作业主体要求。
