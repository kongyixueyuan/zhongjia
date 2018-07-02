package main

import (
	"./BLC"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//区块的序列化与反序列化测试
func blockSerializeTest() {
	block := BLC.NewBlock(1, "trx 100 to zhangsan", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	blockBytes := block.Serialize()
	fmt.Println(blockBytes)
	block = BLC.DeserializeBlock(blockBytes)
	fmt.Printf("%x\n", block.Hash)
	fmt.Printf("%d\n", block.Nonce)
}

func storeBlockToDB() {
	block := BLC.NewBlock(1, "trx 100 to zhangsan", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	db, err := bolt.Open("block.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("blocks"))
		}
		err = bucket.Put([]byte("block1"), block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func readBlockFromDB() {
	db, err := bolt.Open("block.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))
		if bucket != nil {
			blockByte := bucket.Get([]byte("block1"))
			block := BLC.DeserializeBlock(blockByte)
			fmt.Println(block)
		}
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	//blockChainTest()

	//blockSerializeTest()

	//BLC.DBTest()

	//将block对象存储到数据库中
	//storeBlockToDB()

	//从数据库中读取block数据
	//readBlockFromDB()

	/*******************************************************
	//创建创世区块并存储到数据库中
	//添加区块，并存储到数据库中
	//遍历区块链

	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockChain("trx 1000 btc to zhangsan")
	blockchain.AddBlockToBlockChain("trx 200 eth to lisi")
	blockchain.AddBlockToBlockChain("trx 50 eos to wangwu")
	blockchain.AddBlockToBlockChain("trx 88 bch to sunliu")
	blockchain.PrintChain2()
	*/

	cli := BLC.CLI{}
	cli.Run()

}
