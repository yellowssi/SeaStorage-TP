package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/mitchellh/copystructure"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
	"gitlab.com/SeaStorage/SeaStorage-TP/sea"
	"strings"
)

func init() {
	gob.Register(&File{})
	gob.Register(&Directory{})
}

type Root struct {
	Home  *Directory
	Share *Directory
	Keys  map[string]*FileKey
}

type FileInfo struct {
	Name      string
	Size      int64
	Hash      string
	Key       string
	Fragments []*Fragment
}

func NewRoot(home, share *Directory, keys map[string]*FileKey) *Root {
	return &Root{
		Home:  home,
		Share: share,
		Keys:  keys,
	}
}

func NewFileInfo(name string, size int64, hash string, key string, fragments []*Fragment) *FileInfo {
	return &FileInfo{
		Name:      name,
		Size:      size,
		Hash:      hash,
		Key:       key,
		Fragments: fragments,
	}
}

func GenerateRoot() *Root {
	return NewRoot(NewDirectory("home"), NewDirectory("shared"), make(map[string]*FileKey))
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

func validInfo(p, name string) error {
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

func (root *Root) SearchKey(key string, generate, used bool) string {
	keyIndex := crypto.SHA512HexFromHex(key)
	fileKey, ok := root.Keys[string(keyIndex)]
	if ok {
		if used {
			fileKey.Used++
		}
		return string(keyIndex)
	} else if generate {
		fileKey = NewFileKey(key)
		if used {
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
	err = root.Home.CreateFile(p, info.Name, info.Hash, fileKeyIndex, info.Size, info.Fragments)
	if err != nil {
		return err
	}
	root.Home.updateDirectorySize(p)
	return nil
}

func (root *Root) UpdateName(p, name, newName string) error {
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

func (root *Root) UpdateFileData(p string, info FileInfo, userOrGroup bool) (map[string][]sea.Operation, error) {
	err := validInfo(p, info.Name)
	if err != nil {
		return nil, err
	}
	return root.Home.UpdateFileData(p, info.Name, info.Hash, info.Size, info.Fragments, userOrGroup, false)
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

func (root *Root) PublicKey(publicKey, key string) error {
	keyBytes := crypto.AESKeyEncryptedByPublicKey(key, publicKey)
	keyIndex := crypto.SHA512HexFromBytes(keyBytes)
	target, ok := root.Keys[keyIndex]
	if ok {
		if crypto.AESKeyVerify(publicKey, key, target.Key) {
			target.Key = key
			return nil
		}
	}
	return errors.New("invalid key or not exists")
}

func (root *Root) DeleteFile(p, name string, userOrGroup bool) (map[string][]sea.Operation, error) {
	err := validInfo(p, name)
	if err != nil {
		return nil, err
	}
	seaOperations, keyIndex, err := root.Home.DeleteFile(p, name, userOrGroup, false)
	if err != nil {
		return nil, err
	}
	root.updateKeyUsed(map[string]int{keyIndex: 1})
	root.Home.updateDirectorySize(p)
	return seaOperations, nil
}

func (root *Root) CreateDirectory(p string) error {
	err := validPath(p)
	if err != nil {
		return err
	}
	_, err = root.Home.CreateDirectory(p)
	return err
}

func (root *Root) DeleteDirectory(p, name string, userOrGroup bool) (map[string][]sea.Operation, error) {
	err := validInfo(p, name)
	if err != nil {
		return nil, err
	}
	seaOperations, keyUsed, err := root.Home.DeleteDirectory(p, name, userOrGroup, false)
	if err != nil {
		return nil, err
	}
	root.updateKeyUsed(keyUsed)
	root.Home.updateDirectorySize(p)
	return seaOperations, nil
}

func (root *Root) Move(p, name, newPath string) error {
	err := validInfo(p, name)
	if err != nil {
		return err
	}
	err = validPath(newPath)
	if err != nil {
		return err
	}
	return root.Home.Move(p, name, newPath)
}

func (root *Root) GetFile(p, name string) (file FileInfo, err error) {
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

func (root *Root) GetINode(p, name string) (INode, error) {
	return root.Home.checkINodeExists(p, name)
}

func (root *Root) ListDirectory(p string) (iNodes []INodeInfo, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Home.List(p)
}

func (root *Root) GetSharedFile(p, name string) (file FileInfo, err error) {
	err = validInfo(p, name)
	if err != nil {
		return
	}
	f, err := root.Share.checkFileExists(p, name)
	if err != nil {
		return
	}
	key := root.Keys[f.KeyIndex]
	return *NewFileInfo(f.Name, f.Size, f.Hash, key.Key, f.Fragments), nil
}

func (root *Root) GetSharedDirectory(p string) (dir *Directory, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Share.checkPathExists(p)
}

func (root *Root) GetSharedINode(p, name string) (INode, error) {
	return root.Share.checkINodeExists(p, name)
}

func (root *Root) ListSharedDirectory(p string) (iNodes []INodeInfo, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Share.List(p)
}

func (root *Root) AddSea(p, name, hash string, sea *FragmentSea) error {
	return root.Home.AddSea(p, name, hash, sea)
}

func (root *Root) ShareFiles(p, name, dst string, userOrGroup bool) (map[string][]sea.Operation, map[string]string, error) {
	iNode, err := root.GetINode(p, name)
	if err != nil {
		return nil, nil, err
	}
	target, err := copystructure.Copy(iNode)
	if err != nil {
		return nil, nil, err
	}
	var seaOperations map[string][]sea.Operation
	if userOrGroup {
		seaOperations = iNode.GenerateSeaOperations(sea.ActionUserShared, true)
	} else {
		seaOperations = iNode.GenerateSeaOperations(sea.ActionGroupShared, true)
	}
	destination, _ := root.Share.CreateDirectory(p)
	destination.INodes = append(destination.INodes, target.(INode))
	var keys = make(map[string]string)
	keyIndexes := iNode.GetKeys()
	for _, keyIndex := range keyIndexes {
		fileKey := root.Keys[keyIndex]
		fileKey.Used++
		keys[keyIndex] = fileKey.Key
	}
	return seaOperations, keys, nil
}

func (root *Root) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(root)
	return buf.Bytes()
}

func RootFromBytes(data []byte) (*Root, error) {
	root := &Root{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(root)
	return root, err
}
