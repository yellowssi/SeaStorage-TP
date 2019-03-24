package storage

import (
	"errors"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/internal/crypto"
	"strings"
)

type Hash string
type FileKey struct {
	used uint
	key  crypto.Key
}

type INode interface {
	GetName() string
	GetSize() uint
	GetHash() Hash
	GetSharedPath() string
	SetSharedPath(path string)
}

type File struct {
	name       string
	size       uint
	hash       Hash
	keyIndex   Hash
	fragments  []*Fragment
	sharedPath string
}

type Directory struct {
	name       string
	size       uint
	hash       Hash
	iNodes     []INode
	sharedPath string
}

type Fragment struct {
	hash Hash
	seas []*FragmentSea
}

type FragmentSea struct {
	address crypto.Address
	weight  int8
}

type INodeInfo struct {
	name string
	size uint
	hash Hash
}

func NewFileKey(key crypto.Key) *FileKey {
	return &FileKey{key: key, used: 0}
}

func NewFile(name string, size uint, hash Hash, key Hash, fragments []*Fragment) *File {
	return &File{name: name, size: size, hash: hash, keyIndex: key, fragments: fragments, sharedPath: ""}
}

func NewDirectory(name string) *Directory {
	return &Directory{name: name, size: 0, hash: nil, iNodes: make([]INode, 0), sharedPath: ""}
}

func NewFragment(hash Hash, seas []*FragmentSea) *Fragment {
	return &Fragment{hash: hash, seas: seas}
}

func NewFragmentSea(address crypto.Address) *FragmentSea {
	return &FragmentSea{address: address, weight: 0}
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetSize() uint {
	return f.size
}

func (f *File) GetHash() Hash {
	return f.hash
}

func (f *File) GetSharedPath() string {
	return f.sharedPath
}

func (f *File) SetSharedPath(path string) {
	f.sharedPath = path
}

func (d *Directory) GetName() string {
	return d.name
}

func (d *Directory) GetSize() uint {
	return d.size
}

func (d *Directory) GetHash() Hash {
	return d.hash
}

func (d *Directory) GetSharedPath() string {
	return d.sharedPath
}

func (d *Directory) SetSharedPath(path string) {
	d.sharedPath = path
}

func GenerateINodeInfos(iNodes []INode) []*INodeInfo {
	var infos = make([]*INodeInfo, len(iNodes))
	for i := 0; i < len(iNodes); i++ {
		infos[i].name = iNodes[i].GetName()
		infos[i].size = iNodes[i].GetSize()
		infos[i].hash = iNodes[i].GetHash()
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
		if len(dir.iNodes) == 0 {
			return errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/"), nil
		}
		for j := 0; j < len(dir.iNodes); j++ {
			switch dir.iNodes[j].(type) {
			case *Directory:
				if dir.iNodes[j].GetName() == pathParams[i] {
					dir = dir.iNodes[j].(*Directory)
				}
			default:
				if j == len(dir.iNodes)-1 {
					return errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/"), nil
				}
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
	for _, iNode := range dir.iNodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				return nil, iNode.(*File)
			}
		default:
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
	for _, iNode := range dir.iNodes {
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
					return errors.New("The same name file exists: " + strings.Join(pathParams[:i], "/")), nil
				}
			} else if j == len(dir.iNodes)-1 {
				newDir = NewDirectory(pathParams[j])
				dir.iNodes = append(dir.iNodes, newDir)
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
				d.iNodes[i].(*Directory).UpdateDirectorySize(subPath)
			}
			d.size += d.iNodes[i].GetSize()
		default:
		}
	}
}

// Update the name of directory finding by the path.
func (d *Directory) UpdateDirectoryName(path string, name string) (err error) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err
	}
	dir.name = name
	return
}

// Delete directory key.
func (d *Directory) DeleteDirectoryKey() (operations map[Hash]uint) {
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
func (d *Directory) DeleteDirectory(path string, name string) (err error, operations map[Hash]uint) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, nil
	}
	for i := 0; i < len(dir.iNodes); i++ {
		switch dir.iNodes[i].(type) {
		case *Directory:
			if dir.iNodes[i].GetName() == name {
				operations = dir.iNodes[i].(*Directory).DeleteDirectoryKey()
				dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
				d.UpdateDirectorySize(path)
				return nil, operations
			}
		default:
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
	for i := 0; i < len(dir.iNodes); i++ {
		if dir.iNodes[i].GetName() == name {
			return errors.New("The same name file or directory exists: " + path + name)
		}
	}
	dir.iNodes = append(dir.iNodes, NewFile(name, size, hash, keyHash, fragments))
	d.UpdateDirectorySize(path)
	return nil
}

// Update the filename finding by the path and the name.
func (d *Directory) UpdateFileName(path string, name string, newName string) error {
	err, file := d.CheckFileExists(path, name)
	if err != nil {
		return err
	}
	file.name = newName
	return nil
}

// Update the data of file finding by the filename and the path of file.
func (d *Directory) UpdateFileData(path string, name string, size uint, hash Hash, fragments []*Fragment) error {
	err, file := d.CheckFileExists(path, name)
	if err != nil {
		return err
	}
	file.size = size
	file.hash = hash
	file.fragments = fragments
	return nil
}

// Update the key of file
func (d *Directory) UpdateFileKey(path string, name string, keyHash Hash, hash Hash, fragments []*Fragment) (err error, operations map[Hash]int) {
	err, file := d.CheckFileExists(path, name)
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

// Delete the file finding by the name under the path.
func (d *Directory) DeleteFile(path string, name string) (err error, hash Hash) {
	err, dir := d.CheckPathExists(path)
	if err != nil {
		return err, hash
	}
	for i := 0; i < len(dir.iNodes); i++ {
		switch dir.iNodes[i].(type) {
		case *File:
			file := dir.iNodes[i].(*File)
			if file.GetName() == name {
				dir.iNodes = append(dir.iNodes[:i], dir.iNodes[i+1:]...)
				d.UpdateDirectorySize(path)
				return nil, file.keyIndex
			}
		default:
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
	return nil, GenerateINodeInfos(dir.iNodes)
}
