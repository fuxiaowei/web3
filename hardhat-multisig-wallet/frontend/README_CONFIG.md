# 配置文件使用说明

## 快速开始

### 1. 配置环境变量

复制 `.env.example` 为 `.env.local`：

```bash
cp .env.example .env.local
```

### 2. 填入合约地址

部署合约后，在 `.env.local` 中填入合约地址：

```env
# Sepolia 测试网
NEXT_PUBLIC_CONTRACT_ADDRESS_SEPOLIA=0x你的合约地址

# 或本地网络
NEXT_PUBLIC_CONTRACT_ADDRESS_LOCALHOST=0x你的合约地址
```

### 3. 配置网络（可选）

如果需要使用自定义 RPC，修改 `.env.local`：

```env
NEXT_PUBLIC_DEFAULT_NETWORK=sepolia
NEXT_PUBLIC_RPC_URL=https://your-custom-rpc-url
```

## 配置文件说明

### `lib/config.ts`

主配置文件，包含：
- **网络配置**：支持的区块链网络（localhost、sepolia、mainnet）
- **合约地址**：各网络的合约地址
- **工具函数**：获取网络配置、合约地址等

### `lib/contracts.ts`

合约工具函数：
- `getMultiSigWalletContract()` - 创建合约实例
- `isValidContractAddress()` - 验证地址
- `getContractExplorerUrl()` - 获取区块浏览器链接
- `formatAddress()` - 格式化地址显示

### `lib/network.ts`

网络管理工具：
- `switchNetwork()` - 切换网络
- `getCurrentNetwork()` - 获取当前网络
- `isNetworkSupported()` - 验证网络是否支持

## 使用示例

### 在组件中使用配置

```typescript
import { getContractAddress, getNetworkConfig } from '@/lib/config';
import { getMultiSigWalletContract } from '@/lib/contracts';
import { useWeb3 } from '@/components/Web3Provider';

function MyComponent() {
  const { signer, provider } = useWeb3();
  const network = 'sepolia';
  
  // 获取合约地址
  const contractAddress = getContractAddress(network);
  
  // 创建合约实例
  const contract = getMultiSigWalletContract(contractAddress!, signer!);
  
  // 使用合约...
}
```

### 切换网络

```typescript
import { switchNetwork } from '@/lib/network';

async function handleSwitchNetwork() {
  try {
    await switchNetwork('sepolia');
    console.log('网络切换成功');
  } catch (error) {
    console.error('网络切换失败:', error);
  }
}
```

### 获取区块浏览器链接

```typescript
import { getContractExplorerUrl } from '@/lib/contracts';
import { useWeb3 } from '@/components/Web3Provider';

function ContractLink({ address }: { address: string }) {
  const { provider } = useWeb3();
  const [chainId, setChainId] = useState<number | null>(null);
  
  useEffect(() => {
    provider?.getNetwork().then(network => {
      setChainId(Number(network.chainId));
    });
  }, [provider]);
  
  const explorerUrl = chainId ? getContractExplorerUrl(address, chainId) : null;
  
  return explorerUrl ? (
    <a href={explorerUrl} target="_blank" rel="noopener noreferrer">
      查看合约
    </a>
  ) : null;
}
```

## 部署后配置步骤

1. **部署合约**：使用 Hardhat 部署多签钱包合约
2. **获取地址**：记录部署的合约地址
3. **更新配置**：在 `.env.local` 中填入合约地址
4. **重启应用**：重启开发服务器使配置生效

## 环境变量说明

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `NEXT_PUBLIC_DEFAULT_NETWORK` | 默认网络 | `sepolia` |
| `NEXT_PUBLIC_RPC_URL` | 自定义 RPC URL | `https://...` |
| `NEXT_PUBLIC_CONTRACT_ADDRESS_LOCALHOST` | 本地网络合约地址 | `0x...` |
| `NEXT_PUBLIC_CONTRACT_ADDRESS_SEPOLIA` | Sepolia 测试网合约地址 | `0x...` |
| `NEXT_PUBLIC_CONTRACT_ADDRESS_MAINNET` | 主网合约地址 | `0x...` |

## 注意事项

1. **环境变量前缀**：所有前端环境变量必须以 `NEXT_PUBLIC_` 开头
2. **合约地址格式**：必须是有效的以太坊地址（0x 开头，42 字符）
3. **网络匹配**：确保合约地址与当前连接的网络匹配
4. **生产环境**：生产环境建议使用环境变量而非硬编码

## 故障排查

### 问题：找不到合约地址

**解决方案**：
1. 检查 `.env.local` 文件是否存在
2. 确认环境变量名称正确
3. 重启开发服务器

### 问题：网络不匹配

**解决方案**：
1. 确认当前连接的网络
2. 检查对应网络的合约地址是否配置
3. 使用 `switchNetwork()` 切换到正确网络

### 问题：RPC 连接失败

**解决方案**：
1. 检查 RPC URL 是否正确
2. 确认网络是否可访问
3. 尝试使用其他 RPC 提供商

