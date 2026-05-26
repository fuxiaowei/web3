package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"golang-learning-homework/homework-task2/counter" // 这里引入你生成的 Go 包

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	rpcURL := os.Getenv("ETH_RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")

	// 1. 连接 RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 加载私钥
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 获取公钥 / 地址
	publicKey := privKey.Public().(*ecdsa.PublicKey)
	fromAddr := crypto.PubkeyToAddress(*publicKey)

	// 4. 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 5. 获取 gas 费用
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatal(err)
	}
	header, _ := client.HeaderByNumber(ctx, nil)
	baseFee := header.BaseFee
	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	// 6. 创建交易授权
	chainID, _ := client.ChainID(ctx)
	auth, _ := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasTipCap = gasTipCap
	auth.GasFeeCap = gasFeeCap
	auth.GasLimit = uint64(300000)

	// ======================
	// 使用已部署的合约
	// ======================
	contractAddress := common.HexToAddress("0x1e71B91C314FC092A6385edBf8F6312816262ba9")
	instance, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal("连接合约失败：", err)
	}

	// ======================
	// 调用：increment() +1
	// ======================
	tx, err := instance.Incr(auth)
	if err != nil {
		log.Fatal("调用失败：", err)
	}
	fmt.Println("交易哈希：", tx.Hash().Hex())

	// 等待打包
	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal(err)
	}

	// ======================
	// 查询结果
	// ======================
	count, err := instance.Count(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("当前计数：", count)
}
