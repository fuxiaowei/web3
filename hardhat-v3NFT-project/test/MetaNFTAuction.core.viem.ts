import { expect } from "chai";
import { getAddress, parseUnits, zeroAddress } from "viem";
import { deployAuctionFixture, expectRevert, increaseTime } from "./helpers/auctionTestSetup.js";

/**
 * MetaNFTAuction 主流程测试
 *
 * 这份文件关注 V1 的核心业务闭环：
 * - 读取基础信息（版本、价格）
 * - 创建拍卖（start）
 * - 参与竞价（bid）
 * - 到期结算（end）
 *
 * 为什么按“主流程”独立文件：
 * 1) 让拍卖业务测试与升级测试解耦，阅读时上下文更干净。
 * 2) 当业务逻辑频繁改动时，只需要关注本文件，不被 upgrade 噪音干扰。
 * 3) CI 失败定位更快：看到文件名就知道是业务回归还是升级回归。
 */
describe("MetaNFTAuction", function () {
    let env: any;

    beforeEach(async function () {
        // 每个用例都重新部署一套链上状态，避免用例相互影响。
        // 好处：测试可并行执行且可重复运行，定位失败更稳定。
        env = await deployAuctionFixture();
    });

    describe("getVersion", function () {
        it("should return MetaNFTAuctionV1", async function () {
            // 版本号是升级体系里最直观的“实现标识”，先验证代理当前指向 V1。
            const version = await env.auction.read.getVersion();
            expect(version).to.equal("MetaNFTAuctionV1");
        });
    });

    describe("getPriceInDollar", function () {
        it("should return correct prices", async function () {
            // 同时读取 ETH 与 USDC 的喂价，证明 token->oracle 映射生效。
            const ethPrice = await env.auction.read.getPriceInDollar([zeroAddress]);
            const usdcPrice = await env.auction.read.getPriceInDollar([env.usdc.address]);
            // viem 默认返回 bigint，直接使用 bigint 比较可避免精度丢失。
            expect(ethPrice > 0n).to.equal(true);
            expect(usdcPrice > 0n).to.equal(true);
        });
    });

    describe("initialize", function () {
        it("should fail when initialized twice", async function () {
            // 代理构造时已经初始化过一次，再次 initialize 必须失败。
            // 好处：防止管理员被恶意重置，属于升级合约安全底线。
            await expectRevert(
                env.auction.write.initialize([env.admin.account.address], { account: env.admin.account }),
                "already initialized"
            );
        });
    });

    describe("start", function () {
        it("should fail when not called by admin", async function () {
            // 验证权限边界：只有 admin 可创建拍卖，避免任意地址创建恶意拍卖单。
            await expectRevert(
                env.auction.write.start(
                    [env.seller.account.address, 1n, env.nft.address, 1000n, 3600n, env.usdc.address],
                    { account: env.seller.account }
                ),
                "not admin"
            );
        });

        it("should increment auctionId", async function () {
            // 连续创建两场拍卖，验证编号单调递增。
            // 好处：前端索引、事件订阅、链下任务都依赖 auctionId 的稳定递增语义。
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            let auctionId = await env.auction.read.auctionId();
            expect(auctionId).to.equal(1n);

            await env.auction.write.start(
                [env.seller.account.address, 2n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            auctionId = await env.auction.read.auctionId();
            expect(auctionId).to.equal(2n);
        });
    });

    describe("bid", function () {
        it("should fail when auction has ended", async function () {
            // 先创建短拍卖，再推进链上时间模拟“自然结束”。
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;
            const auctionData = await env.auction.read.auctions([currentAuctionId]);
            await increaseTime(env.networkConnection, Number(auctionData[6]) + 1);

            // 结束后继续出价应失败，防止过期竞拍破坏结算公平性。
            await expectRevert(
                env.auction.write.bid([currentAuctionId, parseUnits("1", 18)], {
                    account: env.seller.account,
                    value: parseUnits("1", 18)
                }),
                "ended"
            );
        });

        it("should fail when bid is lower than highest bid", async function () {
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;

            // 先建立有效最高价。
            await env.auction.write.bid([currentAuctionId, parseUnits("2", 18)], {
                account: env.seller.account,
                value: parseUnits("2", 18)
            });

            // 更低价格再次出价必须被拒绝，保证价格发现机制正确。
            await expectRevert(
                env.auction.write.bid([currentAuctionId, parseUnits("1.2", 18)], {
                    account: env.bidder1.account,
                    value: parseUnits("1.2", 18)
                }),
                "invalid highestBid"
            );
        });

        it("should correctly track bidding result", async function () {
            // 多轮交替竞价，模拟真实拍卖中的反超场景。
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;

            await env.auction.write.bid([currentAuctionId, parseUnits("2", 18)], {
                account: env.bidder1.account,
                value: parseUnits("2", 18)
            });
            await env.auction.write.bid([currentAuctionId, parseUnits("3", 18)], {
                account: env.bidder2.account,
                value: parseUnits("3", 18)
            });
            await env.auction.write.bid([currentAuctionId, parseUnits("4", 18)], {
                account: env.bidder1.account,
                value: parseUnits("4", 18)
            });

            // 断言最终领先者和最高价，验证状态写入与覆盖逻辑正确。
            const auctionData = await env.auction.read.auctions([currentAuctionId]);
            expect(auctionData[4]).to.equal(getAddress(env.bidder1.account.address));
            expect(auctionData[8]).to.equal(parseUnits("4", 18));
        });
    });

    describe("end", function () {
        it("should fail when auction is not ended yet", async function () {
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 3600n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;

            // 未到时间提前结算应失败，防止人为提前截胡竞拍。
            await expectRevert(
                env.auction.write.end([currentAuctionId], { account: env.admin.account }),
                "not ended"
            );
        });

        it("should fail when there are no bids after auction ended", async function () {
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;

            await increaseTime(env.networkConnection, 31);

            // 无人出价就结算应失败，避免 NFT 和资金出现异常流转。
            await expectRevert(
                env.auction.write.end([currentAuctionId], { account: env.admin.account }),
                "no bids"
            );
        });

        it("should settle ETH auction: transfer NFT to winner and ETH to seller", async function () {
            // ETH 结算路径：验证“资产换手”和“资金结算”两件事同时成立。
            await env.auction.write.start(
                [env.seller.account.address, 1n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;
            const bidAmount = parseUnits("2", 18);

            const sellerBalanceBefore = await env.networkConnection.ethers.provider.getBalance(env.seller.account.address);
            await env.auction.write.bid([currentAuctionId, bidAmount], {
                account: env.bidder1.account,
                value: bidAmount
            });

            // 到期后任意地址触发 end（合约设计如此），重点检查结算结果而非调用者身份。
            await increaseTime(env.networkConnection, 31);
            await env.auction.write.end([currentAuctionId], { account: env.bidder2.account });

            // 卖家 ETH 余额增加，证明资金结算路径正确。
            const sellerBalanceAfter = await env.networkConnection.ethers.provider.getBalance(env.seller.account.address);
            expect(sellerBalanceAfter - sellerBalanceBefore).to.equal(bidAmount);

            // NFT 所有权切换到最高出价者，证明资产交割路径正确。
            const nftOwner = await env.nft.read.ownerOf([1n]);
            expect(getAddress(nftOwner)).to.equal(getAddress(env.bidder1.account.address));
        });

        it("should use ERC20 path for bidding, refund and settlement", async function () {
            // ERC20 路径覆盖三件关键事：
            // 1) 出价扣款；2) 被反超退款；3) 最终卖家收款。
            await env.auction.write.start(
                [env.seller.account.address, 2n, env.nft.address, 1000n, 30n, env.usdc.address],
                { account: env.admin.account }
            );
            const currentAuctionId = (await env.auction.read.auctionId()) - 1n;

            const bid1 = parseUnits("2000", 6);
            const bid2 = parseUnits("2500", 6);

            await env.usdc.write.mint([env.bidder1.account.address, parseUnits("5000", 6)], { account: env.admin.account });
            await env.usdc.write.mint([env.bidder2.account.address, parseUnits("5000", 6)], { account: env.admin.account });

            await env.usdc.write.approve([env.auction.address, bid1], { account: env.bidder1.account });
            await env.usdc.write.approve([env.auction.address, bid2], { account: env.bidder2.account });

            // 首次出价后余额减少，证明 transferFrom 托管成功。
            const bidder1BalanceBefore = await env.usdc.read.balanceOf([env.bidder1.account.address]);
            await env.auction.write.bid([currentAuctionId, bid1], { account: env.bidder1.account });
            const bidder1BalanceAfterFirstBid = await env.usdc.read.balanceOf([env.bidder1.account.address]);
            expect(bidder1BalanceAfterFirstBid).to.equal(bidder1BalanceBefore - bid1);

            // 被反超后余额恢复，证明“上一名领先者退款”逻辑正确。
            await env.auction.write.bid([currentAuctionId, bid2], { account: env.bidder2.account });
            const bidder1BalanceAfterOutbid = await env.usdc.read.balanceOf([env.bidder1.account.address]);
            expect(bidder1BalanceAfterOutbid).to.equal(bidder1BalanceBefore);

            // 结算后卖家收款并完成 NFT 交割。
            const sellerTokenBefore = await env.usdc.read.balanceOf([env.seller.account.address]);
            await increaseTime(env.networkConnection, 31);
            await env.auction.write.end([currentAuctionId], { account: env.admin.account });
            const sellerTokenAfter = await env.usdc.read.balanceOf([env.seller.account.address]);
            expect(sellerTokenAfter - sellerTokenBefore).to.equal(bid2);

            const nftOwner = await env.nft.read.ownerOf([2n]);
            expect(getAddress(nftOwner)).to.equal(getAddress(env.bidder2.account.address));

            // 附加状态断言：记录的支付代币和最高出价金额都应一致。
            const auctionData = await env.auction.read.auctions([currentAuctionId]);
            expect(getAddress(auctionData[10])).to.equal(getAddress(env.usdc.address));
            expect(auctionData[8]).to.equal(bid2);
        });
    });
});
