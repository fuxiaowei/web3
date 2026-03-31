// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 逻辑合约 V1
contract ImplementationV1 {
    uint256 public value;
    address public owner;
    
    function setValue(uint256 _value) external {
        value = _value;
    }
    
    function getValue() external view returns (uint256) {
        return value;
    }
}

// 逻辑合约 V2（升级版）
contract ImplementationV2 {
    uint256 public value;
    address public owner;
    
    function setValue(uint256 _value) external {
        value = _value * 2; // 新逻辑：值翻倍
    }
    
    function getValue() external view returns (uint256) {
        return value;
    }
    
    // 新增功能
    function reset() external {
        value = 0;
    }
}

// 代理合约
contract Proxy {
    address public implementation;
    uint256 public value; // 与Implementation的存储布局一致
    address public owner;
    
    event Upgraded(address indexed newImplementation);
    
    constructor(address _implementation) {
        implementation = _implementation;
        owner = msg.sender;
    }
    
    // 升级函数
    function upgrade(address newImplementation) external {
        require(msg.sender == owner, "Not owner");
        implementation = newImplementation;
        emit Upgraded(newImplementation);
    }
    
    // fallback函数：将所有调用转发到实现合约
    fallback() external payable {
        address impl = implementation;
        require(impl != address(0), "Implementation not set");
        
        assembly {
            // 复制calldata
            calldatacopy(0, 0, calldatasize())
            
            // delegatecall到实现合约
            let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)
            
            // 复制返回数据
            returndatacopy(0, 0, returndatasize())
            
            // 根据结果返回或revert
            switch result
            case 0 { revert(0, returndatasize()) }
            default { return(0, returndatasize()) }
        }
    }
    
    receive() external payable {}
}