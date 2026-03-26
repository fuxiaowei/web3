import { ethers } from "hardhat";
import assert from "node:assert/strict";

describe("StringReverse", function () {
  let stringReverse;

  beforeEach(async function () {
    const Factory = await ethers.getContractFactory("StringReverse");
    stringReverse = await Factory.deploy();
    await stringReverse.waitForDeployment();
  });

  it("空串与单字符原样返回", async function () {
    assert.equal(await stringReverse.reverse(""), "");
    assert.equal(await stringReverse.reverse("a"), "a");
  });

  it("普通 ASCII 反转", async function () {
    assert.equal(await stringReverse.reverse("abcde"), "edcba");
    assert.equal(await stringReverse.reverse("hello"), "olleh");
  });

  it("回文不变", async function () {
    assert.equal(await stringReverse.reverse("aba"), "aba");
  });
});
