package storage

import (
	"encoding/hex"
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/internal/crypto"
	"strings"
)

type Root struct {
	home   *Directory
	shared *Directory
	keys   map[Hash]*FileKey
}

func NewRoot() *Root {
	return &Root{home: NewDirectory("home"), shared: NewDirectory("shared"), keys: make(map[Hash]*FileKey)}
}

// Check the path whether valid.
// Valid name shouldn't contain '/'
func ValidPath(path string) error {
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

// Check the name whether valid.
// Valid Path
// (1) start and end with '/'
// (2) not contain '//'
func ValidName(name string) error {
	if len(name) == 0 {
		return errors.New("name shouldn't be nil: " + name)
	}
	if strings.Contains(name, "/") {
		return errors.New("name shouldn't contain '/': " + name)
	}
	return nil
}

func ValidFile(path string, name string, fragments []*Fragment) error {
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err = ValidName(name)
	if err != nil {
		return err
	}
	if len(fragments) == 0 {
		return errors.New("File should contain storage address(es) for store data. ")
	}
	for _, fragment := range fragments {
		if len(fragment.seas) == 0 {
			return errors.New("File should contain storage address(es) for store data. ")
		}
	}
	return nil
}

func (root *Root) SearchKey(key crypto.Key, used bool) (hash Hash) {
	keyBytes, _ := hex.DecodeString(string(key))
	keyIndex := crypto.SHA256(keyBytes)
	fileKey := root.keys[Hash(keyIndex)]
	if fileKey != nil {
		return Hash(keyIndex)
	} else if used {
		return Hash(keyIndex)
	}
	return hash
}

func (root *Root) UploadFile(path string, name string, size uint, hash Hash, key crypto.Key, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	fileKeyIndex := root.SearchKey(key, true)
	return root.home.CreateFile(path, name, size, hash, fileKeyIndex, fragments)
}

func (root *Root) UpdateFileName(path string, name string, newName string) error {
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err = ValidName(name)
	if err != nil {
		return err
	}
	err = ValidName(newName)
	if err != nil {
		return err
	}
	return root.home.UpdateFileName(path, name, newName)
}

func (root *Root) UpdateFileData(path string, name string, size uint, hash Hash, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	return root.home.UpdateFileData(path, name, size, hash, fragments)
}

func (root *Root) UpdateFileKey(path string, name string, size uint, hash Hash, key crypto.Key, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	return root.UpdateFileKey(path, name, size, hash, key, fragments)
}

func (root *Root) PublicKey(address crypto.Address, keyHash Hash, key crypto.Key) error {
	target := root.keys[keyHash]
	if target.key.Verify(address, key) {
		target.key = key
	}
	return nil
}

func (root *Root) DeleteFile(path string, name string) error {
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err = ValidName(name)
	if err != nil {
		return err
	}
	err, dir := root.home.CheckPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				sharedPath := iNode.GetSharedPath()
				if sharedPath != "" {
					err, keyIndex := root.shared.DeleteFile(sharedPath, name)
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
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err, _ = root.home.CreateDirectory(path)
	return err
}

func (root *Root) DeleteDirectory(path string, name string) error {
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err = ValidName(name)
	if err != nil {
		return err
	}
	err, dir := root.home.CheckPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *Directory:
			if iNode.GetName() == name {
				sharedPath := iNode.GetSharedPath()
				if sharedPath != "" {
					err, operations := root.shared.DeleteDirectory(sharedPath, name)
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

func (root *Root) ShareFiles(srcPath string, name string, dstPath string) error {
	err := ValidPath(srcPath)
	if err != nil {
		return err
	}
	err = ValidName(name)
	if err != nil {
		return err
	}
	err = ValidPath(dstPath)
	if err != nil {
		return err
	}
	err, dir := root.shared.CheckPathExists(dstPath)
	if err != nil {
		err, dir = root.shared.CreateDirectory(dstPath)
		if err != nil {
			return err
		}
	}
	err, iNode := root.home.CheckINodeExists(srcPath, name)
	if err != nil {
		return err
	}
	if iNode.GetSharedPath() != "" {
		return errors.New("This File or Directory is already shared. ")
	}
	iNode.SetSharedPath(dstPath)
	dir.iNodes = append(dir.iNodes, iNode)
	return nil
}
