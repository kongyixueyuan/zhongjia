package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const TargetBit = 16

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	//创建一个初始值为1的target
	target := big.NewInt(1)

	//向左移动256-targetBit位
	target.Lsh(target, 256-TargetBit)

	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	bytesData := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(pow.block.Height),
		pow.target.Bytes(),
		IntToHex(nonce),
	}, []byte{})

	return bytesData
}

func (pow *ProofOfWork) Run() ([]byte, int64) {
	var nonce int64 = 0
	var hashInt big.Int
	var hash [32]byte

	fmt.Println("开始挖矿...")
	for {
		//1.将block的属性拼接成字节数组
		dataBytes := pow.prepareData(nonce)

		//2.生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		//3.转换成hashInt，进行hash验证
		hashInt.SetBytes(hash[:])

		//4.验证hash的有效性，如果满足条件，跳出循环
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce = nonce + 1
	}

	fmt.Println()
	fmt.Printf("挖矿成功...nonce = %d \n",nonce)
	return hash[:], nonce
}
