package BLC

import (
	"time"
	"strconv"
	"fmt"
	"crypto/sha256"
	"bytes"
)

type Block struct {
	//1.区块高度
	Height int64

	//2.上一个区块的hash
	PrevBlockHash []byte

	//3.交易数据
	Data []byte

	//4.时间戳
	Timestamp int64

	//5.区块的Hash
	Hash []byte

	//6.nonce
	Nonce int64
}

/*
	设置Hash
 */
func (block *Block) SetHash() {
	//1.将Height转换成字节数组
	heightBytes := IntToHex(block.Height)
	fmt.Println("heightBytes = ", heightBytes)

	//2.将时间戳转换成字节数组
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timeBytes := []byte(timeString)
	fmt.Println("timeBytes = ", timeBytes)

	//3.拼接所有属性
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes, block.Hash}, []byte{})

	//4.转换成Hash
	hash := sha256.Sum256(blockBytes)
	fmt.Println("hash = ", hash)
	block.Hash = hash[:]
}

/*
	创建新的区块
 */
func NewBlock(height int64, data string, prevBlockHash []byte) *Block {
	//创建区块
	block := Block{height, prevBlockHash, []byte(data), time.Now().Unix(), nil,0}

	//调用工作量证明方法并返回hansh和nonce
	//创建工作量证明对象
	pow := NewProofOfWork(&block)
	//调用方法
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Println()

	return &block
}

/*
	创建创世区块
 */

func CreateGenesisBlock(data string) *Block {
	return NewBlock(1, data, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
