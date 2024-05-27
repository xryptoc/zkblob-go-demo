package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type batchInString struct {
	Blobs []string
}

type BatchInHash struct {
	BatchNum int64
	Blobs    []common.Hash
}

func getRandomNum() int64 {
	randomSource := rand.NewSource(time.Now().UnixNano())

	return randomSource.Int63()
}

func parseBatch(f string) (*BatchInHash, error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	b := &batchInString{}

	if err := json.Unmarshal(bs, b); err != nil {
		return nil, err
	}

	batch := &BatchInHash{
		BatchNum: getRandomNum(),
	}

	for _, v := range b.Blobs {
		batch.Blobs = append(batch.Blobs, common.HexToHash(v))
	}

	return batch, nil
}
