package storage

import (
	"encoding/hex"
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"strings"
)

type Root struct {
	Home *Directory
	Keys map[crypto.Hash]*FileKey
}

type FileInfo struct {
	Name      string
	Size      uint
	Hash      crypto.Hash
	Key       crypto.Key
	Fragments []*Fragment
}

func NewRoot(home *Directory, keys map[crypto.Hash]*FileKey) *Root {
	return &Root{
		Home: home,
		Keys: keys,
	}
}

func NewFileInfo(name string, size uint, hash crypto.Hash, key crypto.Key, fragments []*Fragment) *FileInfo {
	return &FileInfo{
		Name:      name,
		Size:      size,
		Hash:      hash,
		Key:       key,
		Fragments: fragments,
	}
}

func GenerateRoot() *Root {
	return NewRoot(NewDirectory("home"), make(map[crypto.Hash]*FileKey))
}

// Check the path whether valid.
// Valid Name shouldn't contain '/'
func validPath(path string) error {
	if !strings.HasPrefix(path, "/") {
		return errors.New("Path should start with '/': " + path)
	}
	if !strings.HasSuffix(path, "/") {
		return errors.New("Path should end with '/': " + path)
	}
	pathParams := strings.Split(path, "/")
	for i := 1; i < len(pathParams)-1; i++ {
		if len(pathParams[i]) == 0 {
			return errors.New("Path shouldn't contain '//': " + path)
		}
	}
	return nil
}

// Check the Name whether valid.
// Valid Path
// (1) start and end with '/'
// (2) not contain '//'
func validName(name string) error {
	if len(name) == 0 {
		return errors.New("Name shouldn't be nil: " + name)
	}
	if strings.Contains(name, "/") {
		return errors.New("Name shouldn't contain '/': " + name)
	}
	return nil
}

func validInfo(path string, name string) error {
	err := validPath(path)
	if err != nil {
		return err
	}
	err = validName(name)
	if err != nil {
		return err
	}
	return nil
}

func (root *Root) SearchKey(key crypto.Key, used bool) (hash crypto.Hash) {
	keyBytes, _ := hex.DecodeString(string(key))
	keyIndex := crypto.SHA512(keyBytes)
	_, ok := root.Keys[crypto.Hash(keyIndex)]
	if ok {
		return crypto.Hash(keyIndex)
	} else if used {
		return crypto.Hash(keyIndex)
	}
	return hash
}

func (root *Root) updateKeyUsed(keyUsed map[crypto.Hash]int) {
	for k, v := range keyUsed {
		if v < 0 {
			root.Keys[k].Used -= uint(-v)
			if root.Keys[k].Used == 0 {
				delete(root.Keys, k)
			}
		} else {
			root.Keys[k].Used += uint(v)
		}
	}
}

func (root *Root) CreateFile(path string, info FileInfo) error {
	err := validInfo(path, info.Name)
	if err != nil {
		return err
	}
	fileKeyIndex := root.SearchKey(info.Key, true)
	return root.Home.CreateFile(path, info.Name, info.Size, info.Hash, fileKeyIndex, info.Fragments)
}

func (root *Root) UpdateFileName(path string, name string, newName string) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	err = validName(newName)
	if err != nil {
		return err
	}
	return root.Home.UpdateFileName(path, name, newName)
}

func (root *Root) UpdateFileData(path string, info FileInfo) error {
	err := validInfo(path, info.Name)
	if err != nil {
		return err
	}
	return root.Home.UpdateFileData(path, info.Name, info.Size, info.Hash, info.Fragments)
}

func (root *Root) UpdateFileKey(path string, info FileInfo) error {
	err := validInfo(path, info.Name)
	if err != nil {
		return err
	}
	_ = root.SearchKey(info.Key, true)
	keyUsed, err := root.Home.UpdateFileKey(path, info.Name, crypto.SHA512([]byte(info.Key)), info.Hash, info.Fragments)
	if err != nil {
		return err
	}
	root.updateKeyUsed(keyUsed)
	return nil
}

func (root *Root) PublicKey(address crypto.Address, keyHash crypto.Hash, key crypto.Key) error {
	target, ok := root.Keys[keyHash]
	if ok {
		if target.Key.Verify(address, key) {
			target.Key = key
			return nil
		}
	}
	return errors.New("Key error or not exists. ")
}

func (root *Root) DeleteFile(path string, name string) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	dir, err := root.Home.checkPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.INodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
			}
			return nil
		default:
		}
	}
	return errors.New("File doesn't exists: " + path + name)
}

func (root *Root) CreateDirectory(path string) error {
	err := validPath(path)
	if err != nil {
		return err
	}
	_, err = root.Home.CreateDirectory(path)
	return err
}

func (root *Root) DeleteDirectory(path string, name string) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	dir, err := root.Home.checkPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.INodes {
		switch iNode.(type) {
		case *Directory:
			if iNode.GetName() == name {
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
			}
			return nil
		default:
		}
	}
	return errors.New("Path doesn't exists: " + path + name + "/")
}

func (root *Root) GetFile(path string, name string) (file FileInfo, err error) {
	err = validInfo(path, name)
	if err != nil {
		return
	}
	f, err := root.Home.checkFileExists(path, name)
	if err != nil {
		return
	}
	key := root.Keys[f.KeyIndex]
	return *NewFileInfo(f.Name, f.Size, f.Hash, key.Key, f.Fragments), nil
}

func (root *Root) ListDirectory(path string) (iNodes []INodeInfo, err error) {
	err = validPath(path)
	if err != nil {
		return
	}
	return root.Home.List(path)
}

func (root *Root) GetINode(path string, name string) (INode, error) {
	return root.Home.checkINodeExists(path, name)
}
