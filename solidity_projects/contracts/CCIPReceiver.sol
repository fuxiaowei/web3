// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

abstract contract CCIPReceiver {
    address internal immutable i_router;

    constructor(address router) {
        i_router = router;
    }

    modifier onlyRouter() {
        require(msg.sender == i_router, "Only router can call");
        _;
    }

    function _ccipReceive(
        bytes32 messageId,
        uint64 sourceChainSelector,
        address sender,
        bytes calldata data
    ) internal virtual;
}