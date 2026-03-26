// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract Voting {
    // candidateId => votes（只对“当前轮次”有效）
    mapping(uint256 => uint256) private votes;

    // candidateId => 最后一次被更新到的轮次
    mapping(uint256 => uint256) private lastRound;

    // 全局轮次：resetVotes() 只做轮次递增，避免遍历清零带来的高 gas
    uint256 private currentRound = 1;

    // 允许用户为某个候选人投票
    function vote(uint256 candidateId) external {
        // 如果该候选人还没在本轮次更新过，则视为该轮次得票为 0
        if (lastRound[candidateId] != currentRound) {
            lastRound[candidateId] = currentRound;
            votes[candidateId] = 0;
        }

        unchecked {
            votes[candidateId] += 1;
        }
    }

    // 返回某个候选人的得票数
    function getVotes(uint256 candidateId) external view returns (uint256) {
        return lastRound[candidateId] == currentRound ? votes[candidateId] : 0;
    }

    // 重置所有候选人的得票数（O(1) gas）：通过轮次实现
    function resetVotes() external {
        unchecked {
            currentRound += 1;
        }
    }
}



//
//1. ✅ 创建一个名为Voting的合约，包含以下功能：
//一个mapping来存储候选人的得票数
//一个vote函数，允许用户投票给某个候选人
//一个getVotes函数，返回某个候选人的得票数
//一个resetVotes函数，重置所有候选人的得票数