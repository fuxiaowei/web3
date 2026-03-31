// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 注意：在Remix中使用，需要在GitHub导入
// 在实际项目中，通过npm安装：npm install @openzeppelin/contracts

// 如果无法导入，可以直接复制OpenZeppelin的代码
// 这里演示Strings库的简化实现
library Strings {
    bytes16 private constant _HEX_SYMBOLS = "0123456789abcdef";
    
    function toString(uint256 value) internal pure returns (string memory) {
        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }
    
    function toHexString(uint256 value) internal pure returns (string memory) {
        if (value == 0) {
            return "0x00";
        }
        uint256 temp = value;
        uint256 length = 0;
        while (temp != 0) {
            length++;
            temp >>= 8;
        }
        return toHexString(value, length);
    }
    
    function toHexString(uint256 value, uint256 length) 
        internal 
        pure 
        returns (string memory) 
    {
        bytes memory buffer = new bytes(2 * length + 2);
        buffer[0] = "0";
        buffer[1] = "x";
        for (uint256 i = 2 * length + 1; i > 1; --i) {
            buffer[i] = _HEX_SYMBOLS[value & 0xf];
            value >>= 4;
        }
        require(value == 0, "Strings: hex length insufficient");
        return string(buffer);
    }
}

// 使用Strings库的NFT合约示例
contract SimpleNFT {
    using Strings for uint256;
    using Strings for address;
    
    string public name = "My NFT Collection";
    string public symbol = "MNFT";
    
    mapping(uint256 => address) public ownerOf;
    mapping(address => uint256) public balanceOf;
    uint256 public totalSupply;
    
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    
    // 铸造NFT
    function mint(address to) public returns (uint256) {
        require(to != address(0), "Cannot mint to zero address");
        
        totalSupply++;
        uint256 tokenId = totalSupply;
        
        ownerOf[tokenId] = to;
        balanceOf[to]++;
        
        emit Transfer(address(0), to, tokenId);
        return tokenId;
    }
    
    // 获取Token URI（使用Strings库）
    function tokenURI(uint256 tokenId) public view returns (string memory) {
        require(ownerOf[tokenId] != address(0), "Token does not exist");
        
        // 使用Strings库将tokenId转换为字符串
        string memory tokenIdStr = tokenId.toString();
        
        // 构建URI
        return string(abi.encodePacked(
            "https://api.mynft.com/token/",
            tokenIdStr
        ));
    }
    
    // 获取拥有者信息（使用Strings库）
    function getOwnerInfo(uint256 tokenId) public view returns (string memory) {
        address owner = ownerOf[tokenId];
        require(owner != address(0), "Token does not exist");
        
        // 将地址转换为十六进制字符串
        string memory ownerHex = uint256(uint160(owner)).toHexString(20);
        
        return string(abi.encodePacked(
            "Token #",
            tokenId.toString(),
            " is owned by ",
            ownerHex
        ));
    }
    
    // 批量铸造信息
    function getBatchMintInfo(uint256 count) public pure returns (string memory) {
        return string(abi.encodePacked(
            "Minting ",
            count.toString(),
            " NFTs"
        ));
    }
}