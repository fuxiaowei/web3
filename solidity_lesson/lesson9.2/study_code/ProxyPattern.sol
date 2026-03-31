// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 简单的代理合约
contract SimpleProxy {
    address public implementation;
    address public admin;
    uint256 public value;  // 数据存储
    
    constructor(address _implementation) {
        admin = msg.sender;
        implementation = _implementation;
    }
    
    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin");
        _;
    }
    
    function upgrade(address newImplementation) external onlyAdmin {
        implementation = newImplementation;
    }
    
    fallback() external payable {
        address impl = implementation;
        assembly {
            calldatacopy(0, 0, calldatasize())
            let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)
            returndatacopy(0, 0, returndatasize())
            switch result
            case 0 { revert(0, returndatasize()) }
            default { return(0, returndatasize()) }
        }
    }
}

// V1逻辑合约
contract ImplementationV1 {
    address public implementation;
    address public admin;
    uint256 public value;
    
    function setValue(uint256 _value) public {
        value = _value;
    }
    
    function getValue() public view returns (uint256) {
        return value;
    }
}

// V2逻辑合约 - 新增功能
contract ImplementationV2 {
    address public implementation;
    address public admin;
    uint256 public value;
    uint256 public multiplier;  // 新增变量
    
    function setValue(uint256 _value) public {
        value = _value * (multiplier == 0 ? 1 : multiplier);
    }
    
    function getValue() public view returns (uint256) {
        return value;
    }
    
    function setMultiplier(uint256 _multiplier) public {
        multiplier = _multiplier;
    }
}

