package main

import (
	"context"
	"demo/contract"
	"demo/merkle"
	"demo/solidity"
	"fmt"
	"math/big"
	"os"

	mt "github.com/txaty/go-merkletree"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/siddontang/go/log"
)

func getTransationOpts() *bind.TransactOpts {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), opts.From)
	if err != nil {
		log.Fatal(err)
	}
	opts.Nonce = big.NewInt(int64(nonce))

	return opts
}

func Call() {
	addrBytes, err := os.ReadFile("addr.txt")
	if err != nil {
		panic(err)
	}

	// 读取已部署合约的地址
	contractAddress := common.HexToAddress(string(addrBytes))

	// 实例化合约
	ZkBlobInstance, err := contract.NewZkBlob(contractAddress, client)
	if err != nil {
		panic(err)
	}
	fmt.Println(ZkBlobInstance)

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// test batch output:
	fmt.Println("Test batch data: ", testBatch.BatchNum, testBatch.Blobs)

	// gen tree
	fmt.Println("Gen tree... ", testBatch.Blobs)
	tree, err := merkle.NewMerkleTree(mt.ModeProofGen, testBatch.Blobs...)
	if err != nil {
		panic(err)
	}
	root := tree.Root
	fmt.Println("Tree root: ", hexutil.Encode(tree.Root))

	fmt.Println("Chain ID: ", chainId)

	fmt.Println("Test post sig... ")
	bytesPack, err := solidity.AbiPack([]string{"uint256", "bytes32"}, big.NewInt(testBatch.BatchNum), common.BytesToHash(root))
	if err != nil {
		panic(err)
	}
	// get bytestosign
	bytesToSign, err := solidity.Keccak256(bytesPack)
	if err != nil {
		panic(err)
	}
	fmt.Println("sigtosign: ", bytesToSign)
	// bytesToSign := solsha3.SoliditySHA3(
	// 	// types
	// 	[]string{"uint256", "bytes32"},

	// 	// values
	// 	[]interface{}{
	// 		big.NewInt(testBatch.BatchNum), common.HexToHash("0x68203f90e9d07dc5859259d7536e87a6ba9d345f2552b5b9de2999ddce9ce1bf"),
	// 	},
	// )
	// fmt.Println("sigtosign: ", bytesToSign)

	// sign
	sig, err := solidity.Sign(bytesToSign, privateKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("sig: ", sig)

	// test post
	tx, err := ZkBlobInstance.PostBatch(getTransationOpts(), big.NewInt(testBatch.BatchNum), common.BytesToHash(root), sig)
	if err != nil {
		panic(err)
	}

	// 等待交易被确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Post Success, Transaction mined:", receipt.TxHash.Hex())

	// get proof of hashs[0]
	proof0 := [][32]byte{}
	for _, v := range tree.Proofs[0].Siblings {
		proof0 = append(proof0, common.BytesToHash(v))
	}

	// zkblob verify of hashs[0]
	ok, err := ZkBlobInstance.VerifyBlob(&bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: nil,
		Context:     nil,
	}, big.NewInt(testBatch.BatchNum), testBatch.Blobs[0], proof0)
	if err != nil {
		panic(err)
	}

	fmt.Println("Verify result: ", ok)
}
