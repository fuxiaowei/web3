import { expect } from "chai";
import "@nomicfoundation/hardhat-viem";
import hre from "hardhat";
import { parseUnits, zeroAddress, encodeFunctionData, getAddress } from "viem";

describe("MetaNFTAuction", function () {
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

    // 封装工具函数：增加时间
    const increaseTime = async (seconds: number) => {
        await networkConnection.ethers.provider.send("evm_increaseTime", [seconds]);
        await networkConnection.ethers.provider.send("evm_mine");
    };

    // 封装工具函数：设置指定时间戳
    const setTimestamp = async (target: bigint) => {
        await networkConnection.ethers.provider.send("evm_setNextBlockTimestamp", [Number(target)]);
        await networkConnection.ethers.provider.send("evm_mine");
    };

    beforeEach(async function () {
        networkConnection = await hre.network.connect();
        viem = networkConnection.viem;
        const clients = await viem.getWalletClients();
        [admin, proxyAdminSigner, seller, bidder1, bidder2] = clients;

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

        proxy = await viem.deployContract(
            "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:TransparentUpgradeableProxy",
            [
            impl.address,
            proxyAdmin.address,
            initData
            ]
        );

        auction = await viem.getContractAt("MetaNFTAuction", proxy.address);

        nft = await viem.deployContract("MetaNFT");
        usdc = await viem.deployContract("MockERC20", ["USDC", "USDC", 6n, parseUnits("1000000", 6)]);
        ethOracle = await viem.deployContract("MockOracle", [parseUnits("3000", 8)]);
        usdcOracle = await viem.deployContract("MockOracle", [parseUnits("1", 8)]);

        await auction.write.setTokenOracle([zeroAddress, ethOracle.address], { account: admin.account });
        await auction.write.setTokenOracle([usdc.address, usdcOracle.address], { account: admin.account });

        await nft.write.mint([seller.account.address, 1n], { account: admin.account });
        await nft.write.mint([seller.account.address, 2n], { account: admin.account });
        await nft.write.mint([seller.account.address, 10n], { account: admin.account });
        await nft.write.setApprovalForAll([auction.address, true], { account: seller.account });
    });

    describe("getVersion", function () {
        it("should return MetaNFTAuctionV1", async function () {
            const version = await auction.read.getVersion();
            expect(version).to.equal("MetaNFTAuctionV1");
        });
    });

    describe("getPriceInDollar", function () {
        it("should return correct prices", async function () {
            const ethPrice = await auction.read.getPriceInDollar([zeroAddress]);
            const usdcPrice = await auction.read.getPriceInDollar([usdc.address]);

            expect(ethPrice > 0n).to.equal(true);
            expect(usdcPrice > 0n).to.equal(true);
        });
    });

    describe("initalize", function () {
        it("should fail when initialized twice", async function () {
            await expectRevert(
                auction.write.initalize([admin.account.address], { account: admin.account }),
                "already initialized"
            );
        });
    });

    describe("start", function () {
        it("should fail when not called by admin", async function () {
            await expectRevert(
                auction.write.start([seller.account.address, 1n, nft.address, 1000n, 3600n, usdc.address], { account: seller.account }),
                "not admin"
            );
        });

        it("should increment auctionId", async function () {
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
            const endTime = auctionData[3] + auctionData[6];

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

            await auction.write.bid([currentAuctionId, parseUnits("2", 18)], {
                account: seller.account,
                value: parseUnits("2", 18)
            });

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

            const auctionData = await auction.read.auctions([currentAuctionId]);
            expect(auctionData[4]).to.equal(getAddress(bidder1.account.address));
            expect(auctionData[8]).to.equal(parseUnits("4", 18));
        });
    });

    describe("upgrade", function () {
        it("should upgrade contract successfully", async function () {
            await auction.write.start([seller.account.address, 10n, nft.address, 1000n, 3600n, usdc.address], { account: admin.account });
            const oldAuctionId = await auction.read.auctionId();
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

            await upgradedAuction.write.setTokenOracle([zeroAddress, newEthOracle.address], { account: admin.account });

            const newPrice = await upgradedAuction.read.getPriceInDollar([zeroAddress]);
            expect(newPrice).to.equal(parseUnits("3000", 8));
        });
    });
});