package storage

import (
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/internal/lib/crypto"
	"strings"
)

type Hash string
type FileKey struct {
	Used uint
	Key  crypto.Key
}

type INode interface {
	GetName() string
	GetSize() uint
	GetHash() Hash
	GetSharedPath() string
	SetSharedPath(path string)
}

type File struct {
	Name       string
	Size       uint
	Hash       Hash
	KeyIndex   Hash
	Fragments  []*Fragment
	SharedPath string
}

type Directory struct {
	Name       string
	Size       uint
	Hash       Hash
	INodes     []INode
	SharedPath string
}

type Fragment struct {
	Hash Hash
	Seas []*FragmentSea
}

type FragmentSea struct {
	Address crypto.Address
	Weight  int8
}

type INodeInfo struct {
	Name string
	Size uint
	Hash Hash
}

func NewFileKey(key crypto.Key) *FileKey {
	return &FileKey{Key: key, Used: 0}
}

func NewFile(name string, size uint, hash Hash, key Hash, fragments []*Fragment) *File {
	return &File{Name: name, Size: size, Hash: hash, KeyIndex: key, Fragments: fragments, SharedPath: ""}
}

func NewDirectory(name string) *Directory {
	return &Directory{Name: name, Size: 0, Hash: nil, INodes: make([]INode, 0), SharedPath: ""}
}

func NewFragment(hash Hash, seas []*FragmentSea) *Fragment {
	return &Fragment{Hash: hash, Seas: seas}
}

func NewFragmentSea(address crypto.Address) *FragmentSea {
	return &FragmentSea{Address: address, Weight: 0}
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetSize() uint {
	return f.Size
}

func (f *File) GetHash() Hash {
	return f.Hash
}

func (f *File) GetSharedPath() string {
	return f.SharedPath
}

func (f *File) SetSharedPath(path string) {
	f.SharedPath = path
}

func (d *Directory) GetName() string {
	return d.Name
}

func (d *Directory) GetSize() uint {
	return d.Size
}

func (d *Directory) GetHash() Hash {
	return d.Hash
}

func (d *Directory) GetSharedPath() string {
	return d.SharedPath
}

func (d *Directory) SetSharedPath(path string) {
	d.SharedPath = path
}

func GenerateINodeInfos(iNodes []INode) []*INodeInfo {
	var infos = make([]*INodeInfo, len(iNodes))
	for i := 0; i < len(iNodes); i++ {
		infos[i].Name = iNodes[i].GetName()
		infos[i].Size = iNodes[i].GetSize()
		infos[i].Hash = iNodes[i].GetHash()
	}
	return infos
}

// Check the path whether exists in this Directory INode.
// If exists, return the Directory INode pointer of the path.
// Else, return the error.
func (d *Directory) CheckPathExists(path string) (error, *Directory) {
	pathParams := strings.Split(path, "/")
	dir := d
	for i := 1; i < len(pathParams)-1; i++ {
		if len(dir.INodes) == 0 {
			return errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/"), nil
		}
	L:
		for j := 0; j < len(dir.INodes); j++ {
			switch dir.INodes[j].(type) {
			case *Directory:
				if dir.INodes[j].GetName() == pathParams[i] {
					dir = dir.INodes[j].(*Directory)
				}
				continue L
			default:
				if j == len(dir.INodes)-1 {
					return errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/"), nil
				}
				break
			}
		}
	}
	return nil, dir
}

// Check the file whether exists in this Directory INode.
// If exists, return the pointer of the File INode.
// else, return the error.
func (d *Directory) CheckFileExists(path string, name string) (error, *File) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, nil
	}
	for _, iNode := range dir.INodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				return nil, iNode.(*File)
			}
			break
		default:
			break
		}
	}
	return errors.New("File doesn't exists: " + path + name), nil
}

// Check the file or directory whether exists in this Directory INode.
func (d *Directory) CheckINodeExists(path string, name string) (error, INode) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, nil
	}
	for _, iNode := range dir.INodes {
		if iNode.GetName() == name {
			return nil, iNode
		}
	}
	return errors.New("File or Directory doesn't exists: " + path + name), nil
}

// Create directories recursively
// If there is the same name file exists, it will return error.
// Else, it will return the pointer of the determination directory INode.
func (d *Directory) CreateDirectory(path string) (error, *Directory) {
	var newDir *Directory
	dir := d
	pathParams := strings.Split(path, "/")
	for i := 1; i < len(pathParams)-1; i++ {
		if len(dir.INodes) == 0 {
			newDir = NewDirectory(pathParams[i])
			dir.INodes = append(dir.INodes, newDir)
			dir = newDir
			continue
		}
	L:
		for j := 0; j < len(dir.INodes); j++ {
			if dir.INodes[j].GetName() == pathParams[i] {
				switch dir.INodes[j].(type) {
				case *Directory:
					dir = dir.INodes[j].(*Directory)
					break L
				default:
					return errors.New("The same name file exists: " + strings.Join(pathParams[:i], "/")), nil
				}
			} else if j == len(dir.INodes)-1 {
				newDir = NewDirectory(pathParams[j])
				dir.INodes = append(dir.INodes, newDir)
				dir = newDir
				break
			}
		}
	}
	return nil, dir
}

// Update directories' size in the path recursively.
func (d *Directory) UpdateDirectorySize(path string) {
	if path == "/" {
		d.Size = 0
		for i := 0; i < len(d.INodes); i++ {
			d.Size += d.INodes[i].GetSize()
		}
		return
	}
	pathParams := strings.Split(path, "/")
	d.Size = 0
	for i := 0; i < len(d.INodes); i++ {
		switch d.INodes[i].(type) {
		case *Directory:
			if d.INodes[i].GetName() == pathParams[1] {
				subPath := strings.Join(pathParams[2:], "/")
				subPath = "/" + subPath + "/"
				d.INodes[i].(*Directory).UpdateDirectorySize(subPath)
			}
			d.Size += d.INodes[i].GetSize()
		}
	}
}

// Update the name of directory finding by the path.
func (d *Directory) UpdateDirectoryName(path string, name string) (err error) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err
	}
	dir.Name = name
	return
}

// Delete directory key.
func (d *Directory) DeleteDirectoryKey() (operations map[Hash]uint) {
	for _, iNode := range d.INodes {
		switch iNode.(type) {
		case *Directory:
			for k, v := range iNode.(*Directory).DeleteDirectoryKey() {
				operations[k] += v
			}
			break
		case *File:
			file := iNode.(*File)
			operations[file.KeyIndex]--
			break
		default:
			break
		}
	}
	return operations
}

// Delete iNode of the directory finding by the path.
func (d *Directory) DeleteDirectory(path string, name string) (err error, operations map[Hash]uint) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, nil
	}
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *Directory:
			if dir.INodes[i].GetName() == name {
				operations = dir.INodes[i].(*Directory).DeleteDirectoryKey()
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				d.UpdateDirectorySize(path)
				return nil, operations
			}
			break
		default:
			break
		}
	}
	err = errors.New("Path doesn't exists: " + path + name + "/")
	return err, nil
}

// Store the file into the path.
func (d *Directory) CreateFile(path string, name string, size uint, hash Hash, keyHash Hash, fragments []*Fragment) error {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(dir.INodes); i++ {
		if dir.INodes[i].GetName() == name {
			return errors.New("The same name file or directory exists: " + path + name)
		}
	}
	dir.INodes = append(dir.INodes, NewFile(name, size, hash, keyHash, fragments))
	d.UpdateDirectorySize(path)
	return nil
}

// Update the filename finding by the path and the name.
func (d *Directory) UpdateFileName(path string, name string, newName string) error {
	err, file := d.CheckFileExists(path, name)
	if err != nil {
		return err
	}
	file.Name = newName
	return nil
}

// Update the data of file finding by the filename and the path of file.
func (d *Directory) UpdateFileData(path string, name string, size uint, hash Hash, fragments []*Fragment) error {
	err, file := d.CheckFileExists(path, name)
	if err != nil {
		return err
	}
	file.Size = size
	file.Hash = hash
	file.Fragments = fragments
	return nil
}

// Update the key of file
func (d *Directory) UpdateFileKey(path string, name string, keyHash Hash, hash Hash, fragments []*Fragment) (err error, operations map[Hash]int) {
	err, file := d.CheckFileExists(path, name)
	if err != nil {
		return err, operations
	}
	operations[file.KeyIndex]--
	file.KeyIndex = keyHash
	file.Hash = hash
	file.Fragments = fragments
	operations[keyHash]++
	return nil, operations
}

// Delete the file finding by the name under the path.
func (d *Directory) DeleteFile(path string, name string) (err error, hash Hash) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, hash
	}
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *File:
			file := dir.INodes[i].(*File)
			if file.GetName() == name {
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				d.UpdateDirectorySize(path)
				return nil, file.KeyIndex
			}
			break
		default:
			break
		}
	}
	return errors.New("File doesn't exists: " + path + name), hash
}

// List information of iNodes in the path.
func (d *Directory) List(path string) (error, []*INodeInfo) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, nil
	}
	return nil, GenerateINodeInfos(dir.INodes)
}
