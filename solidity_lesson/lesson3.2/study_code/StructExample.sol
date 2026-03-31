// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StructExample {

    // ========== 定义结构体 ==========

    struct User {
        string name;
        uint age;
        address wallet;
        bool isActive;
    }

    // ========== 使用结构体 ==========

    // 单个结构体变量
    User public admin;

    // 结构体数组
    User[] public users;

    // ========== 创建结构体的三种方式 ==========

    function createUser1() public {
        // 方式1：直接赋值
        admin.name = "Alice";
        admin.age = 25;
        admin.wallet = msg.sender;
        admin.isActive = true;
    }

    function createUser2() public {
        // 方式2：使用构造器语法
        admin = User("Bob", 30, msg.sender, true);
    }

    function createUser3() public {
        // 方式3：使用键值对（推荐）
        admin = User({
        name: "Charlie",
        age: 35,
        wallet: msg.sender,
        isActive: true
        });
    }

    // ========== 添加到数组 ==========

    function addUser(string memory name, uint age) public {
        users.push(User({
        name: name,
        age: age,
        wallet: msg.sender,
        isActive: true
        }));
    }

    // ========== 访问和修改 ==========

    function getUserName(uint index) public view returns (string memory) {
        require(index < users.length, "Index out of bounds");
        return users[index].name;
    }

    function updateUserAge(uint index, uint newAge) public {
        require(index < users.length, "Index out of bounds");
        users[index].age = newAge;
    }

    function deactivateUser(uint index) public {
        require(index < users.length, "Index out of bounds");
        users[index].isActive = false;
    }

    // ========== 结构体作为参数 ==========

    function processUser(User memory user) public pure returns (string memory) {
        if(user.isActive) {
            return user.name;
        } else {
            return "Inactive";
        }
    }

    // ========== 结构体返回值 ==========

    function getUser(uint index) public view returns (User memory) {
        require(index < users.length, "Index out of bounds");
        return users[index];
    }
}
