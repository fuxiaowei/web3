// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

// ERC721：拍卖标的 NFT；ERC20：支付代币；IERC20Metadata：代币小数位等，供价格换算使用
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
// Initializable：与可升级代理配合，用 initializer 做部署后一次性配置
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
// Chainlink Aggregator：按喂价将 ERC20 出价换算为美元等口径
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

/// @title MetaNFTAuction  主拍卖合约，实现拍卖逻辑
/// @notice 可升级代理场景下的 NFT 拍卖合约骨架（`Initializable`）；当前仅定义状态与事件，业务逻辑待实现。
/// @dev 逻辑合约的 `constructor` 中调用 `_disableInitializers()`；代理部署后再通过 `initalize` 设置 `admin`。
///      `tokenToOracle` 供后续按支付代币解析预言机；具体创建拍卖、出价、结算函数待实现。
contract MetaNFTAuction is Initializable {
    /// @notice 管理员地址，通常用于暂停、改参或紧急提款等（具体权限在实现函数中定义）。
    address admin;

    /// @notice ERC20 支付代币地址 → Chainlink 价格预言机（Aggregator）合约地址，用于将出价换算为美元口径。
    /// @dev 字段名 `tokenToOracle` 为历史命名，语义上对应 oracle（价格预言机）。
    mapping(address => address) public tokenToOracle;

    /// @notice 单笔拍卖的链上快照：NFT、卖家、时间与出价信息。
    struct Auction {
        IERC721 nft; // 被拍卖的 NFT 集合合约
        uint256 nftId; // 被拍卖的 tokenId
        address payable seller; // 卖家，落槌后应收款地址
        uint256 startingTime; // 拍卖开始时间戳（秒）
        address highestBidder; // 当前最高出价人；结束时以此地址为买受人
        uint256 startingPriceInDollar; // 起拍价（美元口径，具体精度依赖预言机/换算逻辑）
        uint256 duration; // 拍卖持续时长（秒），通常与 startingTime 一起决定结束时刻
        IERC20 paymentToken; // 本次拍卖接受的 ERC20 支付代币
        uint256 highestBid; // 当前最高出价的代币数量（paymentToken 最小单位）
        uint256 highestBidInDollar; // 当前最高出价折算后的美元口径数值（便于展示与比较）
        address highestBidToken; // 若支持多代币出价，可记录最高出价所用代币；单币场景可与 paymentToken 一致或按实现约定
    }

    /// @notice `auctionId` → 拍卖详情；`auctionId` 为合约内递增的拍卖编号。
    /// @dev `start()` 会在当前 `auctionId` 处写入一条新的 `Auction`，然后递增 `auctionId`。
    mapping(uint256 => Auction) public auctions;

    /// @notice 拍卖开始或起拍价公布时触发。
    /// @dev 事件参数 `startingBid` 在当前实现里实际上会传入“递增后的拍卖编号”（见 `start()` 内的 `emit StartBid(auctionId)`）。
    ///      如果你后续要按“起拍价”触发，建议统一事件语义与参数含义。
    /// @param startingBid 当前合约实现所传入的 starting 信息（此处为拍卖编号）
    event StartBid(uint256 startingBid);

    /// @notice 有人提交有效出价时触发，前端可据此刷新当前最高价与领先地址。
    /// @param sender 出价人地址
    /// @param amount 本次出价的支付代币数量（一般为 `paymentToken` 的最小单位）
    event Bid(address indexed sender, uint256 amount);

    /// @notice 拍卖结束、结算或关闭时触发；监听方可据此触发 UI 更新或链下索引。
    /// @param auctionId 结束的拍卖编号
    event EndBid(uint256 indexed auctionId);

    /// @notice 下一个将分配的拍卖编号。
    /// @dev 从 0 开始（默认值），每次 `start()` 创建拍卖后递增。
    uint256 public auctionId;

    /// @notice 仅允许当前 `admin` 调用；用于治理类操作（如配置预言机）。
    modifier onlyAdmin() {
        require(msg.sender == admin, "not admin");
        _;
    }

    /// @notice 实现合约构造函数：禁用在本合约地址上的初始化器，防止逻辑实现被单独部署时重复初始化；
    ///         真实初始化在代理指向本实现后，通过 `initalize` 完成。
    constructor() {
        _disableInitializers();
    }

    /// @notice 代理部署后的一次性初始化：设置管理员。仅可成功调用一次（`initializer` 修饰符）。
    /// @param _admin 管理员地址，不可为零地址
    /// @dev 函数名为历史拼写 `initalize`；对外 ABI 需与此保持一致。
    function initalize(address _admin) external initializer {
        require(_admin != address(0), "invalid admin");
        admin = _admin;
    }

    /// @notice 为指定 ERC20 支付代币绑定 Chainlink 预言机，用于将代币金额换算为美元（或喂价对约定的报价资产）。
    /// @param token 作为拍卖支付手段的 ERC20 合约地址
    /// @param oracle 该代币对应的 `AggregatorV3Interface` 合约地址
    /// @dev 写入 `tokenToOracle`；具体取价、小数位处理在后续出价/结算逻辑中实现。
    function setTokenOracle(address token, address oracle) external onlyAdmin {
        require(oracle != address(0), "invalid oracle");
        tokenToOracle[token] = oracle;
    }

    /// @notice 在合约中创建一笔新的拍卖，并转移 NFT 至拍卖合约托管。
    /// @dev 这是“拍卖创建”函数：写入拍卖快照、设置初始最高价与出价人为空，
    ///      然后把 `seller` 持有的 `nftId` 转移到本合约。
    ///      该函数被 `onlyAdmin` 修饰，意味着“发起创建拍卖”动作由管理员执行。
    ///
    ///      关于价格：`startingPriceInDollar` 会被乘以 `10**8`，这通常意味着你将“美元价格”按 Chainlink feed 的 8 位小数进行整数化。
    ///      具体还需要你后续出价/结算时也使用同一小数口径进行换算与比较。
    ///
    /// @param seller 拍卖发起方（卖家），落槌后应收款的地址
    /// @param nftId 被拍卖的 tokenId
    /// @param nft 被拍卖的 ERC721 合约地址
    /// @param startingPriceInDollar 起拍价（美元口径，未乘 1e8 前的“人类可读值”）
    /// @param duration 拍卖持续时长（秒）
    /// @param paymentToken 拍卖接受的支付代币（ERC20）
    function start(
        address seller,
        uint256 nftId,
        address nft,
        uint256 startingPriceInDollar,
        uint256 duration,
        address paymentToken
    ) external onlyAdmin {
        require(nft != address(0), "invalid nft");
        require(duration >= 30, "invalid duration");
        require(paymentToken != address(0), "invalid payment token");

        // 在当前 auctionId 处写入拍卖快照，然后创建完成后递增 auctionId。
        Auction storage auction = auctions[auctionId];
        auction.nft = IERC721(nft);
        auction.nftId = nftId;
        auction.seller = payable(seller);
        // 记录拍卖开始时间；结束时间可用 `startingTime + duration` 推导（取决于后续实现）。
        auction.startingTime = block.timestamp;
        // 将“美元起拍价”转为 1e8 精度的整数，便于链上比较计算（通常与 Chainlink 8 decimals 对齐）。
        auction.startingPriceInDollar = startingPriceInDollar * 10 ** 8;
        auction.duration = duration;
        // 设置本笔拍卖的支付代币，并以 ERC20 最小单位计价出价（后续逻辑需保持一致）。
        auction.paymentToken = IERC20(paymentToken);
        // 初始化当前最高出价为空。
        auction.highestBid = 0;
        auction.highestBidder = address(0);
        auction.highestBidInDollar = 0;
        auction.highestBidToken = address(0);
        // 将卖家资产转移到合约托管：要求 seller 已对本合约执行 `approve`。
        IERC721(nft).transferFrom(seller, address(this), nftId);
        auctionId ++;
        // 注意：这里 emit 的是“递增后的 auctionId”，因此监听方拿到的值需要与 `auctions` 的写入编号对应关系对齐。
        emit StartBid(auctionId);
    }

    /// @notice 参与竞价：支持两种出价方式（ETH 或预设 ERC20）。
    /// @dev
    /// - 当 `msg.value > 0` 时按 ETH 出价，`amount` 必须与 `msg.value` 一致；
    /// - 当 `msg.value == 0` 时按 ERC20 出价，`amount` 为 ERC20 最小单位数量；
    /// - 两种方式都会先换算为美元口径（通常按 1e8 精度）再比较高低；
    /// - 若出现新最高价，会退还上一名最高出价人的资金（按其出价资产类型退还）。
    ///
    /// 注意：当前实现中的事件 `emit Bid(msg.sender, msg.value)` 在 ERC20 出价分支会发出 0，
    /// 若前端依赖事件金额，请结合链上状态 `auctions[_auctionId].highestBid` 一并读取。
    ///
    /// @param _auctionId 参与竞价的拍卖编号
    /// @param amount 出价数量；ETH 出价时应等于 `msg.value`，ERC20 出价时为代币最小单位
    function bid(uint256 _auctionId, uint256 amount) external payable {
        Auction storage auction = auctions[_auctionId];
        require(auction.startingTime > 0, "not started");
        require(!isEnded(_auctionId), "ended");

        uint256 bidPrice;
        bool isEthBid = msg.value > 0;
        if (isEthBid) {
            // ETH 出价：使用 token = address(0) 对应的预言机价格进行换算。
            require(amount == msg.value, "amount mismatch");
            uint256 price = getPriceInDollar(address(0));
            bidPrice = _toUsd(msg.value, 18, price);
        } else {
            // ERC20 出价：读取支付代币 decimals，按同一美元口径进行换算。
            require(amount > 0, "invalid amount");
            uint256 price = getPriceInDollar(address(auction.paymentToken));
            uint8 tokenDecimals = IERC20Metadata(address(auction.paymentToken)).decimals();
            bidPrice = _toUsd(amount, tokenDecimals, price);
            // 先把本次代币出价转入合约托管。
            IERC20(address(auction.paymentToken)).transferFrom(msg.sender, address(this), amount);
        }
        // 出价必须同时高于起拍价与当前最高价（均为美元口径）。
        require(auction.startingPriceInDollar < bidPrice, "invalid startingPrice");
        require(auction.highestBidInDollar < bidPrice, "invalid highestBid");
        if (auction.highestBidder != address(0) && auction.highestBidder != msg.sender) {
            // 出现新最高价时，退还上一位领先者的资金。
            uint256 refundAmount = auction.highestBid;
            if (refundAmount > 0) {
                if (auction.highestBidToken == address(0)) {
                    payable(auction.highestBidder).transfer(refundAmount);
                } else {
                    IERC20(address(auction.paymentToken)).transfer(auction.highestBidder, refundAmount);
                }
            }
        }
        if (isEthBid) {
            auction.highestBid = msg.value;
            auction.highestBidToken = address(0);
        } else {
            auction.highestBid = amount;
            auction.highestBidToken = address(auction.paymentToken);
        }
        auction.highestBidder = msg.sender;
        auction.highestBidInDollar = bidPrice;
        // 当前实现记录的是 msg.value（ERC20 分支会是 0）；前端可结合 highestBid 读取真实代币金额。
        emit Bid(msg.sender, msg.value);
    }

    /// @notice 判断拍卖是否已结束。
    /// @param _auctionId 拍卖编号
    /// @return 若当前时间 >= `startingTime + duration` 且拍卖已创建，则返回 true
    function isEnded(uint256 _auctionId) public view returns (bool)  {
        Auction storage auction = auctions[_auctionId];
        return auction.startingTime > 0 &&
        block.timestamp >= auction.startingTime + auction.duration;
    }

    /// @notice 结束拍卖并完成结算：转移 NFT 给最高出价人，并将最高出价转给卖家。
    /// @dev 该函数未限制调用者，任何人都可在拍卖到期后触发结算（只要满足条件）。
    /// @param _auctionId 待结束的拍卖编号
    function end(uint256 _auctionId) external {
        Auction storage auction = auctions[_auctionId];
        require(isEnded(_auctionId), "not ended");
        require(auction.highestBidder != address(0), "no bids");

        // 先把托管 NFT 转给最高出价者。
        auction.nft.transferFrom(address(this), auction.highestBidder, auction.nftId);

        if (auction.highestBid > 0) {
            if (auction.highestBidToken == address(0)) {
                payable(auction.seller).transfer( auction.highestBid);
            } else {
                IERC20(auction.highestBidToken).transfer(auction.seller, auction.highestBid);
            }
        }
        emit EndBid(_auctionId);
    }

    /// @notice 获取某资产对应的链上喂价（通常为美元口径）。
    /// @param token 资产地址；ETH 常用 `address(0)` 作为约定 key
    /// @return 价格数值（精度由预言机决定，常见为 8 位小数）
    function getPriceInDollar(address token) public view returns (uint256) {
        AggregatorV3Interface dataFeed;
        address oracle = tokenToOracle[token];
        require(oracle != address(0), "oracle not set");
        dataFeed = AggregatorV3Interface(oracle);
        (
        /* uint80 roundId */
        ,
        int256 answer,
        /*uint256 startedAt*/
        ,
        /*uint256 updatedAt*/
        ,
        /*uint80 answeredInRound*/
        ) = dataFeed.latestRoundData();
        return uint256(answer);
    }

    /// @notice 按给定资产数量、资产精度和价格，换算成美元口径整数。
    /// @dev 假设 `price` 已包含其自身小数位（例如 Chainlink 常见 8 位）；本函数不再额外处理 priceDecimals。
    /// @param amount 资产数量（最小单位）
    /// @param amountDecimals 资产数量的精度（例如 ETH 为 18）
    /// @param price 预言机价格
    /// @return usd 换算后的美元口径整数（与 `price` 的小数位保持同一口径）
    function _toUsd(uint256 amount, uint256 amountDecimals, uint256 price)
    internal
    pure
    returns (uint256)
    {
        // amount 以最小单位输入，先除以 10^amountDecimals 还原数量，再乘以价格得美元口径值。
        uint256 scale = 10 ** amountDecimals;
        uint256 usd = (amount * price) / scale;
        return usd;
    }

    /// @notice 返回当前实现版本标识，便于升级后做链上区分。
    /// @return 版本字符串（V1）
    function getVersion() external pure virtual returns (string memory) {
        return "MetaNFTAuctionV1";
    }
}
