package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"

	"math/big"
	"crypto/elliptic"
	"time"
)

// UTXO
type Transaction struct {
	//1. 交易hash
	ZjTxHash []byte

	//2. 输入
	ZjVins []*TXInput

	//3. 输出
	ZjVouts []*TXOutput
}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) ZjIsCoinbaseTransaction() bool {

	return len(tx.ZjVins[0].ZjTxHash) == 0 && tx.ZjVins[0].ZjVout == -1
}

//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func ZjNewCoinbaseTransaction(address string) *Transaction {

	//添加时间戳，解决coinbase的hash一致的问题。
	dataBytes := ZjIntToHex(time.Now().Unix())

	//代表消费
	txInput := &TXInput{[]byte{}, -1, nil, dataBytes}

	txOutput := ZjNewTXOutput(10, address)

	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}

	//设置hash值
	txCoinbase.ZjHashTransaction()

	return txCoinbase
}

func (tx *Transaction) ZjHashTransaction() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.ZjTxHash = hash[:]
}

//2. 转账时产生的Transaction
func ZjNewSimpleTransaction(from string, to string, amount int64, utxoSet *UTXOSet) *Transaction {
	wallets, _ := ZjNewWallets()
	wallet := wallets.WalletsMap[from]
	// 通过一个函数，返回
	money, spendableUTXODic := utxoSet.ZjFindSpendableOutputs(from, amount)

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {

		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {

			txInput := &TXInput{txHashBytes, index, nil, wallet.ZjPublicKey}
			txIntputs = append(txIntputs, txInput)
		}
	}

	// 转账
	txOutput := ZjNewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = ZjNewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txIntputs, txOutputs}

	//设置hash值
	tx.ZjHashTransaction()

	//进行签名
	utxoSet.ZjBlc.ZjSignTransaction(tx, wallet.ZjPrivateKey)

	return tx
}

func (tx *Transaction) ZjHash() []byte {
	txCopy := tx

	txCopy.ZjTxHash = []byte{}

	hash := sha256.Sum256(txCopy.ZjSerialize())
	return hash[:]
}

func (tx *Transaction) ZjSerialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.ZjIsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.ZjVins {
		if prevTXs[hex.EncodeToString(vin.ZjTxHash)].ZjTxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.ZjTrimmedCopy()

	for inID, vin := range txCopy.ZjVins {
		prevTx := prevTXs[hex.EncodeToString(vin.ZjTxHash)]
		txCopy.ZjVins[inID].ZjSignature = nil
		txCopy.ZjVins[inID].ZjPublicKey = prevTx.ZjVouts[vin.ZjVout].ZjRipemd160Hash
		txCopy.ZjTxHash = txCopy.ZjHash()
		txCopy.ZjVins[inID].ZjPublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ZjTxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.ZjVins[inID].ZjSignature = signature
	}
}

// 拷贝一份新的Transaction用于签名                                    T
func (tx *Transaction) ZjTrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.ZjVins {
		inputs = append(inputs, &TXInput{vin.ZjTxHash, vin.ZjVout, nil, nil})
	}

	for _, vout := range tx.ZjVouts {
		outputs = append(outputs, &TXOutput{vout.ZjValue, vout.ZjRipemd160Hash})
	}

	txCopy := Transaction{tx.ZjTxHash, inputs, outputs}

	return txCopy
}

// 数字签名验证
func (tx *Transaction) ZjVerify(prevTXs map[string]Transaction) bool {
	if tx.ZjIsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.ZjVins {
		if prevTXs[hex.EncodeToString(vin.ZjTxHash)].ZjTxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.ZjTrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.ZjVins {
		prevTx := prevTXs[hex.EncodeToString(vin.ZjTxHash)]
		txCopy.ZjVins[inID].ZjSignature = nil
		txCopy.ZjVins[inID].ZjPublicKey = prevTx.ZjVouts[vin.ZjVout].ZjRipemd160Hash
		txCopy.ZjTxHash = txCopy.ZjHash()
		txCopy.ZjVins[inID].ZjPublicKey = nil

		//通过私钥和ID产生签名
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.ZjSignature)
		r.SetBytes(vin.ZjSignature[:(sigLen / 2)])
		s.SetBytes(vin.ZjSignature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.ZjPublicKey)
		x.SetBytes(vin.ZjPublicKey[:(keyLen / 2)])
		y.SetBytes(vin.ZjPublicKey[(keyLen / 2):])

		//rawPubKey : 公钥
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ZjTxHash, &r, &s) == false {
			return false
		}
	}

	return true
}
