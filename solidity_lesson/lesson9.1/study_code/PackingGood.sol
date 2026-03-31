// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 未优化版本
contract StoragePackingBad {
    uint8 public a = 1;
    uint256 public b = 2;
    uint8 public c = 3;
    uint256 public d = 4;

    // 部署成本：可以查看Gas消耗
}

// 优化版本
contract StoragePackingGood {
    uint8 public a = 1;
    uint8 public c = 3;
    uint256 public b = 2;
    uint256 public d = 4;

    // 部署成本：对比未优化版本
}