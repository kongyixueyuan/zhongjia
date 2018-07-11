package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的Hash
	ZjTxHash []byte
	// 2. 存储TXOutput在Vout里面的索引
	ZjVout int

	ZjSignature []byte // 数字签名

	ZjPublicKey []byte // 公钥，钱包里面
}

// 判断当前的消费是谁的钱
func (txInput *TXInput) ZjUnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := ZjRipemd160Hash(txInput.ZjPublicKey)

	return bytes.Compare(publicKey, ripemd160Hash) == 0
}
