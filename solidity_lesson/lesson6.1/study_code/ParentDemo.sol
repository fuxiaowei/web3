// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 父合约
contract ParentDemo {
    uint public parentValue;

    // 父合约构造函数，接收一个初始值
    constructor(uint _initialValue) {
        parentValue = _initialValue;
    }

    // 一个可被子合约重写的函数
    function getParentValue() public view virtual returns (uint) {
        return parentValue;
    }

    // 一个可被子合约重写的函数
    function setParentValue(uint _newValue) public virtual {
        parentValue = _newValue;
    }
}

// 子合约继承自 Parent，并向父合约构造函数传递参数
contract Child is ParentDemo(100) {
    uint public childValue;

    // 子合约构造函数
    constructor(uint _childInitialValue) {
        // 父合约构造函数 Parent(100) 会先执行
        childValue = _childInitialValue;
    }

    // 重写父合约的 getParentValue 函数
    function getParentValue() public view override returns (uint) {
        // 使用 super 关键字调用父合约的原始函数
        return super.getParentValue() + childValue;
    }

    // 重写父合约的 setParentValue 函数
    function setParentValue(uint _newValue) public override {
        // 在调用父合约函数前或后添加子合约的逻辑
        super.setParentValue(_newValue * 2); // 将新值乘以2后传递给父合约
        childValue = _newValue; // 子合约自己的值设置为新值
    }

    // 子合约特有的函数
    function getChildValue() public view returns (uint) {
        return childValue;
    }
}