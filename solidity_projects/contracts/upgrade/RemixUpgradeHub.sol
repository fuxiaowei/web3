// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./BoxV1.sol";
import "./BoxV2.sol";
import "./BoxERC1967Proxy.sol";

/**
 * @dev 给 Remix 用：请编译「本文件」，不要在只打开 BoxV1 时单独编译它。
 * Remix 的部署列表只反映当前这次编译涉及的源文件；BoxV2、Proxy 与 BoxV1 互不 import，
 * 只编 BoxV1 时不会出现 BoxV2 / BoxERC1967Proxy。
 */
contract RemixUpgradeHub {}
