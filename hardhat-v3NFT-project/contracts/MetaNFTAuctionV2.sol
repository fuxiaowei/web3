// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./MetaNFTAuction.sol";

/// @title MetaNFTAuctionV2
/// @notice 在 V1 基础上新增“按美元分段动态手续费”能力。
/// @dev
/// - 手续费档位依据 `highestBidInDollar`（1e8 精度）判定，而不是代币数量；
/// - 手续费比例使用 bps（万分比），避免浮点计算；
/// - 所有新增状态变量都追加在 V1 之后，确保透明代理升级的存储兼容性。

contract MetaNFTAuctionV2 is MetaNFTAuction {
    /// @notice bps 分母（10000 = 100%）
    uint16 public constant MAX_BPS = 10_000;

    // 美元分段阈值（1e8 精度，和 highestBidInDollar 保持一致）。
    // 规则：
    // - bidUsd <= tier1UpperUsd        使用 tier1Bps
    // - tier1UpperUsd < bidUsd <= tier2UpperUsd 使用 tier2Bps
    // - bidUsd > tier2UpperUsd         使用 tier3Bps
    uint256 public tier1UpperUsd;
    uint256 public tier2UpperUsd;
    // 各分段手续费（bps，万分比），例如 100 = 1%。
    uint16 public tier1Bps;
    uint16 public tier2Bps;
    uint16 public tier3Bps;
    // 平台手续费接收地址；如果未配置（零地址）则视为不收取手续费。
    address public feeRecipient;

    /// @notice 手续费接收地址更新事件
    event FeeRecipientUpdated(address indexed recipient);
    /// @notice 动态费率配置更新事件
    event FeeConfigUpdated(
        uint256 tier1UpperUsd,
        uint256 tier2UpperUsd,
        uint16 tier1Bps,
        uint16 tier2Bps,
        uint16 tier3Bps
    );
    /// @notice 拍卖结算时手续费明细事件，便于前端展示与链下对账
    event FeeCharged(
        uint256 indexed auctionId,
        address indexed token,
        uint256 grossAmount,
        uint256 feeAmount,
        uint256 sellerAmount,
        address seller,
        address recipient
    );

    function getVersion() external pure override returns (string memory){
        return "MetaNFTAuctionV2";
    }

    function newFeature()external pure returns (string memory){
        return "This is a new feature in V2";
    }

    /// @notice V2 一次性初始化入口（reinitializer(2)）
    /// @dev
    /// - 仅 admin 可调用；
    /// - 只允许执行一次，避免升级后被重复覆盖配置；
    /// - 推荐升级后立即调用，确保手续费逻辑有明确配置。
    function initializeV2(
        address _feeRecipient,
        uint256 _tier1UpperUsd,
        uint256 _tier2UpperUsd,
        uint16 _tier1Bps,
        uint16 _tier2Bps,
        uint16 _tier3Bps
    ) external reinitializer(2) onlyAdmin {
        _setFeeRecipient(_feeRecipient);
        _setDynamicFeeConfig(_tier1UpperUsd, _tier2UpperUsd, _tier1Bps, _tier2Bps, _tier3Bps);
    }

    /// @notice 更新手续费接收地址（管理员）
    function setFeeRecipient(address _feeRecipient) external onlyAdmin {
        _setFeeRecipient(_feeRecipient);
    }

    /// @notice 更新按美元分段的手续费规则（管理员）
    function setDynamicFeeConfig(
        uint256 _tier1UpperUsd,
        uint256 _tier2UpperUsd,
        uint16 _tier1Bps,
        uint16 _tier2Bps,
        uint16 _tier3Bps
    ) external onlyAdmin {
        _setDynamicFeeConfig(_tier1UpperUsd, _tier2UpperUsd, _tier1Bps, _tier2Bps, _tier3Bps);
    }

    /// @notice 根据美元出价档位返回 bps 手续费
    /// @dev 若 feeRecipient 未设置，则返回 0，表示暂不收取手续费
    function getFeeBps(uint256 bidUsd) public view returns (uint16) {
        if (feeRecipient == address(0)) {
            return 0;
        }
        if (bidUsd <= tier1UpperUsd) {
            return tier1Bps;
        }
        if (bidUsd <= tier2UpperUsd) {
            return tier2Bps;
        }
        return tier3Bps;
    }

    /// @notice 预览某拍卖当前应收手续费（不改状态）
    /// @return bps 当前匹配费率（万分比）
    /// @return feeAmount 平台手续费金额（按成交资产最小单位）
    /// @return sellerAmount 卖家净收入（按成交资产最小单位）
    function previewFee(uint256 _auctionId) external view returns (uint16 bps, uint256 feeAmount, uint256 sellerAmount) {
        Auction storage auction = auctions[_auctionId];
        bps = getFeeBps(auction.highestBidInDollar);
        feeAmount = (auction.highestBid * bps) / MAX_BPS;
        sellerAmount = auction.highestBid - feeAmount;
    }

    /// @notice 覆盖 V1 结算逻辑：在结算时按美元档位抽取手续费
    /// @dev
    /// - 先交割 NFT，再分配资金（卖家净收 + 平台手续费）；
    /// - 手续费按成交资产支付：ETH 出价收 ETH，ERC20 出价收 ERC20；
    /// - 费率依据 `highestBidInDollar`，但扣款基于 `highestBid` 资产数量。
    function end(uint256 _auctionId) external override {
        Auction storage auction = auctions[_auctionId];
        require(isEnded(_auctionId), "not ended");
        require(auction.highestBidder != address(0), "no bids");

        auction.nft.transferFrom(address(this), auction.highestBidder, auction.nftId);

        if (auction.highestBid > 0) {
            // 费率按美元档位确定，金额按成交资产计算。
            uint16 bps = getFeeBps(auction.highestBidInDollar);
            uint256 feeAmount = (auction.highestBid * bps) / MAX_BPS;
            uint256 sellerAmount = auction.highestBid - feeAmount;

            if (auction.highestBidToken == address(0)) {
                // ETH 结算路径
                payable(auction.seller).transfer(sellerAmount);
                if (feeAmount > 0) {
                    payable(feeRecipient).transfer(feeAmount);
                }
            } else {
                // ERC20 结算路径
                IERC20(auction.highestBidToken).transfer(auction.seller, sellerAmount);
                if (feeAmount > 0) {
                    IERC20(auction.highestBidToken).transfer(feeRecipient, feeAmount);
                }
            }

            emit FeeCharged(
                _auctionId,
                auction.highestBidToken,
                auction.highestBid,
                feeAmount,
                sellerAmount,
                auction.seller,
                feeRecipient
            );
        }
        emit EndBid(_auctionId);
    }

    /// @dev 内部设置手续费接收地址并做输入校验
    function _setFeeRecipient(address _feeRecipient) internal {
        require(_feeRecipient != address(0), "invalid fee recipient");
        feeRecipient = _feeRecipient;
        emit FeeRecipientUpdated(_feeRecipient);
    }

    /// @dev 内部设置动态费率配置并做一致性校验
    /// 要求：
    /// - 第一档阈值 > 0；
    /// - 第二档阈值必须大于第一档；
    /// - 所有 bps 不可超过 100%。
    function _setDynamicFeeConfig(
        uint256 _tier1UpperUsd,
        uint256 _tier2UpperUsd,
        uint16 _tier1Bps,
        uint16 _tier2Bps,
        uint16 _tier3Bps
    ) internal {
        require(_tier1UpperUsd > 0, "invalid tier1");
        require(_tier2UpperUsd > _tier1UpperUsd, "invalid tier2");
        require(_tier1Bps <= MAX_BPS && _tier2Bps <= MAX_BPS && _tier3Bps <= MAX_BPS, "invalid bps");

        tier1UpperUsd = _tier1UpperUsd;
        tier2UpperUsd = _tier2UpperUsd;
        tier1Bps = _tier1Bps;
        tier2Bps = _tier2Bps;
        tier3Bps = _tier3Bps;

        emit FeeConfigUpdated(_tier1UpperUsd, _tier2UpperUsd, _tier1Bps, _tier2Bps, _tier3Bps);
    }
}
