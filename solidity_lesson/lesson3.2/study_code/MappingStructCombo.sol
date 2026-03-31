// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MappingStructCombo {

    // ========== 定义结构体 ==========

    struct UserInfo {
        string name;
        uint balance;
        uint registeredAt;
        bool exists;  // 重要：标记是否存在
    }

    // ========== Mapping + Struct ==========

    // 核心：mapping存储用户信息
    mapping(address => UserInfo) public users;

    // 辅助：数组存储所有用户地址
    address[] public userAddresses;

    // 用户总数
    uint public userCount;

    // ========== 注册用户 ==========

    function register(string memory name) public {
        require(!users[msg.sender].exists, "Already registered");

        users[msg.sender] = UserInfo({
        name: name,
        balance: 0,
        registeredAt: block.timestamp,
        exists: true  // 标记为存在
        });

        userAddresses.push(msg.sender);
        userCount++;
    }

    // ========== 存款 ==========

    function deposit() public payable {
        require(users[msg.sender].exists, "Not registered");
        require(msg.value > 0, "Send some ETH");

        users[msg.sender].balance += msg.value;
    }

    // ========== 查询 ==========

    function getUserInfo(address user) public view
    returns (string memory name, uint balance, uint registeredAt)
    {
        require(users[user].exists, "User does not exist");

        UserInfo memory info = users[user];
        return (info.name, info.balance, info.registeredAt);
    }

    function isRegistered(address user) public view returns (bool) {
        return users[user].exists;
    }

    // ========== 遍历用户 ==========

    function getAllUsers() public view returns (address[] memory) {
        return userAddresses;
    }

    function getUsersByRange(uint start, uint end)
    public view returns (UserInfo[] memory)
    {
        require(start < end && end <= userAddresses.length, "Invalid range");

        uint count = end - start;
        UserInfo[] memory result = new UserInfo[](count);

        for(uint i = 0; i < count; i++) {
            result[i] = users[userAddresses[start + i]];
        }

        return result;
    }

    // ========== 复杂示例：投票系统 ==========

    struct Proposal {
        string description;
        uint voteCount;
        mapping(address => bool) voters;  // 嵌套mapping
        bool exists;
    }

    mapping(uint => Proposal) public proposals;
    uint public proposalCount;

    function createProposal(string memory description) public returns (uint) {
        uint proposalId = proposalCount;

        Proposal storage p = proposals[proposalId];
        p.description = description;
        p.voteCount = 0;
        p.exists = true;

        proposalCount++;
        return proposalId;
    }

    function vote(uint proposalId) public {
        Proposal storage p = proposals[proposalId];
        require(p.exists, "Proposal does not exist");
        require(!p.voters[msg.sender], "Already voted");

        p.voters[msg.sender] = true;
        p.voteCount++;
    }
}