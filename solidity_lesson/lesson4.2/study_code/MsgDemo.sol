// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MsgDemo {
    mapping(address => uint) public balances;
    
    event PaymentReceived(address indexed sender, uint value, bytes data);
    
    // msg.sender和msg.value演示
    function deposit() public payable {
        require(msg.value > 0, "Must send ETH");
        balances[msg.sender] += msg.value;
    }
    
    // msg.data演示
    function forward(address target) public {
        (bool success, ) = target.call(msg.data);
        require(success, "Forward failed");
    }
    
    // 接收函数，记录msg.data无法在receive中使用，因此只记录sender和value
    receive() external payable {
        emit PaymentReceived(msg.sender, msg.value, "");
    }
}