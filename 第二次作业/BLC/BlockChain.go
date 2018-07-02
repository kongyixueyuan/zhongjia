package BLC

type Blockchain struct {
	Blocks []*Block
}

/**
	创建带有创世区块的区块链
 */
func CreateBlockchainWithGenesisBlock() *Blockchain {
	genesisBlock := CreateGenesisBlock("Genesis block ...")
	return &Blockchain{[]*Block{genesisBlock}}
}

/*
	向区块链中增加区块
 */
func (blc *Blockchain) AddBlockToBlockChain(height int64, data string, prevHash []byte) {
	newBlock := NewBlock(height, data, prevHash)

	blc.Blocks = append(blc.Blocks, newBlock)
}
