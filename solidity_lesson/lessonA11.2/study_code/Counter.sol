// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title Counter
 * @dev 一个简单的计数器合约，用于演示Hardhat 3单元测试
 * @notice 这个合约展示了状态变量、函数和事件的基本用法
 */
contract Counter {
    /**
     * @dev 公开的状态变量，自动生成getter函数
     * @notice 可以通过x()函数访问当前值
     */
    uint256 public x;
    
    /**
     * @dev 递增事件
     * @param by 递增的数值
     */
    event Increment(uint256 by);
    
    /**
     * @dev 构造函数
     * @notice 初始化x为0
     */
    constructor() {
        x = 0;
    }
    
    /**
     * @dev 将数字递增1
     * @notice 调用此函数会将x加1，并触发Increment事件
     */
    function inc() public {
        x++;
        emit Increment(1);
    }
    
    /**
     * @dev 按指定数量递增
     * @param by 要递增的数量
     * @notice 调用此函数会将x增加指定数量，并触发Increment事件
     * @notice 如果by为0，会回退并显示错误消息
     */
    function incBy(uint256 by) public {
        require(by > 0, "incBy: increment should be positive");
        x += by;
        emit Increment(by);
    }
    
    /**
     * @dev 设置数字的值
     * @param _x 新的数字值
     * @notice 任何人都可以调用此函数来设置数字
     */
    function setNumber(uint256 _x) public {
        x = _x;
    }
}

