// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EventVsStorage {
    // 使用数组存储历史
    struct Transaction {
        address from;
        address to;
        uint256 amount;
        uint256 timestamp;
    }
    Transaction[] public transactions;

    // 使用事件记录历史
    event TransactionExecuted(
        address indexed from,
        address indexed to,
        uint256 amount,
        uint256 timestamp
    );

    // 使用存储
    function recordWithStorage(address to, uint256 amount) external {
        transactions.push(Transaction({
        from: msg.sender,
        to: to,
        amount: amount,
        timestamp: block.timestamp
        }));
    }

    // 使用事件
    function recordWithEvent(address to, uint256 amount) external {
        emit TransactionExecuted(msg.sender, to, amount, block.timestamp);
    }
}