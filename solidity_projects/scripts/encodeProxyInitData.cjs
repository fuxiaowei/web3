/**
 * 为 MyTransparentUpgradeableProxy 构造器里的 _data 生成 calldata（仅编码，不广播交易）。
 * 对 MultiSigWalletUpgradeable.initialize(address[], uint256) 做 ABI 编码，供 Remix / 部署使用。
 *
 * 用法:
 *   node scripts/encodeProxyInitData.cjs <owner1,owner2,...> <numConfirmationsRequired>
 *
 * 例（门限 2、两个所有者）:
 *   node scripts/encodeProxyInitData.cjs 0xF5aBB23811113E2d22813DF10f1ac514A82EB031,0x1111111111111111111111111111111111111111 2
 *
 * 若部署代理时不立刻初始化，改用空 bytes:
 *   node scripts/encodeProxyInitData.cjs --empty
 */
const { Interface, getAddress } = require("ethers");

const FRAGMENT =
  "function initialize(address[] _owners, uint256 _numConfirmationsRequired)";

function main() {
  const args = process.argv.slice(2);

  if (args.length === 0) {
    printHelp();
    process.stderr.write("错误: 请传入参数，或见 --help\n");
    process.exit(1);
  }

  if (args[0] === "--empty" || args[0] === "-e") {
    console.log("0x");
    console.log(
      "\n(粘贴到 Remix 的 _data：空 bytes，需之后对代理再调 initialize)"
    );
    return;
  }

  if (args[0] === "--help" || args[0] === "-h") {
    printHelp();
    process.exit(0);
  }

  if (args.length < 2) {
    printHelp();
    process.stderr.write("错误: 需要两个参数: <owners 逗号分隔> <门限正整数>\n");
    process.exit(1);
  }

  const rawOwners = String(args[0]);
  const thresholdStr = String(args[1]);
  const ownerParts = rawOwners
    .split(",")
    .map((s) => s.trim())
    .filter((s) => s.length > 0);

  if (ownerParts.length === 0) {
    process.stderr.write("错误: 至少填一个 owner 地址\n");
    process.exit(1);
  }

  const owners = ownerParts.map((a) => {
    try {
      return getAddress(a);
    } catch (e) {
      process.stderr.write(`错误: 无效地址 "${a}": ${e.message}\n`);
      process.exit(1);
    }
  });

  const threshold = BigInt(thresholdStr);
  if (threshold <= 0n) {
    process.stderr.write("错误: 门限必须 > 0\n");
    process.exit(1);
  }
  if (BigInt(owners.length) < threshold) {
    process.stderr.write("错误: 门限不能大于 owner 数量\n");
    process.exit(1);
  }

  const iface = new Interface([FRAGMENT]);
  const data = iface.encodeFunctionData("initialize", [owners, threshold]);

  console.log(data);
  console.log(
    "\n将上面整行 0x… 粘贴到 Remix 部署 MyTransparentUpgradeableProxy 的 _data 即可。"
  );
}

function printHelp() {
  console.log(`
${require("node:path").basename(__filename)} — 生成透明代理的 initialize calldata

用法:
  node scripts/encodeProxyInitData.cjs <address1,address2,...> <门限>
  node scripts/encodeProxyInitData.cjs --empty

示例:
  node scripts/encodeProxyInitData.cjs 0xAbcdefabcdefabcdefabcdefabcdefabcdefAbcd,0xBbcdefabcdefabcdefabcdefabcdefabcdefBbcd 2
`);
}

main();
