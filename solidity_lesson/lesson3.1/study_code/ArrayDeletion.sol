// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ArrayDeletion {

    uint[] public numbers;

    // 初始化数组
    function initialize() public {
        delete numbers;  // 清空数组
        numbers.push(10);
        numbers.push(20);
        numbers.push(30);
        numbers.push(40);
        numbers.push(50);
    }

    // 方法1：保持顺序删除
    function removeWithOrder(uint index) public {
        require(index < numbers.length, "Index out of bounds");

        // 将后面的元素向前移动
        for(uint i = index; i < numbers.length - 1; i++) {
            numbers[i] = numbers[i + 1];
        }

        // 删除最后一个元素
        numbers.pop();
    }

    // 方法2：快速删除（不保持顺序）
    function quickRemove(uint index) public {
        require(index < numbers.length, "Index out of bounds");

        // 用最后一个元素替换要删除的元素
        numbers[index] = numbers[numbers.length - 1];

        // 删除最后一个元素
        numbers.pop();
    }

    // 获取数组
    function getArray() public view returns (uint[] memory) {
        return numbers;
    }

    // 获取长度
    function getLength() public view returns (uint) {
        return numbers.length;
    }
}