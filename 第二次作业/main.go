package main

import (
	"Projects/MyPublicChain/PublicChain/BLC"
)

func main() {

	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	blockchain.AddBlockToBlockChain(blockchain.Blocks[len(blockchain.Blocks)-1].Height,"trx 100 to zhangsan",blockchain.Blocks[len(blockchain.Blocks)-1].PrevBlockHash)
	blockchain.AddBlockToBlockChain(blockchain.Blocks[len(blockchain.Blocks)-1].Height,"trx 300 to lisi",blockchain.Blocks[len(blockchain.Blocks)-1].PrevBlockHash)
	blockchain.AddBlockToBlockChain(blockchain.Blocks[len(blockchain.Blocks)-1].Height,"trx 200 to wangwu",blockchain.Blocks[len(blockchain.Blocks)-1].PrevBlockHash)

}
