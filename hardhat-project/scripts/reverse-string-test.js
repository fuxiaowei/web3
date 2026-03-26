import hre from "hardhat";
import assert from "node:assert/strict";

async function main() {
    const conn = await hre.network.connect();
    const ethers = conn.ethers;
    if (!ethers) throw new Error("conn.ethers is undefined");

    try {
        console.log("字符串反转脚本测试开始…");
        const StringReverse = await ethers.getContractFactory("StringReverse");
        const stringReverse = await StringReverse.deploy();
        await stringReverse.waitForDeployment();
        console.log("部署者地址:", await stringReverse.getAddress());

        // 脚本里手写断言（与 test/ReverseString.test.js 二选一或并存）
        assert.equal(await stringReverse.reverse(""), "");
        assert.equal(await stringReverse.reverse("a"), "a");
        assert.equal(await stringReverse.reverse("abcde"), "edcba");
        var reverse = await stringReverse.reverse("abcde");
        console.log("字符串反转结果:", reverse);

        console.log("脚本断言全部通过");
    } finally {
        await conn.close();
    }
}

main().catch((e) => {
    console.error(e); // 捕获异常并打印错误堆栈
    process.exitCode = 1; // 非 0 退出码表示失败
});

