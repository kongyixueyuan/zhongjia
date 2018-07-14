package BLC

type Version struct {
	ZjVersion    int64 // 版本
	ZjBestHeight int64 // 当前节点区块的高度
	ZjAddrFrom   string //当前节点的地址
}
