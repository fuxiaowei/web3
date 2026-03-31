// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ArrayBasics {

    // 定长数组
    uint[5] public fixedArray;
    uint[3] public numbers = [1, 2, 3];

    // 动态数组
    uint[] public dynamicArray;

    // 设置定长数组
    function setFixedArray() public {
        fixedArray[0] = 10;
        fixedArray[1] = 20;
        fixedArray[2] = 30;
        fixedArray[3] = 40;
        fixedArray[4] = 50;
    }

    // 添加元素（仅动态数组可用）
    function addElement(uint value) public {
        dynamicArray.push(value);
    }

    // 删除最后元素
    function removeLastElement() public {
        require(dynamicArray.length > 0, "Array is empty");
        dynamicArray.pop();
    }

    // 获取动态数组长度
    function getLength() public view returns (uint) {
        return dynamicArray.length;
    }

    // 获取整个数组
    function getArray() public view returns (uint[] memory) {
        return dynamicArray;
    }

    // 更新指定元素
    function updateElement(uint index, uint value) public {
        require(index < dynamicArray.length, "Index out of bounds");
        dynamicArray[index] = value;
    }

    // delete的效果
    function demonstrateDelete(uint index) public {
        require(index < dynamicArray.length, "Index out of bounds");
        delete dynamicArray[index];  // 注意：只是设为0，不改变length
    }
}