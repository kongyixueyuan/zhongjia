package main

import (
	"./BLC"
)

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
	cli.ZjRun()

}
