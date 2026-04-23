// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.7;

import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract CrossChainTokenTransferReceiver is CCIPReceiver {
    using SafeERC20 for IERC20;

    event TokensReceived(
        address indexed receiver,
        address token,
        uint256 amount
    );

    constructor(address router) CCIPReceiver(router) {}

    function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
        address receiver = abi.decode(message.data, (address));

        // ✅ 修复：发送端用 tokenAmounts，接收端必须用 destTokenAmounts
        for (uint256 i = 0; i < message.destTokenAmounts.length; i++) {
            address token = message.destTokenAmounts[i].token;
            uint256 amount = message.destTokenAmounts[i].amount;

            IERC20(token).safeTransfer(receiver, amount);

            emit TokensReceived(receiver, token, amount);
        }
    }
}