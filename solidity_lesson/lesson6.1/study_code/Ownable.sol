// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Ownable {
    address public owner; // 合约所有者地址

    // 构造函数：部署合约时，将部署者设置为所有者
    constructor() {
        owner = msg.sender;
    }

    // 定义 onlyOwner 修饰符
    modifier onlyOwner() {
        require(msg.sender == owner, "Caller is not the owner"); // 检查调用者是否为所有者
        _; // 执行被修饰的函数体
    }

    // 只有所有者才能调用的函数
    function doSomethingImportant() public onlyOwner returns (string memory) {
        return "This action can only be performed by the owner!";
    }

    // 允许所有者转移合约所有权
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner cannot be the zero address");
        owner = newOwner;
    }

    // 任何人都可以调用的函数
    function getInfo() public pure returns (string memory) {
        return "This is public information.";
    }
}