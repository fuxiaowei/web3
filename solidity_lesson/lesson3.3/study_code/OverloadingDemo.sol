// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract OverloadingDemo {
    event Transfer(address indexed to, uint256 amount, string memo);

    // 版本1：简单转账
    function transfer(address to, uint256 amount) public {
        // 实际转账逻辑省略
        emit Transfer(to, amount, "");
    }

    // 版本2：带备注的转账
    function transfer(address to, uint256 amount, string memory memo) public {
        // 实际转账逻辑省略
        emit Transfer(to, amount, memo);
    }

    // 测试函数
    function testOverloading() public {
        address recipient = 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4;

        // 调用版本1
        transfer(recipient, 100);

        // 调用版本2
        transfer(recipient, 200, "Payment for service");
    }
}