// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ErrorHandlingDemo {
    address public owner;
    mapping(address => uint) public balances;
    uint public totalSupply;
    
    constructor() {
        owner = msg.sender;
        totalSupply = 1000000;
        balances[msg.sender] = totalSupply;
    }
    
    // require - 输入验证
    function transfer(address to, uint amount) public {
        require(to != address(0), "Cannot transfer to zero address");
        require(amount > 0, "Amount must be positive");
        require(balances[msg.sender] >= amount, "Insufficient balance");
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }
    
    // require - 权限检查
    function mint(address to, uint amount) public {
        require(msg.sender == owner, "Only owner can mint");
        require(to != address(0), "Cannot mint to zero address");
        
        balances[to] += amount;
        totalSupply += amount;
    }
    
    // assert - 不变量检查
    function transferWithInvariant(address to, uint amount) public {
        uint balanceBefore = balances[msg.sender] + balances[to];
        
        require(to != address(0), "Invalid address");
        require(balances[msg.sender] >= amount, "Insufficient balance");
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
        
        // 检查不变量：总和应该不变
        assert(balances[msg.sender] + balances[to] == balanceBefore);
    }
    
    // revert - 复杂条件
    function complexTransfer(address to, uint amount, bool urgent) public {
        if (to == address(0)) {
            revert("Invalid recipient");
        }
        
        if (amount == 0) {
            revert("Amount cannot be zero");
        }
        
        if (balances[msg.sender] < amount) {
            revert("Insufficient balance");
        }
        
        if (urgent && amount > 10000) {
            revert("Urgent transfers limited to 10000");
        }
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }
    
    // 自定义错误（推荐）
    error InsufficientBalance(uint requested, uint available);
    error InvalidAddress(address addr);
    error Unauthorized(address caller);
    
    function transferWithCustomError(address to, uint amount) public {
        if (to == address(0)) {
            revert InvalidAddress(to);
        }
        
        if (balances[msg.sender] < amount) {
            revert InsufficientBalance(amount, balances[msg.sender]);
        }
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }
    
    function adminFunction() public {
        if (msg.sender != owner) {
            revert Unauthorized(msg.sender);
        }
        // 管理员操作
    }
}