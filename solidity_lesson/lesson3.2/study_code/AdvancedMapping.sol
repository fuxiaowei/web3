// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AdvancedMapping {

    // ========== 不同类型的mapping ==========

    // address => uint（最常用）
    mapping(address => uint) public balances;

    // uint => bool
    mapping(uint => bool) public isValid;

    // address => string
    mapping(address => string) public names;

    // address => address（推荐人系统）
    mapping(address => address) public referrer;

    // ========== 嵌套mapping ==========

    // mapping套mapping
    // 外层：用户地址
    // 内层：代币地址 => 余额
    mapping(address => mapping(address => uint)) public tokenBalances;

    function setTokenBalance(address token, uint amount) public {
        tokenBalances[msg.sender][token] = amount;
    }

    function getTokenBalance(address user, address token)
    public view returns (uint)
    {
        return tokenBalances[user][token];
    }

    // ERC20授权模式
    // owner => spender => amount
    mapping(address => mapping(address => uint)) public allowance;

    function approve(address spender, uint amount) public {
        allowance[msg.sender][spender] = amount;
    }

    // ========== 删除mapping的值 ==========

    function deleteBalance(address user) public {
        delete balances[user];  // 恢复为默认值0
    }

    // ========== 检查键是否被设置 ==========

    // 问题：无法直接知道键是否被设置过
    // 解决方案：使用额外的bool mapping

    mapping(address => uint) public userData;
    mapping(address => bool) public hasData;

    function setData(address user, uint value) public {
        userData[user] = value;
        hasData[user] = true;  // 标记已设置
    }

    function checkIfSet(address user) public view returns (bool) {
        return hasData[user];
    }

    // 或者使用数组辅助
    address[] public userList;
    mapping(address => bool) public isUser;

    function addUser(address user, uint balance) public {
        require(!isUser[user], "User already exists");

        balances[user] = balance;
        userList.push(user);
        isUser[user] = true;
    }

    function getUserCount() public view returns (uint) {
        return userList.length;
    }

    function getAllUsers() public view returns (address[] memory) {
        return userList;
    }
}