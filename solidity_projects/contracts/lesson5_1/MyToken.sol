// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract MyToken {
    // 代币基本信息
    string public name;  // 中文名
    string public symbol;  // 英文名
    uint8 public decimals; // 小数位
    uint256 public totalSupply;  // 代币数量

    // 状态变量
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    // 所有者
    address public  owner;

    // 事件 用于记录log indexed用于创建索引进行搜索
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approve(address indexed owner, address indexed spender, uint256 value);

    // 修饰符
    modifier onlyOwner(){
        require(msg.sender == owner, "Only owner can call this");
        _;
    }

    // 构造函数
    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _initialSupply
    ){
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _initialSupply * 10 ** _decimals;
        owner = msg.sender;
        balanceOf[msg.sender] = totalSupply;
        emit Transfer(address(0), msg.sender, totalSupply);
    }

    // 转账函数
    function transfer(address to, uint256 amount) public returns (bool) {
        require(to != address(0), "Can not transfer to zero address");
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");

        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;

        emit Transfer(msg.sender, to, amount);

        return true;
    }

    // 授权函数
    function approve(address spender, uint256 amount) public returns (bool){
        require(spender != address(0), "Can not approve zero address");

        allowance[msg.sender][spender] = amount;

        emit Approve(msg.sender, spender, amount);

        return true;
    }

    // 授权转账函数
    function transferFrom(address from, address to, uint256 amount) public returns (bool){
        require(from != address(0), "Can not transfer from zero address");
        require(to != address(0), "Can not transfer to zero address");
        require(balanceOf[from] >= amount, "Insufficient balance");
        require(allowance[from][msg.sender] >= amount, "Insufficient allowance");

        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        allowance[from][msg.sender] -= amount;

        emit Transfer(from, to, amount);

        return true;
    }

    // 铸造函数 （可选，需要onlyOwner修饰符）
    function mint(address to, uint256 amount) public onlyOwner {
        require(to != address(0), "Can not mint to zero address");

        totalSupply += amount;
        balanceOf[to] += amount;

        emit Transfer(address(0), to, amount);
    }

    // 销毁函数 （可选）
    function burn(uint256 amount) public {
        require(balanceOf[msg.sender] > amount, "Insufficient balance to burn");

        totalSupply -= amount;
        balanceOf[msg.sender] -= amount;

        emit Transfer(msg.sender, address(0), amount);
    }
}