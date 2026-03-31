// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UserManager {

    address[] public users;
    mapping(address => bool) public isUser;
    mapping(address => uint) public userIndex;

    uint public constant MAX_USERS = 1000;

    // 添加用户
    function addUser(address user) public {
        require(user != address(0), "Invalid address");
        require(!isUser[user], "User already exists");
        require(users.length < MAX_USERS, "Maximum users reached");

        users.push(user);
        isUser[user] = true;
        userIndex[user] = users.length - 1;
    }

    // 删除用户（快速删除）
    function removeUser(address user) public {
        require(isUser[user], "User does not exist");

        uint index = userIndex[user];
        uint lastIndex = users.length - 1;

        // 如果不是最后一个元素，用最后一个替换
        if(index != lastIndex) {
            address lastUser = users[lastIndex];
            users[index] = lastUser;
            userIndex[lastUser] = index;
        }

        // 删除最后一个元素
        users.pop();
        delete isUser[user];
        delete userIndex[user];
    }

    // 获取用户数量
    function getUserCount() public view returns (uint) {
        return users.length;
    }

    // 获取所有用户
    function getAllUsers() public view returns (address[] memory) {
        return users;
    }

    // 检查用户是否存在（O(1)复杂度）
    function checkUser(address user) public view returns (bool) {
        return isUser[user];
    }

    // 获取指定范围的用户
    function getUsersRange(uint start, uint end) public view returns (address[] memory) {
        require(start < end, "Invalid range");
        require(end <= users.length, "End exceeds array length");

        address[] memory result = new address[](end - start);
        for(uint i = start; i < end; i++) {
            result[i - start] = users[i];
        }
        return result;
    }
}