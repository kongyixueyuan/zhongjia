package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
)

const DBName = "blockchain.db"
const blockTableName = "blocks"
const currentHash = "currentHash"

type Blockchain struct {
	Tip []byte   //最新区块的hash
	DB  *bolt.DB //存储区块数据的数据库
}

func DBExists() bool {
	if _,err := os.Stat(DBName); os.IsNotExist(err){
		return false
	}
	return true
}

/**
	创建带有创世区块的区块链
 */
func CreateBlockchainWithGenesisBlock(data string) *Blockchain {
	var blockHash []byte

	if DBExists(){
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	}

	//创建或开启数据库
	db, err := bolt.Open(DBName, 0060, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		//创建表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			//创建创世区块
			genesisBlock := CreateGenesisBlock(data)
			//将创世区块存储到数据库中
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//存储当前的hash值到数据库中
			err = b.Put([]byte(currentHash), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash
		}

		return nil

	})

	return &Blockchain{blockHash, db}
}

/*
	向区块链中增加区块
 */
func (blc *Blockchain) AddBlockToBlockChain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			//获取最新的区块
			blockBytes := b.Get(blc.Tip)
			block := DeserializeBlock(blockBytes)

			//创建新区块
			newBlock := NewBlock(block.Height+1, data, block.Hash)
			//存储新区块
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//存储最新的hash值
			err = b.Put([]byte(currentHash), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			//更新tip
			blc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (blc *Blockchain) PrintChain2() {
	blockChainIterator := blc.Iterator()

	for{
		block := blockChainIterator.Next()

		fmt.Println("-----printchain2----")
		fmt.Printf("Height : %d\n", block.Height)
		fmt.Printf("prevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("data : %x\n", block.Data)
		fmt.Printf("time : %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0{
			break
		}
	}
}

//遍历输出所有区块的信息
func (blc *Blockchain) PrintChain() {

	var block *Block
	var currentHash = blc.Tip

	for {
		err := blc.DB.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(blockTableName))
			if b != nil {
				blockBytes := b.Get(currentHash)
				block = DeserializeBlock(blockBytes)

				fmt.Println("-----printchain----")
				fmt.Printf("Height : %d\n", block.Height)
				fmt.Printf("prevBlockHash : %x\n", block.PrevBlockHash)
				fmt.Printf("data : %x\n", block.Data)
				fmt.Printf("time : %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
				fmt.Printf("Hash : %x\n", block.Hash)
				fmt.Printf("Nonce : %d\n", block.Nonce)

			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

		currentHash = block.PrevBlockHash
	}

}

// 返回Blockchain对象
func BlockchainObject() *Blockchain {

	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte(currentHash))
		}

		return nil
	})

	return &Blockchain{tip,db}
}