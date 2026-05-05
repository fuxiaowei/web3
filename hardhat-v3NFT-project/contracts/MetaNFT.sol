// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

/// @title MetaNFTAuction
/// @notice 可升级代理场景下的 NFT 拍卖合约骨架（`Initializable`）；当前仅定义状态与事件，业务逻辑待实现。
/// @dev 配合 OpenZeppelin 代理使用时，构造函数不可用，需用 `initialize` 等方式完成一次性初始化（本文件尚未写出）。
contract MetaNFTAuction is Initializable {
    /// @notice 管理员地址，通常用于暂停、改参或紧急提款等（具体权限在实现函数中定义）。
    address admin;

    /// @notice ERC20 支付代币地址 → Chainlink 价格预言机（Aggregator）合约地址，用于将出价换算为美元口径。
    /// @dev 字段名 `tokenToOrace` 为历史命名，语义上对应 oracle。
    mapping(address => address) public tokenToOrace;

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

    /// @notice auctionId → 拍卖详情；`auctionId` 为合约内递增的拍卖编号。
    mapping(uint256 => Auction) public auctions;

    /// @notice 拍卖开始或起拍价公布时触发（参数含义以实现为准）。
    event StartBid(uint256 startingBid);

    /// @notice 有人出价时触发。
    event Bid(address indexed sender, uint256 amount);

    /// @notice 拍卖结束、结算或关闭时触发。
    event EndBid(uint256 indexed auctionId);

    /// @notice 下一个将分配的拍卖编号；创建新拍卖后应递增。
    uint256 public auctionId;
}


