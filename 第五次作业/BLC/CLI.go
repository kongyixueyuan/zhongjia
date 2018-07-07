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
	fmt.Println("\taddressLists -- 输出所有钱包地址.")
	fmt.Println("\tcreateWallet -- 创建钱包.")
	fmt.Println("\tcreateBlockChain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintBlockChain -- 输出区块信息.")
	fmt.Println("\tgetBalance -address -- 获取地址余额.")
}

func isVaildArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

// 转账
func (cli *CLI) send(from []string, to []string, amount []string) {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}

/**
	打印区块链
 */
func (cli *CLI) printBlockChain() {

	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()

}

/**
	创建创世区块
 */
func (cli *CLI) createGenesisBlockchain(data string) {

	CreateBlockchainWithGenesisBlock(data)
}

/**
	获取对应账户的余额
 */
func (cli *CLI) getBalance(address string) {

	fmt.Println("地址：" + address)

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)
}

func (cli *CLI) createWallet() {

	wallets, _ := NewWallets()

	wallets.CreateNewWallet()

	fmt.Println(len(wallets.WalletsMap))
}

// 打印所有的钱包地址
func (cli *CLI) addressLists() {

	fmt.Println("打印所有的钱包地址:")

	wallets, _ := NewWallets()

	for address, _ := range wallets.WalletsMap {

		fmt.Println(address)
	}
}

func (cli *CLI) Run() {
	isVaildArgs()

	addresslistsCmd := flag.NewFlagSet("addressLists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printBlockChain", flag.ExitOnError)

	flagCreateBlokcChainData := createBlockChainCmd.String("address", "", "创建创世区块的地址")
	flagFrom := sendBlockCmd.String("from", "", "转账源地址......")
	flagTo := sendBlockCmd.String("to", "", "转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额......")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addressLists":
		err := addresslistsCmd.Parse(os.Args[2:])
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
	case "getBalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		for index, fromAdress := range from {
			if IsValidForAdress([]byte(fromAdress)) == false || IsValidForAdress([]byte(to[index])) == false {
				fmt.Printf("地址无效......")
				printUsage()
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flagAmount)
		cli.send(from, to, amount)
	}

	if printBlockChainCmd.Parsed() {

		cli.printBlockChain()
	}

	if addresslistsCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.addressLists()
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet()
	}

	if createBlockChainCmd.Parsed() {

		if *flagCreateBlokcChainData == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlokcChainData)
	}

	if getbalanceCmd.Parsed() {

		if IsValidForAdress([]byte(*getbalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getbalanceWithAdress)
	}
}
