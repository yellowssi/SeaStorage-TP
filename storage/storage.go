package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
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
	return NewRoot(NewDirectory("root"), make(map[string]*FileKey))
}

// Check the path whether valid.
// Valid Name shouldn't contain '/'
func validPath(p string) error {
	if !strings.HasPrefix(p, "/") {
		return errors.New("Path should start with '/': " + p)
	}
	if !strings.HasSuffix(p, "/") {
		return errors.New("Path should end with '/': " + p)
	}
	pParams := strings.Split(p, "/")
	for i := 1; i < len(pParams)-1; i++ {
		if len(pParams[i]) == 0 {
			return errors.New("Path shouldn't contain '//': " + p)
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

func validInfo(p string, name string) error {
	err := validPath(p)
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

func (root *Root) CreateFile(p string, info FileInfo) error {
	err := validInfo(p, info.Name)
	if err != nil {
		return err
	}
	fileKeyIndex := root.SearchKey(info.Key, true, true)
	err = root.Home.CreateFile(p, info.Name, info.Size, info.Hash, fileKeyIndex, info.Fragments)
	if err != nil {
		return err
	}
	root.Home.updateDirectorySize(p)
	return nil
}

func (root *Root) UpdateName(p string, name string, newName string) error {
	err := validInfo(p, name)
	if err != nil {
		return err
	}
	err = validName(newName)
	if err != nil {
		return err
	}
	return root.Home.UpdateName(p, name, newName)
}

func (root *Root) UpdateFileData(p string, info FileInfo) error {
	err := validInfo(p, info.Name)
	if err != nil {
		return err
	}
	return root.Home.UpdateFileData(p, info.Name, info.Size, info.Hash, info.Fragments)
}

func (root *Root) UpdateFileKey(p string, info FileInfo) error {
	err := validInfo(p, info.Name)
	if err != nil {
		return err
	}
	_ = root.SearchKey(info.Key, true, false)
	keyUsed, err := root.Home.UpdateFileKey(p, info.Name, crypto.SHA512HexFromHex(info.Key), info.Hash, info.Fragments)
	if err != nil {
		return err
	}
	root.updateKeyUsed(keyUsed)
	return nil
}

func (root *Root) PublicKey(publicKey string, key string) error {
	keyBytes := crypto.AESKeyEncryptedByPublicKey(key, publicKey)
	keyIndex := crypto.SHA512HexFromBytes(keyBytes)
	target, ok := root.Keys[keyIndex]
	if ok {
		if crypto.AESKeyVerify(publicKey, key, target.Key) {
			target.Key = key
			return nil
		}
	}
	return errors.New("Key error or not exists. ")
}

func (root *Root) DeleteFile(p string, name string) error {
	err := validInfo(p, name)
	if err != nil {
		return err
	}
	keyIndex, err := root.Home.DeleteFile(p, name)
	if err != nil {
		return err
	}
	root.updateKeyUsed(map[string]int{keyIndex: 1})
	root.Home.updateDirectorySize(p)
	return nil
}

func (root *Root) CreateDirectory(p string) error {
	err := validPath(p)
	if err != nil {
		return err
	}
	_, err = root.Home.CreateDirectory(p)
	return err
}

func (root *Root) DeleteDirectory(p string, name string) error {
	err := validInfo(p, name)
	if err != nil {
		return err
	}
	operations, err := root.Home.DeleteDirectory(p, name)
	if err != nil {
		return err
	}
	root.updateKeyUsed(operations)
	root.Home.updateDirectorySize(p)
	return errors.New("Path doesn't exists: " + p + name + "/")
}

func (root *Root) GetFile(p string, name string) (file FileInfo, err error) {
	err = validInfo(p, name)
	if err != nil {
		return
	}
	f, err := root.Home.checkFileExists(p, name)
	if err != nil {
		return
	}
	key := root.Keys[f.KeyIndex]
	return *NewFileInfo(f.Name, f.Size, f.Hash, key.Key, f.Fragments), nil
}

func (root *Root) GetDirectory(p string) (dir *Directory, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Home.checkPathExists(p)
}

func (root *Root) GetINode(p string, name string) (INode, error) {
	return root.Home.checkINodeExists(p, name)
}

func (root *Root) ListDirectory(p string) (iNodes []INodeInfo, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Home.List(p)
}

func (root *Root) AddSea(p string, name string, hash string, sea *FragmentSea) error {
	return root.Home.AddSea(p, name, hash, sea)
}

func RootFromBytes(data []byte) (*Root, error) {
	root := &Root{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(root)
	return root, err
}
