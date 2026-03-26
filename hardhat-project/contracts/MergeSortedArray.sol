// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

/// @title MergeSortedArray
/// @notice 合并两个已按非降序排列的整数数组，得到新的非降序数组（经典「双指针 / 归并」子过程）。
/// @dev 时间复杂度 O(m+n)，额外空间 O(m+n)。使用 `<=` 在相等时优先取 `nums1` 的元素，合并结果稳定（同值相对顺序与 nums1 在前一致）。
contract MergeSortedArray {
    /*
     * =============================================================================
     * 5. 合并两个有序数组 (Merge Sorted Array)
     * =============================================================================
     * 题目：给定两个各自升序的数组，合并为一个升序数组。
     * 示例：nums1 = [1,3,5], nums2 = [2,4,6] → [1,2,3,4,5,6]
     * =============================================================================
     */

    /// @notice 将两个升序数组合并为一个升序数组
    /// @dev 算法步骤：
    ///      1. 用指针 i、j 分别指向 nums1、nums2 的当前候选元素。
    ///      2. 在「两个数组都还有元素」时，比较 nums1[i] 与 nums2[j]，将较小者写入 result[k]，
    ///         相等时取 nums1[i]（`<=`），然后移动对应指针与 k。
    ///      3. 当某一数组先耗尽时，另一数组剩余部分已整体大于 result 中已有最大值，
    ///         可直接依次尾插，无需再比较。
    /// @param nums1 升序数组（允许为空）
    /// @param nums2 升序数组（允许为空）
    /// @return result 长度为 len(nums1)+len(nums2) 的升序数组
    function merge(int[] memory nums1, int[] memory nums2) public pure returns (int[] memory) {
        // 双指针：i 扫描 nums1，j 扫描 nums2（均从 0 开始）
        uint i = 0;
        uint j = 0;

        uint len1 = nums1.length;
        uint len2 = nums2.length;

        // 结果一次性分配足量长度；k 为写入下标
        int[] memory result = new int[](len1 + len2);
        uint k = 0;

        // 主循环：两个数组都还有未合并元素时，每次选当前较小的头元素推进
        while (i < len1 && j < len2) {
            if (nums1[i] <= nums2[j]) {
                result[k] = nums1[i];
                i++;
            } else {
                result[k] = nums2[j];
                j++;
            }
            k++;
        }

        // 仅 nums1 还有剩余：已全部 ≥ 已写入部分，顺序拷贝即可
        while (i < len1) {
            result[k] = nums1[i];
            i++;
            k++;
        }

        // 仅 nums2 还有剩余：同上
        while (j < len2) {
            result[k] = nums2[j];
            j++;
            k++;
        }

        return result;
    }
}
