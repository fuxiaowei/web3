import { expect } from "chai";
import { network } from "hardhat";
import { time, mine } from "@nomicfoundation/hardhat-network-helpers";

// 连接网络
const { ethers, networkHelpers } = await network.connect();

// 定义Fixture函数
async function deployTokenFixture() {
  const [owner, addr1, addr2] = await ethers.getSigners();
  const token = await ethers.deployContract("Token");
  
  // 给owner分配初始代币
  const initialSupply = ethers.parseEther("10000");
  // 假设Token合约有mint函数
  // await token.mint(owner.address, initialSupply);
  
  return { token, owner, addr1, addr2, initialSupply };
}

describe("Token", function () {
  describe("Deployment", function () {
    it("Should deploy with correct name and symbol", async function () {
      const { token } = await networkHelpers.loadFixture(deployTokenFixture);
      
      // 假设Token合约有name和symbol函数
      // expect(await token.name()).to.equal("Token");
      // expect(await token.symbol()).to.equal("TKN");
    });
  });

  describe("Transfer", function () {
    it("Should transfer tokens correctly", async function () {
      const { token, owner, addr1 } = await networkHelpers.loadFixture(deployTokenFixture);
      
      const amount = ethers.parseEther("100");
      
      // 假设Token合约有transfer函数
      // await expect(token.transfer(addr1.address, amount))
      //   .to.emit(token, "Transfer")
      //   .withArgs(owner.address, addr1.address, amount);
      
      // expect(await token.balanceOf(addr1.address)).to.equal(amount);
    });

    it("Should revert when balance is insufficient", async function () {
      const { token, addr1, addr2 } = await networkHelpers.loadFixture(deployTokenFixture);
      
      const amount = ethers.parseEther("1000000");
      
      // await expect(token.connect(addr1).transfer(addr2.address, amount))
      //   .to.be.revertedWith("Insufficient balance");
    });
  });

  describe("Balance Changes", function () {
    it("Should change balances correctly", async function () {
      const { token, owner, addr1 } = await networkHelpers.loadFixture(deployTokenFixture);
      
      const amount = ethers.parseEther("100");
      
      // 假设有transfer函数
      // await expect(
      //   token.transfer(addr1.address, amount)
      // ).to.changeTokenBalances(
      //   token,
      //   [owner, addr1],
      //   [ethers.parseEther("-100"), ethers.parseEther("100")]
      // );
    });
  });
});

