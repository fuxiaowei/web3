// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract BoxV1 is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    // 状态变量
    uint256 public number;

    // 构造函数：禁止初始化，防止被恶意调用
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // 初始化函数（替代构造函数）
    function initialize(uint256 _initNum) external initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
        number = _initNum;
    }

    // 业务方法
    function setNumber(uint256 _num) external virtual {
        number = _num;
    }

    // UUPS 必须：授权升级权限，仅管理员
    function _authorizeUpgrade(address) internal override onlyOwner {}
}