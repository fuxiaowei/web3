# 多签钱包前端应用

一个现代化、功能完整的多签钱包管理前端应用，基于 Next.js 14 和 Tailwind CSS 构建。

## 功能特性

- ✅ **钱包连接**：支持 MetaMask 钱包连接
- ✅ **钱包管理**：查看钱包余额、所有者列表、确认阈值
- ✅ **交易管理**：
  - 创建交易提案
  - 确认/撤销确认交易
  - 执行交易
  - 查看交易详情和状态
- ✅ **所有者管理**：
  - 添加所有者
  - 删除所有者
  - 修改确认阈值
- ✅ **实时更新**：监听链上事件，自动刷新数据
- ✅ **现代化 UI**：美观的界面设计，支持深色模式
- ✅ **响应式设计**：适配各种屏幕尺寸

## 技术栈

- **框架**: Next.js 14 (App Router)
- **语言**: TypeScript
- **样式**: Tailwind CSS
- **Web3**: Ethers.js v6
- **图标**: Lucide React
- **日期处理**: date-fns

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 配置环境变量（可选）

如果需要自定义网络配置，可以创建 `.env.local` 文件：

```env
NEXT_PUBLIC_NETWORK=sepolia
NEXT_PUBLIC_RPC_URL=your_rpc_url
```

### 3. 启动开发服务器

```bash
npm run dev
```

应用将在 [http://localhost:3000](http://localhost:3000) 启动。

### 4. 构建生产版本

```bash
npm run build
npm start
```

## 使用说明

### 连接钱包

1. 确保已安装 MetaMask 浏览器扩展
2. 点击"连接钱包"按钮
3. 在 MetaMask 中确认连接
4. 确保连接到 Sepolia 测试网络

### 管理多签钱包

1. **输入钱包地址**：在首页输入已部署的多签钱包合约地址
2. **查看钱包信息**：查看余额、所有者数量、确认阈值等
3. **创建交易**：点击"创建交易"按钮，填写接收地址和金额
4. **确认交易**：作为所有者，可以确认或撤销确认交易
5. **执行交易**：当确认数达到阈值时，可以执行交易
6. **管理所有者**：在"所有者"标签页中添加或删除所有者，修改阈值

## 项目结构

```
frontend/
├── app/
│   ├── globals.css          # 全局样式
│   ├── layout.tsx           # 根布局
│   └── page.tsx             # 首页
├── components/
│   ├── AccountDetailsModal.tsx    # 账户详情模态框
│   ├── AlertProvider.tsx          # 通知提供者
│   ├── CreateTransactionModal.tsx # 创建交易模态框
│   ├── Header.tsx                 # 头部组件
│   ├── OwnersManagement.tsx       # 所有者管理组件
│   ├── TransactionCard.tsx         # 交易卡片
│   ├── TransactionList.tsx         # 交易列表
│   ├── WalletConnectPrompt.tsx    # 钱包连接提示
│   ├── WalletDashboard.tsx        # 钱包仪表板
│   ├── WalletInfo.tsx             # 钱包信息卡片
│   ├── WalletInput.tsx            # 钱包地址输入
│   └── Web3Provider.tsx            # Web3 提供者
├── hooks/
│   ├── useMultiSigWallet.ts        # 多签钱包 Hook
│   └── useTransactions.ts          # 交易 Hook
├── lib/
│   └── abis.ts                     # 合约 ABI
└── types/
    └── window.d.ts                  # Window 类型定义
```

## 设计特点

- **现代化设计**：采用渐变背景、毛玻璃效果、流畅动画
- **深色模式支持**：自动适配系统主题
- **响应式布局**：完美适配桌面和移动设备
- **交互反馈**：按钮悬停效果、加载状态、成功/错误提示
- **信息层次**：清晰的信息架构和视觉层次

## 注意事项

1. **网络要求**：默认配置为 Sepolia 测试网络，如需使用其他网络，请修改 `Web3Provider.tsx` 中的网络配置
2. **合约地址**：确保输入的多签钱包地址是已部署的合约地址
3. **权限要求**：只有钱包的所有者才能执行管理操作（创建交易、确认交易、管理所有者等）
4. **Gas 费用**：所有链上操作都需要支付 Gas 费用

## 开发建议

- 使用 TypeScript 严格模式进行类型检查
- 遵循 React Hooks 最佳实践
- 使用 Tailwind CSS 工具类进行样式设计
- 保持组件的单一职责原则
- 及时处理错误和边界情况

## 许可证

MIT License

