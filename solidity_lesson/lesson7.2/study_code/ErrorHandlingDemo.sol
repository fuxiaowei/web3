// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract ErrorHandlingDemo {
    uint256 public balance = 100;
    uint256 public totalSupply = 1000;

    // 自定义错误
    error InsufficientBalance(uint256 available, uint256 required);
    error Unauthorized();
    error InvalidAmount();

    // require示例 - 输入验证
    function withdrawRequire(uint256 amount) public {
        require(amount > 0, "金额必须大于0");
        require(balance >= amount, "余额不足");
        balance -= amount;
    }

    // require + 自定义错误
    function withdrawCustomError(uint256 amount) public {
        if (amount == 0) revert InvalidAmount();
        if (balance < amount) revert InsufficientBalance(balance, amount);
        balance -= amount;
    }

    // assert示例 - 不变量检查
    function checkInvariant() public view {
        assert(balance <= totalSupply);
    }

    // revert示例 - 自定义错误条件
    function transferRevert(address to, uint256 amount) public {
        if (to == address(0)) revert InvalidRecipient();
        if (balance < amount) revert InsufficientBalance(balance, amount);
        balance -= amount;
    }

    // 触发assert失败（用于演示）
    function triggerAssertFailure() public {
        totalSupply = 50; // 故意设置不满足不变量
        assert(balance <= totalSupply); // 这会失败
    }

    error InvalidRecipient();
}