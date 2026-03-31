// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 注意：这是简化版本，实际使用时应导入 @openzeppelin/contracts
// import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
// import "@openzeppelin/contracts/access/Ownable.sol";
// import "@openzeppelin/contracts/security/Pausable.sol";

// 简化的ERC20实现（用于演示，实际应使用OpenZeppelin）
contract ERC20 {
    string public name;
    string public symbol;
    uint8 public decimals = 18;
    uint256 public totalSupply;
    
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
    }
    
    function transfer(address to, uint256 amount) public virtual returns (bool) {
        _transfer(msg.sender, to, amount);
        return true;
    }
    
    function _transfer(address from, address to, uint256 amount) internal {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");
        
        uint256 fromBalance = balanceOf[from];
        require(fromBalance >= amount, "ERC20: transfer amount exceeds balance");
        
        balanceOf[from] = fromBalance - amount;
        balanceOf[to] += amount;
        
        emit Transfer(from, to, amount);
    }
    
    function _mint(address to, uint256 amount) internal {
        totalSupply += amount;
        balanceOf[to] += amount;
        emit Transfer(address(0), to, amount);
    }
}

// 简化的Ownable实现
contract Ownable {
    address public owner;
    
    constructor() {
        owner = msg.sender;
    }
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Ownable: caller is not the owner");
        _;
    }
    
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "Ownable: new owner is the zero address");
        owner = newOwner;
    }
}

// 简化的Pausable实现
contract Pausable {
    bool public paused;
    
    modifier whenNotPaused() {
        require(!paused, "Pausable: paused");
        _;
    }
    
    function _pause() internal {
        paused = true;
    }
    
    function _unpause() internal {
        paused = false;
    }
}

// 使用OpenZeppelin风格的代币合约
contract MyToken is ERC20, Ownable, Pausable {
    constructor() ERC20("My Token", "MTK") {
        _mint(msg.sender, 1000 * 10**decimals());
    }

    // 重写transfer函数，添加暂停功能
    function transfer(address to, uint256 amount) public override whenNotPaused returns (bool) {
        return super.transfer(to, amount);
    }

    // onlyOwner修饰符确保只有所有者可以调用
    function pause() public onlyOwner {
        _pause();
    }
    
    function unpause() public onlyOwner {
        _unpause();
    }
}