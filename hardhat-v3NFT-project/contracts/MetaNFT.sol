// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

/// @title MetaNFT
/// @notice 基于 OpenZeppelin `ERC721` 的简单 NFT 合约：支持按指定 tokenId 铸造、按序号自增铸造，以及持有者销毁。
/// @dev 铸造函数未做权限控制（任何人可调用），若用于生产环境请按需增加 `onlyOwner` / 白名单等。
contract MetaNFT is ERC721 {
    /// @notice `mintNext` 使用的下一个 tokenId，从 1 起递增；`mint(address, id)` 不受此计数器约束。
    uint256 private _nextId = 1;

    /// @notice 部署时设置集合名称与符号（Name: MetaNFT, Symbol: MFT）。
    constructor() ERC721("MetaNFT", "MFT") {}

    /// @notice 向 `to` 铸造指定 `id` 的 NFT。
    /// @param to 接收方地址
    /// @param id 要铸造的 tokenId（若已存在则会因 ERC721 规则 revert）
    function mint(address to, uint256 id) external {
        _mint(to, id);
    }

    /// @notice 向 `to` 铸造当前计数器对应的 tokenId，并返回该 id，随后将 `_nextId` 加一。
    /// @param to 接收方地址
    /// @return 本次铸造使用的 tokenId
    function mintNext(address to) external returns (uint256) {
        _mint(to, _nextId);
        uint256 id = _nextId;
        _nextId++;
        return id;
    }

    /// @notice 销毁指定 tokenId；仅当前持有者（`ownerOf`）可调用。
    /// @param id 要销毁的 tokenId
    function burn(uint256 id) external {
        require(msg.sender == ownerOf(id), "not owner");
        _burn(id);
    }
}


