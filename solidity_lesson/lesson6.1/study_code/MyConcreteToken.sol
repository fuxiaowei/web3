// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 接口：定义了代币合约必须实现的功能
interface IBasicToken {
    function getName() external view returns (string memory);
    function getSymbol() external view returns (string memory);
    function getTotalSupply() external view returns (uint256);
    function getBalance(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
}

// 抽象合约：实现了部分接口功能，但仍有抽象函数
abstract contract AbstractToken is IBasicToken {
    string public override name;
    string public override symbol;
    uint256 public override totalSupply;
    mapping(address => uint256) private _balances;

    constructor(string memory _name, string memory _symbol, uint256 _initialSupply) {
        name = _name;
        symbol = _symbol;
        totalSupply = _initialSupply;
        _balances[msg.sender] = _initialSupply; // 将初始供应量分配给部署者
    }

    // 实现了接口中的 getBalance 函数
    function getBalance(address account) public view override returns (uint256) {
        return _balances[account];
    }

    // 抽象函数：transfer，必须由继承者实现
    function transfer(address recipient, uint256 amount) public virtual override returns (bool);

    // 辅助函数，供子合约使用
    function _transfer(address sender, address recipient, uint256 amount) internal virtual {
        require(sender != address(0), "Transfer from zero address");
        require(recipient != address(0), "Transfer to zero address");
        require(_balances[sender] >= amount, "Insufficient balance");

        _balances[sender] -= amount;
        _balances[recipient] += amount;
    }

    // 实现了接口中的 getName 和 getSymbol
    function getName() public view override returns (string memory) {
        return name;
    }

    function getSymbol() public view override returns (string memory) {
        return symbol;
    }

    function getTotalSupply() public view override returns (uint256) {
        return totalSupply;
    }
}

// 具体合约：继承抽象合约并实现所有抽象函数
contract MyConcreteToken is AbstractToken("MyToken", "MTK", 1000000) {
    constructor() {} // 继承了父合约的构造函数，并传递了参数

    // 实现抽象合约中的 transfer 函数
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        _transfer(msg.sender, recipient, amount); // 调用抽象合约中的辅助函数
        return true;
    }
}