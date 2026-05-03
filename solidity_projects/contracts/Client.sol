// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library Client {
    struct EVMTokenAmount {
        address token;
        uint256 amount;
    }

    struct Message {
        uint64 destinationChainSelector;
        address receiver;
        bytes data;
        EVMTokenAmount[] tokenAmounts;
        address feeToken;
        bytes extraArgs;
    }
}