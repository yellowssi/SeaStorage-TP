package storage

import (
	"encoding/hex"
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"strings"
)

type Root struct {
	home   *Directory
	shared *Directory
	keys   map[crypto.Hash]*FileKey
}

type FileInfo struct {
	Name      string
	Size      uint
	Hash      crypto.Hash
	Key       crypto.Key
	Fragments []*Fragment
}

func NewRoot() *Root {
	return &Root{home: NewDirectory("home"), shared: NewDirectory("shared"), keys: make(map[crypto.Hash]*FileKey)}
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
	_, ok := root.keys[crypto.Hash(keyIndex)]
	if ok {
		return crypto.Hash(keyIndex)
	} else if used {
		return crypto.Hash(keyIndex)
	}
	return hash
}

func (root *Root) UploadFile(path string, name string, size uint, hash crypto.Hash, key crypto.Key, fragments []*Fragment) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	fileKeyIndex := root.SearchKey(key, true)
	return root.home.CreateFile(path, name, size, hash, fileKeyIndex, fragments)
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
	return root.home.UpdateFileName(path, name, newName)
}

func (root *Root) UpdateFileData(path string, name string, size uint, hash crypto.Hash, fragments []*Fragment) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	return root.home.UpdateFileData(path, name, size, hash, fragments)
}

func (root *Root) UpdateFileKey(path string, name string, size uint, hash crypto.Hash, key crypto.Key, fragments []*Fragment) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	return root.UpdateFileKey(path, name, size, hash, key, fragments)
}

func (root *Root) PublicKey(address crypto.Address, keyHash crypto.Hash, key crypto.Key) error {
	target, ok := root.keys[keyHash]
	if ok {
		if target.key.Verify(address, key) {
			target.key = key
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
	dir, err := root.home.checkPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				if iNode.GetShared() {
					keyIndex, err := root.shared.DeleteFile("/", name)
					if err != nil {
						return err
					}
					root.keys[keyIndex].used--
					if root.keys[keyIndex].used == 0 {
						delete(root.keys, keyIndex)
					}
				}
			}
			dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
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
	_, err = root.home.CreateDirectory(path)
	return err
}

func (root *Root) DeleteDirectory(path string, name string) error {
	err := validInfo(path, name)
	if err != nil {
		return err
	}
	dir, err := root.home.checkPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *Directory:
			if iNode.GetName() == name {
				if iNode.GetShared() {
					operations, err := root.shared.DeleteDirectory("/", name)
					if err != nil {
						return err
					}
					for k, v := range operations {
						root.keys[k].used -= v
					}
				}
			}
			dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
			return nil
		default:
		}
	}
	return errors.New("Path doesn't exists: " + path + name + "/")
}

func (root *Root) ShareFiles(srcPath string, name string) error {
	err := validInfo(srcPath, name)
	if err != nil {
		return err
	}
	iNode, err := root.home.checkINodeExists(srcPath, name)
	if err != nil {
		return err
	}
	if iNode.GetShared() {
		return errors.New("This File or Directory is already shared. ")
	}
	iNode.SetShared(true)
	root.shared.iNodes = append(root.shared.iNodes, iNode)
	return nil
}

func (root *Root) CancelShare(name string) error {
	err := validName(name)
	if err != nil {
		return err
	}
	for i, iNode := range root.shared.iNodes {
		if iNode.GetName() == name {
			iNode.SetShared(false)
			root.shared.iNodes = append(root.shared.iNodes[:i], root.shared.iNodes[i+1:]...)
			return nil
		}
	}
	return errors.New("The file or directory is not shared. ")
}

func (root *Root) GetFile(path string, name string) (file FileInfo, err error) {
	err = validInfo(path, name)
	if err != nil {
		return
	}
	f, err := root.home.checkFileExists(path, name)
	if err != nil {
		return
	}
	key := root.keys[f.keyIndex]
	return *NewFileInfo(f.name, f.size, f.hash, key.key, f.fragments), nil
}

func (root *Root) ListDirectory(path string) (iNodes []INodeInfo, err error) {
	err = validPath(path)
	if err != nil {
		return
	}
	return root.home.List(path)
}

func (root *Root) ListSharedDirectory(path string) (iNodes []INodeInfo, err error) {
	err = validPath(path)
	if err != nil {
		return
	}
	return root.shared.List(path)
}

func (root *Root) GetSharedFile(path string, name string) (file FileInfo, err error) {
	err = validInfo(path, name)
	if err != nil {
		return
	}
	f, err := root.shared.checkFileExists(path, name)
	if err != nil {
		return
	}
	key := root.keys[f.keyIndex]
	return *NewFileInfo(f.name, f.size, f.hash, key.key, f.fragments), nil
}
