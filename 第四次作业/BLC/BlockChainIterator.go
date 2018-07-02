package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (blc *Blockchain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{blc.Tip, blc.DB}
}

func (blockChainIterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := blockChainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取当前迭代器中hash对应的区块
			currentBlockBytes := b.Get(blockChainIterator.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			//更新迭代器中的当前hash值
			blockChainIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}
