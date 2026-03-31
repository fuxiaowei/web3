// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UncheckedDemo {
    uint256[] public data;

    constructor() {
        // 初始化100个元素
        for (uint256 i = 0; i < 100; i++) {
            data.push(i);
        }
    }

    // 未优化：有溢出检查
    function sumWithCheck() external view returns (uint256) {
        uint256 total = 0;
        uint256 length = data.length;
        for (uint256 i = 0; i < length; i++) {
            total += data[i];
        }
        return total;
    }

    // 优化：unchecked
    function sumWithoutCheck() external view returns (uint256) {
        uint256 total = 0;
        uint256 length = data.length;
        for (uint256 i = 0; i < length;) {
            total += data[i];
        unchecked { i++; }
        }
        return total;
    }
}