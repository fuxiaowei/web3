// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract LoopDemo {
    
    // for循环 - 计算1到n的和
    function sum(uint n) public pure returns (uint) {
        uint total = 0;
        for (uint i = 0; i <= n; i++) {
            total += i;
        }
        return total;
    }
    
    // for循环 - 数组求和
    function sumArray(uint[] memory numbers) public pure returns (uint) {
        uint total = 0;
        for (uint i = 0; i < numbers.length; i++) {
            total += numbers[i];
        }
        return total;
    }
    
    // while循环 - 倒计数
    function countdown(uint start) public pure returns (uint) {
        uint count = start;
        while (count > 0) {
            count--;
        }
        return count;
    }
    
    // do-while循环
    function doWhileDemo(uint n) public pure returns (uint) {
        uint i = 0;
        uint result = 0;
        do {
            result += i;
            i++;
        } while (i < n);
        return result;
    }
    
    // break - 查找数组中的元素
    function findTarget(uint[] memory arr, uint target) 
        public pure returns (bool, uint) 
    {
        for (uint i = 0; i < arr.length; i++) {
            if (arr[i] == target) {
                return (true, i);  // 找到就break（通过return）
            }
        }
        return (false, 0);
    }
    
    // continue - 只累加偶数
    function sumEven(uint n) public pure returns (uint) {
        uint total = 0;
        for (uint i = 0; i <= n; i++) {
            if (i % 2 != 0) {
                continue;  // 跳过奇数
            }
            total += i;
        }
        return total;
    }
    
    // 嵌套循环 - 乘法表
    function multiplicationTable(uint n) public pure returns (uint[][] memory) {
        uint[][] memory table = new uint[][](n);
        for (uint i = 0; i < n; i++) {
            table[i] = new uint[](n);
            for (uint j = 0; j < n; j++) {
                table[i][j] = (i + 1) * (j + 1);
            }
        }
        return table;
    }
}