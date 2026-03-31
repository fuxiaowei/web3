// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/interfaces/IERC2981.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165.sol";

/**
 * @title MyNFTWithRoyalty
 * @dev 支持ERC2981版税标准的NFT合约
 * @notice 继承ERC2981接口，实现版税功能
 */
contract MyNFTWithRoyalty is ERC721, ERC721URIStorage, Ownable, IERC2981, ERC165 {
    uint256 private _tokenIdCounter;
    uint256 public constant MAX_SUPPLY = 10000;
    uint256 public mintPrice = 0.01 ether;
    
    // 版税接收地址
    address private _royaltyReceiver;
    
    // 版税比例（基点，10000 = 100%）
    uint96 private _royaltyBps = 1000; // 10%
    
    event NFTMinted(
        address indexed minter, 
        uint256 indexed tokenId, 
        string uri
    );
    
    /**
     * @dev 构造函数
     * @param royaltyReceiver 版税接收地址
     * @param royaltyBps 版税比例（基点）
     */
    constructor(
        address royaltyReceiver,
        uint96 royaltyBps
    ) ERC721("MyNFTWithRoyalty", "MNFR") Ownable(msg.sender) {
        require(royaltyReceiver != address(0), "Invalid royalty receiver");
        require(royaltyBps <= 1000, "Royalty too high"); // 最大10%
        
        _royaltyReceiver = royaltyReceiver;
        _royaltyBps = royaltyBps;
    }
    
    /**
     * @dev 铸造NFT
     */
    function mint(string memory uri) public payable returns (uint256) {
        require(_tokenIdCounter < MAX_SUPPLY, "Max supply reached");
        require(msg.value >= mintPrice, "Insufficient payment");
        
        _tokenIdCounter++;
        uint256 newTokenId = _tokenIdCounter;
        
        _safeMint(msg.sender, newTokenId);
        _setTokenURI(newTokenId, uri);
        
        emit NFTMinted(msg.sender, newTokenId, uri);
        
        return newTokenId;
    }
    
    /**
     * @dev 实现ERC2981标准：获取版税信息
     * @param tokenId Token ID
     * @param salePrice 售价
     * @return receiver 版税接收地址
     * @return royaltyAmount 版税金额
     */
    function royaltyInfo(
        uint256 /* tokenId */,
        uint256 salePrice
    ) external view override returns (
        address receiver,
        uint256 royaltyAmount
    ) {
        receiver = _royaltyReceiver;
        royaltyAmount = (salePrice * _royaltyBps) / 10000;
    }
    
    /**
     * @dev 设置版税信息
     * @param receiver 新的版税接收地址
     * @param bps 新的版税比例（基点）
     * @notice 只有合约所有者可以调用
     */
    function setRoyaltyInfo(address receiver, uint96 bps) external onlyOwner {
        require(receiver != address(0), "Invalid receiver");
        require(bps <= 1000, "Royalty too high");
        
        _royaltyReceiver = receiver;
        _royaltyBps = bps;
    }
    
    /**
     * @dev 查询版税接收地址
     */
    function royaltyReceiver() external view returns (address) {
        return _royaltyReceiver;
    }
    
    /**
     * @dev 查询版税比例
     */
    function royaltyBps() external view returns (uint96) {
        return _royaltyBps;
    }
    
    /**
     * @dev 重写tokenURI函数
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    /**
     * @dev 检查接口支持
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage, IERC165)
        returns (bool)
    {
        return 
            interfaceId == type(IERC2981).interfaceId ||
            super.supportsInterface(interfaceId);
    }
    
    /**
     * @dev 查询总供应量
     */
    function totalSupply() public view returns (uint256) {
        return _tokenIdCounter;
    }
    
    /**
     * @dev 提取铸造费用
     */
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance to withdraw");
        payable(owner()).transfer(balance);
    }
}

