// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract CommitRevealVoting {
    // 投票选项
    enum Choice { None, OptionA, OptionB, OptionC }
    
    // 投票阶段
    enum Phase { Commit, Reveal, Ended }
    
    Phase public currentPhase;
    uint256 public commitDeadline;
    uint256 public revealDeadline;
    
    mapping(address => bytes32) public commits;
    mapping(address => bool) public hasRevealed;
    mapping(Choice => uint256) public votes;
    
    event Committed(address indexed voter);
    event Revealed(address indexed voter, Choice choice);
    event PhaseChanged(Phase newPhase);
    
    constructor(uint256 _commitDuration, uint256 _revealDuration) {
        currentPhase = Phase.Commit;
        commitDeadline = block.timestamp + _commitDuration;
        revealDeadline = commitDeadline + _revealDuration;
    }
    
    // 计算承诺哈希（链下调用）
    function getCommitHash(
        Choice choice, 
        bytes32 nonce
    ) external view returns (bytes32) {
        return keccak256(abi.encodePacked(choice, nonce, msg.sender));
    }
    
    // 阶段1：提交承诺
    function commit(bytes32 hash) external {
        require(currentPhase == Phase.Commit, "Not in commit phase");
        require(block.timestamp < commitDeadline, "Commit deadline passed");
        require(commits[msg.sender] == bytes32(0), "Already committed");
        
        commits[msg.sender] = hash;
        emit Committed(msg.sender);
    }
    
    // 切换到揭示阶段
    function startRevealPhase() external {
        require(currentPhase == Phase.Commit, "Not in commit phase");
        require(block.timestamp >= commitDeadline, "Commit phase not ended");
        
        currentPhase = Phase.Reveal;
        emit PhaseChanged(Phase.Reveal);
    }
    
    // 阶段2：揭示投票
    function reveal(Choice choice, bytes32 nonce) external {
        require(currentPhase == Phase.Reveal, "Not in reveal phase");
        require(block.timestamp < revealDeadline, "Reveal deadline passed");
        require(!hasRevealed[msg.sender], "Already revealed");
        require(commits[msg.sender] != bytes32(0), "No commit found");
        
        // 验证承诺
        bytes32 hash = keccak256(abi.encodePacked(choice, nonce, msg.sender));
        require(hash == commits[msg.sender], "Invalid reveal");
        
        hasRevealed[msg.sender] = true;
        votes[choice]++;
        
        emit Revealed(msg.sender, choice);
    }
    
    // 结束投票
    function endVoting() external {
        require(currentPhase == Phase.Reveal, "Not in reveal phase");
        require(block.timestamp >= revealDeadline, "Reveal phase not ended");
        
        currentPhase = Phase.Ended;
        emit PhaseChanged(Phase.Ended);
    }
    
    // 获取结果
    function getResults() external view returns (
        uint256 optionA,
        uint256 optionB,
        uint256 optionC
    ) {
        require(currentPhase == Phase.Ended, "Voting not ended");
        return (votes[Choice.OptionA], votes[Choice.OptionB], votes[Choice.OptionC]);
    }
}