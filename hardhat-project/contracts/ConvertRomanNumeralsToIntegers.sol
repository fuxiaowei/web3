// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "hardhat/console.sol";

contract ConvertRomanNumeralsToIntegers {
    // 罗马字符 => 数字 映射
    function getValue(bytes1 char) internal pure returns (uint) {
        if (char == "I") return 1;
        if (char == "V") return 5;
        if (char == "X") return 10;
        if (char == "L") return 50;
        if (char == "C") return 100;
        if (char == "D") return 500;
        if (char == "M") return 1000;
        return 0;
    }

    /// @notice 将合法罗马数字字符串转为整数（题目约定 [1, 3999]，从左到右一次扫描）
    /// @dev 算法要点：从左向右看每一位的「面值」；若当前位面值 **严格小于** 右边一位，说明当前位在「减法组合」里
    ///      （如 IV：先遇到 I，右边 V 更大 → I 表示 -1，随后 V 表示 +5，合计 4；IX、XL、XC、CD、CM 同理）。
    ///      否则当前位直接累加（如 III 三个 I 都是 +1）。
    function romanToInt(string memory s) public pure returns (uint) {
        // string 在 Solidity 中不可下标访问，先转为可索引的 bytes（每个 ASCII 罗马字母占 1 字节）
        bytes memory arr = bytes(s);
        uint len = arr.length;
        uint result = 0;

        // 从左到右遍历每一位罗马符号（中间结果只用 uint 累加，避免「先减后加」导致 uint 下溢）
        for (uint i = 0; i < len; i++) {
            uint value = getValue(arr[i]);
            console.log("value", value);

            if (i < len - 1) {
                uint nextVal = getValue(arr[i + 1]);
                if (value < nextVal) {
                    // 减法组合：一次加上 (大 − 小)，并跳过右邻（已在本次处理）
                    result += nextVal - value;
                    i++;
                    console.log("pair result:", result);
                } else {
                    result += value;
                    console.log("else  result:", result);
                }
            } else {
                result += value;
                console.log("else  result:", result);
            }
        }

        return result;
    }
}

///*
// * =============================================================================
// * 13. 罗马数字转整数
// * =============================================================================
// *
// * 【字符与数值】
// *   罗马数字由以下七种字符组成：I、V、X、L、C、D、M。
// *
// *     字符     数值
// *     ----     ----
// *     I         1
// *     V         5
// *     X        10
// *     L        50
// *     C       100
// *     D       500
// *     M      1000
// *
// * 【书写规则】
// *   - 一般规则：较小的数字写在较大数字的右侧，表示相加。
// *     例：II = 2；XII = 12；XXVII = 27。
// *   - 特例（减法记法）：较小数字写在较大数字左侧时，表示「大 − 小」。
// *     例：4 写作 IV（5−1），9 写作 IX（10−1）；不能写作 IIII。
// *   - 减法组合仅六种：
// *       I 在 V、X 左侧 → 4、9
// *       X 在 L、C 左侧 → 40、90
// *       C 在 D、M 左侧 → 400、900
// *
// * 【任务】给定罗马数字字符串 s，转换为整数。
// *
// * 【示例】
// *   示例 1  输入 "III"      输出 3
// *   示例 2  输入 "IV"       输出 4
// *   示例 3  输入 "IX"       输出 9
// *   示例 4  输入 "LVIII"    输出 58   （L=50, V=5, III=3）
// *   示例 5  输入 "MCMXCIV"  输出 1994 （M=1000, CM=900, XC=90, IV=4）
// *
// * 【提示与数据范围】
// *   - 1 ≤ s.length ≤ 15
// *   - s 仅含 'I','V','X','L','C','D','M'
// *   - 保证 s 合法，对应整数 ∈ [1, 3999]
// *   - 测试用例符合规范，无非法跨位写法；例如 49 为 XLIX、999 为 CMXCIX，
// *     不会出现 IL、IM 等不符合规则的串。
// * =============================================================================
// */