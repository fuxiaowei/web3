# 众筹平台前端

基于 Next.js 14 和 React 18 构建的现代化众筹平台前端界面。

## 技术栈

- **Next.js 14** - React 框架
- **TypeScript** - 类型安全
- **Tailwind CSS** - 样式框架
- **Ethers.js v6** - 以太坊交互
- **Lucide React** - 图标库
- **date-fns** - 日期处理

## 快速开始

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `.env.local.example` 为 `.env.local` 并填入配置：

```bash
cp .env.local.example .env.local
```

编辑 `.env.local`：

```env
NEXT_PUBLIC_FACTORY_ADDRESS=你的工厂合约地址
NEXT_PUBLIC_CHAIN_ID=31337
NEXT_PUBLIC_RPC_URL=http://127.0.0.1:8545
```

### 运行开发服务器

```bash
npm run dev	
```

打开 [http://localhost:3000](http://localhost:3000) 查看应用。

## 项目结构

```
frontend/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # 根布局
│   ├── page.tsx           # 首页
│   ├── campaign/          # 项目详情页
│   └── globals.css        # 全局样式
├── components/            # React 组件
│   ├── Header.tsx         # 头部导航
│   ├── CampaignCard.tsx   # 项目卡片
│   ├── CreateCampaignModal.tsx  # 创建项目弹窗
│   └── Web3Provider.tsx  # Web3 上下文
├── hooks/                 # 自定义 Hooks
│   ├── useWeb3.ts         # Web3 连接
│   ├── useCampaigns.ts    # 项目列表
│   ├── useCampaign.ts     # 项目详情
│   └── useFactory.ts      # 工厂合约
└── lib/                   # 工具函数
    └── abis.ts            # 合约 ABI
```

## 功能说明

### 首页

- 显示所有众筹项目
- 项目卡片展示进度、状态、贡献者数量
- 支持创建新项目
- 实时数据更新

### 项目详情页

- 显示项目完整信息
- 进度条可视化
- 支持贡献资金
- 支持提取资金（项目创建者）
- 支持退款（失败项目）

### 钱包连接

- MetaMask 集成
- 自动检测网络
- 账户切换监听
- 连接状态持久化

## 部署

### 构建生产版本

```bash
npm run build
```

### 启动生产服务器

```bash
npm start
```

### Vercel 部署

项目可以直接部署到 Vercel：

1. 推送代码到 GitHub
2. 在 Vercel 中导入项目
3. 配置环境变量
4. 部署

## 注意事项

1. **合约地址配置**：确保 `.env.local` 中的 `NEXT_PUBLIC_FACTORY_ADDRESS` 正确
2. **网络配置**：确保 RPC URL 和 Chain ID 正确
3. **MetaMask**：用户需要安装 MetaMask 浏览器扩展
4. **测试网 ETH**：在测试网部署时需要测试网 ETH

## 开发建议

- 使用 TypeScript 严格模式
- 遵循 React Hooks 最佳实践
- 使用 Tailwind CSS 工具类
- 保持组件小而专注
- 使用自定义 Hooks 复用逻辑

## 许可证

MIT License

