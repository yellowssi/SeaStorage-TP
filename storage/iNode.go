package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
	"time"
)

type FileKey struct {
	Used uint
	Key  string
}

type INode interface {
	GetName() string
	GetSize() uint
	GetHash() string
}

type File struct {
	Name      string
	Size      uint
	Hash      string
	KeyIndex  string
	Fragments []*Fragment
}

type Directory struct {
	Name   string
	Size   uint
	Hash   string
	INodes []INode
}

type Fragment struct {
	Hash string
	Seas []*FragmentSea
}

type FragmentSea struct {
	PublicKey string
	Weight    int8
	Timestamp time.Time
}

type INodeInfo struct {
	IsDir bool
	Name  string
	Size  uint
}

func NewFileKey(key string) *FileKey {
	return &FileKey{Key: key, Used: 0}
}

func NewFile(name string, size uint, hash string, key string, fragments []*Fragment) *File {
	return &File{Name: name, Size: size, Hash: hash, KeyIndex: key, Fragments: fragments}
}

func NewDirectory(name string) *Directory {
	return &Directory{Name: name, Size: 0, Hash: "", INodes: make([]INode, 0)}
}

func NewFragment(hash string, seas []*FragmentSea) *Fragment {
	return &Fragment{Hash: hash, Seas: seas}
}

func NewFragmentSea(publicKey string) *FragmentSea {
	return &FragmentSea{PublicKey: publicKey, Weight: 0}
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetSize() uint {
	return f.Size
}

func (f *File) GetHash() string {
	return f.Hash
}

func (d *Directory) GetName() string {
	return d.Name
}

func (d *Directory) GetSize() uint {
	return d.Size
}

func (d *Directory) GetHash() string {
	return d.Hash
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
		if len(dir.INodes) == 0 {
			return nil, errors.New("Path doesn't exists: " + strings.Join(pathParams[:i], "/") + "/")
		}
		for j := 0; j < len(dir.INodes); j++ {
			switch dir.INodes[j].(type) {
			case *Directory:
				if dir.INodes[j].GetName() == pathParams[i] {
					dir = dir.INodes[j].(*Directory)
				}
			default:
				if j == len(dir.INodes)-1 {
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
	for _, iNode := range dir.INodes {
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
	for _, iNode := range dir.INodes {
		if iNode.GetName() == name {
			return iNode, nil
		}
	}
	return nil, errors.New("File or Directory doesn't exists: " + path + name)
}

// Target directories recursively
// If there is the same Name file exists, it will return error.
// Else, it will return the pointer of the determination directory INode.
func (d *Directory) CreateDirectory(path string) (*Directory, error) {
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
					return nil, errors.New("The same Name file exists: " + strings.Join(pathParams[:i], "/"))
				}
			} else if j == len(dir.INodes)-1 {
				newDir = NewDirectory(pathParams[j])
				dir.INodes = append(dir.INodes, newDir)
				dir = newDir
				break
			}
		}
	}
	return dir, nil
}

// Update directories' Size in the path recursively.
func (d *Directory) updateDirectorySize(path string) {
	pathParams := strings.Split(path, "/")
	d.Size = 0
	for i := 0; i < len(d.INodes); i++ {
		switch d.INodes[i].(type) {
		case *Directory:
			if d.INodes[i].GetName() == pathParams[1] {
				subPath := strings.Join(pathParams[2:], "/")
				subPath = "/" + subPath + "/"
				d.INodes[i].(*Directory).updateDirectorySize(subPath)
			}
			d.Size += d.INodes[i].GetSize()
		case *File:
			d.Size += d.INodes[i].GetSize()
		}
	}
}

// Update the Name of directory finding by the path.
func (d *Directory) UpdateName(path string, name string, newName string) error {
	iNode, err := d.checkINodeExists(path, name)
	if err != nil {
		return err
	}
	switch iNode.(type) {
	case *File:
		iNode.(*File).Name = newName
	case *Directory:
		iNode.(*Directory).Name = newName
	}
	return nil
}

// Delete directory Key.
func (d *Directory) DeleteDirectoryKey() map[string]uint {
	operations := make(map[string]uint)
	for _, iNode := range d.INodes {
		switch iNode.(type) {
		case *Directory:
			for k, v := range iNode.(*Directory).DeleteDirectoryKey() {
				operations[k] += v
			}
		case *File:
			file := iNode.(*File)
			operations[file.KeyIndex]--
		default:
		}
	}
	return operations
}

// Delete iNode of the directory finding by the path.
func (d *Directory) DeleteDirectory(path string, name string) (operations map[string]uint, err error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *Directory:
			if dir.INodes[i].GetName() == name {
				operations = dir.INodes[i].(*Directory).DeleteDirectoryKey()
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				d.updateDirectorySize(path)
				return operations, nil
			}
		default:
		}
	}
	return nil, errors.New("Path doesn't exists: " + path + name + "/")
}

// Store the file into the path.
func (d *Directory) CreateFile(path string, name string, size uint, hash string, keyHash string, fragments []*Fragment) error {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(dir.INodes); i++ {
		if dir.INodes[i].GetName() == name {
			return errors.New("The same Name file or directory exists: " + path + name)
		}
	}
	dir.INodes = append(dir.INodes, NewFile(name, size, hash, keyHash, fragments))
	d.updateDirectorySize(path)
	return nil
}

// Update the data of file finding by the filename and the path of file.
func (d *Directory) UpdateFileData(path string, name string, size uint, hash string, fragments []*Fragment) error {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return err
	}
	file.Size = size
	file.Hash = hash
	file.Fragments = fragments
	return nil
}

// Update the Key of file
func (d *Directory) UpdateFileKey(path string, name string, keyHash string, hash string, fragments []*Fragment) (operations map[string]int, err error) {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return operations, err
	}
	operations[file.KeyIndex]--
	file.KeyIndex = keyHash
	file.Hash = hash
	file.Fragments = fragments
	operations[keyHash]++
	return operations, nil
}

// Delete the file finding by the Name under the path.
func (d *Directory) DeleteFile(path string, name string) (string, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *File:
			file := dir.INodes[i].(*File)
			if file.GetName() == name {
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				d.updateDirectorySize(path)
				return file.KeyIndex, nil
			}
		default:
		}
	}
	return "", errors.New("File doesn't exists: " + path + name)
}

func (d Directory) AddSea(path string, name string, hash string, sea *FragmentSea) error {
	file, err := d.checkFileExists(path, name)
	if err != nil {
		return err
	}
	for _, fragment := range file.Fragments {
		if fragment.Hash == hash {
			for _, s := range fragment.Seas {
				if s.PublicKey == sea.PublicKey {
					return errors.New("fragment stored")
				}
			}
			fragment.Seas = append(fragment.Seas, sea)
			return nil
		}
	}
	return errors.New("fragment is not valid")
}

// List information of INodes in the path.
func (d *Directory) List(path string) ([]INodeInfo, error) {
	dir, err := d.checkPathExists(path)
	if err != nil {
		return nil, err
	}
	return generateINodeInfos(dir.INodes), nil
}

func (d *Directory) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(d)
	return buf.Bytes()
}

func DirectoryFromBytes(data []byte) (*Directory, error) {
	d := &Directory{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(d)
	return d, err
}
