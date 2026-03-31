// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract SimpleToken {
    string public name;
    string public symbol;
    address public creator;
    uint256 public totalSupply;
    
    mapping(address => uint256) public balances;
    
    constructor(string memory _name, string memory _symbol, uint256 _supply) {
        name = _name;
        symbol = _symbol;
        creator = msg.sender;
        totalSupply = _supply;
        balances[msg.sender] = _supply;
    }
    
    function transfer(address to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }
}

contract TokenFactory {
    SimpleToken[] public tokens;
    mapping(address => address[]) public userTokens;
    
    event TokenCreated(address tokenAddress, string name, string symbol);
    
    function createToken(
        string memory name,
        string memory symbol,
        uint256 initialSupply
    ) public returns (address) {
        SimpleToken newToken = new SimpleToken(name, symbol, initialSupply);
        
        tokens.push(newToken);
        userTokens[msg.sender].push(address(newToken));
        
        emit TokenCreated(address(newToken), name, symbol);
        return address(newToken);
    }
    
    function getTokenCount() public view returns (uint256) {
        return tokens.length;
    }
    
    function getUserTokens(address user) public view returns (address[] memory) {
        return userTokens[user];
    }
}

