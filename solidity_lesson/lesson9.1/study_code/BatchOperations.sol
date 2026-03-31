// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BatchOperations {
    mapping(address => uint256) public balances;

    constructor() {
        // 给当前账户初始余额
        balances[msg.sender] = 1000000;
    }

    // 单个转账
    function transfer(address to, uint256 amount) external {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }

    // 批量转账
    function batchTransfer(address[] calldata recipients, uint256[] calldata amounts) external {
        require(recipients.length == amounts.length, "Length mismatch");
        require(recipients.length <= 100, "Too many recipients");

        // 先计算总额
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < amounts.length;) {
            totalAmount += amounts[i];
        unchecked { i++; }
        }

        // 检查余额
        require(balances[msg.sender] >= totalAmount, "Insufficient balance");

        // 执行批量转账
        balances[msg.sender] -= totalAmount;
        for (uint256 i = 0; i < recipients.length;) {
            balances[recipients[i]] += amounts[i];
        unchecked { i++; }
        }
    }
}