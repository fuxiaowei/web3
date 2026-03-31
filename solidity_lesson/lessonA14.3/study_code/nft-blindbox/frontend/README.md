# NFT盲盒前端项目

这是一个基于 Next.js 14 构建的 NFT 盲盒前端应用，与 NFTBlindBoxUpgradeable 合约完全匹配。

## 功能特性

- **购买盲盒** - 支持白名单和公售阶段
- **NFT展示** - 精美的卡片式展示，支持已揭示和未揭示状态
- **稀有度系统** - 展示 Common、Rare、Epic、Legendary 四种稀有度
- **实时统计** - 显示销售进度、已售数量等
- **钱包连接** - 支持 MetaMask 钱包连接
- **多网络支持** - 支持 localhost、Sepolia、Mainnet
- **现代化UI** - 使用 Tailwind CSS 和 Framer Motion 构建

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 配置环境变量

复制 `env.example` 为 `.env.local`：

```bash
cp env.example .env.local
```

编辑 `.env.local`，填入合约地址：

```env
NEXT_PUBLIC_CONTRACT_ADDRESS_SEPOLIA=0x你的合约地址
NEXT_PUBLIC_DEFAULT_NETWORK=sepolia
```

### 3. 运行开发服务器

```bash
npm run dev
```

访问 [http://localhost:3000](http://localhost:3000)

## 环境变量配置

### 必需配置

- `NEXT_PUBLIC_CONTRACT_ADDRESS_SEPOLIA` - Sepolia 测试网合约地址
- `NEXT_PUBLIC_CONTRACT_ADDRESS_MAINNET` - 主网合约地址（可选）
- `NEXT_PUBLIC_CONTRACT_ADDRESS_LOCALHOST` - 本地网络合约地址（可选）

### 可选配置

- `NEXT_PUBLIC_DEFAULT_NETWORK` - 默认网络（localhost | sepolia | mainnet）
- `NEXT_PUBLIC_RPC_URL` - 自定义 RPC URL
- `NEXT_PUBLIC_APP_NAME` - 应用名称
- `NEXT_PUBLIC_APP_DESCRIPTION` - 应用描述
- `NEXT_PUBLIC_IPFS_GATEWAY` - IPFS 网关地址
- `NEXT_PUBLIC_THEME_COLOR` - 主题色（十六进制，不带#）

## 项目结构

```
frontend/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # 根布局
│   ├── page.tsx           # 主页面
│   └── globals.css        # 全局样式
├── components/             # React 组件
│   ├── Header.tsx         # 头部导航
│   ├── Web3Provider.tsx   # Web3 上下文
│   ├── BlindBoxCard.tsx  # NFT 卡片
│   ├── BlindBoxShowcase.tsx # 盲盒展示区
│   ├── MyNFTs.tsx        # 我的NFT列表
│   ├── PurchaseModal.tsx  # 购买模态框
│   └── WalletConnectPrompt.tsx # 钱包连接提示
├── hooks/                  # 自定义 Hooks
│   └── useNFTBlindBox.ts  # NFT盲盒交互Hook
├── lib/                    # 工具函数
│   ├── config.ts          # 配置文件
│   ├── contracts.ts       # 合约工具
│   ├── abis.ts            # 合约ABI
│   └── format.ts          # 格式化工具
└── types/                  # TypeScript 类型
    └── window.d.ts        # Window 类型扩展
```

## 主要功能

### 购买盲盒

1. 连接钱包
2. 查看销售信息和价格
3. 点击"立即购买"按钮
4. 确认交易

### 查看我的NFT

- 自动加载用户拥有的所有NFT
- 显示已揭示和未揭示状态
- 展示稀有度信息

### 销售状态

- 实时显示销售进度
- 显示已售数量和剩余数量
- 支持白名单和公售阶段

## 样式定制

### 主题色

在 `.env.local` 中设置 `NEXT_PUBLIC_THEME_COLOR` 可以自定义主题色：

```env
NEXT_PUBLIC_THEME_COLOR=a855f7
```

### Tailwind 配置

编辑 `tailwind.config.js` 可以自定义颜色、动画等。

## 构建生产版本

```bash
npm run build
npm start
```

## 技术栈

- **Next.js 14** - React 框架
- **TypeScript** - 类型安全
- **Tailwind CSS** - 样式框架
- **Framer Motion** - 动画库
- **Ethers.js v6** - 以太坊交互
- **Lucide React** - 图标库

## 注意事项

1. 确保合约已正确部署
2. 确保网络配置正确
3. 生产环境建议使用环境变量管理敏感信息
4. IPFS 图片需要配置正确的网关地址

## 许可证

MIT

