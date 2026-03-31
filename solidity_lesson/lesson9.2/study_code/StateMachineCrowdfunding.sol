// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract SimpleCrowdfunding {
    enum State { Preparing, Funding, Success, Failed }
    
    State public state;
    address public owner;
    uint256 public goal;
    uint256 public raised;
    uint256 public deadline;
    
    mapping(address => uint256) public contributions;
    
    event StateChanged(State newState);
    event Contributed(address contributor, uint256 amount);
    
    constructor(uint256 _goal, uint256 _durationMinutes) {
        owner = msg.sender;
        goal = _goal;
        deadline = block.timestamp + (_durationMinutes * 1 minutes);
        state = State.Preparing;
    }
    
    modifier inState(State _state) {
        require(state == _state, "Wrong state");
        _;
    }
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    function startFunding() public onlyOwner inState(State.Preparing) {
        state = State.Funding;
        emit StateChanged(State.Funding);
    }
    
    function contribute() public payable inState(State.Funding) {
        require(block.timestamp < deadline, "Funding ended");
        require(msg.value > 0, "Must send ETH");
        
        contributions[msg.sender] += msg.value;
        raised += msg.value;
        emit Contributed(msg.sender, msg.value);
    }
    
    function finalize() public inState(State.Funding) {
        require(block.timestamp >= deadline, "Funding not ended");
        
        if (raised >= goal) {
            state = State.Success;
            emit StateChanged(State.Success);
        } else {
            state = State.Failed;
            emit StateChanged(State.Failed);
        }
    }
    
    function withdrawFunds() public onlyOwner inState(State.Success) {
        payable(owner).transfer(address(this).balance);
    }
    
    function refund() public inState(State.Failed) {
        uint256 amount = contributions[msg.sender];
        require(amount > 0, "No contribution");
        
        contributions[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
    }
}

