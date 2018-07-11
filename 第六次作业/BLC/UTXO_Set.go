package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	ZjBlc *Blockchain
}

/**
	重置UTXO Set
 */
func (utxoSet *UTXOSet) ZjResetUTXOSet() {
	db := utxoSet.ZjBlc.ZjDB
	bucketName := []byte(utxoBucket)

	//先删除原先的表，再新建表
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//存储
	utxos := utxoSet.ZjBlc.ZjFindUTXO()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txID, outputs := range utxos {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outputs.ZjSerialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (u *UTXOSet) ZjFindSpendableOutputs(address string, amount int64) (int64, map[string][]int) {
	db := u.ZjBlc.ZjDB
	var accumulated int64 = 0
	unspendOutputs := make(map[string][]int)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := ZjDeserializeOutputs(v)

			for index, out := range outs.ZjOutputs {
				if out.ZjUnLockScriptPubKeyWithAddress(address) && accumulated < amount {
					accumulated += out.ZjValue
					unspendOutputs[txID] = append(unspendOutputs[txID], index)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//额度不足
	if accumulated < amount {
		log.Panic("Insufficient amount")
	}

	return accumulated, unspendOutputs
}

/**
	查找某一指定地址所有未发费的输出
 */
func (u *UTXOSet) ZjFindTXOutputs(address string) []*TXOutput {
	db := u.ZjBlc.ZjDB
	var utxos []*TXOutput

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := ZjDeserializeOutputs(v)
			for _, out := range outs.ZjOutputs {
				if out.ZjUnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return utxos
}

// 查询余额
func (u *UTXOSet) ZjGetBalance(address string) int64 {

	utxos := u.ZjFindTXOutputs(address)
	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.ZjValue
	}
	return amount
}

//更新UTXOSet
/*
	增：将transaction中的vouts添加到UXTOSet中。
	删：遍历transaction中的vins，根据input删除UXTOSet中对应的output
 */
func (u *UTXOSet) ZjUTXOSetUpdate() {
	db := u.ZjBlc.ZjDB
	block := u.ZjBlc.ZjIterator().ZjNext()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.ZjTxs {

			if tx.ZjIsCoinbaseTransaction() == false {

				for _, in := range tx.ZjVins {
					updateOuts := TXOutputs{}
					outsByte := b.Get(in.ZjTxHash)
					outputs := ZjDeserializeOutputs(outsByte)
					for index, out := range outputs.ZjOutputs {
						if index != in.ZjVout {
							//如果是不相等，则是未花费的输出，需要添加到updateOuts中
							updateOuts.ZjOutputs = append(updateOuts.ZjOutputs, out)
						}
					}
					if len(updateOuts.ZjOutputs) == 0 { //说明该hash对应的outputs全被花费了
						err := b.Delete(in.ZjTxHash)
						if err != nil {
							log.Panic(err)
						}
					} else { //将更新过的outputs重新保存到数据表中
						err := b.Put(in.ZjTxHash, updateOuts.ZjSerialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}
			newOutputs := TXOutputs{}
			for _, out := range tx.ZjVouts {
				newOutputs.ZjOutputs = append(newOutputs.ZjOutputs, out)
			}
			err := b.Put(tx.ZjTxHash, newOutputs.ZjSerialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
