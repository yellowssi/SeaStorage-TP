package merklePatriciaTree

import (
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/common/dataStructure"
	"time"
)

/**
 * Merkle Patricia Tree
 */
type MPT struct {
	RootNode NodeInterface
}

func (mpt *MPT) GetHash() []byte {
	mpt.RootRecursiveHashUpdate()
	return mpt.RootNode.GetNodeFlag().Hash
}

/**
 * Generate Root Node If Root Node is Nil
 */
func (mpt *MPT) GenerateRoot(value HashInterface) bool {
	if mpt.RootNode != nil {
		return false
	}
	key := value.CalculateHash()
	mpt.RootNode = NewNode(EvenLeafNode, *dataStructure.FromByteArray(key), value)
	return true
}

/**
 * Get Value By Key (Hash)
 */
func (mpt *MPT) Get(key []byte) HashInterface {
	if mpt.RootNode == nil {
		return nil
	}
	keyBitArray := dataStructure.FromByteArray(key)
	nodeInterface := mpt.RootNode
	var node *Node
	var branchNode *BranchNode
	var indexKey *dataStructure.BitArray
	var keyLength int
	index := 0
	NibblesLen := NIBBLES[0].GetLength()
	NibblesNum := len(NIBBLES)
	for {
		switch nodeInterface.(type) {
		case *Node:
			node = nodeInterface.(*Node)
			if dataStructure.CompareBitArray(node.Prefix, EvenExtensionNode) == 0 ||
				dataStructure.CompareBitArray(node.Prefix, OddExtensionNode) == 0 {
				keyLength = node.Key.GetLength()
				if dataStructure.CompareBitArray(*keyBitArray.SubBitArray(index, index+keyLength), node.Key) == 0 {
					index += keyLength
					nodeInterface = node.Value.(*BranchNode)
				} else {
					return nil
				}
			} else {
				if dataStructure.CompareBitArray(*keyBitArray.SubBitArray(index, -1), node.Key) == 0 {
					return node.Value
				} else {
					return nil
				}
			}
			break
		case *BranchNode:
			branchNode = nodeInterface.(*BranchNode)
			indexKey = keyBitArray.SubBitArray(index, index+NibblesLen)
			for i, nibble := range NIBBLES {
				if dataStructure.CompareBitArray(*indexKey, nibble) == 0 {
					index += NibblesLen
					nodeInterface = branchNode.Nodes[i]
					break
				} else if i == NibblesNum-1 {
					return nil
				}
			}
			break
		default:
			return nil
		}
	}
}

/**
 * Insert Value
 */
func (mpt *MPT) Insert(value HashInterface) bool {
	if mpt.RootNode == nil {
		return mpt.GenerateRoot(value)
	} else {
		key := dataStructure.FromByteArray(value.CalculateHash())
		node := InsertNode(mpt.RootNode, *key, value)
		if node != nil {
			mpt.RootNode = node
			return true
		} else {
			return false
		}
	}
}

/**
 * Insert Node (Recursive)
 */
func InsertNode(nodeInterface NodeInterface, key dataStructure.BitArray, value HashInterface) NodeInterface {
	switch nodeInterface.(type) {
	case *Node:
		node := nodeInterface.(*Node)
		node.GenerateNodeFlag()
		keyLength := node.Key.GetLength()
		if dataStructure.CompareBitArray(node.Prefix, EvenExtensionNode) == 0 ||
			dataStructure.CompareBitArray(node.Prefix, OddExtensionNode) == 0 {
			if dataStructure.CompareBitArray(*key.SubBitArray(0, keyLength), node.Key) == 0 {
				result := InsertNode(node.Value.(*BranchNode), *key.SubBitArray(keyLength, -1), value)
				if result != nil {
					node.Value = result
					return node
				} else {
					return nil
				}
			} else {
				i := 0
				for ; i < keyLength; i++ {
					if key.GetBit(i) != node.Key.GetBit(i) {
						break
					}
				}
				if i < NibblesLen {
					newBranchNode := NewBranchNode()
					var newLeafNodePrefix dataStructure.BitArray
					if (key.GetLength()-NibblesLen)%8 == 0 {
						newLeafNodePrefix = EvenLeafNode
					} else {
						newLeafNodePrefix = OddLeafNode
					}
					newNodeIndexKey := key.SubBitArray(0, NibblesLen)
					oldNodeIndexKey := node.Key.SubBitArray(0, NibblesLen)
					for j := 0; j < NibblesNum; j++ {
						if dataStructure.CompareBitArray(*newNodeIndexKey, NIBBLES[j]) == 0 {
							newBranchNode.Nodes[j] =
								NewNode(newLeafNodePrefix, *key.SubBitArray(NibblesLen, -1), value)
						} else if dataStructure.CompareBitArray(*oldNodeIndexKey, NIBBLES[j]) == 0 {
							if node.Key.GetLength() == NibblesLen {
								newBranchNode.Nodes[j] = node.Value.(*BranchNode)
							} else {
								var newExtensionNodePrefix dataStructure.BitArray
								if (node.Key.GetLength()-NibblesLen)%8 == 0 {
									newExtensionNodePrefix = EvenExtensionNode
								} else {
									newExtensionNodePrefix = OddExtensionNode
								}
								newBranchNode.Nodes[j] = NewNode(newExtensionNodePrefix,
									*node.Key.SubBitArray(NibblesLen, -1), node.Value)
							}
						}
					}
					return newBranchNode
				} else {
					diff := i / NibblesLen
					var newExtensionNodePrefix dataStructure.BitArray
					if diff%2 == 0 {
						newExtensionNodePrefix = EvenExtensionNode
					} else {
						newExtensionNodePrefix = OddExtensionNode
					}
					var oldExtensionNodePrefix dataStructure.BitArray
					if (node.Key.GetLength()-diff*4)%8 == 0 {
						oldExtensionNodePrefix = EvenExtensionNode
					} else {
						oldExtensionNodePrefix = OddExtensionNode
					}
					var newLeafNodePrefix dataStructure.BitArray
					if (key.GetLength()-diff*NibblesLen-NibblesLen)%8 == 0 {
						newLeafNodePrefix = EvenLeafNode
					} else {
						newLeafNodePrefix = OddLeafNode
					}
					newBranchNode := NewBranchNode()
					for j, nibble := range NIBBLES {
						if dataStructure.CompareBitArray(
							*key.SubBitArray(diff*NibblesLen, diff*NibblesLen+NibblesLen), nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(newLeafNodePrefix,
								*key.SubBitArray(diff*NibblesLen+NibblesLen, -1), value)
						} else if dataStructure.CompareBitArray(
							*node.Key.SubBitArray(diff*NibblesLen, diff*NibblesLen+NibblesLen), nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(oldExtensionNodePrefix,
								*node.Key.SubBitArray(diff*NibblesLen+NibblesLen, -1), node.Value)
						}
					}
					return NewNode(newExtensionNodePrefix, *key.SubBitArray(0, diff*NibblesLen), newBranchNode)
				}
			}
		} else {
			if dataStructure.CompareBitArray(key, node.Key) == 0 {
				return nil
			} else {
				i := 0
				for ; i < keyLength; i++ {
					if key.GetBit(i) != node.Key.GetBit(i) {
						break
					}
				}
				newBranchNode := NewBranchNode()
				if i < NibblesLen {
					newNodeIndexKey := key.SubBitArray(0, NibblesLen)
					oldNodeIndexKey := node.Key.SubBitArray(0, NibblesLen)
					var newLeafNodePrefix dataStructure.BitArray
					if (key.GetLength()-NibblesLen)%8 == 0 {
						newLeafNodePrefix = EvenLeafNode
					} else {
						newLeafNodePrefix = OddLeafNode
					}
					for j, nibble := range NIBBLES {
						if dataStructure.CompareBitArray(*newNodeIndexKey, nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(newLeafNodePrefix,
								*key.SubBitArray(NibblesLen, -1), value)
						} else if dataStructure.CompareBitArray(*oldNodeIndexKey, nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(newLeafNodePrefix,
								*node.Key.SubBitArray(NibblesLen, -1), value)
						}
					}
					return newBranchNode
				} else {
					diff := i / NibblesLen
					var newExtensionNodePrefix dataStructure.BitArray
					if diff%2 == 0 {
						newExtensionNodePrefix = EvenExtensionNode
					} else {
						newExtensionNodePrefix = OddExtensionNode
					}
					var newLeafNodePrefix dataStructure.BitArray
					if (key.GetLength()-diff*NibblesLen-NibblesLen)%8 == 0 {
						newLeafNodePrefix = EvenLeafNode
					} else {
						newLeafNodePrefix = OddLeafNode
					}
					newNodeIndexKey := key.SubBitArray(diff*NibblesLen, diff*NibblesLen+NibblesLen)
					oldNodeIndexKey := node.Key.SubBitArray(diff*NibblesLen, diff*NibblesLen+NibblesLen)
					for j, nibble := range NIBBLES {
						if dataStructure.CompareBitArray(*newNodeIndexKey, nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(newLeafNodePrefix,
								*key.SubBitArray(diff*NibblesLen+NibblesLen, -1), value)
						} else if dataStructure.CompareBitArray(*oldNodeIndexKey, nibble) == 0 {
							newBranchNode.Nodes[j] = NewNode(newLeafNodePrefix,
								*node.Key.SubBitArray(diff*NibblesLen+NibblesLen, -1), node.Value)
						}
					}
					return NewNode(newExtensionNodePrefix, *key.SubBitArray(0, diff*NibblesLen), newBranchNode)
				}
			}
		}
	case *BranchNode:
		branchNode := nodeInterface.(*BranchNode)
		indexKey := key.SubBitArray(0, NibblesLen)
		for i, nibble := range NIBBLES {
			if dataStructure.CompareBitArray(*indexKey, nibble) == 0 {
				node := InsertNode(branchNode.Nodes[i], *key.SubBitArray(NibblesLen, -1), value)
				if node != nil {
					branchNode.Nodes[i] = node
					return branchNode
				} else {
					return nil
				}
			}
		}
		return nil
	default:
		var newLeafNodePrefix dataStructure.BitArray
		if key.GetLength()%8 == 0 {
			newLeafNodePrefix = EvenLeafNode
		} else {
			newLeafNodePrefix = OddLeafNode
		}
		return NewNode(newLeafNodePrefix, key, value)
	}
}

/**
 * Remove Value By Key (Hash)
 */
func (mpt *MPT) Remove(key []byte) bool {
	if mpt.RootNode == nil {
		return false
	} else {
		return RemoveNode(mpt.RootNode, dataStructure.FromByteArray(key))
	}
}

/**
 * Remove Node (Recursive)
 */
func RemoveNode(nodeInterface NodeInterface, key *dataStructure.BitArray) bool {
	switch nodeInterface.(type) {
	case *Node:
		node := nodeInterface.(*Node)
		if dataStructure.CompareBitArray(node.Prefix, EvenExtensionNode) == 0 ||
			dataStructure.CompareBitArray(node.Prefix, OddExtensionNode) == 0 {
			if dataStructure.CompareBitArray(*key.SubBitArray(0, node.Key.GetLength()), node.Key) == 0 {
				if RemoveNode(node.Value.(*BranchNode), key.SubBitArray(node.Key.GetLength(), -1)) {
					node.GenerateNodeFlag()
					node.NodeFlag.Flag++
					node.NodeFlag.Timestamp = uint32(time.Now().Unix())
					for _, childNode := range node.Value.(*BranchNode).Nodes {
						if childNode != nil {
							return true
						}
					}
					node.Value = nil
					return true
				}
			}
		} else {
			if dataStructure.CompareBitArray(*key, node.Key) == 0 {
				node.Value = nil
				return true
			}
		}
		break
	case *BranchNode:
		branchNode := nodeInterface.(*BranchNode)
		for i, childNode := range branchNode.Nodes {
			if childNode != nil && RemoveNode(childNode, key.SubBitArray(NIBBLES[0].GetLength(), -1)) {
				branchNode.GenerateNodeFlag()
				branchNode.NodeFlag.Flag++
				branchNode.NodeFlag.Timestamp = uint32(time.Now().Unix())
				switch childNode.(type) {
				case *Node:
					if childNode.(*Node).Value == nil {
						branchNode.Nodes[i] = nil
						return true
					}
					break
				case *BranchNode:
					for _, childChildNode := range childNode.(*BranchNode).Nodes {
						if childChildNode != nil {
							return true
						}
					}
					branchNode.Nodes[i] = nil
					break
				}
				return true
			}
		}
		break
	}
	return false
}

/**
 * Recursive Hash Update From Root Node
 */
func (mpt *MPT) RootRecursiveHashUpdate() {
	RecursiveHashUpdate(mpt.RootNode)
}

/**
 * Recursive Hash Update
 */
func RecursiveHashUpdate(nodeInterface NodeInterface) {
	switch nodeInterface.(type) {
	case *Node:
		switch nodeInterface.(*Node).Value.(type) {
		case *Node:
			childNode := nodeInterface.(*Node).Value.(*Node)
			if childNode.NodeFlag != nil {
				RecursiveHashUpdate(childNode)
			}
			break
		case *BranchNode:
			childNodes := nodeInterface.(*Node).Value.(*BranchNode).Nodes
			for _, childNode := range childNodes {
				if childNode != nil {
					switch childNode.(type) {
					case *Node:
						if childNode.(*Node).NodeFlag != nil {
							RecursiveHashUpdate(childNode)
						}
						break
					case *BranchNode:
						if childNode.(*Node).NodeFlag != nil {
							RecursiveHashUpdate(childNode)
						}
						break
					}
				}
			}
			break
		}
		break
	case *BranchNode:
		childNodes := nodeInterface.(*BranchNode).Nodes
		for _, childNode := range childNodes {
			if childNode != nil {
				RecursiveHashUpdate(childNode)
			}
		}
		break
	}
	nodeInterface.UpdateNodeFlag()
}

/**
 * Store Changed Node from Memory to File
 */
func (mpt *MPT) Store(path string) bool {
	return false
}

/**
 * Load Node from File to Memory
 */
func Load(path string) *MPT {
	return &MPT{}
}
