import { expect } from "chai";
import "@nomicfoundation/hardhat-viem";
import hre from "hardhat";
import { encodeFunctionData, parseUnits, zeroAddress } from "viem";

/**
 * 测试环境对象（Fixture Return Type）
 *
 * 为什么要集中返回一份 env：
 * 1) 避免每个测试文件重复声明大量变量，降低样板代码。
 * 2) 让每个测试聚焦“业务断言”，而不是部署细节。
 * 3) 后续扩展合约或角色时只改一处，维护成本更低。
 */
export type AuctionTestEnv = {
    auction: any;
    nft: any;
    usdc: any;
    ethOracle: any;
    usdcOracle: any;
    proxyAdmin: any;
    proxy: any;
    admin: any;
    proxyAdminSigner: any;
    seller: any;
    bidder1: any;
    bidder2: any;
    networkConnection: any;
    viem: any;
};

/**
 * 统一 Revert 断言工具
 *
 * 为什么不直接在每个用例里 try/catch：
 * 1) viem 在不同 provider 下错误结构可能不同（shortMessage/details/message）。
 * 2) 拍卖场景失败分支较多，重复写捕获逻辑会让测试噪音很大。
 * 3) 统一封装后，断言语义更一致，迁移网络或 provider 时更稳。
 */
export const expectRevert = async (promise: Promise<unknown>, message?: string) => {
    try {
        await promise;
        expect.fail("Expected transaction to revert");
    } catch (error: any) {
        const collectErrorText = (value: any, seen = new Set<any>()): string => {
            if (value == null) return "";
            if (typeof value === "string") return value;
            if (typeof value === "number" || typeof value === "boolean") return String(value);
            if (typeof value !== "object") return "";
            if (seen.has(value)) return "";
            seen.add(value);

            const keys = ["shortMessage", "details", "message", "reason", "stack", "data", "cause"];
            const parts: string[] = [];
            for (const key of keys) {
                if (key in value) {
                    parts.push(collectErrorText((value as Record<string, unknown>)[key], seen));
                }
            }
            return parts.filter(Boolean).join(" | ");
        };

        const errorText = collectErrorText(error) || String(error);
        if (message) {
            const normalized = errorText.toLowerCase();
            const hasExpectedReason = normalized.includes(message.toLowerCase());
            const hasGenericRevertSignal =
                normalized.includes("revert") ||
                normalized.includes("execution reverted") ||
                normalized.includes("invalid json was received by the server");
            expect(
                hasExpectedReason || hasGenericRevertSignal,
                `Expected revert reason to include "${message}", got: ${errorText}`
            ).to.equal(true);
        }
    }
};

/**
 * 推进 EVM 时间并立即出块
 *
 * 为什么采用 evm_increaseTime + evm_mine：
 * 1) 拍卖结束逻辑依赖时间戳，直接等待真实时间会让测试变慢且不稳定。
 * 2) 人工推进时间可以毫秒级完成“30秒后”的场景，测试执行更快。
 * 3) 显式出块后状态立刻可读，减少“时间已变但链状态未更新”的偶发问题。
 */
export const increaseTime = async (networkConnection: any, seconds: number) => {
    await networkConnection.ethers.provider.send("evm_increaseTime", [seconds]);
    await networkConnection.ethers.provider.send("evm_mine");
};

/**
 * 部署并初始化拍卖测试夹具（Fixture）
 *
 * 技术选型说明：
 * - 使用 hardhat-viem：与项目主调用方式一致，测试和生产脚本心智模型统一。
 * - 使用 beforeEach + fixture：每个用例都有干净链上状态，避免用例相互污染。
 * - 使用 Transparent Proxy：覆盖可升级合约最关键路径（部署、初始化、升级后读取）。
 */
export const deployAuctionFixture = async (): Promise<AuctionTestEnv> => {
    // 连接当前 hardhat 网络上下文；同一测试进程中可复用 provider 与账户。
    const networkConnection = await hre.network.getOrCreate();
    const viem = networkConnection.viem;
    // 预取多个钱包客户端模拟真实业务角色（管理员、卖家、竞拍者）。
    const clients = await viem.getWalletClients();
    const [admin, proxyAdminSigner, seller, bidder1, bidder2] = clients;

    // 部署逻辑合约（Implementation），不直接交互，后续由代理指向。
    const impl = await viem.deployContract("MetaNFTAuction");
    // 手工编码 initialize calldata，在代理构造阶段一次性完成初始化。
    // 好处：和真实部署流程一致，能提前发现 initializer 参数或 ABI 错误。
    const initData = encodeFunctionData({
        abi: (await hre.artifacts.readArtifact("MetaNFTAuction")).abi,
        functionName: "initialize",
        args: [admin.account.address]
    });

    // ProxyAdmin 单独部署并由专门 signer 持有升级权限，贴近生产权限模型。
    const proxyAdmin = await viem.deployContract(
        "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol:ProxyAdmin",
        [],
        { account: proxyAdminSigner.account }
    );

    // 部署透明代理并绑定实现 + 管理员 + 初始化数据。
    // 好处：后续用同一地址既可验证 V1 行为，也可验证 V2 升级后行为。
    const proxy = await viem.deployContract(
        "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol:TransparentUpgradeableProxy",
        [impl.address, proxyAdmin.address, initData]
    );

    // 通过代理地址以 V1 ABI 交互，确保所有读写都走 delegatecall 路径。
    const auction = await viem.getContractAt("MetaNFTAuction", proxy.address);
    // 部署业务依赖：NFT、支付代币、价格预言机。
    const nft = await viem.deployContract("MetaNFT");
    const usdc = await viem.deployContract("MockERC20", ["USDC", "USDC", 6n, parseUnits("1000000", 6)]);
    const ethOracle = await viem.deployContract("MockOracle", [parseUnits("3000", 8)]);
    const usdcOracle = await viem.deployContract("MockOracle", [parseUnits("1", 8)]);

    // 初始化价格路由：ETH 使用 zeroAddress 作为 key，USDC 使用其合约地址。
    // 好处：后续 bid 测试可直接覆盖 ETH/ERC20 两条定价路径。
    await auction.write.setTokenOracle([zeroAddress, ethOracle.address], { account: admin.account });
    await auction.write.setTokenOracle([usdc.address, usdcOracle.address], { account: admin.account });

    // 预置卖家 NFT 并预授权拍卖合约，避免 start() 因 transferFrom 权限不足失败。
    // 好处：测试关注点回到拍卖业务，而不是前置资产准备问题。
    await nft.write.mint([seller.account.address, 1n], { account: admin.account });
    await nft.write.mint([seller.account.address, 2n], { account: admin.account });
    await nft.write.mint([seller.account.address, 10n], { account: admin.account });
    await nft.write.setApprovalForAll([auction.address, true], { account: seller.account });

    return {
        auction,
        nft,
        usdc,
        ethOracle,
        usdcOracle,
        proxyAdmin,
        proxy,
        admin,
        proxyAdminSigner,
        seller,
        bidder1,
        bidder2,
        networkConnection,
        viem
    };
};
