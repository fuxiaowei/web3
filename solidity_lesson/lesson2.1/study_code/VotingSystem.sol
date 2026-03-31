// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract VotingSystem {
    bool public votingOpen = true;
    mapping(address => bool) public hasVoted;

    function vote() public {
        require(votingOpen, "Voting is closed");           // 检查投票是否开放
        require(!hasVoted[msg.sender], "Already voted");   // 检查是否已投票

        hasVoted[msg.sender] = true;
        // 记录投票...
    }

    function closeVoting() public {
        votingOpen = false;
    }
}