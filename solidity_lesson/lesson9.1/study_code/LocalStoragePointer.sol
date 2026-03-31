// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract LocalStoragePointer {
    struct User {
        uint256 balance;
        uint256 score;
        uint256 timestamp;
    }

    mapping(address => User) public users;

    // 未优化：重复SLOAD
    function updateUserBad(address user, uint256 balance, uint256 score) external {
        users[user].balance = balance;
        users[user].score = score;
        users[user].timestamp = block.timestamp;
    }

    // 优化：使用局部指针
    function updateUserGood(address user, uint256 balance, uint256 score) external {
        User storage u = users[user];
        u.balance = balance;
        u.score = score;
        u.timestamp = block.timestamp;
    }
}