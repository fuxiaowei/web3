import { expect } from "chai";
import { network } from "hardhat";

// 连接网络
const { ethers, networkHelpers } = await network.connect();

async function deployCounterFixture() {
  const [owner] = await ethers.getSigners();
  const counter = await ethers.deployContract("Counter");
  return { counter, owner };
}

describe("Error Handling", function () {
  describe("Require Revert", function () {
    it("Should revert with error message", async function () {
      const { counter } = await networkHelpers.loadFixture(deployCounterFixture);
      
      await expect(counter.incBy(0))
        .to.be.revertedWith("incBy: increment should be positive");
    });

    it("Should revert without reason", async function () {
      // 演示无原因回退
      // 假设有某个函数会无条件回退
      // await expect(contract.fail())
      //   .to.be.reverted;
    });
  });

  describe("Custom Error", function () {
    it("Should revert with custom error", async function () {
      // 演示自定义错误
      // 假设合约有自定义错误
      // error InsufficientBalance(uint256 required, uint256 available);
      
      // await expect(contract.withdraw(amount))
      //   .to.be.revertedWithCustomError(contract, "InsufficientBalance")
      //   .withArgs(required, available);
    });
  });

  describe("Panic Error", function () {
    it("Should revert with panic on overflow", async function () {
      const { counter } = await networkHelpers.loadFixture(deployCounterFixture);
      
      // 设置最大值
      const maxValue = 2n ** 256n - 1n;
      await counter.setNumber(maxValue);
      
      // 在Solidity 0.8+中，溢出会触发panic
      await expect(counter.inc())
        .to.be.revertedWithPanic(0x11); // 算术溢出
    });
  });
});

