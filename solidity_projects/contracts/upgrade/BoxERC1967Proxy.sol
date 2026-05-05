// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/// @dev UUPS 与透明代理均使用 EIP-1967 实现槽；此类便于部署时固定代理字节码。
contract BoxERC1967Proxy is ERC1967Proxy {
    constructor(address implementation, bytes memory initData) payable ERC1967Proxy(implementation, initData) {}
}