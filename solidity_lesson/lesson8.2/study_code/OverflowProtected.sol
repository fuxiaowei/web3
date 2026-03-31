// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// Solidity 0.8.19 自动保护
contract OverflowProtected {
    uint8 public value = 255;
    uint8 public value2 = 0;
    
    // 会自动检查溢出并回滚
    function testOverflow() external {
        value++; // 会panic回滚
    }
    
    function testUnderflow() external {
        value2--; // 会panic回滚
    }
    
    // 安全的增加
    function safeIncrement() external {
        require(value < type(uint8).max, "Would overflow");
        value++;
    }
}

// unchecked块示例
contract UncheckedExample {
    uint256 public counter;
    
    // 不安全：跳过溢出检查
    function unsafeIncrement() external {
        unchecked {
            counter++;
        }
    }
    
    // 安全：确定不会溢出的场景
    function safeLoop(uint256[] calldata data) external returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < data.length; ) {
            sum += data[i];
            unchecked {
                i++; // 循环计数器确定安全
            }
        }
        return sum;
    }
    
    // 类型转换检查
    function safeCast(uint256 value) external pure returns (uint8) {
        require(value <= type(uint8).max, "Value too large");
        return uint8(value);
    }
}