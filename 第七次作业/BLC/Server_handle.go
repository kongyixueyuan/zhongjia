package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"fmt"
)

func zjHandleVersion(request []byte, bc *Blockchain) {

	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	bestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.ZjBestHeight

	if bestHeight > foreignerBestHeight {
		zjSendVersion(payload.ZjAddrFrom, bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		zjSendGetBlocks(payload.ZjAddrFrom)
	}

}

func zjHandleAddr(request []byte, bc *Blockchain) {

}

func zjHandleGetblocks(request []byte, bc *Blockchain) {

	var buff bytes.Buffer
	var payload GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()

	zjSendInv(payload.ZjAddrFrom, BLOCK_TYPE, blocks)

}

func zjHandleGetData(request []byte, bc *Blockchain) {

	var buff bytes.Buffer
	var payload GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.ZjType == BLOCK_TYPE {

		block, err := bc.GetBlock([]byte(payload.ZjHash))
		if err != nil {
			return
		}

		zjSendBlock(payload.ZjAddrFrom, block)
	}

	if payload.ZjType == "tx" {

	}
}

func zjHandleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := payload.ZjBlock

	block := ZjDeserializeBlock(blockBytes)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.ZjHash)

	if len(transactionArray) > 0 {
		blockHash := transactionArray[0]
		zjSendGetData(payload.ZjAddrFrom, "block", blockHash)

		transactionArray = transactionArray[1:]
	} else {

		fmt.Println("数据库重置......")
		UTXOSet := &UTXOSet{bc}
		UTXOSet.ZjResetUTXOSet()

	}

}

func zjHandleTx(request []byte, bc *Blockchain) {

}

func zjHandleInv(request []byte, bc *Blockchain) {

	var buff bytes.Buffer
	var payload Inv

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.ZjType == BLOCK_TYPE {

		blockHash := payload.ZjItems[0]
		zjSendGetData(payload.ZjAddrFrom, BLOCK_TYPE, blockHash)

		if len(payload.ZjItems) >= 1 {
			transactionArray = payload.ZjItems[1:]
		}
	}

	if payload.ZjType == TX_TYPE {

	}
}
