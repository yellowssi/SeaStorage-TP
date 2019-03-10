package storage

import (
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/internal/lib/crypto"
	"strings"
)

type Root struct {
	Home   *Directory
	Shared *Directory
	Keys []*FileKey
}

// Check the path whether valid.
// Valid Name shouldn't contain '/'
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
		return errors.New("Name shouldn't be nil: " + name)
	}
	if strings.Contains(name, "/") {
		return errors.New("Name shouldn't contain '/': " + name)
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
		if len(fragment.Seas) == 0 {
			return errors.New("File should contain storage address(es) for store data. ")
		}
	}
	return nil
}

func (root *Root) SearchKey(key crypto.Key, create bool) FileKeyIndex {
	for i, fileKey := range root.Keys {
		if fileKey.Key == key {
			return FileKeyIndex(i)
		}
	}
	if create {
		root.Keys = append(root.Keys, NewFileKey(key))
		return FileKeyIndex(len(root.Keys) - 1)
	}
	return FileKeyIndex(-1)
}

func (root *Root) UploadFile(path string, name string, size uint, hash Hash, key crypto.Key, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	fileKeyIndex := root.SearchKey(key, true)
	return root.Home.CreateFile(path, name, size, hash, fileKeyIndex, fragments)
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
	return root.Home.UpdateFileName(path, name, newName)
}

func (root *Root) UpdateFileData(path string, name string, size uint, hash Hash, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	return root.Home.UpdateFileData(path, name, size, hash, fragments)
}

func (root *Root) UpdateFileKey(path string, name string, size uint, hash Hash, key crypto.Key, fragments []*Fragment) error {
	err := ValidFile(path, name, fragments)
	if err != nil {
		return err
	}
	return root.UpdateFileKey(path, name, size, hash, key, fragments)
}

func (root *Root) PublicKey(address crypto.Address, index FileKeyIndex, key crypto.Key) error {
	target := root.Keys[index]
	if target.Key.Verify(address, key) {
		target.Key = key
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
	err, dir := root.Home.CheckPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.INodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				sharedPath := iNode.GetSharedPath()
				if sharedPath != "" {
					err = root.Shared.DeleteFile(sharedPath, name)
					if err != nil {
						return err
					}
				}
			}
			dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
			return nil
		default:
			break
		}
	}
	return errors.New("File doesn't exists: " + path + name)
}

func (root *Root) CreateDirectory(path string) error {
	err := ValidPath(path)
	if err != nil {
		return err
	}
	err, _ = root.Home.CreateDirectory(path)
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
	err, dir := root.Home.CheckPathExists(path)
	if err != nil {
		return err
	}
	for i, iNode := range dir.INodes {
		switch iNode.(type) {
		case *Directory:
			if iNode.GetName() == name {
				sharedPath := iNode.GetSharedPath()
				if sharedPath != "" {
					err = root.Shared.DeleteDirectory(sharedPath, name)
					if err != nil {
						return err
					}
				}
			}
			dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
			return nil
		default:
			break
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
	err, dir := root.Shared.CheckPathExists(dstPath)
	if err != nil {
		err, dir = root.Shared.CreateDirectory(dstPath)
		if err != nil {
			return err
		}
	}
	err, iNode := root.Home.CheckINodeExists(srcPath, name)
	if err != nil {
		return err
	}
	if iNode.GetSharedPath() != "" {
		return errors.New("This File or Directory is already shared. ")
	}
	iNode.SetSharedPath(dstPath)
	dir.INodes = append(dir.INodes, iNode)
	return nil
}
