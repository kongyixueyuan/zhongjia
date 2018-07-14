package BLC

import (
	"io"
	"bytes"
	"log"
	"net"
)

//COMMAND_VERSION
func zjSendVersion(toAddress string, bc *Blockchain) {

	bestHeight := bc.GetBestHeight()

	payload := zjGobEncode(Version{NODE_VERSION, bestHeight, nodeAddress})

	//version
	request := append(zjCommandToBytes(COMMAND_VERSION), payload...)

	zjSendData(toAddress, request)

}

//COMMAND_GETBLOCKS
func zjSendGetBlocks(toAddress string) {

	payload := zjGobEncode(GetBlocks{nodeAddress})

	request := append(zjCommandToBytes(COMMAND_GETBLOCKS), payload...)

	zjSendData(toAddress, request)

}

// 主节点将自己的所有的区块hash发送给钱包节点
//COMMAND_BLOCK
//
func zjSendInv(toAddress string, kind string, hashes [][]byte) {

	payload := zjGobEncode(Inv{nodeAddress, kind, hashes})

	request := append(zjCommandToBytes(COMMAND_INV), payload...)

	zjSendData(toAddress, request)

}

func zjSendGetData(toAddress string, kind string, blockHash []byte) {

	payload := zjGobEncode(GetData{nodeAddress, kind, blockHash})

	request := append(zjCommandToBytes(COMMAND_GETDATA), payload...)

	zjSendData(toAddress, request)
}

func zjSendBlock(toAddress string, block []byte) {

	payload := zjGobEncode(BlockData{nodeAddress, block})

	request := append(zjCommandToBytes(COMMAND_BLOCK), payload...)

	zjSendData(toAddress, request)
}

func zjSendData(to string, data []byte) {

	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
