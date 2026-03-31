// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract VaultWithPause {
    mapping(address => uint256) public balances;
    bool public paused;
    address public admin;
    
    event Paused(address admin);
    event Unpaused(address admin);
    event Deposited(address user, uint256 amount);
    event Withdrawn(address user, uint256 amount);
    event EmergencyWithdrawal(address user, uint256 amount);
    
    constructor() {
        admin = msg.sender;
        paused = false;
    }
    
    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }
    
    modifier whenPaused() {
        require(paused, "Contract is not paused");
        _;
    }
    
    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin");
        _;
    }
    
    function deposit() public payable whenNotPaused {
        require(msg.value > 0, "Must deposit something");
        balances[msg.sender] += msg.value;
        emit Deposited(msg.sender, msg.value);
    }
    
    function withdraw(uint256 amount) public whenNotPaused {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        payable(msg.sender).transfer(amount);
        emit Withdrawn(msg.sender, amount);
    }
    
    // 紧急提现 - 只能在暂停时调用
    function emergencyWithdraw() public whenPaused {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");
        
        balances[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
        emit EmergencyWithdrawal(msg.sender, amount);
    }
    
    function pause() public onlyAdmin whenNotPaused {
        paused = true;
        emit Paused(admin);
    }
    
    function unpause() public onlyAdmin whenPaused {
        paused = false;
        emit Unpaused(admin);
    }
}

