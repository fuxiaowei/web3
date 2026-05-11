import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
// 引入 Ignition 模块构建函数，用来组织代理合约部署步骤。

// 第一层模块：部署实现合约 + 透明代理 + ProxyAdmin。
const metaNFTAuctionProxyModule = buildModule(
  "MetaNFTAuctionProxyModule",
  (m) => {
    // 取第 0 个账户作为代理管理员（Transparent Proxy admin）。
    // 注意：该地址不能直接调用实现合约函数。
    const proxyAdminOwner = m.getAccount(0);
    // 取第 1 个账户作为业务管理员（MetaNFTAuction.admin）。
    // 该地址用于调用 onlyAdmin 的业务方法（setTokenOracle/start）。
    const appAdmin = m.getAccount(1);

    // 部署拍卖合约实现（逻辑合约，不直接对外交互）。
    const auctionImpl = m.contract("MetaNFTAuction");

    // 编码 initialize(appAdmin) 初始化调用，供代理构造时执行。
    const encodedFunctionCall = m.encodeFunctionCall(
      auctionImpl,
      "initialize",
      [appAdmin],
    );

    // 部署透明代理：参数分别是实现地址、管理员地址、初始化 calldata。
    const proxy = m.contract("TransparentUpgradeableProxy", [
      auctionImpl,
      proxyAdminOwner,
      encodedFunctionCall,
    ]);

    // 从代理部署触发的 AdminChanged 事件中读取新管理员地址。
    const proxyAdminAddress = m.readEventArgument(
      proxy,
      "AdminChanged",
      "newAdmin",
    );

    // 通过地址绑定 ProxyAdmin 合约实例，后续用于升级逻辑。
    const proxyAdmin = m.contractAt("ProxyAdmin", proxyAdminAddress);

    // 导出代理和管理员实例，给上层模块复用。
    return { proxyAdmin, proxy };
  },
);

// 第二层模块：基于上面的代理模块，返回业务合约类型化实例。
const metaNFTAuctionModule = buildModule("MetaNFTAuctionModule", (m) => {
  // 复用已定义的代理部署模块，拿到 proxy 与 proxyAdmin。
  const { proxy, proxyAdmin } = m.useModule(metaNFTAuctionProxyModule);

  // 将代理地址按 MetaNFTAuction ABI 进行绑定，便于后续函数调用。
  const auction = m.contractAt("MetaNFTAuction", proxy);

  // 导出业务实例、代理实例、管理员实例。
  return { auction, proxy, proxyAdmin };
});

// 默认导出上层模块，一般部署时直接使用它。
export default metaNFTAuctionModule;
