package main

import (
	"context"
	"demo/contract"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/siddontang/go/log"
)

func Deploy() {
	var err error

	_, err = os.ReadFile("addr.txt")
	if err == nil {
		fmt.Println("ZkBlob contract already deployed, no need to deploy again")
		return
	}
	ctx := context.Background()

	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))

	addr, tx, _, err := contract.DeployZkBlob(auth, client, common.HexToAddress("0x69eC1D6E8D4EDf341F1a4ab6d6141506A2E38acC"))
	if err != nil {
		log.Fatal(err)
	}

	// 等待交易被确认
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	auth.Nonce = nil
	auth.Value = nil
	fmt.Printf("Deploy tx: %v\n", tx.Hash().Hex())
	fmt.Printf("ZkBlob Addr: %v\n", addr.Hex())
	if err = os.WriteFile("addr.txt", []byte(addr.Hex()), 0644); err != nil {
		panic(err)
	}
}
