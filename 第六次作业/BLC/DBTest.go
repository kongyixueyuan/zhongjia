package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

/*
	开启一个数据库，如不存在则创建
 */
func openOrCreateDB() {
	db, err := bolt.Open("myDB.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("open or create db sucessfully")
	defer db.Close()
}

/**
	创建一个数据表
 */
func openOrCreateTable() {
	//1.创建/开启数据库
	db, err := bolt.Open("myDB.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("open db sucessfully")
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		//2.创建表
		b, err := tx.CreateBucket([]byte("block"))
		if err != nil {
			return fmt.Errorf("create failed: %s", err)
		}
		fmt.Println("create table sucessfully")
		if b != nil {
			//3.存储数据
			err = b.Put([]byte("hash1"), []byte("tx 1000 eth to lisi"))
			if err != nil {
				log.Panic(err)
			}
			fmt.Println("update data sucessfully")
		}
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}
}

/**
	更新数据
 */
func updateTableData() {
	//1.创建/开启数据库
	db, err := bolt.Open("myDB.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	fmt.Println("open db sucessfully")

	err = db.Update(func(tx *bolt.Tx) error {
		//2.开启表
		b := tx.Bucket([]byte("block"))
		fmt.Println("open table sucessfully")
		if b != nil {
			//3.存储数据
			err = b.Put([]byte("hash2"), []byte("tx 1000 btc to zhangsan"))
			if err != nil {
				log.Panic(err)
			}
		}
		fmt.Println("update data sucessfully")
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}
}

/**
	读取数据
 */
func viewTableData() {
	//1.创建/开启数据库
	db, err := bolt.Open("myDB.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		//2.开启表
		b := tx.Bucket([]byte("block"))
		if b != nil {
			//3.读取数据
			value1 := b.Get([]byte("hash1"))
			fmt.Printf("%s\n", value1)
			value2 := b.Get([]byte("hash2"))
			fmt.Printf("%s\n", value2)
		}
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}
}

func DBTest() {
	//1.数据库安装
	//在命令行中 输入 go get github.com/boltdb/bolt/
	//导入数据库

	//2.创建/开启数据库
	//openOrCreateDB()

	//3.创建/开启表
	//openOrCreateTable()

	//4.更新表数据
	//updateTableData()

	//5.查看表数据
	viewTableData()
}
