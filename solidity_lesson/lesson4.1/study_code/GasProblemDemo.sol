// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract GasProblemDemo {
    
    // 危险：无限制的循环
    function dangerousLoop(uint[] memory data) public pure returns (uint) {
        uint total = 0;
        // 如果data很大，可能gas耗尽
        for (uint i = 0; i < data.length; i++) {
            total += data[i];
        }
        return total;
    }
    
    // 安全：限制循环次数
    function safeLoop(uint[] memory data) public pure returns (uint) {
        require(data.length <= 100, "Array too large");
        uint total = 0;
        for (uint i = 0; i < data.length; i++) {
            total += data[i];
        }
        return total;
    }
    
    // 更好：使用mapping替代数组遍历
    mapping(address => uint) public balances;
    
    function setBalance(address user, uint amount) public {
        balances[user] = amount;  // O(1) 操作，不需要循环
    }
    
    function getBalance(address user) public view returns (uint) {
        return balances[user];  // O(1) 查询
    }
    
    // 分批处理示例
    uint[] public largeArray;
    
    function processBatch(uint startIndex, uint batchSize) public {
        require(batchSize <= 50, "Batch too large");
        uint endIndex = startIndex + batchSize;
        require(endIndex <= largeArray.length, "Out of bounds");
        
        for (uint i = startIndex; i < endIndex; i++) {
            // 处理每个元素
            largeArray[i] = largeArray[i] * 2;
        }
    }
    
    // 使用事件记录，链下查询
    event DataAdded(address indexed user, uint value, uint timestamp);
    
    function addData(uint value) public {
        emit DataAdded(msg.sender, value, block.timestamp);
        // 不需要在链上存储所有历史记录
        // 可以通过事件日志查询
    }
}