// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 不安全：推送模式
contract UnsafeReward {
    address[] public users;
    uint256 public rewardAmount = 0.1 ether;
    
    function addUser(address user) external {
        users.push(user);
    }
    
    // 危险：如果用户太多会Gas耗尽
    function distributeRewards() external {
        for (uint256 i = 0; i < users.length; i++) {
            (bool success, ) = payable(users[i]).call{value: rewardAmount}("");
            require(success, "Transfer failed");
        }
    }
    
    receive() external payable {}
}

// 可以拒绝接收ETH的恶意合约
contract MaliciousUser {
    receive() external payable {
        revert("I don't want your money!");
    }
}

// 安全：拉取模式
contract SafeReward {
    mapping(address => uint256) public rewards;
    
    event RewardSet(address indexed user, uint256 amount);
    event RewardClaimed(address indexed user, uint256 amount);
    
    function setReward(address user, uint256 amount) external {
        rewards[user] = amount;
        emit RewardSet(user, amount);
    }
    
    // 批量设置奖励（分页处理）
    function batchSetRewards(
        address[] calldata users, 
        uint256[] calldata amounts
    ) external {
        require(users.length == amounts.length, "Length mismatch");
        require(users.length <= 50, "Too many users"); // 限制批量大小
        
        for (uint256 i = 0; i < users.length; i++) {
            rewards[users[i]] = amounts[i];
        }
    }
    
    // 用户自己领取
    function claimReward() external {
        uint256 amount = rewards[msg.sender];
        require(amount > 0, "No rewards");
        
        rewards[msg.sender] = 0;
        
        (bool success, ) = payable(msg.sender).call{value: amount}("");
        require(success, "Transfer failed");
        
        emit RewardClaimed(msg.sender, amount);
    }
    
    receive() external payable {}
}