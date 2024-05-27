package main

import (
	"crypto/ecdsa"
	"demo/mne"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/siddontang/go/log"
)

var (
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey

	testBatch *BatchInHash
)

func init() {
	var err error
	// 获取私钥
	privateKey, err = mne.ToPrivateKey("inner twenty cricket science clever crime typical thumb there hospital monitor summer")
	if err != nil {
		panic(err)
	}

	client, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	// parse test batch data
	testBatch, err = parseBatch("./batch.json")
	if err != nil {
		log.Fatal(err)
	}
}
