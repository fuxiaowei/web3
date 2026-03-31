# NFT市场合约代码

本目录包含完整的NFT市场实现代码，包括NFT合约、市场合约和版税支持。

## 文件说明

### 1. MyNFT.sol
基础的NFT合约实现，支持：
- ERC721标准功能
- NFT铸造
- 供应量控制
- 价格设置
- 元数据管理

### 2. NFTMarketplace.sol
完整的NFT交易市场合约，支持：
- NFT上架和下架
- 固定价格购买
- 英式拍卖
- ERC2981版税支持
- 安全防护机制

### 3. MyNFTWithRoyalty.sol
支持ERC2981版税标准的NFT合约，支持：
- 所有MyNFT.sol的功能
- ERC2981版税标准
- 版税自动分配

## 部署说明

### 环境要求
- Solidity 0.8.20
- OpenZeppelin Contracts 5.0.0+

### 部署步骤

1. **部署NFT合约**
   ```solidity
   // MyNFT.sol
   constructor() // 无需参数
   
   // MyNFTWithRoyalty.sol
   constructor(address royaltyReceiver, uint96 royaltyBps)
   ```

2. **部署市场合约**
   ```solidity
   constructor(address _feeRecipient)
   ```

3. **授权市场合约**
   ```solidity
   // 在NFT合约上调用
   setApprovalForAll(marketplaceAddress, true)
   ```

## 使用流程

### 铸造NFT
```solidity
// 支付0.01 ETH
mint("ipfs://QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx")
```

### 上架NFT
```solidity
// 在市场合约上调用
listNFT(nftContractAddress, tokenId, priceInWei)
```

### 购买NFT
```solidity
// 支付足够的ETH
buyNFT(listingId)
```

### 创建拍卖
```solidity
createAuction(nftContractAddress, tokenId, startPrice, durationHours)
```

### 出价
```solidity
// 支付足够的ETH
placeBid(auctionId)
```

## 安全特性

1. **重入攻击防护**：使用ReentrancyGuard
2. **CEI原则**：遵循检查-效果-交互模式
3. **权限验证**：严格验证所有权和授权
4. **资金安全**：使用低级call并检查返回值

## 测试建议

1. 在Remix IDE中测试基础功能
2. 在测试网（Sepolia/Goerli）上部署测试
3. 测试所有边界情况和异常场景
4. 验证版税分配是否正确
5. 测试拍卖流程的完整性

## 注意事项

1. 确保NFT合约已授权市场合约
2. 购买时支付足够的ETH
3. 拍卖结束后需要调用endAuction进行结算
4. 被超越的出价者需要调用withdrawBid提取资金

## Gas消耗参考

- 铸造NFT：~80,000 Gas
- 上架NFT：~50,000 Gas
- 购买NFT：~150,000 Gas
- 创建拍卖：~80,000 Gas
- 出价：~60,000 Gas

## 扩展功能建议

1. 批量操作（批量上架、批量购买）
2. 要约系统（买家主动出价）
3. 白名单销售
4. 租赁功能（ERC4907）
5. 盲盒机制（集成Chainlink VRF）

