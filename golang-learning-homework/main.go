package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Sepolia 链 ID（也可通过 client.ChainID 读取，作业里写死便于对照文档）
const sepoliaChainID = int64(11155111)

func main() {
	ctx := context.Background()

	rpcURL := strings.TrimSpace(os.Getenv("SEPOLIA_RPC_URL"))
	if rpcURL == "" {
		log.Fatal("请设置环境变量 SEPOLIA_RPC_URL，例如 https://sepolia.infura.io/v3/<YOUR_KEY>")
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("连接 RPC 失败: %v", err)
	}
	defer client.Close()

	if err := queryBlock(ctx, client); err != nil {
		log.Fatalf("查询区块: %v", err)
	}

	if err := maybeSendTx(ctx, client); err != nil {
		log.Fatalf("发送交易: %v", err)
	}
}

// queryBlock 查询指定高度区块；未设置 SEPOLIA_BLOCK_NUMBER 时查最新块。
func queryBlock(ctx context.Context, client *ethclient.Client) error {
	var blockNumber *big.Int
	if s := strings.TrimSpace(os.Getenv("SEPOLIA_BLOCK_NUMBER")); s != "" {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("解析 SEPOLIA_BLOCK_NUMBER: %w", err)
		}
		blockNumber = big.NewInt(n)
	}

	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return err
	}

	ts := time.Unix(int64(block.Time()), 0).UTC()
	fmt.Println("========== 区块信息 ==========")
	fmt.Printf("区块号:     %s\n", block.Number().String())
	fmt.Printf("区块哈希:   %s\n", block.Hash().Hex())
	fmt.Printf("父哈希:     %s\n", block.ParentHash().Hex())
	fmt.Printf("时间戳(UTC): %s\n", ts.Format(time.RFC3339))
	fmt.Printf("交易数量:   %d\n", len(block.Transactions()))
	fmt.Println("==============================")
	return nil
}

// maybeSendTx 在设置了 SEPOLIA_PRIVATE_KEY 与 SEPOLIA_TO_ADDRESS 时发送一笔原生币转账。
func maybeSendTx(ctx context.Context, client *ethclient.Client) error {
	keyHex := strings.TrimSpace(os.Getenv("SEPOLIA_PRIVATE_KEY"))
	toHex := strings.TrimSpace(os.Getenv("SEPOLIA_TO_ADDRESS"))
	if keyHex == "" || toHex == "" {
		fmt.Println("未设置 SEPOLIA_PRIVATE_KEY 或 SEPOLIA_TO_ADDRESS，跳过发送交易。")
		return nil
	}

	priv, err := crypto.HexToECDSA(strings.TrimPrefix(keyHex, "0x"))
	if err != nil {
		return fmt.Errorf("解析私钥: %w", err)
	}

	from := crypto.PubkeyToAddress(priv.PublicKey)
	to := common.HexToAddress(toHex)

	value := big.NewInt(1) // 默认 1 wei，可通过 SEPOLIA_VALUE_WEI 覆盖
	if s := strings.TrimSpace(os.Getenv("SEPOLIA_VALUE_WEI")); s != "" {
		v, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return fmt.Errorf("无效的 SEPOLIA_VALUE_WEI: %s", s)
		}
		value = v
	}

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("nonce: %w", err)
	}

	gasTip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("gas tip: %w", err)
	}
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("最新块头: %w", err)
	}
	// 常见写法：feeCap = tip + 2*baseFee
	gasFeeCap := new(big.Int).Add(gasTip, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))

	chainID := big.NewInt(sepoliaChainID)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTip,
		GasFeeCap: gasFeeCap,
		Gas:       21000,
		To:        &to,
		Value:     value,
		Data:      nil,
	})

	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), priv)
	if err != nil {
		return fmt.Errorf("签名: %w", err)
	}

	if err := client.SendTransaction(ctx, signed); err != nil {
		return fmt.Errorf("广播: %w", err)
	}

	fmt.Println("========== 已广播交易 ==========")
	fmt.Printf("交易哈希: %s\n", signed.Hash().Hex())
	fmt.Printf("发送方:   %s\n", from.Hex())
	fmt.Printf("接收方:   %s\n", to.Hex())
	fmt.Printf("金额(wei): %s\n", value.String())
	fmt.Println("================================")
	return nil
}
