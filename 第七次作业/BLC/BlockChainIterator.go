package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	ZjCurrentHash []byte
	ZjDB          *bolt.DB
}

func (blc *Blockchain) ZjIterator() *BlockChainIterator {
	return &BlockChainIterator{blc.ZjTip, blc.ZjDB}
}

func (blockChainIterator *BlockChainIterator) ZjNext() *Block {
	var block *Block

	err := blockChainIterator.ZjDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取当前迭代器中hash对应的区块
			currentBlockBytes := b.Get(blockChainIterator.ZjCurrentHash)
			block = ZjDeserializeBlock(currentBlockBytes)
			//更新迭代器中的当前hash值
			blockChainIterator.ZjCurrentHash = block.ZjPrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	blockChainIterator.ZjCurrentHash = block.ZjPrevBlockHash

	return block
}
