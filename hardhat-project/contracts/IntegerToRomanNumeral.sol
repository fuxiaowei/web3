// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract IntegerToRomanNumeral {
    // 核心：从大到小排列所有可能的数值 + 对应罗马字符（包含6种减法组合）
    uint[] private values = [
        1000,
        900,
        500,
        400,
        100,
        90,
        50,
        40,
        10,
        9,
        5,
        4,
        1
    ];

    string[] private symbols = [
        "M",
        "CM",
        "D",
        "CD",
        "C",
        "XC",
        "L",
        "XL",
        "X",
        "IX",
        "V",
        "IV",
        "I"
    ];

    // 主函数：整数转罗马数字
    function intToRoman(uint num) public view returns (string memory) {
        // 复制数值，避免修改原始入参
        uint n = num;
        // 动态拼接结果
        string memory roman = "";

        // 遍历所有数值组合
        for (uint i = 0; i < values.length; i++) {
            // 当前数值能减多少次
            while (n >= values[i]) {
                roman = string.concat(roman, symbols[i]);
                n -= values[i];
            }
            // 提前退出优化
            if (n == 0) break;
        }

        return roman;
    }
}

/*
 * =============================================================================
 * 12. 整数转罗马数字（难度：中等）
 * =============================================================================
 *
 * 【任务】
 *   给定整数 num，将其转为罗马数字字符串。
 *
 * 【符号与数值】
 *   字符    数值
 *   ----    ----
 *   I        1
 *   V        5
 *   X       10
 *   L       50
 *   C      100
 *   D      500
 *   M     1000
 *
 * 【转换规则（按数位从高到低）】
 *   1. 若当前剩余数值的「这一位」不是 4 或 9：
 *      选出能减去的最大符号，接到结果末尾，减去该符号值，对剩余部分继续转换。
 *   2. 若这一位是 4 或 9：必须用减法形式（小数在大数左侧），且只能用下列六种：
 *        4 → IV    9 → IX    40 → XL    90 → XC    400 → CD    900 → CM
 *      例如：4 = V−I 写作 IV；9 = X−I 写作 IX。
 *   3. 仅 I、X、C、M（10 的幂）最多连续出现 3 次；V、L、D 不能连续重复。
 *      若同一符号需写 4 次，应改为减法形式（如 40 用 XL，不用 XXXX）。
 *   4. 注意：49 写作 XLIX（按位处理），不能写成 IL（非法减法）。
 *
 * 【示例】
 *   示例 1  输入 num = 3749   输出 "MMMDCCXLIX"
 *            3000 = MMM；700 = DCC；40 = XL；9 = IX
 *
 *   示例 2  输入 num = 58     输出 "LVIII"
 *            50 = L；8 = VIII
 *
 *   示例 3  输入 num = 1994   输出 "MCMXCIV"
 *            1000 = M；900 = CM；90 = XC；4 = IV
 *
 * 【数据范围】
 *   1 ≤ num ≤ 3999
 * =============================================================================
 */