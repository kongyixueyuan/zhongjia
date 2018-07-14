package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)


func zjStartServer(nodeID string,minerAdd string)  {

	// 当前节点的IP地址
	nodeAddress = fmt.Sprintf("localhost:%s",nodeID)

	ln,err := net.Listen(PROTOCOL,nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	bc := ZjBlockchainObject(nodeID)

	defer bc.ZjDB.Close()

	if nodeAddress != knowNodes[0]{
		 // 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		 zjSendVersion(knowNodes[0],bc)
	}

	for {
		// 收到的数据的格式是固定的，12字节+结构体字节数组
		// 接收客户端发送过来的数据
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go zjHandleConnection(conn,bc)
	}
}


func zjHandleConnection(conn net.Conn,bc *Blockchain) {

	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n",request[:COMMANDLENGTH])

	command := zjBytesToCommand(request[:COMMANDLENGTH])

	// 12字节 + 某个结构体序列化以后的字节数组
	switch command {
		case COMMAND_VERSION:
			zjHandleVersion(request, bc)
		case COMMAND_ADDR:
			zjHandleAddr(request, bc)
		case COMMAND_BLOCK:
			zjHandleBlock(request, bc)
		case COMMAND_GETBLOCKS:
			zjHandleGetblocks(request, bc)
		case COMMAND_GETDATA:
			zjHandleGetData(request, bc)
		case COMMAND_INV:
			zjHandleInv(request, bc)
		case COMMAND_TX:
			zjHandleTx(request, bc)
		default:
			fmt.Println("Unknown command!")
	}

	conn.Close()
}
