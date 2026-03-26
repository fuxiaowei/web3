import hre from "hardhat";
import assert from "node:assert/strict";

async function main() {
    const conn = await hre.network.connect();
    const ethers = conn.ethers;
    if (!ethers) throw new Error("conn.ethers is undefined");

    try {
        console.log("罗马数字转阿拉伯数字本地测试开始…");
        const ConvertRomanNumeralsToIntegers = await ethers.getContractFactory("ConvertRomanNumeralsToIntegers");
        const romanToNumber = await ConvertRomanNumeralsToIntegers.deploy();
        await romanToNumber.waitForDeployment();
        console.log("部署者地址:", await romanToNumber.getAddress());


        const result = await romanToNumber.romanToInt.staticCall("III");
        console.log("罗马数字 III 转整数结果:", result.toString());
        
        // 测试更多用例验证合约功能
        const testCases = ["IV", "IX", "LVIII", "MCMXCIV"];
        for (const roman of testCases) {
            const value = await romanToNumber.romanToInt.staticCall(roman);
            console.log(`罗马数字 ${roman} 转整数结果:`, value.toString());
        }
    } finally {
        await conn.close();
    }
}

main().catch((e) => {
    console.error(e); // 捕获异常并打印错误堆栈
    process.exitCode = 1; // 非 0 退出码表示失败
});

