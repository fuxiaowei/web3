// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./BoxV1.sol";

// 继承 V1
contract BoxV2 is BoxV1 {
    // 新增状态变量（只能往后加，不能插中间）
    string public message;

    // 重写旧方法
    function setNumber(uint256 _num) external override {
        number = _num * 2;
    }

    // 新增方法
    function setMessage(string memory _msg) external {
        message = _msg;
    }
}