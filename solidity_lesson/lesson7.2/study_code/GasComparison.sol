// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract GasComparison {
    uint256 public value = 100;

    // 自定义错误
    error ValueTooHigh(uint256 current, uint256 max);
    error ValueTooLow();

    // require + 字符串
    function testRequireString(uint256 amount) public {
        require(amount <= value, "Value too high");
        value = amount;
    }

    // require + 自定义错误
    function testRequireCustomError(uint256 amount) public {
        if (amount > value) revert ValueTooHigh(value, amount);
        value = amount;
    }

    // assert
    function testAssert(uint256 amount) public {
        value = amount;
        assert(value <= 1000);
    }

    // revert + 字符串
    function testRevertString(uint256 amount) public {
        if (amount > value) revert("Value too high");
        value = amount;
    }

    // revert + 自定义错误
    function testRevertCustomError(uint256 amount) public {
        if (amount > value) revert ValueTooHigh(value, amount);
        value = amount;
    }

    // 触发错误的辅助函数
    function triggerRequireString() public {
        testRequireString(200);
    }

    function triggerRequireCustomError() public {
        testRequireCustomError(200);
    }

    function triggerAssert() public {
        testAssert(2000);
    }

    function triggerRevertString() public {
        testRevertString(200);
    }

    function triggerRevertCustomError() public {
        testRevertCustomError(200);
    }
}