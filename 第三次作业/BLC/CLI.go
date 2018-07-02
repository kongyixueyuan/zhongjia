package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
)

type CLI struct{}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockChain -data -- 创建区块链.")
	fmt.Println("\taddBlock -data DATA -- 交易数据.")
	fmt.Println("\tprintBlockChain -- 输出区块信息.")
}

func isVaildArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string)  {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	//获取区块链对象
	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockChain(data)
}

/**
	打印区块链
 */
func (cli *CLI) printBlockChain()  {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.PrintChain2()

}

/**
	创建创世区块
 */
func (cli *CLI) createGenesisBlockchain(data string)  {

	CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run() {
	isVaildArgs()

	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printBlockChain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "trx 1000 btc to zhagnsan", "交易数据...")
	flagCreateBlokcChainData := createBlockChainCmd.String("data", "Genesis block data......", "创世区块交易数据...")

	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printBlockChain":
		err := printBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		//fmt.Println(*flagAddBlockData)
		cli.addBlock(*flagAddBlockData)
	}

	if printBlockChainCmd.Parsed() {
		//fmt.Println("输出所有区块的数据........")
		cli.printBlockChain()
	}

	if createBlockChainCmd.Parsed() {
		if *flagCreateBlokcChainData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlokcChainData)
	}
}
