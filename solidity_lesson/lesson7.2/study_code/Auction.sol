// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

contract Auction {
    // 自定义错误
    error BidTooLow(uint256 currentBid, uint256 newBid);
    error AuctionEnded();
    error AuctionNotEnded();
    error BidAlreadyPlaced();
    error WithdrawalFailed();
    error InvalidToken();

    IERC20 public token;
    address public owner;
    uint256 public auctionEnd;
    uint256 public highestBid;
    address public highestBidder;

    mapping(address => uint256) public bids;
    mapping(address => bool) public hasBid;

    event NewBid(address indexed bidder, uint256 amount);
    event BidRefunded(address indexed bidder, uint256 amount);
    event AuctionWon(address indexed winner, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }

    modifier auctionActive() {
        require(block.timestamp < auctionEnd, "Auction ended");
        _;
    }

    constructor(address _token, uint256 _duration) {
        token = IERC20(_token);
        owner = msg.sender;
        auctionEnd = block.timestamp + _duration;
        require(_token != address(0), "Invalid token address");
    }

    // 出价函数 - 使用require和自定义错误
    function bid(uint256 amount) public auctionActive {
        // 输入验证
        require(amount > 0, "Bid amount must be greater than 0");

        // 业务规则验证
        if (amount <= highestBid) {
            revert BidTooLow(highestBid, amount);
        }

        // 不变量检查
        assert(auctionEnd > block.timestamp);

        // 处理之前的出价
        if (hasBid[msg.sender]) {
            uint256 previousBid = bids[msg.sender];
            bids[msg.sender] = 0;
            _refund(msg.sender, previousBid);
        }

        // 转移新出价
        require(
            token.transferFrom(msg.sender, address(this), amount),
            "Token transfer failed"
        );

        // 更新状态
        bids[msg.sender] = amount;
        hasBid[msg.sender] = true;
        highestBid = amount;
        highestBidder = msg.sender;

        emit NewBid(msg.sender, amount);
    }

    // 退款函数 - 使用try-catch
    function _refund(address to, uint256 amount) private {
        if (amount == 0) return;

        try token.transfer(to, amount) returns (bool success) {
            if (!success) {
                revert WithdrawalFailed();
            }
            emit BidRefunded(to, amount);
        } catch Error(string memory) {
            revert WithdrawalFailed();
        } catch {
            revert WithdrawalFailed();
        }
    }

    // 结束拍卖
    function endAuction() public onlyOwner {
        require(block.timestamp >= auctionEnd, "Auction not ended");

        if (highestBidder != address(0)) {
            emit AuctionWon(highestBidder, highestBid);
        }
    }

    // 提取出价
    function withdraw() public {
        require(block.timestamp >= auctionEnd, "Auction not ended");
        require(hasBid[msg.sender], "No bid to withdraw");
        require(msg.sender != highestBidder, "Winner cannot withdraw");

        uint256 amount = bids[msg.sender];
        bids[msg.sender] = 0;
        hasBid[msg.sender] = false;

        _refund(msg.sender, amount);
    }

    // 获取当前最高出价
    function getHighestBid() public view returns (uint256, address) {
        return (highestBid, highestBidder);
    }
}