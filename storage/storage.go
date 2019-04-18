package storage

import (
	"encoding/gob"
	"errors"
	"gitlab.com/SeaStorage/SeaStorage/crypto"
	"strings"
)

func init() {
	gob.Register(&File{})
	gob.Register(&Directory{})
}

type Root struct {
	Home *Directory
	Keys map[string]*FileKey
}

type FileInfo struct {
	Name      string
	Size      uint
	Hash      string
	Key       string
	Fragments []*Fragment
}

func NewRoot(home *Directory, keys map[string]*FileKey) *Root {
	return &Root{
		Home: home,
		Keys: keys,
	}
}

func NewFileInfo(name string, size uint, hash string, key string, fragments []*Fragment) *FileInfo {
	return &FileInfo{
		Name:      name,
		Size:      size,
		Hash:      hash,
		Key:       key,
		Fragments: fragments,
	}
}

func GenerateRoot() *Root {
	return NewRoot(NewDirectory("home"), make(map[string]*FileKey))
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

func (root *Root) SearchKey(key string, generate bool, creation bool) string {
	keyIndex := crypto.SHA512HexFromHex(key)
	fileKey, ok := root.Keys[string(keyIndex)]
	if ok {
		if creation {
			fileKey.Used++
		}
		return string(keyIndex)
	} else if generate {
		fileKey = NewFileKey(key)
		if creation {
			fileKey.Used++
		}
		root.Keys[string(keyIndex)] = fileKey
		return string(keyIndex)
	}
	return ""
}

func (root *Root) updateKeyUsed(keyUsed map[string]int) {
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
	fileKeyIndex := root.SearchKey(info.Key, true, true)
	return root.Home.CreateFile(path, info.Name, info.Size, info.Hash, fileKeyIndex, info.Fragments)
}

func (root *Root) UpdateName(path string, name string, newName string) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	err = validName(newName)
	if err != nil {
		return err
	}
	return root.Home.UpdateName(path, name, newName)
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
	_ = root.SearchKey(info.Key, true, false)
	keyUsed, err := root.Home.UpdateFileKey(path, info.Name, crypto.SHA512HexFromHex(info.Key), info.Hash, info.Fragments)
	if err != nil {
		return err
	}
	root.updateKeyUsed(keyUsed)
	return nil
}

func (root *Root) PublicKey(publicKey string, key string) error {
	keyBytes := crypto.AESKeyEncryptedByPublicKey(key, publicKey)
	keyIndex := crypto.SHA512HexFromHex(crypto.BytesToHex(keyBytes))
	target, ok := root.Keys[keyIndex]
	if ok {
		if crypto.AESKeyVerify(publicKey, key, target.Key) {
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

func (root *Root) AddSea(path string, name string, hash string, sea *FragmentSea) error {
	return root.Home.AddSea(path, name, hash, sea)
}
