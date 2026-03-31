// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract SignatureExample {
    // 存储消息哈希
    mapping(bytes32 => bool) public processedMessages;

    function processMessage(
        bytes32 messageHash,
        bytes memory signature
    ) public {
        require(!processedMessages[messageHash], "Already processed");

        // 验证签名...
        // address signer = recoverSigner(messageHash, signature);

        processedMessages[messageHash] = true;
    }

    function hashMessage(string memory message) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(message));
    }
}