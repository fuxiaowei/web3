// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract A {
    function foo() public pure virtual returns (string memory) {
        return "A.foo";
    }
}

contract B is A {
    function foo() public pure virtual override returns (string memory) {
        return string.concat("B.foo -> ", super.foo()); // super.foo() calls A.foo
    }
}

contract C is A {
    function foo() public pure virtual override returns (string memory) {
        return string.concat("C.foo -> ", super.foo()); // super.foo() calls A.foo
    }
}

// D 继承 B 和 C。Solidity 使用 C3 线性化，并且对基类列表按“从右到左”处理。
// 对于 `contract D is B, C`，线性化顺序是：D, C, B, A。
// 因此，在 D 中调用 super.foo() 时，会先跳到线性化中的下一个合约 C。
contract D is B, C {
    // D 必须 override foo，因为它继承了 B 和 C，而 B 和 C 都提供了 foo 的实现。
    // 这里需要明确指定 override 哪些父合约的 foo。
    function foo() public pure override(B, C) returns (string memory) {
        // super.foo() 会调用线性化顺序中紧跟在 D 之后的实现。
        // 在 `D is B, C` 的情况下，顺序是 D, C, B, A，所以这里会先调用 C.foo。
        // 最终返回的字符串为："D.foo -> C.foo -> A.foo"。
        return string.concat("D.foo -> ", super.foo());
    }
}