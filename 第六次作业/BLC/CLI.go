package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
)

type CLI struct{}

func zjPrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddressLists -- 输出所有钱包地址.")
	fmt.Println("\tcreateWallet -- 创建钱包.")
	fmt.Println("\tcreateBlockChain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintBlockChain -- 输出区块信息.")
	fmt.Println("\tgetBalance -address -- 获取地址余额.")
}

func zjIsVaildArgs() {
	if len(os.Args) < 2 {
		zjPrintUsage()
		os.Exit(1)
	}
}

// 转账
func (cli *CLI) zjSend(from []string, to []string, amount []string) {

	//main send -from "[\"1KqnEpCnCxosZYyN2aXYBpbBtJe9CtpeZc\"]" -to "[\"1CbUCq1oQuQJRjkSzGiAS1CW9e7oqFvGfN\"]" -amount "[\"4\"]"

	if ZjDBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := ZjBlockchainObject()
	defer blockchain.ZjDB.Close()

	utxoSet := UTXOSet{blockchain}

	blockchain.ZjMineNewBlock(from, to, amount, &utxoSet)
	//update
	utxoSet.ZjUTXOSetUpdate()
}

/**
	打印区块链
 */
func (cli *CLI) zjPrintBlockChain() {

	if ZjDBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := ZjBlockchainObject()

	defer blockchain.ZjDB.Close()

	blockchain.ZjPrintchain()
}

/**
	创建创世区块
 */
func (cli *CLI) zjCreateGenesisBlockchain(data string) {

	blockChain := ZjCreateBlockchainWithGenesisBlock(data)
	defer blockChain.ZjDB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ZjResetUTXOSet()
}

/**
	获取对应账户的余额
 */
func (cli *CLI) zjGetBalance(address string) {

	fmt.Println("地址：" + address)

	blockchain := ZjBlockchainObject()
	defer blockchain.ZjDB.Close()

	utxoSet := UTXOSet{blockchain}

	amount := utxoSet.ZjGetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)
}

func (cli *CLI) zjCreateWallet() {

	wallets, _ := ZjNewWallets()

	wallets.ZjCreateNewWallet()

	fmt.Println(len(wallets.WalletsMap))
}

// 打印所有的钱包地址
func (cli *CLI) zjAddressLists() {

	fmt.Println("打印所有的钱包地址:")

	wallets, _ := ZjNewWallets()

	for address, _ := range wallets.WalletsMap {

		fmt.Println(address)
	}
}

func (cli *CLI) ZjRun() {
	zjIsVaildArgs()

	addresslistsCmd := flag.NewFlagSet("addressLists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("zjSend", flag.ExitOnError)
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
		zjPrintUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			zjPrintUsage()
			os.Exit(1)
		}

		from := ZjJSONToArray(*flagFrom)
		to := ZjJSONToArray(*flagTo)
		for index, fromAdress := range from {
			if ZjIsValidForAdress([]byte(fromAdress)) == false || ZjIsValidForAdress([]byte(to[index])) == false {
				fmt.Printf("地址无效......")
				zjPrintUsage()
				os.Exit(1)
			}
		}

		amount := ZjJSONToArray(*flagAmount)
		cli.zjSend(from, to, amount)
	}

	if printBlockChainCmd.Parsed() {

		cli.zjPrintBlockChain()
	}

	if addresslistsCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.zjAddressLists()
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.zjCreateWallet()
	}

	if createBlockChainCmd.Parsed() {

		if *flagCreateBlokcChainData == "" {
			fmt.Println("地址不能为空....")
			zjPrintUsage()
			os.Exit(1)
		}

		cli.zjCreateGenesisBlockchain(*flagCreateBlokcChainData)
	}

	if getbalanceCmd.Parsed() {

		if ZjIsValidForAdress([]byte(*getbalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			zjPrintUsage()
			os.Exit(1)
		}

		cli.zjGetBalance(*getbalanceWithAdress)
	}
}
