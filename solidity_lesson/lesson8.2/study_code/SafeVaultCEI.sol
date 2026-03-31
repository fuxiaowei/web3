// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 安全版本：使用CEI模式
contract SafeVaultCEI {
    mapping(address => uint256) public balances;
    
    event Deposit(address indexed user, uint256 amount);
    event Withdrawal(address indexed user, uint256 amount);
    
    function deposit() external payable {
        balances[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }
    
    // 安全的提款函数（CEI模式）
    function withdraw() external {
        // 1. Checks - 检查
        uint256 balance = balances[msg.sender];
        require(balance > 0, "Insufficient balance");
        
        // 2. Effects - 更新状态（先更新！）
        balances[msg.sender] = 0;
        
        // 3. Interactions - 外部调用（最后执行）
        (bool success, ) = msg.sender.call{value: balance}("");
        require(success, "Transfer failed");
        
        emit Withdrawal(msg.sender, balance);
    }
    
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}