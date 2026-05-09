import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
// 引入 Ignition 模块构建函数。

import MetaNFTAuctionModule from "./MetaNFTAuctionProxyModule.js";
// 引入已有的代理部署模块（包含 proxy 与 proxyAdmin）。

// 定义升级模块：将代理从 V1 实现升级到 V2 实现。
const metaNFTAuctionUpgradeModule = buildModule(
  "MetaNFTAuctionUpgradeModule",
  (m) => {
    // 使用第 0 个账户作为升级交易的发送者（应具备管理员权限）。
    const proxyAdminOwner = m.getAccount(0);

    // 复用代理部署模块，拿到代理地址与 ProxyAdmin 实例。
    const { proxyAdmin, proxy } = m.useModule(MetaNFTAuctionModule);

    // 部署新的实现合约 V2。
    const auctionV2 = m.contract("MetaNFTAuctionV2");

    // 调用 ProxyAdmin.upgradeAndCall 执行升级；"0x" 表示升级后不额外调用初始化数据。
    m.call(proxyAdmin, "upgradeAndCall", [proxy, auctionV2, "0x"], {
      from: proxyAdminOwner,
    });

    // 将代理地址按 V2 ABI 绑定，方便后续以 V2 接口访问同一个代理地址。
    const auction = m.contractAt("MetaNFTAuctionV2", proxy, {
      // 显式指定部署对象 ID，便于在部署产物中识别这个绑定实例。
      id: "MetaNFTAuctionV2AtProxy",
    });

    // 导出升级后的业务实例及相关管理实例。
    return { auction, proxyAdmin, proxy };
  },
);

// 默认导出升级模块。
export default metaNFTAuctionUpgradeModule;
