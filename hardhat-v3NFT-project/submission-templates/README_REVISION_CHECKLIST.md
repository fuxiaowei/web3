# README 修订清单（提交前逐项确认）

## 1. 版本与依赖一致性
- [ ] README 中 Solidity 版本与实际配置一致（当前应为 `0.8.19`）
- [ ] README 中 OpenZeppelin 版本与 `package.json` 一致（当前 `@openzeppelin/contracts ^4.9.6`）
- [ ] README 中 Hardhat 版本与 `package.json` 一致（当前 `hardhat ^3.2.0`）

## 2. 测试文件与命令一致性
- [ ] 测试文件名更新为：
  - `test/MetaNFTAuction.core.viem.ts`
  - `test/MetaNFTAuctionV2.viem.ts`
- [ ] 测试命令可直接执行（建议统一写 `npx hardhat test`）
- [ ] 删除 README 中不存在或过期的测试文件引用

## 3. 部署说明一致性
- [ ] 明确实际部署工具为 `Hardhat Ignition`
- [ ] 补充透明代理部署流程与升级流程
- [ ] 引用当前可用模块文件名：
  - `ignition/modules/MetaNFTAuctionProxyModule.ts`
  - `ignition/modules/MetaNFTAuctionUpgradeModule.ts`

## 4. 地址与验收证据
- [ ] 增加“Sepolia 部署地址”章节（可引用 `submission-templates/DEPLOYMENTS.md`）
- [ ] 写明地址来源文件：`ignition/deployments/chain-11155111/deployed_addresses.json`
- [ ] 提供关键地址的 Etherscan 链接

## 5. 任务要求对照
- [ ] 增加“任务要求 -> 当前实现 -> 证据文件”对照表
- [ ] 标注“已完成 / 部分完成 / 待完成”状态
- [ ] 明确当前与任务原文存在差异的点（如 `start` 当前是 onlyAdmin）

## 6. 结果可复现性
- [ ] README 中保留从安装到测试的最短复现路径
- [ ] 命令与当前目录结构一致
- [ ] 避免出现已删除文件、旧命令、旧截图

## 7. 提交前最终检查
- [ ] `npx hardhat test` 通过
- [ ] 文档中的所有路径、命令、文件名可在项目中找到
- [ ] 模板占位符（如 `<待填写>`、`YYYY-MM-DD`）已替换为真实值
