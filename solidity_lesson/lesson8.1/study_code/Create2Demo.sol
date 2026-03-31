// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 简单的计数器合约
contract Counter {
    uint256 public count;
    address public owner;
    
    constructor(address _owner) {
        owner = _owner;
        count = 0;
    }
    
    function increment() external {
        count++;
    }
}

// 工厂合约
contract CounterFactory {
    event CounterCreated(address indexed counterAddress, bytes32 salt);
    
    // 使用new创建（地址不可预测）
    function createWithNew() external returns (address) {
        Counter counter = new Counter(msg.sender);
        return address(counter);
    }
    
    // 使用create2创建（地址可预测）
    function createWithCreate2(bytes32 salt) external returns (address) {
        Counter counter = new Counter{salt: salt}(msg.sender);
        emit CounterCreated(address(counter), salt);
        return address(counter);
    }
    
    // 预计算create2地址
    function computeAddress(bytes32 salt, address deployer) external view returns (address) {
        bytes memory bytecode = abi.encodePacked(
            type(Counter).creationCode,
            abi.encode(deployer)
        );
        
        bytes32 hash = keccak256(
            abi.encodePacked(
                bytes1(0xff),
                address(this),
                salt,
                keccak256(bytecode)
            )
        );
        
        return address(uint160(uint256(hash)));
    }
}