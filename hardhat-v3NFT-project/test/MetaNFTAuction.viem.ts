import { expect } from "chai";
import "@nomicfoundation/hardhat-viem";
import hre from "hardhat";
import { parseUnits, zeroAddress, encodeFunctionData, getAddress } from "viem";

describe("MetaNFTAuction", function () {
    // 统一的环境对象：在 beforeEach 中重新部署，确保用例互不污染
    let auction: any;
    let nft: any;
    let usdc: any;
    let ethOracle: any;
    let usdcOracle: any;
    let proxyAdmin: any;
    let proxy: any;

    let admin: any;
    let proxyAdminSigner: any;
    let seller: any;
    let bidder1: any;
    let bidder2: any;
    let networkConnection: any;
    let viem: any;

    // viem 抛错结构较深（shortMessage/details/message 分散在不同字段），
    // 这里统一做 revert 捕获与信息断言，避免每个测试重复写 try/catch。
    // - promise: 预期会失败的交易 Promise
    // - message: 可选，校验 revert 关键字（如 "not admin"、"ended"）
    const expectRevert = async (promise: Promise<unknown>, message?: string) => {
        try {
            await promise;
            expect.fail("Expected transaction to revert");
        } catch (error: any) {
            const errorText =
                error?.shortMessage ??
                error?.details ??
                error?.message ??
                String(error);
            if (message) {
                expect(errorText).to.include(message);
            }
        }
    };

    // 封装工具函数：在 EVM 层推进区块时间（秒）
    // 注意只改变本地测试链，不影响真实网络。
    const increaseTime = async (seconds: number) => {
        await networkConnection.ethers.provider.send("evm_increaseTime", [seconds]);
        await networkConnection.ethers.provider.send("evm_mine");
    };

    // 封装工具函数：直接把下一块时间设置为指定时间戳
    // 用于精确命中“拍卖结束时刻”等边界条件。
    const setTimestamp = async (target: bigint) => {
        await networkConnection.ethers.provider.send("evm_setNextBlockTimestamp", [Number(target)]);
        await networkConnection.ethers.provider.send("evm_mine");
    };

    beforeEach(async function () {
        // 连接本地测试网络并获取默认测试账户
        networkConnection = await hre.network.connect();
        viem = networkConnection.viem;
        const clients = await viem.getWalletClients();
        [admin, proxyAdminSigner, seller, bidder1, bidder2] = clients;

        // 部署逻辑合约（Implementation），后续由 Transparent Proxy 指向它
        // 并通过 initData 在代理构造时执行初始化。
        // 注意：合约里函数名是历史拼写 initalize（不是 initialize）。
        const impl = await viem.deployContract("MetaNFTAuction");

        const initData = encodeFunctionData({
            abi: (await hre.artifacts.readArtifact("MetaNFTAuction")).abi,
            functionName: "initalize",
            args: [admin.account.address]
        });

        proxyAdmin = await viem.deployContract(
            "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol:ProxyAdmin",
            [],
            { account: proxyAdminSigner.account }
        );

        // 部署透明代理：
        // - admin 为 ProxyAdmin 合约地址（不是 EOA）
        // - data 为初始化 calldata，会在部署时 delegatecall 到 impl 执行
        proxy = await viem.deployContract(
            "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:TransparentUpgradeableProxy",
            [
            impl.address,
            proxyAdmin.address,
            initData
            ]
        );

        auction = await viem.getContractAt("MetaNFTAuction", proxy.address);

        // 部署测试依赖：NFT、支付代币、喂价
        nft = await viem.deployContract("MetaNFT");
        usdc = await viem.deployContract("MockERC20", ["USDC", "USDC", 6n, parseUnits("1000000", 6)]);
        ethOracle = await viem.deployContract("MockOracle", [parseUnits("3000", 8)]);
        usdcOracle = await viem.deployContract("MockOracle", [parseUnits("1", 8)]);

        await auction.write.setTokenOracle([zeroAddress, ethOracle.address], { account: admin.account });
        await auction.write.setTokenOracle([usdc.address, usdcOracle.address], { account: admin.account });

        // 预置 NFT 给 seller，并提前授权拍卖合约托管 NFT
        // 否则 start() 内部 transferFrom 会失败。
        await nft.write.mint([seller.account.address, 1n], { account: admin.account });
        await nft.write.mint([seller.account.address, 2n], { account: admin.account });
        await nft.write.mint([seller.account.address, 10n], { account: admin.account });
        await nft.write.setApprovalForAll([auction.address, true], { account: seller.account });
    });

    describe("getVersion", function () {
        it("should return MetaNFTAuctionV1", async function () {
            // V1 版本标识用于验证当前代理实现确实是旧逻辑
            const version = await auction.read.getVersion();
            expect(version).to.equal("MetaNFTAuctionV1");
        });
    });

    describe("getPriceInDollar", function () {
        it("should return correct prices", async function () {
            // zeroAddress 在本项目中约定表示 ETH 喂价 key
            const ethPrice = await auction.read.getPriceInDollar([zeroAddress]);
            const usdcPrice = await auction.read.getPriceInDollar([usdc.address]);

            // viem 返回 bigint，这里直接做 bigint 比较
            expect(ethPrice > 0n).to.equal(true);
            expect(usdcPrice > 0n).to.equal(true);
        });
    });

    describe("initalize", function () {
        it("should fail when initialized twice", async function () {
            // 代理部署时已经执行过一次 initalize，再调应触发 Initializable 防重入
            await expectRevert(
                auction.write.initalize([admin.account.address], { account: admin.account }),
                "already initialized"
            );
        });
    });

    describe("start", function () {
        it("should fail when not called by admin", async function () {
            // start 受 onlyAdmin 保护，seller 直接调用必须失败
            await expectRevert(
                auction.write.start([seller.account.address, 1n, nft.address, 1000n, 3600n, usdc.address], { account: seller.account }),
                "not admin"
            );
        });

        it("should increment auctionId", async function () {
            // 每创建一场拍卖，auctionId 自增 1
            await auction.write.start([seller.account.address, 1n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            let auctionId = await auction.read.auctionId();
            expect(auctionId).to.equal(1n);

            await auction.write.start([seller.account.address, 2n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            auctionId = await auction.read.auctionId();
            expect(auctionId).to.equal(2n);
        });
    });

    describe("bid", function () {
        it("should fail when auction has ended", async function () {
            await auction.write.start([seller.account.address, 1n, nft.address, 1000n, 30n, usdc.address], { account: admin.account });
            const currentAuctionId = (await auction.read.auctionId()) - 1n;
            const auctionData = await auction.read.auctions([currentAuctionId]);
            // auctions 结构中 [3]=startingTime，[6]=duration
            const endTime = auctionData[3] + auctionData[6];

            // 把链上时间推进到结束点，验证 bid 被拒绝
            await setTimestamp(endTime);

            await expectRevert(
                auction.write.bid([currentAuctionId, parseUnits("1", 18)], {
                    account: seller.account,
                    value: parseUnits("1", 18)
                }),
                "ended"
            );
        });

        it("should fail when bid is lower than highest bid", async function () {
            await auction.write.start([seller.account.address, 1n, nft.address, 1000n, 30n, usdc.address], { account: admin.account });
            const currentAuctionId = (await auction.read.auctionId()) - 1n;

            // 首次有效出价建立最高价
            await auction.write.bid([currentAuctionId, parseUnits("2", 18)], {
                account: seller.account,
                value: parseUnits("2", 18)
            });

            // 第二次出价低于当前最高价，应被拒绝
            await expectRevert(
                auction.write.bid([currentAuctionId, parseUnits("1.2", 18)], {
                    account: bidder1.account,
                    value: parseUnits("1.2", 18)
                }),
                "invalid highestBid"
            );
        });

        it("should correctly track bidding result", async function () {
            await auction.write.start([seller.account.address, 1n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            const currentAuctionId = (await auction.read.auctionId()) - 1n;

            await auction.write.bid([currentAuctionId, parseUnits("2", 18)], {
                account: bidder1.account,
                value: parseUnits("2", 18)
            });
            await auction.write.bid([currentAuctionId, parseUnits("3", 18)], {
                account: bidder2.account,
                value: parseUnits("3", 18)
            });
            await auction.write.bid([currentAuctionId, parseUnits("4", 18)], {
                account: bidder1.account,
                value: parseUnits("4", 18)
            });

            // 多轮竞价后，校验链上记录的最高出价人与最高价金额
            const auctionData = await auction.read.auctions([currentAuctionId]);
            expect(auctionData[4]).to.equal(getAddress(bidder1.account.address));
            expect(auctionData[8]).to.equal(parseUnits("4", 18));
        });
    });

    describe("upgrade", function () {
        it("should upgrade contract successfully", async function () {
            // 升级前先做一次业务写入，验证升级后存储不丢失
            await auction.write.start([seller.account.address, 10n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            const oldAuctionId = await auction.read.auctionId();
            const newImpl = await viem.deployContract("MetaNFTAuctionV2");

            const proxyAsV2 = await viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                proxy.address
            );

            // OZ v4 ProxyAdmin: 无需附带 call data 时使用 upgrade
            await proxyAdmin.write.upgrade(
                [proxyAsV2.address, newImpl.address],
                { account: proxyAdminSigner.account }
            );

            // 升级后同地址以 V2 ABI 读取：应保留旧状态 + 返回新逻辑
            const upgradedAuction = await viem.getContractAt("MetaNFTAuctionV2", proxy.address);

            const newAuctionId = await upgradedAuction.read.auctionId();
            expect(newAuctionId).to.equal(oldAuctionId);

            const version = await upgradedAuction.read.getVersion();
            expect(version).to.equal("MetaNFTAuctionV2");

            const newFeature = await upgradedAuction.read.newFeature();
            expect(newFeature).to.equal("This is a new feature in V2");
        });

        it("should fail when non-admin tries to upgrade", async function () {
            await auction.write.start([seller.account.address, 10n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            const newImpl = await viem.deployContract("MetaNFTAuctionV2");

            const proxyAsV2 = await viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                proxy.address
            );

            // 升级权限由 ProxyAdmin owner 控制，seller 调用应失败
            await expectRevert(
                proxyAdmin.write.upgrade(
                    [proxyAsV2.address, newImpl.address],
                    { account: seller.account }
                )
            );
        });

        it("should change oracle after upgrade", async function () {
            await auction.write.start([seller.account.address, 10n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            const newEthOracle = await viem.deployContract("MockOracle", [parseUnits("3000", 8)]);
            const newImpl = await viem.deployContract("MetaNFTAuctionV2");

            const proxyAsV2 = await viem.getContractAt(
                "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:ITransparentUpgradeableProxy",
                proxy.address
            );

            await proxyAdmin.write.upgrade(
                [proxyAsV2.address, newImpl.address],
                { account: proxyAdminSigner.account }
            );

            const upgradedAuction = await viem.getContractAt("MetaNFTAuctionV2", proxy.address);

            // 验证升级后仍可调用管理员函数并更新状态（证明新实现可用）
            await upgradedAuction.write.setTokenOracle([zeroAddress, newEthOracle.address], { account: admin.account });

            const newPrice = await upgradedAuction.read.getPriceInDollar([zeroAddress]);
            expect(newPrice).to.equal(parseUnits("3000", 8));
        });
    });
});