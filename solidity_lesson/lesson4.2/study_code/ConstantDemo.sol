// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ConstantDemo {
    // constant常量
    uint public constant MAX_SUPPLY = 1_000_000;
    uint public constant DECIMALS = 18;
    address public constant FEE_RECIPIENT = 0x0000000000000000000000000000000000000000;
    
    // immutable变量
    address public immutable owner;
    uint public immutable creationTime;
    address public immutable factory;
    
    // 普通变量
    uint public totalSupply;
    
    constructor(address _factory) {
        owner = msg.sender;
        creationTime = block.timestamp;
        factory = _factory;
    }
    
    // 使用常量
    function mint(uint amount) public {
        require(totalSupply + amount <= MAX_SUPPLY, "Exceeds max supply");
        totalSupply += amount;
    }
    
    // 计算函数
    function calculateWithDecimals(uint value) public pure returns (uint) {
        return value * 10**DECIMALS;
    }
}