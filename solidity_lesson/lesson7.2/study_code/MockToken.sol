// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// 模拟外部ERC20代币合约
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

// 一个可能失败的代币合约
contract MockToken {
    mapping(address => uint256) public balanceOf;
    bool public paused = false;

    constructor() {
        balanceOf[msg.sender] = 1000;
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        require(!paused, "Token transfer paused");
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        return true;
    }

    function pause() external {
        paused = true;
    }
}

// 使用try-catch的合约
contract TokenHandler {
    IERC20 public token;

    event TransferSuccess(address indexed to, uint256 amount);
    event TransferFailed(address indexed to, uint256 amount, string reason);
    event UnknownError(address indexed to, uint256 amount);

    constructor(address _token) {
        token = IERC20(_token);
    }

    // 使用try-catch处理代币转账
    function safeTransfer(address to, uint256 amount) public {
        try token.transfer(to, amount) returns (bool success) {
            if (success) {
                emit TransferSuccess(to, amount);
            } else {
                emit TransferFailed(to, amount, "Transfer returned false");
            }
        } catch Error(string memory reason) {
            // 捕获字符串错误
            emit TransferFailed(to, amount, reason);
        } catch (bytes memory lowLevelData) {
            // 捕获未知错误
            emit UnknownError(to, amount);
        }
    }

    // 批量安全转账
    function batchSafeTransfer(
        address[] memory recipients,
        uint256[] memory amounts
    ) public {
        require(recipients.length == amounts.length, "Array length mismatch");

        for (uint256 i = 0; i < recipients.length; i++) {
            try token.transfer(recipients[i], amounts[i]) returns (bool success) {
                if (success) {
                    emit TransferSuccess(recipients[i], amounts[i]);
                } else {
                    emit TransferFailed(recipients[i], amounts[i], "Transfer returned false");
                }
            } catch Error(string memory reason) {
                emit TransferFailed(recipients[i], amounts[i], reason);
                // 继续处理下一个，不中断整个批次
            } catch (bytes memory) {
                emit UnknownError(recipients[i], amounts[i]);
                // 继续处理下一个
            }
        }
    }

    // 检查余额并转账
    function transferIfSufficient(address to, uint256 amount) public {
        uint256 balance = token.balanceOf(address(this));

        if (balance < amount) {
            revert("Insufficient contract balance");
        }

        try token.transfer(to, amount) {
            emit TransferSuccess(to, amount);
        } catch {
            revert("Transfer failed");
        }
    }
}