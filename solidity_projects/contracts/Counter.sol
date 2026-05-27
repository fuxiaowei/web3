// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Counter{
    uint public count;

    function incr() external {
        count +=1;
    } 

    function dec()external {
        count -=1;
    }

    // 重置计数
    function reset() public {
        count = 0;
    }
}