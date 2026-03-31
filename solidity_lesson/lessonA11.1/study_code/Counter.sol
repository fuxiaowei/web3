// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

/**
 * @title Counter
 * @dev 一个简单的计数器合约，演示Hardhat 3的基础功能
 * @notice 这个合约展示了状态变量、函数和事件的基本用法
 */
contract Counter {
    /**
     * @dev 公开的状态变量，自动生成getter函数
     * @notice 可以通过number()函数访问当前值
     */
    uint256 public number;
    
    /**
     * @dev 递增事件
     * @param by 递增的数值
     */
    event Increment(uint256 by);

    /**
     * @dev 设置数字的值
     * @param newNumber 新的数字值
     * @notice 任何人都可以调用此函数来设置数字
     */
    function setNumber(uint256 newNumber) public {
        number = newNumber;
    }

    /**
     * @dev 将数字递增1
     * @notice 调用此函数会将number加1，并触发Increment事件
     */
    function increment() public {
        number++;
        emit Increment(1);
    }

    /**
     * @dev 按指定数量递增
     * @param amount 要递增的数量
     * @notice 调用此函数会将number增加指定数量，并触发Increment事件
     */
    function incBy(uint256 amount) public {
        number += amount;
        emit Increment(amount);
    }
}

