// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract VotingSystem {
    
    struct Proposal {
        string description;
        uint voteCount;
        uint deadline;
        bool exists;
    }
    
    address public owner;
    uint public proposalCount;
    uint public constant MAX_PROPOSALS = 100;
    
    mapping(uint => Proposal) public proposals;
    mapping(uint => mapping(address => bool)) public hasVoted;
    
    event ProposalCreated(uint indexed proposalId, string description, uint deadline);
    event Voted(uint indexed proposalId, address indexed voter);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    // 创建提案
    function createProposal(string memory description, uint durationInDays) 
        public onlyOwner 
    {
        require(bytes(description).length > 0, "Description cannot be empty");
        require(durationInDays >= 1 && durationInDays <= 30, "Duration must be 1-30 days");
        require(proposalCount < MAX_PROPOSALS, "Too many proposals");
        
        uint proposalId = proposalCount;
        uint deadline = block.timestamp + (durationInDays * 1 days);
        
        proposals[proposalId] = Proposal({
            description: description,
            voteCount: 0,
            deadline: deadline,
            exists: true
        });
        
        proposalCount++;
        emit ProposalCreated(proposalId, description, deadline);
    }
    
    // 投票
    function vote(uint proposalId) public {
        require(proposals[proposalId].exists, "Proposal does not exist");
        require(block.timestamp <= proposals[proposalId].deadline, "Voting has ended");
        require(!hasVoted[proposalId][msg.sender], "Already voted");
        
        hasVoted[proposalId][msg.sender] = true;
        proposals[proposalId].voteCount++;
        
        emit Voted(proposalId, msg.sender);
    }
    
    // 获取提案详情
    function getProposal(uint proposalId) 
        public view 
        returns (string memory description, uint voteCount, uint deadline, bool isActive) 
    {
        require(proposals[proposalId].exists, "Proposal does not exist");
        Proposal memory p = proposals[proposalId];
        return (
            p.description,
            p.voteCount,
            p.deadline,
            block.timestamp <= p.deadline
        );
    }
    
    // 获取获胜提案
    function getWinner() public view returns (uint winningProposalId, uint winningVoteCount) {
        require(proposalCount > 0, "No proposals");
        
        uint maxVotes = 0;
        uint winnerId = 0;
        
        for (uint i = 0; i < proposalCount; i++) {
            if (proposals[i].voteCount > maxVotes) {
                maxVotes = proposals[i].voteCount;
                winnerId = i;
            }
        }
        
        return (winnerId, maxVotes);
    }
    
    // 检查是否已投票
    function checkIfVoted(uint proposalId, address voter) public view returns (bool) {
        return hasVoted[proposalId][voter];
    }
}