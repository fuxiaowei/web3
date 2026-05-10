# 部署地址清单

## 1. 网络信息
- 网络名称：`Sepolia`
- Chain ID：`11155111`
- 部署工具：`Hardhat Ignition`
- 部署时间：`YYYY-MM-DD HH:mm`

## 2. 合约地址（Sepolia）

### 2.1 拍卖代理体系
- MetaNFTAuction 实现合约（V1）：
  - `0x67f081cdf945042F3979E52f0189541F2bd78F8e`
- TransparentUpgradeableProxy（对外拍卖地址）：
  - `0x0D46eEa7cd0E5aB73bdDC9E8a273CF1b2CEC027D`
- ProxyAdmin：
  - `0x39eE4a3E15058D4C9a70Fdf98bdbA11C15ad0982`
- MetaNFTAuctionV2 实现合约：
  - `0xa52799A91DF6fd7F115cE9F2804290066142c2A3`

### 2.2 其他依赖合约（如已部署）
- MetaNFT：
  - `<待填写>`
- MockERC20：
  - `<待填写>`
- MockOracle(ETH/USD)：
  - `<待填写>`
- MockOracle(ERC20/USD)：
  - `<待填写>`

## 3. 地址来源
- 文件：`ignition/deployments/chain-11155111/deployed_addresses.json`

## 4. 区块浏览器链接
- Proxy：
  - `https://sepolia.etherscan.io/address/0x0D46eEa7cd0E5aB73bdDC9E8a273CF1b2CEC027D`
- V1 Implementation：
  - `https://sepolia.etherscan.io/address/0x67f081cdf945042F3979E52f0189541F2bd78F8e`
- V2 Implementation：
  - `https://sepolia.etherscan.io/address/0xa52799A91DF6fd7F115cE9F2804290066142c2A3`
- ProxyAdmin：
  - `https://sepolia.etherscan.io/address/0x39eE4a3E15058D4C9a70Fdf98bdbA11C15ad0982`

## 5. 升级记录
- 升级前版本：`MetaNFTAuctionV1`
- 升级后版本：`MetaNFTAuctionV2`
- 升级方式：`Transparent Proxy + ProxyAdmin.upgrade/upgradeAndCall`
- 升级验证：
  - [x] 状态未丢失（auctionId 保持）
  - [x] 新功能可调用（newFeature）
