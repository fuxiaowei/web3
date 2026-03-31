// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 被调用的目标合约
contract Target {
    uint256 public value;
    address public sender;
    
    event ValueChanged(uint256 newValue, address caller);
    
    function setValue(uint256 _value) external {
        value = _value;
        sender = msg.sender;
        emit ValueChanged(_value, msg.sender);
    }
    
    function getValue() external view returns (uint256) {
        return value;
    }
}

// 调用者合约
contract Caller {
    uint256 public value;
    address public sender;
    
    event CallResult(bool success, bytes data);
    
    // 使用call调用
    function testCall(address target, uint256 newValue) external {
        (bool success, bytes memory data) = target.call(
            abi.encodeWithSignature("setValue(uint256)", newValue)
        );
        require(success, "Call failed");
        emit CallResult(success, data);
    }
    
    // 使用delegatecall调用
    function testDelegatecall(address target, uint256 newValue) external {
        (bool success, bytes memory data) = target.delegatecall(
            abi.encodeWithSignature("setValue(uint256)", newValue)
        );
        require(success, "Delegatecall failed");
        emit CallResult(success, data);
    }
    
    // 使用staticcall调用
    function testStaticcall(address target) external {
        (bool success, bytes memory data) = target.staticcall(
            abi.encodeWithSignature("getValue()")
        );
        require(success, "Staticcall failed");
        emit CallResult(success, data);
    }
}