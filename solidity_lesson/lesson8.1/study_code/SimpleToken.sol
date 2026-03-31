// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 定义ERC20接口
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

// 简单的ERC20实现
contract SimpleToken is IERC20 {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    
    constructor() {
        _balances[msg.sender] = 1000000 * 10**18;
    }
    
    function transfer(address to, uint256 amount) external returns (bool) {
        require(_balances[msg.sender] >= amount, "Insufficient balance");
        _balances[msg.sender] -= amount;
        _balances[to] += amount;
        return true;
    }
    
    function balanceOf(address account) external view returns (uint256) {
        return _balances[account];
    }
    
    function approve(address spender, uint256 amount) external returns (bool) {
        _allowances[msg.sender][spender] = amount;
        return true;
    }
    
    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        require(_allowances[from][msg.sender] >= amount, "Insufficient allowance");
        require(_balances[from] >= amount, "Insufficient balance");
        
        _allowances[from][msg.sender] -= amount;
        _balances[from] -= amount;
        _balances[to] += amount;
        return true;
    }
}

// 代币交换合约（使用接口调用）
contract TokenSwap {
    IERC20 public tokenA;
    IERC20 public tokenB;
    
    event Swap(address indexed user, uint256 amountA, uint256 amountB);
    
    constructor(address _tokenA, address _tokenB) {
        tokenA = IERC20(_tokenA);
        tokenB = IERC20(_tokenB);
    }
    
    function swap(uint256 amountA) external {
        // 使用接口调用，类型安全
        require(tokenA.transferFrom(msg.sender, address(this), amountA), "Transfer A failed");
        
        uint256 amountB = amountA; // 简化的1:1兑换
        require(tokenB.transfer(msg.sender, amountB), "Transfer B failed");
        
        emit Swap(msg.sender, amountA, amountB);
    }
}