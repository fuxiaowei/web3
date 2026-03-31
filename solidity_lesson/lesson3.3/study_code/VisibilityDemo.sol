// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract VisibilityDemo {
    uint256 private secretNumber = 42;

    // public: 外部、内部、继承都可调用
    function publicFunction() public pure returns (string memory) {
        return "This is public";
    }

    // external: 只能外部调用
    function externalFunction() external pure returns (string memory) {
        return "This is external";
    }

    // internal: 内部和继承可调用
    function internalFunction() internal pure returns (string memory) {
        return "This is internal";
    }

    // private: 只能本合约内部调用
    function privateFunction() private pure returns (string memory) {
        return "This is private";
    }

    // 测试内部调用
    function testInternalCall() public pure returns (string memory) {
        return internalFunction();  // 可以调用
    }

    // 测试external不能内部调用
    // function testExternalCall() public pure returns (string memory) {
    //     return externalFunction();  // 编译错误
    // }
}

// 测试继承
contract ChildContract is VisibilityDemo {
    function testInheritance() public pure returns (string memory) {
        return internalFunction();  // 可以调用internal
        // return privateFunction();  // 不能调用private
    }
}