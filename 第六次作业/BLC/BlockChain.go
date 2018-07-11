package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
	"encoding/hex"
	"strconv"
	"bytes"
	"crypto/ecdsa"
)

const DBName = "blockchain.db"
const blockTableName = "blocks"
const currentHash = "currentHash"

type Blockchain struct {
	ZjTip []byte   //最新区块的hash
	ZjDB  *bolt.DB //存储区块数据的数据库
}

func ZjDBExists() bool {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 遍历输出所有区块的信息
func (blc *Blockchain) ZjPrintchain() {

	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
	blockchainIterator := blc.ZjIterator()

	for {
		block := blockchainIterator.ZjNext()

		fmt.Printf("Height：%d\n", block.ZjHeight)
		fmt.Printf("PrevBlockHash：%x\n", block.ZjPrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.ZjTimestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.ZjHash)
		fmt.Printf("Nonce：%d\n", block.ZjNonce)
		fmt.Println("Txs:")
		for _, tx := range block.ZjTxs {

			fmt.Printf("%x\n", tx.ZjTxHash)
			fmt.Println("Vins:")
			for _, in := range tx.ZjVins {
				fmt.Printf("%x\n", in.ZjTxHash)
				fmt.Printf("%d\n", in.ZjVout)
				fmt.Printf("%s\n", in.ZjPublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.ZjVouts {
				fmt.Println(out.ZjValue)
				fmt.Println(out.ZjRipemd160Hash)
			}
		}

		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.ZjPrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

/**
	创建带有创世区块的区块链
 */
func ZjCreateBlockchainWithGenesisBlock(address string) *Blockchain {

	// 判断数据库是否存在
	if ZjDBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	// 关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := ZjNewCoinbaseTransaction(address)

			genesisBlock := ZjCreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.ZjHash, genesisBlock.ZjSerialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte(currentHash), genesisBlock.ZjHash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.ZjHash
		}
		return nil
	})
	return &Blockchain{genesisHash, db}
}

//// 增加区块到区块链里面
func (blc *Blockchain) ZjAddBlockToBlockChain(txs []*Transaction) {

	err := blc.ZjDB.Update(func(tx *bolt.Tx) error {

		//1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		//2. 创建新区块
		if b != nil {

			// 先获取最新区块
			blockBytes := b.Get(blc.ZjTip)
			// 反序列化
			block := ZjDeserializeBlock(blockBytes)

			//3. 将区块序列化并且存储到数据库中
			newBlock := ZjNewBlock(txs, block.ZjHeight+1, block.ZjHash)
			err := b.Put(newBlock.ZjHash, newBlock.ZjSerialize())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"currentHash"对应的hash
			err = b.Put([]byte(currentHash), newBlock.ZjHash)
			if err != nil {
				log.Panic(err)
			}
			//5. 更新blockchain的Tip
			blc.ZjTip = newBlock.ZjHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 返回Blockchain对象
func ZjBlockchainObject() *Blockchain {

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

	fmt.Println("ZjBlockchainObject")
	return &Blockchain{tip, db}
}

// 挖掘新的区块
func (blockchain *Blockchain) ZjMineNewBlock(from []string, to []string, amount []string, utxoSet *UTXOSet) {
	//1.建立一笔交易
	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := ZjNewSimpleTransaction(address, to[index], int64(value), utxoSet)
		txs = append(txs, tx)
	}

	//奖励
	tx := ZjNewCoinbaseTransaction(from[0])
	txs = append(txs, tx)

	//2. 在建立新区块之前对txs进行签名验证
	for _, tx := range txs {
		if blockchain.ZjVerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	//3. 建立新的区块
	blockchain.ZjAddBlockToBlockChain(txs)
}

func (bclockchain *Blockchain) ZjSignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {

	if tx.ZjIsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.ZjVins {
		prevTX, err := bclockchain.ZjFindTransaction(vin.ZjTxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ZjTxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) ZjFindTransaction(ID []byte) (Transaction, error) {

	bci := bc.ZjIterator()

	for {
		block := bci.ZjNext()

		for _, tx := range block.ZjTxs {
			if bytes.Compare(tx.ZjTxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.ZjPrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

	return Transaction{}, nil
}

// 验证数字签名
func (bc *Blockchain) ZjVerifyTransaction(tx *Transaction) bool {

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.ZjVins {
		prevTX, err := bc.ZjFindTransaction(vin.ZjTxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ZjTxHash)] = prevTX
	}

	return tx.ZjVerify(prevTXs)
}

/*
	map[string]TXOutputs
	string ---> Transcation Hash
	TXOutputs --> 对应的未花费的输出

	遍历区块链，找出所有的未花费的输出
 */
func (blockchain *Blockchain) ZjFindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentUTXOs := make(map[string][]int)
	blockIterator := blockchain.ZjIterator()

	for {
		block := blockIterator.ZjNext()
	works:
		for _, tx := range block.ZjTxs {
			txID := hex.EncodeToString(tx.ZjTxHash)
			//先遍历vouts，因为vin总是指向之前的区块的outputs
			for idx, out := range tx.ZjVouts {
				if spentUTXOs[txID] != nil {
					for _, spentOutInx := range spentUTXOs[txID] {
						if spentOutInx == idx {
							continue works
						}
					}
				}
				outs := UTXO[txID]
				outs.ZjOutputs = append(outs.ZjOutputs, out)
				UTXO[txID] = outs
			}

			//遍历vins，存储到spentUTXOS中。
			if tx.ZjIsCoinbaseTransaction() == false {
				for _, in := range tx.ZjVins {
					inID := hex.EncodeToString(in.ZjTxHash)
					spentUTXOs[inID] = append(spentUTXOs[inID], in.ZjVout)
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.ZjPrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return UTXO
}
