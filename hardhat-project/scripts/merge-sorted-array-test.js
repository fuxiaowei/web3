import hre from "hardhat";
import assert from "node:assert/strict";

async function main() {
    console.log("合并两个有序数组 本地测试开始…");

    const conn = await hre.network.connect();
    const ethers = conn.ethers;
    if (!ethers) throw new Error("conn.ethers is undefined");
    console.log("部署者地址:", ethers.address);

    // 部署合约
    const Merge = await ethers.getContractFactory("MergeSortedArray");
    const contract = await Merge.deploy();
    await contract.waitForDeployment();

    // 测试用例
    const nums1 = [1, 3, 5];
    const nums2 = [2, 4, 6];

    const merged = await contract.merge(nums1, nums2);
    console.log(`数组1: ${nums1}`);
    console.log(`数组2: ${nums2}`);
    console.log(`合并后: ${merged}`);
}

main()
    .catch((err) => {
        console.error(err);
        process.exit(1);
    });