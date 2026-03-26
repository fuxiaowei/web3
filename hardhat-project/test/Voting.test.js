import { ethers } from "hardhat";
import assert from "node:assert/strict";
import { log } from "node:console";

describe("Voting", function () {
  it("vote & getVotes should work", async function () {
    log("vote & getVotes should work");
    const Voting = await ethers.getContractFactory("Voting");
    const voting = await Voting.deploy();
    await voting.waitForDeployment();

    assert.equal(await voting.getVotes(1n), 0n);
    assert.equal(await voting.getVotes(2n), 0n);

    await voting.vote(1n);
    await voting.vote(1n);
    await voting.vote(2n);

    assert.equal(await voting.getVotes(1n), 2n);
    assert.equal(await voting.getVotes(2n), 1n);
  });

  it("resetVotes should reset visible votes (epoch reset)", async function () {
    const Voting = await ethers.getContractFactory("Voting");
    const voting = await Voting.deploy();
    await voting.waitForDeployment();

    await voting.vote(1n);
    await voting.vote(2n);

    assert.equal(await voting.getVotes(1n), 1n);
    assert.equal(await voting.getVotes(2n), 1n);

    await voting.resetVotes();

    // 当前 epoch 的 getVotes 视为 0
    assert.equal(await voting.getVotes(1n), 0n);
    assert.equal(await voting.getVotes(2n), 0n);

    // 新 epoch 再投票：同一个 candidateId 会懒重置后重新计数
    await voting.vote(1n);
    assert.equal(await voting.getVotes(1n), 1n);
  });
});

