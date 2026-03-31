// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RoleManagement {
    // 定义角色枚举
    enum Role { None, Owner, Admin, User }

    // 存储用户角色
    mapping(address => Role) public roles;

    // 事件
    event RoleGranted(address indexed account, Role role);

    constructor() {
        roles[msg.sender] = Role.Owner;
        emit RoleGranted(msg.sender, Role.Owner);
    }

    // Modifier: 只有Owner
    modifier onlyOwner() {
        require(roles[msg.sender] == Role.Owner, "Not owner");
        _;
    }

    // Modifier: 只有Admin
    modifier onlyAdmin() {
        require(roles[msg.sender] == Role.Admin, "Not admin");
        _;
    }

    // Modifier: 只有User
    modifier onlyUser() {
        require(roles[msg.sender] == Role.User, "Not user");
        _;
    }

    // 只有Owner可以添加Admin
    function addAdmin(address account) public onlyOwner {
        require(account != address(0), "Invalid address");
        roles[account] = Role.Admin;
        emit RoleGranted(account, Role.Admin);
    }

    // 只有Admin可以添加User
    function addUser(address account) public onlyAdmin {
        require(account != address(0), "Invalid address");
        roles[account] = Role.User;
        emit RoleGranted(account, Role.User);
    }

    // 所有人可以查询角色
    function getRole(address account) public view returns (Role) {
        return roles[account];
    }

    // 检查是否是Owner
    function isOwner(address account) public view returns (bool) {
        return roles[account] == Role.Owner;
    }

    // 检查是否是Admin
    function isAdmin(address account) public view returns (bool) {
        return roles[account] == Role.Admin;
    }

    // 检查是否是User
    function isUser(address account) public view returns (bool) {
        return roles[account] == Role.User;
    }
}