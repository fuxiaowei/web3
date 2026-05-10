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
    // 把“人类可读美元值”转换为合约内部使用的 1e8 精度，
    // 与 highestBidInDollar / Chainlink 常见 8 位小数口径保持一致。
    const usd = (value: string) => parseUnits(value, 8);

    beforeEach(async function () {
        // 每个升级用例都从全新 V1 状态开始，防止升级副作用污染其他用例。
        env = await deployAuctionFixture();
    });

    const upgradeToV2 = async () => {
        // 统一升级步骤：
        // 1) 部署 V2 实现
        // 2) 通过 ProxyAdmin 升级代理
        // 3) 以 V2 ABI 绑定同一代理地址返回
        // 这样可以避免每个用例重复写升级样板代码。
        const newImpl = await env.viem.deployContract("MetaNFTAuctionV2");
        const proxyAsV2 = await env.viem.getContractAt(
            "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
            env.proxy.address
        );

        await env.proxyAdmin.write.upgrade(
            [proxyAsV2.address, newImpl.address],
            { account: env.proxyAdminSigner.account }
        );

        return env.viem.getContractAt("MetaNFTAuctionV2", env.proxy.address);
    };

    describe("upgrade", function () {
        it("should upgrade contract successfully", async function () {
            // 升级前先写入业务状态，后续用于验证“升级不丢数据”。
            await env.auction.write.start(
                [env.seller.account.address, 10n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const oldAuctionId = await env.auction.read.auctionId();
            // 部署新实现合约（V2）。
            // 升级后仍是同一代理地址，但需要切到 V2 ABI 读取新接口。
            const upgradedAuction = await upgradeToV2();

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
            const upgradedAuction = await upgradeToV2();

            // 升级后更新 ETH 预言机地址，验证管理员写路径仍正常。
            await upgradedAuction.write.setTokenOracle([zeroAddress, newEthOracle.address], { account: env.admin.account });

            // 读回价格验证更新生效，确保新实现与旧存储配合正常。
            const newPrice = await upgradedAuction.read.getPriceInDollar([zeroAddress]);
            expect(newPrice).to.equal(parseUnits("3000", 8));
        });
    });

    describe("dynamic fee by USD tiers", function () {
        it("should apply ETH fee based on highestBidInDollar", async function () {
            // 先创建拍卖，再升级到 V2，确保“升级前创建的数据”也能使用新结算逻辑。
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;
            const upgradedAuction = await upgradeToV2();

            // 分段规则：<=2000 美元收 3%，<=5000 美元收 2%，>5000 美元收 1%
            await upgradedAuction.write.setFeeRecipient([env.bidder2.account.address], { account: env.admin.account });
            await upgradedAuction.write.setDynamicFeeConfig(
                [usd("2000"), usd("5000"), 300, 200, 100],
                { account: env.admin.account }
            );

            const bidAmount = parseUnits("2", 18); // 2 ETH * 3000 = 6000 USD，命中第 3 档 1%
            const expectedFee = parseUnits("0.02", 18);

            const sellerBalanceBefore = await env.networkConnection.ethers.provider.getBalance(env.seller.account.address);
            const feeRecipientBalanceBefore = await env.networkConnection.ethers.provider.getBalance(env.bidder2.account.address);

            await upgradedAuction.write.bid([currentAuctionId, bidAmount], {
                account: env.bidder1.account,
                value: bidAmount
            });

            // 拍卖到期后执行结算，验证“卖家净收 + 平台抽佣”两条资金流。
            await env.networkConnection.ethers.provider.send("evm_increaseTime", [31]);
            await env.networkConnection.ethers.provider.send("evm_mine");
            await upgradedAuction.write.end([currentAuctionId], { account: env.admin.account });

            const sellerBalanceAfter = await env.networkConnection.ethers.provider.getBalance(env.seller.account.address);
            const feeRecipientBalanceAfter = await env.networkConnection.ethers.provider.getBalance(env.bidder2.account.address);

            expect(sellerBalanceAfter - sellerBalanceBefore).to.equal(bidAmount - expectedFee);
            expect(feeRecipientBalanceAfter - feeRecipientBalanceBefore).to.equal(expectedFee);
        });

        it("should apply ERC20 fee based on highestBidInDollar", async function () {
            // 与 ETH 场景同理，这里验证 ERC20 结算路径的动态手续费。
            await env.auction.write.start(
                [env.seller.account.address, 2n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;
            const upgradedAuction = await upgradeToV2();

            await upgradedAuction.write.setFeeRecipient([env.bidder2.account.address], { account: env.admin.account });
            // 2500 USDC 会落在第 2 档（2%）
            await upgradedAuction.write.setDynamicFeeConfig(
                [usd("2000"), usd("4000"), 300, 200, 100],
                { account: env.admin.account }
            );

            const bidAmount = parseUnits("2500", 6);
            const expectedFee = parseUnits("50", 6);

            // 准备 ERC20 资金和授权，确保 bid() 可成功 transferFrom。
            await env.usdc.write.mint([env.bidder1.account.address, parseUnits("5000", 6)], { account: env.admin.account });
            await env.usdc.write.approve([upgradedAuction.address, bidAmount], { account: env.bidder1.account });

            const sellerBefore = await env.usdc.read.balanceOf([env.seller.account.address]);
            const recipientBefore = await env.usdc.read.balanceOf([env.bidder2.account.address]);

            await upgradedAuction.write.bid([currentAuctionId, bidAmount], { account: env.bidder1.account });

            // 到期后结算，核对卖家与手续费接收者余额变化。
            await env.networkConnection.ethers.provider.send("evm_increaseTime", [31]);
            await env.networkConnection.ethers.provider.send("evm_mine");
            await upgradedAuction.write.end([currentAuctionId], { account: env.admin.account });

            const sellerAfter = await env.usdc.read.balanceOf([env.seller.account.address]);
            const recipientAfter = await env.usdc.read.balanceOf([env.bidder2.account.address]);

            expect(sellerAfter - sellerBefore).to.equal(bidAmount - expectedFee);
            expect(recipientAfter - recipientBefore).to.equal(expectedFee);
        });

        it("should fail when non-admin updates fee config", async function () {
            const upgradedAuction = await upgradeToV2();
            // 配置类接口受 onlyAdmin 保护，非管理员调用必须失败。
            await expectRevert(
                upgradedAuction.write.setDynamicFeeConfig(
                    [usd("2000"), usd("5000"), 300, 200, 100],
                    { account: env.seller.account }
                ),
                "not admin"
            );
        });
    });
});
