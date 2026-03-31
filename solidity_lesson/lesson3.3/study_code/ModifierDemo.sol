// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ModifierDemo {
    address public owner;
    bool public paused = false;
    uint256 public value;
    uint256 public totalDeposits;

    constructor() {
        owner = msg.sender;
    }

    // Modifier 1: 权限控制
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    // Modifier 2: 状态检查
    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    // Modifier 3: 参数验证
    modifier validAddress(address _addr) {
        require(_addr != address(0), "Invalid address");
        _;
    }

    // Modifier 4: 最小值检查
    modifier minValue(uint256 _min) {
        require(msg.value >= _min, "Insufficient value");
        _;
    }

    // 单个modifier
    function setValue(uint256 _value) public onlyOwner {
        value = _value;
    }

    // 组合多个modifier
    function setValueWithChecks(uint256 _value)
    public
    onlyOwner
    whenNotPaused
    {
        value = _value;
    }

    // 带参数的modifier
    function transferOwnership(address newOwner)
    public
    onlyOwner
    validAddress(newOwner)
    {
        owner = newOwner;
    }

    // payable + modifier
    function deposit()
    public
    payable
    whenNotPaused
    minValue(0.01 ether)
    {
        totalDeposits += msg.value;
    }

    // 管理函数
    function pause() public onlyOwner {
        paused = true;
    }

    function unpause() public onlyOwner {
        paused = false;
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}