// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract SafeCounter {
    uint256 public count;

    function increment() public {
        count = count + 1;  // 安全！溢出会自动回退
        // 虽然uint256的最大值是天文数字，
        // 但万一真的溢出了，交易会回退，不会出问题
    }

    function decrement() public {
        count = count - 1;  // 安全！下溢会自动回退
        // 如果count是0，这个操作会失败
    }

    function safeDecrement() public {
        require(count > 0, "Count is already zero");
        count = count - 1;  // 更明确的错误信息
    }
}