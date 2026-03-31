// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Shape {
    // 定义一个虚函数，表示它可以被子合约重写
    function getArea() public view virtual returns (uint) {
        return 0; // 默认面积为0
    }

    function getDescription() public pure returns (string memory) {
        return "This is a generic shape.";
    }
}

contract Rectangle is Shape {
    uint public width;
    uint public height;

    constructor(uint _width, uint _height) {
        width = _width;
        height = _height;
    }

    // 重写 getArea 函数
    function getArea() public view override returns (uint) {
        return width * height;
    }
}

contract Circle is Shape {
    uint public radius;

    constructor(uint _radius) {
        radius = _radius;
    }

    // 重写 getArea 函数
    function getArea() public view override returns (uint) {
        // 简化计算，实际应为 PI * r * r
        return radius * radius * 3; // 假设 PI 为 3
    }
}