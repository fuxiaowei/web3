import hre from "hardhat";
import assert from "node:assert/strict";

async function main() {
    console.log("整数转罗马数字 本地测试开始…");
    const conn = await hre.network.connect();
    const ethers = conn.ethers;
    if (!ethers) throw new Error("conn.ethers is undefined");

    try {
        const IntegerToRoman = await ethers.getContractFactory("IntegerToRomanNumeral");
        const contract = await IntegerToRoman.deploy();
        await contract.waitForDeployment();


        // 测试用例
        const testCases = [3749, 58, 1994, 1, 4, 9, 58, 400, 900];

        for (const num of testCases) {
            const result = awai

            t contract.intToRoman(num);
            console.log(`数字 ${num} → 罗马数字: ${result}`);
        }
    } finally {
        console.log("整数转罗马数字 本地测试结束");
        await conn.close();
    }
}

main()
    .catch((err) => {
        console.error(err);
        process.exit(1);
    });