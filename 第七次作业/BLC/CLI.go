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
	fmt.Println("\tresetUTXO -- 重置.")
	fmt.Println("\tstartnode -miner ADDRESS -- 启动节点服务器，并且指定挖矿奖励的地址.")
}

func zjIsVaildArgs() {
	if len(os.Args) < 2 {
		zjPrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) zjSend(from []string, to []string, amount []string, nodeID string, mineNow bool) {
	//main send -from "[\"1KqnEpCnCxosZYyN2aXYBpbBtJe9CtpeZc\"]" -to "[\"1CbUCq1oQuQJRjkSzGiAS1CW9e7oqFvGfN\"]" -amount "[\"4\"]"

	blockchain := ZjBlockchainObject(nodeID)
	defer blockchain.ZjDB.Close()

	if mineNow {
		blockchain.ZjMineNewBlock(from, to, amount, nodeID)

		utxoSet := &UTXOSet{blockchain}

		//转账成功以后，需要更新一下
		utxoSet.ZjUTXOSetUpdate()
	} else {
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}

}

/**
	打印区块链
 */
func (cli *CLI) zjPrintBlockChain(nodeID string) {

	blockchain := ZjBlockchainObject(nodeID)

	defer blockchain.ZjDB.Close()

	blockchain.ZjPrintchain()
}

/**
	创建创世区块
 */
func (cli *CLI) zjCreateGenesisBlockchain(data string, nodeID string) {

	blockChain := ZjCreateBlockchainWithGenesisBlock(nodeID, data)
	defer blockChain.ZjDB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ZjResetUTXOSet()
}

/**
	获取对应账户的余额
 */
func (cli *CLI) zjGetBalance(address string, nodeID string) {

	fmt.Println("地址：" + address)

	blockchain := ZjBlockchainObject(nodeID)
	defer blockchain.ZjDB.Close()

	utxoSet := UTXOSet{blockchain}

	amount := utxoSet.ZjGetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)
}



func (cli *CLI) zjCreateWallet(nodeID string) {

	wallets, _ := ZjNewWallets(nodeID)

	wallets.ZjCreateNewWallet(nodeID)

	fmt.Println(len(wallets.WalletsMap))
}

// 打印所有的钱包地址
func (cli *CLI) zjAddressLists(nodeID string) {

	fmt.Println("打印所有的钱包地址:")

	wallets, _ := ZjNewWallets(nodeID)

	for address, _ := range wallets.WalletsMap {

		fmt.Println(address)
	}
}

func (cli *CLI) zjResetUTXOSet(nodeID string) {

	blockchain := ZjBlockchainObject(nodeID)

	defer blockchain.ZjDB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ZjResetUTXOSet()

}

func (cli *CLI) zjStartNode(nodeID string,minerAdd string)  {

	// 启动服务器
	if minerAdd == "" || ZjIsValidForAdress([]byte(minerAdd))  {
		//  启动服务器
		fmt.Printf("启动服务器:localhost:%s\n",nodeID)
		zjStartServer(nodeID,minerAdd)

	} else {

		fmt.Println("指定的地址无效....")
		os.Exit(0)
	}

}
func (cli *CLI) ZjRun() {
	zjIsVaildArgs()

	//获取节点ID
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!\n")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID:%s\n", nodeID)

	resetUTXOCMD := flag.NewFlagSet("resetUTXO", flag.ExitOnError)
	addresslistsCmd := flag.NewFlagSet("addressLists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("zjSend", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printBlockChain", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	flagCreateBlokcChainData := createBlockChainCmd.String("address", "", "创建创世区块的地址")
	flagFrom := sendBlockCmd.String("from", "", "转账源地址......")
	flagTo := sendBlockCmd.String("to", "", "转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额......")
	flagMiner := startNodeCmd.String("miner", "", "定义挖矿奖励的地址......")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")
	flagMine := sendBlockCmd.Bool("mine",false,"是否在当前节点中立即验证....")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "resetUTXO":
		err := resetUTXOCMD.Parse(os.Args[2:])
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
		cli.zjSend(from, to, amount, nodeID, *flagMine)
	}

	if resetUTXOCMD.Parsed() {

		fmt.Println("重置UTXO表单......")
		cli.zjResetUTXOSet(nodeID)
	}

	if printBlockChainCmd.Parsed() {

		cli.zjPrintBlockChain(nodeID)
	}

	if addresslistsCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.zjAddressLists(nodeID)
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.zjCreateWallet(nodeID)
	}

	if createBlockChainCmd.Parsed() {

		if *flagCreateBlokcChainData == "" {
			fmt.Println("地址不能为空....")
			zjPrintUsage()
			os.Exit(1)
		}

		cli.zjCreateGenesisBlockchain(*flagCreateBlokcChainData, nodeID)
	}

	if getbalanceCmd.Parsed() {

		if ZjIsValidForAdress([]byte(*getbalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			zjPrintUsage()
			os.Exit(1)
		}

		cli.zjGetBalance(*getbalanceWithAdress, nodeID)
	}

	if startNodeCmd.Parsed() {

		cli.zjStartNode(nodeID, *flagMiner)
	}

}
