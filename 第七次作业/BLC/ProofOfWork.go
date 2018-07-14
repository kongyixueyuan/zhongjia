package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const TargetBit = 16

type ProofOfWork struct {
	zjblock  *Block
	zjtarget *big.Int
}

func ZjNewProofOfWork(block *Block) *ProofOfWork {
	//创建一个初始值为1的target
	target := big.NewInt(1)

	//向左移动256-targetBit位
	target.Lsh(target, 256-TargetBit)

	return &ProofOfWork{block, target}
}

/*
	验证hash是否有效
 */
func (pow *ProofOfWork) ZjIsVaild() bool {
	//将pow中block的hash与target 进行对比
	var hashInt big.Int

	//将[]byte 类型的hash转换成 big.Int类型。
	hashInt.SetBytes(pow.zjblock.ZjHash)

	//进行对比
	if pow.zjtarget.Cmp(&hashInt) == 1 {
		return false
	}

	return true
}

func (pow *ProofOfWork) zjprepareData(nonce int64) []byte {
	bytesData := bytes.Join([][]byte{
		pow.zjblock.ZjPrevBlockHash,
		pow.zjblock.ZjHashTransactions(),
		ZjIntToHex(pow.zjblock.ZjTimestamp),
		ZjIntToHex(pow.zjblock.ZjHeight),
		pow.zjtarget.Bytes(),
		ZjIntToHex(nonce),
	}, []byte{})

	return bytesData
}

func (pow *ProofOfWork) ZjRun() ([]byte, int64) {
	var nonce int64 = 0
	var hashInt big.Int
	var hash [32]byte

	fmt.Println("开始挖矿...")
	for {
		//1.将block的属性拼接成字节数组
		dataBytes := pow.zjprepareData(nonce)

		//2.生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		//3.转换成hashInt，进行hash验证
		hashInt.SetBytes(hash[:])

		//4.验证hash的有效性，如果满足条件，跳出循环
		if pow.zjtarget.Cmp(&hashInt) == 1 {
			break
		}
		nonce = nonce + 1
	}

	fmt.Println()
	fmt.Printf("挖矿成功...nonce = %d \n", nonce)
	return hash[:], nonce
}
