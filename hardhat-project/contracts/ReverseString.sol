// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringReverse {
    /// @notice 原地反转 UTF-8 字节序列（多字节码点会被字节级反转，与常见 LeetCode 题一致）
    function reverse(string memory input) public pure returns (string memory) {
        bytes memory b = bytes(input);
        uint256 n = b.length;
        // 避免 n==0 时 n-1 下溢；长度 0/1 无需交换
        if (n < 2) {
            return input;
        }

        unchecked {
            uint256 left = 0;
            uint256 right = n - 1;
            while (left < right) {
                (b[left], b[right]) = (b[right], b[left]);
                ++left;
                --right;
            }
        }

        return string(b);
    }
}

//2. ✅ 反转字符串 (Reverse String)
//题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"