// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract GasOptimization {

    uint[] public data;
    uint public constant MAX_SIZE = 100;

    // 初始化测试数据
    function initializeData(uint count) public {
        require(count <= MAX_SIZE, "Exceeds max size");
        delete data;
        for(uint i = 0; i < count; i++) {
            data.push(i + 1);
        }
    }

    //  未优化：每次读取length
    function sumBad() public view returns (uint) {
        uint total = 0;
        for(uint i = 0; i < data.length; i++) {
            total += data[i];
        }
        return total;
    }

    // 优化：缓存length
    function sumGood() public view returns (uint) {
        uint total = 0;
        uint len = data.length;  // 缓存length
        for(uint i = 0; i < len; i++) {
            total += data[i];
        }
        return total;
    }

    // 分批求和
    function sumRange(uint start, uint end) public view returns (uint) {
        require(start < end, "Invalid range");
        require(end <= data.length, "End exceeds array length");

        uint total = 0;
        for(uint i = start; i < end; i++) {
            total += data[i];
        }
        return total;
    }

    // 安全添加元素
    function safePush(uint value) public {
        require(data.length < MAX_SIZE, "Array is full");
        data.push(value);
    }

    // 使用calldata（外部函数）
    function processArrayCalldata(uint[] calldata arr) external pure returns (uint) {
        uint total = 0;
        uint len = arr.length;
        for(uint i = 0; i < len; i++) {
            total += arr[i];
        }
        return total;
    }

    // 使用memory（对比）
    function processArrayMemory(uint[] memory arr) public pure returns (uint) {
        uint total = 0;
        uint len = arr.length;
        for(uint i = 0; i < len; i++) {
            total += arr[i];
        }
        return total;
    }
}