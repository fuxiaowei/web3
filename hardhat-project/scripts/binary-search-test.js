import hre from "hardhat";
import assert from "node:assert/strict";

async function main() {
    console.log("二分查找 本地测试开始…");

    const conn = await hre.network.connect();
    const ethers = conn.ethers;
    if (!ethers) throw new Error("conn.ethers is undefined");
    console.log("部署者地址:", ethers.address);

    try {


        // 部署合约
        const BinarySearch = await ethers.getContractFactory("BinarySearch");
        const contract = await BinarySearch.deploy();
        await contract.waitForDeployment();

        // 有序数组
        const arr = [-5, 0, 2, 7, 9, 11, 15];

        // 测试用例
        const testTargets = [2, 7, -5, 15, 100];

        for (const target of testTargets) {
            const index = await contract.search(arr, target);
            if (index.toString() === ethers.MaxUint256.toString()) {
                console.log(`目标值 ${target} → 未找到`);
            } else {
                console.log(`目标值 ${target} → 找到，索引为: ${index.toString()}`);
            }
        }
    } finally {
        console.log("二分查找 本地测试结束。");
        await conn.close();
    }
}

main().catch((err) => {
    console.error(err);
    process.exit(1);
});