import hre from "hardhat"; // 引入 Hardhat 运行时（用于获取网络、编译产物等能力）

async function main() {
  const conn = await hre.network.connect(); // 连接到 Hardhat 内置的开发链/网络（本地模拟）
  const ethers = conn.ethers; // 从连接对象上拿到 ethers 相关能力（signer/contract factory 等）
  if (!ethers) throw new Error("conn.ethers is undefined"); // 防御性检查：如果插件没生效就直接报错

  try {
    console.log("vote & getVotes should work 执行开始......."); // 脚本开始提示
    const [deployer] = await ethers.getSigners(); // 拿到默认账户列表，并取第一个作为部署者/发送者
    console.log("Deployer:", deployer.address); // 打印部署者地址，方便你确认调用的账户是谁

    const Voting = await ethers.getContractFactory("Voting"); // 获取合约工厂（用于 deploy/创建合约实例）
    const voting = await Voting.deploy(); // 部署合约 Voting（构造函数无参数）
    await voting.waitForDeployment(); // 等待部署交易被打包并完成
    console.log("Voting deployed to:", await voting.getAddress()); // 打印部署后的合约地址

    const c1 = 1n; // 候选人 ID = 1（BigInt 类型）
    const c2 = 2n; // 候选人 ID = 2（BigInt 类型）

    console.log("Initial votes:", {
      c1: await voting.getVotes(c1), // 查询候选人 1 当前得票
      c2: await voting.getVotes(c2), // 查询候选人 2 当前得票
    });

    await (await voting.vote(c1)).wait(); // 对候选人 1 投票一次（wait 等交易确认）
    await (await voting.vote(c1)).wait(); // 对候选人 1 再投一次
    await (await voting.vote(c2)).wait(); // 对候选人 2 投票一次

    console.log("After voting:", {
      c1: await voting.getVotes(c1), // 再次查询候选人 1 得票
      c2: await voting.getVotes(c2), // 再次查询候选人 2 得票
    });

    await (await voting.resetVotes()).wait(); // 重置投票（你当前实现是 epoch/轮次重置，O(1)）
    console.log("After resetVotes:", {
      c1: await voting.getVotes(c1), // 重置后查询候选人 1 得票是否为 0
      c2: await voting.getVotes(c2), // 重置后查询候选人 2 得票是否为 0
    });

    await (await voting.vote(c1)).wait(); // 重置后再给候选人 1 投一次，验证“新轮次”能重新累计
    console.log("After voting again:", {
      c1: await voting.getVotes(c1), // 新轮次下候选人 1 得票
      c2: await voting.getVotes(c2), // 新轮次下候选人 2 得票应仍为 0
    });
  } finally {
    await conn.close(); // 关闭网络连接/资源（避免进程无法退出）
  }
}

main().catch((e) => {
  console.error(e); // 捕获异常并打印错误堆栈
  process.exitCode = 1; // 非 0 退出码表示失败
});

