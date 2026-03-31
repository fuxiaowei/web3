// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import "./Token.sol";

/**
 * @title Vault
 * @dev 一个保险库合约，依赖于Token合约
 * @notice 这个合约演示了合约之间的依赖关系
 */
contract Vault {
    /**
     * @dev Token合约实例
     * @notice 存储Token合约的地址
     */
    Token public token;
    
    /**
     * @dev 存款映射
     * @notice 记录每个地址的存款数量
     */
    mapping(address => uint256) public deposits;
    
    /**
     * @dev 构造函数
     * @param _token Token合约地址
     * @notice 初始化时设置Token合约地址
     */
    constructor(address _token) {
        require(_token != address(0), "Invalid token address");
        token = Token(_token);
    }
    
    /**
     * @dev 存款函数
     * @param amount 存款数量
     * @notice 将Token从调用者账户转移到Vault，并记录存款
     */
    function deposit(uint256 amount) public {
        require(amount > 0, "Amount must be greater than 0");
        
        // 从调用者账户转移Token到Vault
        token.transfer(address(this), amount);
        
        // 记录存款
        deposits[msg.sender] += amount;
    }
    
    /**
     * @dev 查询Vault中的Token余额
     * @return Vault合约持有的Token数量
     */
    function getVaultBalance() public view returns (uint256) {
        return token.balanceOf(address(this));
    }
}

