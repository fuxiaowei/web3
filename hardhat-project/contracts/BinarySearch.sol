// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BinarySearch {
    // 二分查找：在升序有序数组中找到目标值的索引
    // 找到：返回下标 index
    // 未找到：返回 type(uint).max (表示无效值)
    function search(int[] memory arr, int target) public pure returns (uint) {
        // 定义左右指针
        uint left = 0;
        uint right = arr.length - 1;

        // 核心循环：left <= right
        while (left <= right) {
            // 求中间索引（防止溢出写法）
            uint mid = left + (right - left) / 2;

            if (arr[mid] == target) {
                // 找到目标，直接返回下标
                return mid;
            } else if (arr[mid] < target) {
                // 目标在右半部分 → 左指针右移
                left = mid + 1;
            } else {
                // 目标在左半部分 → 右指针左移
                // 安全判断：避免 right 下溢
                if (mid == 0) break;
                right = mid - 1;
            }
        }

        // 未找到 → 返回最大值表示不存在
        return type(uint).max;
    }
}

//6. ✅ 二分查找 (Binary Search)
//题目描述：在一个有序数组中查找目标值。