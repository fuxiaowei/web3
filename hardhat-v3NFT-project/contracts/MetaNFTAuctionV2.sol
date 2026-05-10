// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./MetaNFTAuction.sol";

/// @notice 升级版拍卖合约占位文件，用于演示代理升级时替换实现合约。
/// @dev 实际逻辑应在与 `MetaNFTAuction` 兼容的存储布局下实现，并配合 UUPS/Transparent Proxy 等模式部署。

contract MetaNFTAuctionV2 is MetaNFTAuction {
    function getVersion() external pure override returns (string memory){
        return "MetaNFTAuctionV2";
    }

    function newFeature()external pure returns (string memory){
        return "This is a new feature in V2";
    }
}
