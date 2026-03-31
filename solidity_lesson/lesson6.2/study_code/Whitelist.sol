// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 地址集合库
library AddressSet {
    struct Set {
        address[] values;
        mapping(address => uint256) indexes;
    }
    
    // 添加地址
    function add(Set storage set, address value) internal returns (bool) {
        if (contains(set, value)) {
            return false;
        }
        
        set.values.push(value);
        set.indexes[value] = set.values.length;
        return true;
    }
    
    // 移除地址
    function remove(Set storage set, address value) internal returns (bool) {
        uint256 index = set.indexes[value];
        
        if (index == 0) {
            return false;
        }
        
        uint256 toDeleteIndex = index - 1;
        uint256 lastIndex = set.values.length - 1;
        
        if (toDeleteIndex != lastIndex) {
            address lastValue = set.values[lastIndex];
            set.values[toDeleteIndex] = lastValue;
            set.indexes[lastValue] = index;
        }
        
        set.values.pop();
        delete set.indexes[value];
        
        return true;
    }
    
    // 检查是否包含
    function contains(Set storage set, address value) 
        internal 
        view 
        returns (bool) 
    {
        return set.indexes[value] != 0;
    }
    
    // 获取长度
    function length(Set storage set) internal view returns (uint256) {
        return set.values.length;
    }
    
    // 根据索引获取地址
    function at(Set storage set, uint256 index) 
        internal 
        view 
        returns (address) 
    {
        require(index < set.values.length, "AddressSet: index out of bounds");
        return set.values[index];
    }
}

// 使用地址集合库的白名单合约
contract Whitelist {
    using AddressSet for AddressSet.Set;
    
    AddressSet.Set private whitelist;
    address public owner;
    
    event AddedToWhitelist(address indexed account);
    event RemovedFromWhitelist(address indexed account);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    // 添加到白名单
    function addToWhitelist(address account) public onlyOwner {
        require(account != address(0), "Cannot add zero address");
        require(whitelist.add(account), "Address already in whitelist");
        emit AddedToWhitelist(account);
    }
    
    // 从白名单移除
    function removeFromWhitelist(address account) public onlyOwner {
        require(whitelist.remove(account), "Address not in whitelist");
        emit RemovedFromWhitelist(account);
    }
    
    // 检查是否在白名单中
    function isWhitelisted(address account) public view returns (bool) {
        return whitelist.contains(account);
    }
    
    // 获取白名单大小
    function getWhitelistSize() public view returns (uint256) {
        return whitelist.length();
    }
    
    // 获取指定索引的地址
    function getWhitelistAddress(uint256 index) public view returns (address) {
        return whitelist.at(index);
    }
    
    // 获取所有白名单地址
    function getAllWhitelistAddresses() public view returns (address[] memory) {
        uint256 size = whitelist.length();
        address[] memory addresses = new address[](size);
        
        for (uint256 i = 0; i < size; i++) {
            addresses[i] = whitelist.at(i);
        }
        
        return addresses;
    }
    
    // 只允许白名单用户访问的函数示例
    function restrictedFunction() public view returns (string memory) {
        require(isWhitelisted(msg.sender), "Not whitelisted");
        return "Welcome! You are whitelisted.";
    }
}