package mne

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	bip39 "github.com/tyler-smith/go-bip39"
)

func ToPrivateKey(mne string) (*ecdsa.PrivateKey, error) {
	// 助记词转换成种子
	seed := bip39.NewSeed(mne, "")

	// 从种子生成主扩展密钥
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// 按照BIP44路径派生出以太坊钱包子密钥 m/44'/60'/0'/0/0
	childKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}
	childKey, err = childKey.Derive(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		return nil, err
	}
	childKey, err = childKey.Derive(hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	childKey, err = childKey.Derive(0)
	if err != nil {
		return nil, err
	}
	childKey, err = childKey.Derive(0)
	if err != nil {
		return nil, err
	}

	// 获取ECDSA私钥
	privKeyBytes, err := childKey.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.ToECDSA(privKeyBytes.Serialize())
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
