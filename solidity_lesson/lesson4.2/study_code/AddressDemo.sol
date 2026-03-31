// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AddressDemo {
    address public owner;
    address payable public recipient;
    
    constructor() {
        owner = msg.sender;
    }
    
    // 查询余额
    function getBalance(address addr) public view returns (uint) {
        return addr.balance;
    }
    
    // 检查是否是合约
    function isContract(address addr) public view returns (bool) {
        return addr.code.length > 0;
    }
    
    // 接收以太币
    receive() external payable {}
    
    // 转账演示
    function transferDemo(address payable to) public payable {
        require(msg.value > 0, "Must send ETH");
        to.transfer(msg.value);
    }
    
    // call方法演示
    function callDemo(address payable to) public payable {
        (bool success, ) = to.call{value: msg.value}("");
        require(success, "Call failed");
    }
}