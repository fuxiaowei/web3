// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FunctionVisibility {
    uint256 public value;

    // public函数
    function setValuePublic(uint256 _value) public {
        value = _value;
    }

    // external函数
    function setValueExternal(uint256 _value) external {
        value = _value;
    }
}