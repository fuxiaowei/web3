# 新手流程指南（Ganache）

本指南对应脚本：`scripts/beginner.flow.viem.ts`

## 1) 先准备配置

1. 复制配置模板：
   - `.env.beginner.example` -> `.env`
2. 打开 `.env`，至少替换这几项：
   - `ADMIN_PRIVATE_KEY`
   - `SELLER_PRIVATE_KEY`
   - `BIDDER_PRIVATE_KEY`
   - `ENDER_PRIVATE_KEY`（可选）
   - `SELLER_ADDRESS`（可选，不填自动推导）

> 其余地址你当前项目已经是 Ganache 可用值。

## 2) 一键跑完整流程（推荐）

```powershell
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- all
```

该命令会依次执行：

1. `set-oracle`
2. `mint-nft`
3. `approve-nft`
4. `start`
5. `bid-eth`
6. `end`
7. `status`

## 3) 分步执行（排错时用）

```powershell
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- set-oracle
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- mint-nft
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- approve-nft
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- start
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- bid-eth
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- end
npx hardhat run scripts/beginner.flow.viem.ts --network ganache -- status
```

## 4) 常见问题

- `not admin`：
  - `ADMIN_PRIVATE_KEY` 不是拍卖合约 admin。
- `ERC721: caller is not token owner or approved`：
  - `approve-nft` 没做成功，或 `SELLER_ADDRESS` / `SELLER_PRIVATE_KEY` 不匹配。
- `not ended`：
  - 拍卖还没到期。默认 `AUTO_TIME_TRAVEL=true` 会在本地链自动快进时间。

