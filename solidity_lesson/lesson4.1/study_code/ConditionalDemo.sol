// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ConditionalDemo {
    
    // 基本if语句
    function checkAge(uint age) public pure returns (bool) {
        if (age >= 18) {
            return true;
        }
        return false;
    }
    
    // if-else语句
    function checkValue(uint value) public pure returns (string memory) {
        if (value > 100) {
            return "High";
        } else {
            return "Low";
        }
    }
    
    // else if链 - 分数分级
    function getGrade(uint score) public pure returns (string memory) {
        if (score >= 90) {
            return "A";
        } else if (score >= 80) {
            return "B";
        } else if (score >= 70) {
            return "C";
        } else if (score >= 60) {
            return "D";
        } else {
            return "F";
        }
    }
    
    // 嵌套if
    function checkEligibility(uint age, uint score) public pure returns (string memory) {
        if (age >= 18) {
            if (score >= 60) {
                return "Eligible and Passed";
            } else {
                return "Eligible but Failed";
            }
        } else {
            return "Not Eligible";
        }
    }
    
    // 三元运算符
    function max(uint a, uint b) public pure returns (uint) {
        return a > b ? a : b;
    }
    
    // 三元运算符嵌套（不推荐过度使用）
    function classify(uint value) public pure returns (string memory) {
        return value > 100 ? "High" : value > 50 ? "Medium" : "Low";
    }
    
    // 组合条件判断
    function checkRange(uint value) public pure returns (bool) {
        return (value >= 10 && value <= 100);
    }
}