// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Crowdfunding {
    
    enum State { Preparing, Active, Success, Failed }
    
    State public currentState;
    address public owner;
    uint public goal;
    uint public deadline;
    uint public totalFunded;
    
    mapping(address => uint) public contributions;
    
    event StateChanged(State newState);
    event Contributed(address indexed contributor, uint amount);
    event Refunded(address indexed contributor, uint amount);
    
    modifier inState(State expectedState) {
        require(currentState == expectedState, "Invalid state");
        _;
    }
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }
    
    constructor() {
        owner = msg.sender;
        currentState = State.Preparing;
    }
    
    // 开始众筹
    function start(uint _goal, uint durationInDays) 
        public 
        onlyOwner 
        inState(State.Preparing) 
    {
        require(_goal > 0, "Goal must be positive");
        require(durationInDays > 0 && durationInDays <= 90, "Invalid duration");
        
        goal = _goal;
        deadline = block.timestamp + (durationInDays * 1 days);
        currentState = State.Active;
        
        emit StateChanged(State.Active);
    }
    
    // 贡献资金
    function contribute() public payable inState(State.Active) {
        require(block.timestamp <= deadline, "Campaign ended");
        require(msg.value > 0, "Must send ETH");
        
        contributions[msg.sender] += msg.value;
        totalFunded += msg.value;
        
        emit Contributed(msg.sender, msg.value);
    }
    
    // 结束众筹
    function finalize() public {
        require(currentState == State.Active, "Not active");
        require(block.timestamp > deadline, "Campaign not ended");
        
        if (totalFunded >= goal) {
            currentState = State.Success;
            // 转账给owner
            payable(owner).transfer(totalFunded);
        } else {
            currentState = State.Failed;
        }
        
        emit StateChanged(currentState);
    }
    
    // 退款（众筹失败时）
    function refund() public inState(State.Failed) {
        uint amount = contributions[msg.sender];
        require(amount > 0, "No contribution");
        
        contributions[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
        
        emit Refunded(msg.sender, amount);
    }
    
    // 查询合约状态
    function getStatus() public view returns (
        State state,
        uint currentFunding,
        uint targetGoal,
        uint timeLeft
    ) {
        uint remaining = 0;
        if (block.timestamp < deadline) {
            remaining = deadline - block.timestamp;
        }
        
        return (currentState, totalFunded, goal, remaining);
    }
}