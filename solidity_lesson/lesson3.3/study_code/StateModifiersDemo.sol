// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StateModifiersDemo {
    uint256 public counter = 0;

    // 默认：可以修改状态
    function incrementCounter() public {
        counter++;
    }

    // view：只读状态
    function getCounter() public view returns (uint256) {
        return counter;
    }

    function getDoubleCounter() public view returns (uint256) {
        return counter * 2;  // 可以读取和计算
    }

    // 错误示例：view不能修改状态
    // function badView() public view {
    //     counter++;  // 编译错误
    // }

    // pure：不读不写状态
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }

    // 错误示例：pure不能读状态
    // function badPure() public pure returns (uint256) {
    //     return counter;  // 编译错误
    // }

    // payable：可接收ETH
    function deposit() public payable {
        require(msg.value > 0, "Send some ETH");
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }

    // 错误示例：没有payable不能收ETH
    function normalFunction() public {
        // 如果调用时发送ETH，交易会失败
    }
}