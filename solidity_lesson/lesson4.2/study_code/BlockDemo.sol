// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BlockDemo {
    uint public deploymentTime;
    uint public deploymentBlock;
    
    constructor() {
        deploymentTime = block.timestamp;
        deploymentBlock = block.number;
    }
    
    // 时间相关函数
    function timeSinceDeployment() public view returns (uint) {
        return block.timestamp - deploymentTime;
    }
    
    // 区块相关函数
    function blocksSinceDeployment() public view returns (uint) {
        return block.number - deploymentBlock;
    }
    
    // 时间锁演示
    function timeLock(uint lockDuration) public view returns (uint) {
        return block.timestamp + lockDuration;
    }
    
    // blockhash演示
    function getBlockHash(uint blockNumber) public view returns (bytes32) {
        return blockhash(blockNumber);
    }
}