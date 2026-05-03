// SPDX-License-Identifier: UNLICENSED
pragma solidity >=0.8.2 <0.9.0;

contract ContractStructure{
    uint256 public balance;

    constructor(uint256 init){
        balance = init;
    }

    event BalanceChanged(uint256 old, uint256 incr);

    modifier IncrimentRange(uint256 incr){
        // 修饰器对函数输入条件进行约束
        require(incr>100,"too small!");
        // 执行被修饰函数的逻辑
        _;
    }

    function addBalance(uint256 incr) public IncrimentRange(incr) {
        uint256 old = balance;
        balance += incr;
        emit BalanceChanged(old, incr);
    }
}




