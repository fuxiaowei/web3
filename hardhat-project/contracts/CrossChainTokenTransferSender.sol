// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.7;

import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/src/v0.8/ccip/interfaces/IRouterClient.sol";
import {OwnerIsCreator} from "@chainlink/contracts-ccip/src/v0.8/shared/access/OwnerIsCreator.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract CrossChainTokenTransferSender is OwnerIsCreator {
    using SafeERC20 for IERC20;

    event TokensSent(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address receiver,
        address token,
        uint256 tokenAmount,
        address feeToken,
        uint256 fees
    );

    address immutable i_router;

    constructor(address router) {
        i_router = router;
    }

    // 👉 这里加了 payable，就解决了！
    function sendTokens(
        uint64 destinationChainSelector,
        address receiver,
        address token,
        uint256 amount,
        address feeToken
    ) external payable onlyOwner returns (bytes32 messageId)
    {
        Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
        tokenAmounts[0] = Client.EVMTokenAmount({
            token: token,
            amount: amount
        });

        Client.EVM2AnyMessage memory evm2AnyMessage = Client.EVM2AnyMessage({
            receiver: abi.encode(receiver),
            data: "",
            tokenAmounts: tokenAmounts,
            extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200000})),
            feeToken: feeToken
        });

        IRouterClient router = IRouterClient(i_router);
        uint256 fees = router.getFee(destinationChainSelector, evm2AnyMessage);

        if (feeToken != address(0)) {
            IERC20(feeToken).safeIncreaseAllowance(address(router), fees);
        }

        messageId = router.ccipSend{value: msg.value}(
            destinationChainSelector,
            evm2AnyMessage
        );

        emit TokensSent(
            messageId,
            destinationChainSelector,
            receiver,
            token,
            amount,
            feeToken,
            fees
        );
    }
}