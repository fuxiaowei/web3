// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract EventDemo {
    // 定义一个简单的Transfer事件
    event Transfer(address indexed from, address indexed to, uint256 value);

    // 定义一个包含更多信息的事件
    event DataUpdate(
        address indexed user,
        uint256 indexed id,
        string data,
        uint256 timestamp
    );

    mapping(address => uint256) public balances;

    constructor() {
        balances[msg.sender] = 1000;
    }

    // 转账函数，触发Transfer事件
    function transfer(address to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        balances[msg.sender] -= amount;
        balances[to] += amount;

        // 触发事件
        emit Transfer(msg.sender, to, amount);
    }

    // 更新数据函数，触发DataUpdate事件
    function updateData(uint256 id, string memory data) public {
        emit DataUpdate(msg.sender, id, data, block.timestamp);
    }
}