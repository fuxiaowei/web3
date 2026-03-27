// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

/// @title BeggingContract（讨饭 / 捐赠合约）
/// @notice 接收用户自愿转入的 ETH，按地址累计金额；仅部署者可提空合约余额。
/// @dev
/// - `donations` 与 `allDonors` 配合：`allDonors` 仅在某地址**首次**产生正额捐赠时 `push` 一次（见 `donate` 中 `donations[msg.sender] == 0` 判断）。
/// - 因 `donate` 要求 `msg.value > 0`，「从未捐赠」与「累计额为 0」在未捐赠时一致；若将来允许 0 元调用，需改用单独 mapping 标记是否已登记捐赠者。
/// - `withdraw` 使用 `transfer`：向 EOA 通常足够；若 `owner` 为会耗气的合约地址，可能因 2300 gas 失败，生产环境可改为 `call` 并处理返回值。
contract BeggingContract {
    /// @dev 合约部署者，唯一有权调用 `withdraw` 的地址（构造函数中设为 `msg.sender`）
    address public owner;

    /// @dev 各地址向本合约累计转入的 wei（多次 `donate` 会累加）
    mapping(address => uint256) public donations;

    /// @dev 所有至少捐赠过一次（且当时 `donations` 由 0 变为正）的地址列表，供 `getTop3Donors` 遍历
    address[] public allDonors;

    /// @notice 每次成功捐赠时触发，链下可索引 `donor` 做统计或前端展示
    /// @param donor 捐赠者地址（`msg.sender`）
    /// @param amount 本次捐赠的 wei（`msg.value`）
    //indexed = 给事件参数 建索引 / 加数据库索引
    //作用只有两个：
    //  可以被后端精准过滤、查询（web3j /ethers 必备）
    //  链上会对这个参数建立索引，检索更快
    event Donation(address indexed donor, uint256 amount);

    /// @notice 所有者成功提款后触发，`amount` 为本次转出的合约余额（通常等于提款前全余额）
    event Withdraw(uint256 amount);

    /// @notice 部署合约并将 `owner` 设为部署账户
    constructor() {
        owner = msg.sender;
    }

    /// @dev 将 `msg.sender` 限制为 `owner`，否则 revert
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }

    /// @notice 向合约捐赠 ETH，并更新累计金额与（首次捐赠时）捐赠者列表
    /// @dev
    /// - 必须通过交易附带 `msg.value > 0` 的 wei。
    /// - 首次捐赠：`donations[sender]` 仍为 0 → 将 `sender` 加入 `allDonors`；再执行 `+= msg.value`。
    /// - 非首次：仅累加 `donations`，不重复 `push`。
    function donate() external payable {
        require(msg.value > 0, "Donation must > 0");

        if (donations[msg.sender] == 0) {
            allDonors.push(msg.sender);
        }

        donations[msg.sender] += msg.value;

        emit Donation(msg.sender, msg.value);
    }

    /// @notice 将所有合约余额一次性转给 `owner`（仅所有者可调）
    /// @dev 先读 `address(this).balance`，再 `transfer` 全额；若余额为 0 则 revert
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance");

        payable(owner).transfer(balance);
        emit Withdraw(balance);
    }

    /// @notice 查询任意地址在本合约中的累计捐赠额（wei）
    /// @param _account 要查询的地址
    /// @return 该地址 `donations[_account]`，未捐赠过则为 0
    function getDonation(address _account) external view returns (uint256) {
        return donations[_account];
    }

    /// @notice 返回当前统计下捐赠总额最高的 3 个地址及其金额（降序：下标 0 最大）
    /// @dev
    /// - 遍历 `allDonors`，对每个 `amount` 尝试插入「有序 TOP3」数组：
    ///   自 `j = 0` 到 `2`，若 `amount > topAmounts[j]`，则应在第 `j` 槽插入；
    ///   先将 `j..2` 原有元素整体右移一位（从右端 `k=2` 向左挪到 `k=j`），再写入新地址与金额，并 `break`。
    /// - 若捐赠者少于 3 人或金额不足，对应槽保持 `address(0)` 与 `0`。
    /// - Gas 随捐赠者人数线性增长；捐赠者很多时可改为链下索引或 Merkle 证明等方案。
    /// @return topAddresses 前三名地址（可能含 `address(0)`）
    /// @return topAmounts 与 `topAddresses` 对齐的累计捐赠额
    function getTop3Donors() external view returns (address[3] memory, uint256[3] memory) {
        address[3] memory topAddresses;
        uint256[3] memory topAmounts;

        for (uint i = 0; i < allDonors.length; i++) {
            address addr = allDonors[i];
            uint256 amount = donations[addr];

            for (uint j = 0; j < 3; j++) {
                if (amount > topAmounts[j]) {
                    // 将 j 号位及右侧元素各向右挤一格，腾出 j 号位给当前 (addr, amount)
                    for (uint k = 2; k > j; k--) {
                        topAmounts[k] = topAmounts[k - 1];
                        topAddresses[k] = topAddresses[k - 1];
                    }
                    topAmounts[j] = amount;
                    topAddresses[j] = addr;
                    break;
                }
            }
        }
        return (topAddresses, topAmounts);
    }

    /// @notice 当前合约持有的 ETH 总余额（wei），等于所有捐赠之和减去已 `withdraw` 转出的部分
    function getTotalBalance() external view returns (uint256) {
        return address(this).balance;
    }
}

/*
 * =============================================================================
 * 作业 3：讨饭合约（BeggingContract）
 * =============================================================================
 *
 * 【任务目标】
 *   - 允许用户向合约发送以太币（捐赠）。
 *   - 记录每位捐赠者的地址与累计捐赠金额。
 *   - 仅合约所有者可一次性提走合约内全部余额。
 *
 * 【合约须实现】
 *   - mapping：address → 累计捐赠金额（wei）。
 *   - donate()：payable，接收转账并更新记录（题目要求可用 transfer 等实现资金流）。
 *   - withdraw()：仅 owner 可调，将合约余额转给所有者。
 *   - getDonation(address)：view，查询某地址累计捐赠额。
 *   - 建议配合 onlyOwner 修饰符限制提款。
 *
 * 【推荐流程（作业说明）】
 *   编写：在 Remix 或本地环境编译 BeggingContract。
 *   部署：部署到测试网（如原题中的 Goerli / Sepolia；现多使用 Sepolia 等）。
 *   测试：
 *     - 用钱包向合约 donate；
 *     - owner 调用 withdraw；
 *     - 调用 getDonation 核对金额。
 *
 * 【提交内容（按课程要求）】
 *   - Solidity 源文件（如 BeggingContract.sol）。
 *   - 测试网上的合约地址。
 *   - Remix / Etherscan 等测试过程截图（若教师要求）。
 *
 * 【任务要求摘要】
 *   - 用 mapping 记录捐赠金额；donate / 资金流使用 payable 与转账 API。
 *   - withdraw 仅所有者可调用；donate、withdraw、getDonation 行为可测、可验证。
 *
 * 【额外挑战（可选）】
 *   - 捐赠事件：每次捐赠触发事件（记录捐赠者与金额）。
 *   - 排行榜：展示捐赠额最高的前 3 个地址。
 *   - 时间窗：仅在指定时间段内允许捐赠。
 * =============================================================================
 */

// 部署测试网成功 =============================================================================
//PS D:\wisest\web3\web3\hardhat-project> npx hardhat run scripts/deploy-begging-contract-sepolia.js --network sepolia
//
//部署账户: 0x4E2b10e8E740B73673cCFA1F4f651a8F74f33B7E
//BeggingContract 已部署到 Sepolia，合约地址: 0x0899604061c6f029DfD415728507b3a42eCdB612
//浏览器: https://sepolia.etherscan.io/address/0x0899604061c6f029DfD415728507b3a42eCdB612