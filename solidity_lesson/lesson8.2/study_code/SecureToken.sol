// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";

// 使用OpenZeppelin Ownable
contract SecureToken is Ownable {
    string public name = "Secure Token";
    string public symbol = "STKN";
    uint8 public decimals = 18;
    
    mapping(address => uint256) public balances;
    uint256 public totalSupply;
    
    event Minted(address indexed to, uint256 amount);
    event Burned(address indexed from, uint256 amount);
    event Transfer(address indexed from, address indexed to, uint256 amount);
    
    constructor() Ownable(msg.sender) {}
    
    // 只有owner可以mint
    function mint(address to, uint256 amount) external onlyOwner {
        require(to != address(0), "Invalid address");
        balances[to] += amount;
        totalSupply += amount;
        emit Minted(to, amount);
    }
    
    // 任何人可以burn自己的代币
    function burn(uint256 amount) external {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        totalSupply -= amount;
        emit Burned(msg.sender, amount);
    }
    
    function transfer(address to, uint256 amount) external {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        require(to != address(0), "Invalid address");
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
        emit Transfer(msg.sender, to, amount);
    }
}