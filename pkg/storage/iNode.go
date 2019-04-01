package storage

import (
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"strings"
)

type FileKey struct {
	used uint
	key  crypto.Key
}

type INode interface {
	GetName() string
	GetSize() uint
	GetHash() crypto.Hash
	GetShared() bool
	SetShared(shared bool)
}

type File struct {
	name      string
	size      uint
	hash      crypto.Hash
	keyIndex  crypto.Hash
	fragments []*Fragment
	shared    bool
}

type Directory struct {
	name   string
	size   uint
	hash   crypto.Hash
	iNodes []INode
	shared bool
}

type Fragment struct {
	Hash crypto.Hash
	Seas []*FragmentSea
}

type FragmentSea struct {
	Address crypto.Address
	Weight  int8
}

type INodeInfo struct {
	IsDir bool
	Name  string
	Size  uint
}

func NewFileKey(key crypto.Key) *FileKey {
	return &FileKey{key: key, used: 0}
}

func NewFile(name string, size uint, hash crypto.Hash, key crypto.Hash, fragments []*Fragment) *File {
	return &File{name: name, size: size, hash: hash, keyIndex: key, fragments: fragments, shared: false}
}

func NewDirectory(name string) *Directory {
	return &Directory{name: name, size: 0, hash: "", iNodes: make([]INode, 0), shared: false}
}

func NewFragment(hash crypto.Hash, seas []*FragmentSea) *Fragment {
	return &Fragment{Hash: hash, Seas: seas}
}

func NewFragmentSea(address crypto.Address) *FragmentSea {
	return &FragmentSea{Address: address, Weight: 0}
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetSize() uint {
	return f.size
}

func (f *File) GetHash() crypto.Hash {
	return f.hash
}

func (f *File) GetShared() bool {
	return f.shared
}

func (f *File) SetShared(shared bool) {
	f.shared = shared
}

func (d *Directory) GetName() string {
	return d.name
}

func (d *Directory) GetSize() uint {
	return d.size
}

func (d *Directory) GetHash() crypto.Hash {
	return d.hash
}

func (d *Directory) GetShared() bool {
	return d.shared
}

func (d *Directory) SetShared(shared bool) {
	d.shared = shared
}

func generateINodeInfos(iNodes []INode) []INodeInfo {
	var infos = make([]INodeInfo, len(iNodes))
	for i := 0; i < len(iNodes); i++ {
		switch iNodes[i].(type) {
		case *Directory:
			infos[i].IsDir = true
		case *File:
			infos[i].IsDir = false
		}
		infos[i].Name = iNodes[i].GetName()
		infos[i].Size = iNodes[i].GetSize()
	}
	return infos
}

// Check the path whether exists in this Directory INode.
// If exists, return the Directory INode pointer of the path.
// Else, return the error.
func (d *Directory) checkPathExists(path string) (*Directory, error) {
	pathParams := strings.Split(path, "/")
	dir := d
	for i := 1; i < len(pathParams)-1; i++ {
		if len(dir.iNodes) == 0 {
			return nil, errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/")
		}
		for j := 0; j < len(dir.iNodes); j++ {
			switch dir.iNodes[j].(type) {
			case *Directory:
				if dir.iNodes[j].GetName() == pathParams[i] {
					dir = dir.iNodes[j].(*Directory)
				}
			default:
				if j == len(dir.iNodes)-1 {
					return nil, errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/")
				}
			}
		}
	}
	return dir, nil
}

// Check the file whether exists in this Directory INode.
// If exists, return the pointer of the File INode.
// else, return the error.
func (d *Directory) checkFileExists(path string, name string) (*File, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	for _, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				return iNode.(*File), nil
			}
		default:
		}
	}
	return nil, errors.New("File doesn't exists: " + path + name)
}

// Check the file or directory whether exists in this Directory INode.
func (d *Directory) checkINodeExists(path string, name string) (INode, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	for _, iNode := range dir.iNodes {
		if iNode.GetName() == name {
			return iNode, nil
		}
	}
	return nil, errors.New("File or Directory doesn't exists: " + path + name)
}

// Create directories recursively
// If there is the same Name file exists, it will return error.
// Else, it will return the pointer of the determination directory INode.
func (d *Directory) CreateDirectory(path string) (*Directory, error) {
	var newDir *Directory
	dir := d
	pathParams := strings.Split(path, "/")
	for i := 1; i < len(pathParams)-1; i++ {
		if len(dir.iNodes) == 0 {
			newDir = NewDirectory(pathParams[i])
			dir.iNodes = append(dir.iNodes, newDir)
			dir = newDir
			continue
		}
	L:
		for j := 0; j < len(dir.iNodes); j++ {
			if dir.iNodes[j].GetName() == pathParams[i] {
				switch dir.iNodes[j].(type) {
				case *Directory:
					dir = dir.iNodes[j].(*Directory)
					break L
				default:
					return nil, errors.New("The same Name file exists: " + strings.Join(pathParams[:i], "/"))
				}
			} else if j == len(dir.iNodes)-1 {
				newDir = NewDirectory(pathParams[j])
				dir.iNodes = append(dir.iNodes, newDir)
				dir = newDir
				break
			}
		}
	}
	return dir, nil
}

// Update directories' Size in the path recursively.
func (d *Directory) updateDirectorySize(path string) {
	if path == "/" {
		d.size = 0
		for i := 0; i < len(d.iNodes); i++ {
			d.size += d.iNodes[i].GetSize()
		}
		return
	}
	pathParams := strings.Split(path, "/")
	d.size = 0
	for i := 0; i < len(d.iNodes); i++ {
		switch d.iNodes[i].(type) {
		case *Directory:
			if d.iNodes[i].GetName() == pathParams[1] {
				subPath := strings.Join(pathParams[2:], "/")
				subPath = "/" + subPath + "/"
				d.iNodes[i].(*Directory).updateDirectorySize(subPath)
			}
			d.size += d.iNodes[i].GetSize()
		default:
		}
	}
}

// Update the Name of directory finding by the path.
func (d *Directory) UpdateDirectoryName(path string, name string) error {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return err
	}
	dir.name = name
	return nil
}

// Delete directory Key.
func (d *Directory) DeleteDirectoryKey() (operations map[crypto.Hash]uint) {
	for _, iNode := range d.iNodes {
		switch iNode.(type) {
		case *Directory:
			for k, v := range iNode.(*Directory).DeleteDirectoryKey() {
				operations[k] += v
			}
		case *File:
			file := iNode.(*File)
			operations[file.keyIndex]--
		default:
		}
	}
	return operations
}

// Delete iNode of the directory finding by the path.
func (d *Directory) DeleteDirectory(path string, name string) (operations map[crypto.Hash]uint, err error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dir.iNodes); i++ {
		switch dir.iNodes[i].(type) {
		case *Directory:
			if dir.iNodes[i].GetName() == name {
				operations = dir.iNodes[i].(*Directory).DeleteDirectoryKey()
				dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
				d.updateDirectorySize(path)
				return operations, nil
			}
		default:
		}
	}
	return nil, errors.New("Path doesn't exists: " + path + name + "/")
}

// Store the file into the path.
func (d *Directory) CreateFile(path string, name string, size uint, hash crypto.Hash, keyHash crypto.Hash, fragments []*Fragment) error {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(dir.iNodes); i++ {
		if dir.iNodes[i].GetName() == name {
			return errors.New("The same Name file or directory exists: " + path + name)
		}
	}
	dir.iNodes = append(dir.iNodes, NewFile(name, size, hash, keyHash, fragments))
	d.updateDirectorySize(path)
	return nil
}

// Update the filename finding by the path and the Name.
func (d *Directory) UpdateFileName(path string, name string, newName string) error {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return err
	}
	file.name = newName
	return nil
}

// Update the data of file finding by the filename and the path of file.
func (d *Directory) UpdateFileData(path string, name string, size uint, hash crypto.Hash, fragments []*Fragment) error {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return err
	}
	file.size = size
	file.hash = hash
	file.fragments = fragments
	return nil
}

// Update the Key of file
func (d *Directory) UpdateFileKey(path string, name string, keyHash crypto.Hash, hash crypto.Hash, fragments []*Fragment) (err error, operations map[crypto.Hash]int) {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return err, operations
	}
	operations[file.keyIndex]--
	file.keyIndex = keyHash
	file.hash = hash
	file.fragments = fragments
	operations[keyHash]++
	return nil, operations
}

// Delete the file finding by the Name under the path.
func (d *Directory) DeleteFile(path string, name string) (crypto.Hash, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(dir.iNodes); i++ {
		switch dir.iNodes[i].(type) {
		case *File:
			file := dir.iNodes[i].(*File)
			if file.GetName() == name {
				dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
				d.updateDirectorySize(path)
				return file.keyIndex, nil
			}
		default:
		}
	}
	return "", errors.New("File doesn't exists: " + path + name)
}

// List information of iNodes in the path.
func (d *Directory) List(path string) ([]INodeInfo, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	return generateINodeInfos(dir.iNodes), nil
}
