// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Client.sol";

interface IRouterClient {
    function isChainSupported(uint64 chainSelector) external view returns (bool);
    function getFee(uint64 destinationChainSelector, Client.Message calldata message, bytes calldata extraArgs) external view returns (uint256);
    function ccipSend(uint64 destinationChainSelector, Client.Message calldata message) external payable returns (bytes32);
}