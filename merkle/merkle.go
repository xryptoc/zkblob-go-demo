package merkle

import (
	"demo/solidity"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	mt "github.com/txaty/go-merkletree"
)

type dataBlock struct {
	value common.Hash
}

func (d *dataBlock) Serialize() ([]byte, error) {
	if len(d.value) == 0 {
		return nil, errors.New("datablock should not beempty value")
	}
	return d.value[:], nil
}

func NewMerkleTree(mode mt.TypeConfigMode, hashs ...common.Hash) (*mt.MerkleTree, error) {
	dataBlocs := []mt.DataBlock{}
	for _, hash := range hashs {
		dataBlocs = append(dataBlocs, &dataBlock{
			value: hash,
		})
	}
	cfg := &mt.Config{
		Mode:               mode,
		DisableLeafHashing: true,
		SortSiblingPairs:   true,
		HashFunc:           solidity.Keccak256,
	}

	tree, err := mt.New(cfg, dataBlocs)
	if err != nil {
		return nil, err
	}

	return tree, nil
}
