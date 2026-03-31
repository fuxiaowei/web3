// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract AnonymousEventDemo {
    // 普通事件：最多3个indexed参数
    event RegularEvent(
        address indexed a,
        address indexed b,
        address indexed c,
        uint256 value
    );

    // 匿名事件：最多4个indexed参数
    event MyEvent(
        address indexed a,
        address indexed b,
        address indexed c,
        address indexed d,
        uint256 value
    ) anonymous;

    function triggerRegularEvent() public {
        emit RegularEvent(
            address(0x1),
            address(0x2),
            address(0x3),
            100
        );
    }

    function triggerAnonymousEvent() public {
        emit MyEvent(
            address(0x1),
            address(0x2),
            address(0x3),
            address(0x4),
            200
        );
    }
}