// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MappingBasics {

    // ========== 基本mapping ==========

    // 语法：mapping(keyType => valueType) 变量名;
    mapping(address => uint) public balances;

    // 设置余额
    function setBalance(address user, uint amount) public {
        balances[user] = amount;
    }

    // 获取余额
    function getBalance(address user) public view returns (uint) {
        return balances[user];  // 如果不存在，返回0
    }

    // 增加余额
    function addBalance(address user, uint amount) public {
        balances[user] += amount;
    }

    // ========== mapping的特性 ==========

    mapping(uint => string) public idToName;

    function setName(uint id, string memory name) public {
        idToName[id] = name;
    }

    // 测试默认值
    function testDefaultValue() public view returns (string memory) {
        return idToName[999];  // 从未设置过的键，返回空字符串
    }
}