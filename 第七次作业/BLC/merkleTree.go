package BLC

import (
	"crypto/sha256"
)

type MerkleTree struct {
	ZjRootNode *MerkleNode
}

type MerkleNode struct {
	ZjLeft  *MerkleNode
	ZjRight *MerkleNode
	ZjData  []byte
}

func ZjNewMerkleTree(data [][]byte) *MerkleTree {

	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])

	}

	for _, datum := range data {
		node := ZjNewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {

		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := ZjNewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

func ZjNewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.ZjData = hash[:]
	} else {
		prevHashes := append(left.ZjData, right.ZjData...)
		hash := sha256.Sum256(prevHashes)
		mNode.ZjData = hash[:]
	}

	mNode.ZjLeft = left
	mNode.ZjRight = right

	return &mNode
}
