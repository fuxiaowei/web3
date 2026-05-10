import { expect } from "chai";
import { parseUnits, zeroAddress } from "viem";
import { deployAuctionFixture, expectRevert } from "./helpers/auctionTestSetup.js";

/**
 * MetaNFTAuction 升级测试（V2）
 *
 * 为什么单独拆出升级测试文件：
 * 1) 升级流程涉及 ProxyAdmin / Implementation / ABI 切换，和日常业务断言关注点不同。
 * 2) 当升级失败时可以快速定位是“治理权限问题”还是“业务逻辑问题”。
 * 3) 团队评审时可独立审查“存储兼容性 + 新功能生效”这类高风险项。
 */
describe("MetaNFTAuctionV2", function () {
    let env: any;

    beforeEach(async function () {
        // 每个升级用例都从全新 V1 状态开始，防止升级副作用污染其他用例。
        env = await deployAuctionFixture();
    });

    describe("upgrade", function () {
        it("should upgrade contract successfully", async function () {
            // 升级前先写入业务状态，后续用于验证“升级不丢数据”。
            await env.auction.write.start(
                [env.seller.account.address, 10n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const oldAuctionId = await env.auction.read.auctionId();
            // 部署新实现合约（V2）。
            const newImpl = await env.viem.deployContract("MetaNFTAuctionV2");

            // 以 ITransparentUpgradeableProxy ABI 访问代理，调用 upgrade 所需接口。
            const proxyAsV2 = await env.viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                env.proxy.address
            );

            // 由 ProxyAdmin owner 执行升级，贴近生产治理路径。
            await env.proxyAdmin.write.upgrade(
                [proxyAsV2.address, newImpl.address],
                { account: env.proxyAdminSigner.account }
            );

            // 升级后仍是同一代理地址，但需要切到 V2 ABI 读取新接口。
            const upgradedAuction = await env.viem.getContractAt("MetaNFTAuctionV2", env.proxy.address);

            // 核心回归1：旧状态不丢（存储布局兼容）。
            const newAuctionId = await upgradedAuction.read.auctionId();
            expect(newAuctionId).to.equal(oldAuctionId);

            // 核心回归2：版本标识变化，证明实现切换成功。
            const version = await upgradedAuction.read.getVersion();
            expect(version).to.equal("MetaNFTAuctionV2");

            // 核心回归3：V2 新增功能可调用，证明新逻辑对外可见。
            const newFeature = await upgradedAuction.read.newFeature();
            expect(newFeature).to.equal("This is a new feature in V2");
        });

        it("should fail when non-admin tries to upgrade", async function () {
            await env.auction.write.start(
                [env.seller.account.address, 10n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const newImpl = await env.viem.deployContract("MetaNFTAuctionV2");

            const proxyAsV2 = await env.viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                env.proxy.address
            );

            // 权限负例：非 ProxyAdmin owner 不应拥有升级能力。
            // 好处：防止恶意实现替换导致资产或逻辑失控。
            await expectRevert(
                env.proxyAdmin.write.upgrade(
                    [proxyAsV2.address, newImpl.address],
                    { account: env.seller.account }
                )
            );
        });

        it("should change oracle after upgrade", async function () {
            // 升级后继续调用管理函数，验证“业务可用性”而不只是“版本可读性”。
            await env.auction.write.start(
                [env.seller.account.address, 10n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const newEthOracle = await env.viem.deployContract("MockOracle", [parseUnits("3000", 8)]);
            const newImpl = await env.viem.deployContract("MetaNFTAuctionV2");

            const proxyAsV2 = await env.viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                env.proxy.address
            );

            await env.proxyAdmin.write.upgrade(
                [proxyAsV2.address, newImpl.address],
                { account: env.proxyAdminSigner.account }
            );

            const upgradedAuction = await env.viem.getContractAt("MetaNFTAuctionV2", env.proxy.address);

            // 升级后更新 ETH 预言机地址，验证管理员写路径仍正常。
            await upgradedAuction.write.setTokenOracle([zeroAddress, newEthOracle.address], { account: env.admin.account });

            // 读回价格验证更新生效，确保新实现与旧存储配合正常。
            const newPrice = await upgradedAuction.read.getPriceInDollar([zeroAddress]);
            expect(newPrice).to.equal(parseUnits("3000", 8));
        });
    });
});
