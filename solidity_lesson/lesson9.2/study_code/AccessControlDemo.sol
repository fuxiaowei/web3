// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 没有访问控制的危险合约
contract UnsafeToken {
    mapping(address => uint256) public balances;
    uint256 public totalSupply;
    
    // 任何人都可以铸造代币！
    function mint(address to, uint256 amount) public {
        balances[to] += amount;
        totalSupply += amount;
    }
    
    // 任何人都可以销毁合约！
    function destroy() public {
        selfdestruct(payable(msg.sender));
    }
}

// Ownable模式
contract OwnableToken {
    address public owner;
    mapping(address => uint256) public balances;
    uint256 public totalSupply;
    
    constructor() {
        owner = msg.sender;
    }
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    // 只有owner可以铸造
    function mint(address to, uint256 amount) public onlyOwner {
        balances[to] += amount;
        totalSupply += amount;
    }
    
    function transferOwnership(address newOwner) public onlyOwner {
        owner = newOwner;
    }
}

// RBAC模式
contract RBACToken {
    mapping(bytes32 => mapping(address => bool)) private roles;
    mapping(address => uint256) public balances;
    uint256 public totalSupply;
    bool public paused;
    
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    constructor() {
        roles[ADMIN_ROLE][msg.sender] = true;
    }
    
    modifier onlyRole(bytes32 role) {
        require(roles[role][msg.sender], "Access denied");
        _;
    }
    
    // 只有MINTER可以铸造
    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        balances[to] += amount;
        totalSupply += amount;
    }
    
    // 只有PAUSER可以暂停
    function pause() public onlyRole(PAUSER_ROLE) {
        paused = true;
    }
    
    // 只有ADMIN可以授予角色
    function grantRole(bytes32 role, address account) public onlyRole(ADMIN_ROLE) {
        roles[role][account] = true;
    }
}

