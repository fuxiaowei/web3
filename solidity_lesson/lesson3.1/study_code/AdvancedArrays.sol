// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AdvancedArrays {

    // 二维数组
    uint[][] public matrix;

    // Memory数组
    function createMemoryArray(uint size) public pure returns (uint[] memory) {
        uint[] memory arr = new uint[](size);
        for(uint i = 0; i < size; i++) {
            arr[i] = i * 10;
        }
        return arr;
    }

    // 添加一行到矩阵
    function addRow(uint[] memory row) public {
        matrix.push(row);
    }

    // 获取矩阵元素
    function getElement(uint row, uint col) public view returns (uint) {
        require(row < matrix.length, "Row out of bounds");
        require(col < matrix[row].length, "Column out of bounds");
        return matrix[row][col];
    }

    // 获取整个矩阵
    function getMatrix() public view returns (uint[][] memory) {
        return matrix;
    }

    // 计算矩阵所有元素的和
    function sumMatrix() public view returns (uint) {
        uint total = 0;
        uint rows = matrix.length;
        for(uint i = 0; i < rows; i++) {
            uint cols = matrix[i].length;
            for(uint j = 0; j < cols; j++) {
                total += matrix[i][j];
            }
        }
        return total;
    }

    // 定长二维数组
    uint[3][4] public fixedMatrix;

    function initFixedMatrix() public {
        // fixedMatrix 是 4 个长度为 3 的数组
        fixedMatrix[0] = [1, 2, 3];
        fixedMatrix[1] = [4, 5, 6];
        fixedMatrix[2] = [7, 8, 9];
        fixedMatrix[3] = [10, 11, 12];
    }

    function getFixedElement(uint row, uint col) public view returns (uint) {
        return fixedMatrix[row][col];
    }
}