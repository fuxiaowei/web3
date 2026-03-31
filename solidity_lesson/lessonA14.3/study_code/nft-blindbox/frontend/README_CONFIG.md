# NFT盲盒前端配置说明

## 环境变量配置

### 必需配置

在 `.env.local` 文件中配置以下变量：

```env
# 合约地址（至少配置一个网络）
NEXT_PUBLIC_CONTRACT_ADDRESS_SEPOLIA=0x你的合约地址
NEXT_PUBLIC_CONTRACT_ADDRESS_MAINNET=0x你的合约地址
NEXT_PUBLIC_CONTRACT_ADDRESS_LOCALHOST=0x你的合约地址

# 默认网络
NEXT_PUBLIC_DEFAULT_NETWORK=sepolia
```

### 可选配置

```env
# 自定义RPC URL
NEXT_PUBLIC_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/your-api-key

# 应用信息
NEXT_PUBLIC_APP_NAME=Mystery NFT
NEXT_PUBLIC_APP_DESCRIPTION=神秘NFT盲盒，开启你的收藏之旅

# IPFS网关
NEXT_PUBLIC_IPFS_GATEWAY=https://ipfs.io/ipfs/

# 主题色（十六进制，不带#）
NEXT_PUBLIC_THEME_COLOR=a855f7
```

## 部署步骤

1. **部署合约**
   ```bash
   # 在合约项目目录
   npx hardhat run scripts/deployModules.ts --network sepolia
   npx hardhat run scripts/deployWithUUPS.ts --network sepolia
   ```

2. **配置前端环境变量**
   ```bash
   # 在frontend目录
   cp env.example .env.local
   # 编辑 .env.local，填入合约地址
   ```

3. **安装依赖并运行**
   ```bash
   npm install
   npm run dev
   ```

## 样式定制

### 主题色

通过环境变量 `NEXT_PUBLIC_THEME_COLOR` 可以自定义主题色。默认值为 `a855f7`（紫色）。

### 稀有度颜色

在 `lib/config.ts` 中的 `RARITY_CONFIG` 可以自定义稀有度颜色和样式。

### Tailwind配置

编辑 `tailwind.config.js` 可以自定义：
- 颜色主题
- 动画效果
- 背景图案

## 网络配置

支持的网络：
- `localhost` - 本地开发网络（Chain ID: 31337）
- `sepolia` - Sepolia 测试网（Chain ID: 11155111）
- `mainnet` - 以太坊主网（Chain ID: 1）

在 `lib/config.ts` 中可以添加更多网络配置。

## 功能说明

### 购买盲盒
- 自动检测销售状态
- 支持白名单和公售阶段
- 实时显示价格和剩余数量

### NFT展示
- 已揭示：显示NFT图片和稀有度
- 未揭示：显示盲盒样式
- 稀有度颜色区分

### 我的收藏
- 自动加载用户拥有的NFT
- 显示每个NFT的状态和稀有度

