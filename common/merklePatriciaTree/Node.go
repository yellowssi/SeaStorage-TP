package merklePatriciaTree

import (
	"bytes"
	"crypto/sha256"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/common/dataStructure"
	"time"
)

var NIBBLES = []dataStructure.BitArray{
	{Value: []bool{false, false, false, false}},
	{Value: []bool{false, false, false, true}},
	{Value: []bool{false, false, true, false}},
	{Value: []bool{false, false, true, true}},
	{Value: []bool{false, true, false, false}},
	{Value: []bool{false, true, false, true}},
	{Value: []bool{false, true, true, false}},
	{Value: []bool{false, true, true, true}},
	{Value: []bool{true, false, false, false}},
	{Value: []bool{true, false, false, true}},
	{Value: []bool{true, false, true, false}},
	{Value: []bool{true, false, true, true}},
	{Value: []bool{true, true, false, false}},
	{Value: []bool{true, true, false, true}},
	{Value: []bool{true, true, true, false}},
	{Value: []bool{true, true, true, true}},
}

var NibblesLen = NIBBLES[0].GetLength()
var NibblesNum = len(NIBBLES)

var (
	EvenExtensionNode = dataStructure.BitArray{Value: []bool{false, false, false, false}}
	OddExtensionNode  = dataStructure.BitArray{Value: []bool{false, false, false, true}}
	EvenLeafNode      = dataStructure.BitArray{Value: []bool{false, false, true, false}}
	OddLeafNode       = dataStructure.BitArray{Value: []bool{false, false, true, true}}
)

type NodeFlag struct {
	Hash      []byte
	Timestamp uint32
	Flag      uint32
}

// hash Interface: 含有Hash类的接口
type HashInterface interface {
	CalculateHash() []byte
}

// Merkle Patricia Tree Interface: 树节点接口
type NodeInterface interface {
	GetNodeFlag() *NodeFlag
	GenerateNodeFlag()
	UpdateNodeFlag()
	CalculateHash() []byte
}

// Node: 树节点（Leaf Node & Extension Node）
type Node struct {
	Prefix   dataStructure.BitArray
	Key      dataStructure.BitArray
	Value    HashInterface
	NodeFlag *NodeFlag
}

func NewNode(prefix dataStructure.BitArray, key dataStructure.BitArray, value HashInterface) *Node {
	node := &Node{Prefix: prefix, Key: key, Value: value}
	node.GenerateNodeFlag()
	node.NodeFlag.Flag++
	return node
}

func (n *Node) GetNodeFlag() *NodeFlag {
	return n.NodeFlag
}

func (n *Node) GenerateNodeFlag() {
	if n.NodeFlag == nil {
		n.NodeFlag = &NodeFlag{Hash: n.CalculateHash(), Timestamp: uint32(time.Now().Unix()), Flag: 0}
	}
}

func (n *Node) UpdateNodeFlag() {
	n.NodeFlag.Hash = n.CalculateHash()
	n.NodeFlag.Timestamp = uint32(time.Now().Unix())
}

func (n *Node) CalculateHash() []byte {
	var valueHash []byte
	switch n.Value.(type) {
	case *Node:
		childNode := n.Value.(*Node)
		if childNode.NodeFlag != nil {
			valueHash = childNode.NodeFlag.Hash
		} else {
			valueHash = childNode.CalculateHash()
		}
		break
	case *BranchNode:
		childNode := n.Value.(*BranchNode)
		if childNode.NodeFlag != nil {
			valueHash = childNode.NodeFlag.Hash
		} else {
			valueHash = childNode.CalculateHash()
		}
		break
	default:
		valueHash = n.Value.CalculateHash()
		break
	}
	headers := bytes.Join([][]byte{n.Prefix.ToByteArray(), n.Key.ToByteArray(), valueHash}, []byte{})
	hash := sha256.Sum256(headers)
	return hash[:]
}

// BranchNode: 扩展节点
type BranchNode struct {
	Nodes    [16]NodeInterface
	NodeFlag *NodeFlag
}

func NewBranchNode() *BranchNode {
	branchNode := &BranchNode{Nodes: [16]NodeInterface{}}
	branchNode.GenerateNodeFlag()
	branchNode.NodeFlag.Flag++
	return branchNode
}

func (bn *BranchNode) GetNodeFlag() *NodeFlag {
	return bn.NodeFlag
}

func (bn *BranchNode) GenerateNodeFlag() {
	bn.NodeFlag = &NodeFlag{Hash: bn.CalculateHash(), Timestamp: uint32(time.Now().Unix()), Flag: 0}
}

func (bn *BranchNode) UpdateNodeFlag() {
	bn.NodeFlag.Hash = bn.CalculateHash()
	bn.NodeFlag.Timestamp = uint32(time.Now().Unix())
}

func (bn *BranchNode) CalculateHash() []byte {
	var nodes []byte
	var valueHash []byte
	for i := 0; i < 16; i++ {
		switch bn.Nodes[i].(type) {
		case *Node:
			childNode := bn.Nodes[i].(*Node)
			if childNode.NodeFlag != nil {
				valueHash = childNode.NodeFlag.Hash
			} else {
				valueHash = childNode.CalculateHash()
			}
			break
		case *BranchNode:
			childNode := bn.Nodes[i].(*BranchNode)
			if childNode.NodeFlag != nil {
				valueHash = childNode.NodeFlag.Hash
			} else {
				valueHash = childNode.CalculateHash()
			}
			break
		}
		nodes = bytes.Join([][]byte{nodes, NIBBLES[i].ToByteArray(), valueHash}, []byte{})
	}
	hash := sha256.Sum256(nodes)
	return hash[:]
}
