// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 内部库示例
library InternalStringLib {
    // internal函数，代码会被嵌入调用合约
    function concat(string memory a, string memory b) 
        internal 
        pure 
        returns (string memory) 
    {
        return string(abi.encodePacked(a, b));
    }
    
    function length(string memory str) internal pure returns (uint256) {
        return bytes(str).length;
    }
}

// 外部库示例
library ExternalStringLib {
    // public函数，需要独立部署
    function toUpperCase(string memory str) public pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        bytes memory result = new bytes(strBytes.length);
        
        for (uint i = 0; i < strBytes.length; i++) {
            // 将小写字母转换为大写
            if (strBytes[i] >= 0x61 && strBytes[i] <= 0x7A) {
                result[i] = bytes1(uint8(strBytes[i]) - 32);
            } else {
                result[i] = strBytes[i];
            }
        }
        
        return string(result);
    }
}

// 使用内部库的合约
contract InternalLibUser {
    using InternalStringLib for string;
    
    function combineStrings(string memory a, string memory b) 
        public 
        pure 
        returns (string memory) 
    {
        return a.concat(b);
    }
    
    function getLength(string memory str) public pure returns (uint256) {
        return str.length();
    }
}

// 使用外部库的合约
contract ExternalLibUser {
    function convertToUpper(string memory str) 
        public 
        pure 
        returns (string memory) 
    {
        // 需要显式调用外部库
        return ExternalStringLib.toUpperCase(str);
    }
}