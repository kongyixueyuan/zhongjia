package BLC

import (
	"bytes"
)

type TXOutput struct {
	ZjValue         int64
	ZjRipemd160Hash []byte //用户名
}

func (txOutput *TXOutput) ZjLock(address string) {

	publicKeyHash := Base58Decode([]byte(address))

	txOutput.ZjRipemd160Hash = publicKeyHash[1 : len(publicKeyHash)-4]
}

func ZjNewTXOutput(value int64, address string) *TXOutput {

	txOutput := &TXOutput{value, nil}

	// 设置Ripemd160Hash
	txOutput.ZjLock(address)

	return txOutput
}

// 解锁
func (txOutput *TXOutput) ZjUnLockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1 : len(publicKeyHash)-4]

	return bytes.Compare(txOutput.ZjRipemd160Hash, hash160) == 0
}



