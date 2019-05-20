package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type FileKey struct {
	Used uint
	Key  string
}

type INode interface {
	GetName() string
	GetSize() int64
	GetHash() string
	ToBytes() []byte
	ToJson() string
	lock()
	unlock()
}

type File struct {
	mutex     sync.Mutex
	Name      string
	Size      int64
	Hash      string
	KeyIndex  string
	Fragments []*Fragment
}

type Directory struct {
	mutex  sync.Mutex
	Name   string
	Size   int64
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
	Size  int64
}

func NewFileKey(key string) *FileKey {
	return &FileKey{Key: key, Used: 0}
}

func NewFile(name string, size int64, hash string, key string, fragments []*Fragment) *File {
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

func (f *File) lock() {
	f.mutex.Lock()
}

func (f *File) unlock() {
	f.mutex.Unlock()
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetSize() int64 {
	return f.Size
}

func (f *File) GetHash() string {
	return f.Hash
}

func (d *Directory) lock() {
	d.mutex.Lock()
}

func (d *Directory) unlock() {
	d.mutex.Unlock()
}

func (d *Directory) GetName() string {
	return d.Name
}

func (d *Directory) GetSize() int64 {
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
func (d *Directory) checkPathExists(p string) (*Directory, error) {
	pathParams := strings.Split(p, "/")
	dir := d
L:
	for i := 1; i < len(pathParams)-1; i++ {
		if len(dir.INodes) == 0 {
			return nil, errors.New("Path doesn't exists: " + strings.Join(pathParams[:i+1], "/") + "/")
		}
		for j := 0; j < len(dir.INodes); j++ {
			switch dir.INodes[j].(type) {
			case *Directory:
				if dir.INodes[j].GetName() == pathParams[i] {
					dir = dir.INodes[j].(*Directory)
					continue L
				}
			}
			if j == len(dir.INodes)-1 {
				return nil, errors.New("Path doesn't exists: " + strings.Join(pathParams[:i+1], "/") + "/")
			}
		}
	}
	return dir, nil
}

// Check the file whether exists in this Directory INode.
// If exists, return the pointer of the File INode.
// else, return the error.
func (d *Directory) checkFileExists(p string, name string) (*File, error) {
	dir, err := d.checkPathExists(p)
	if err != nil {
		return nil, err
	}
	for _, iNode := range dir.INodes {
		switch iNode.(type) {
		case *File:
			if iNode.GetName() == name {
				return iNode.(*File), nil
			}
		}
	}
	return nil, errors.New("File doesn't exists: " + p + name)
}

// Check the file or directory whether exists in this Directory INode.
func (d *Directory) checkINodeExists(p string, name string) (INode, error) {
	dir, err := d.checkPathExists(p)
	if err != nil {
		return nil, err
	}
	for _, iNode := range dir.INodes {
		if iNode.GetName() == name {
			return iNode, nil
		}
	}
	return nil, errors.New("File or Directory doesn't exists: " + p + name)
}

// Target directories recursively
// If there is the same Name file exists, it will return error.
// Else, it will return the pointer of the determination directory INode.
func (d *Directory) CreateDirectory(p string) (*Directory, error) {
	var newDir *Directory
	dir := d
	pathParams := strings.Split(p, "/")
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
				newDir = NewDirectory(pathParams[i])
				dir.INodes = append(dir.INodes, newDir)
				dir = newDir
				break
			}
		}
	}
	return dir, nil
}

// Update directories' Size in the path recursively.
func (d *Directory) updateDirectorySize(p string) {
	pathParams := strings.Split(p, "/")
	d.Size = 0
	d.lock()
	defer d.unlock()
	for i := 0; i < len(d.INodes); i++ {
		switch d.INodes[i].(type) {
		case *Directory:
			if d.INodes[i].GetName() == pathParams[1] {
				subPath := strings.Join(pathParams[2:], "/")
				subPath = "/" + subPath
				d.INodes[i].(*Directory).updateDirectorySize(subPath)
			}
			d.Size += d.INodes[i].GetSize()
		case *File:
			d.Size += d.INodes[i].GetSize()
		}
	}
}

// Update the Name of directory finding by the path.
func (d *Directory) UpdateName(p string, name string, newName string) error {
	iNode, err := d.checkINodeExists(p, name)
	if err != nil {
		return err
	}
	iNode.lock()
	defer iNode.unlock()
	switch iNode.(type) {
	case *File:
		iNode.(*File).Name = newName
	case *Directory:
		iNode.(*Directory).Name = newName
	}
	return nil
}

// Delete directory Key.
func (d *Directory) DeleteDirectoryKey() map[string]int {
	operations := make(map[string]int)
	for _, iNode := range d.INodes {
		switch iNode.(type) {
		case *Directory:
			for k, v := range iNode.(*Directory).DeleteDirectoryKey() {
				operations[k] -= v
			}
		case *File:
			file := iNode.(*File)
			operations[file.KeyIndex]--
		}
	}
	return operations
}

// Delete iNode of the directory finding by the path.
func (d *Directory) DeleteDirectory(p string, name string) (operations map[string]int, err error) {
	dir, err := d.checkPathExists(p)
	if err != nil {
		return nil, err
	}
	d.lock()
	defer d.unlock()
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *Directory:
			if dir.INodes[i].GetName() == name {
				operations = dir.INodes[i].(*Directory).DeleteDirectoryKey()
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				return operations, nil
			}
		}
	}
	return nil, errors.New("Path doesn't exists: " + p + name + "/")
}

// Store the file into the path.
func (d *Directory) CreateFile(p string, name string, size int64, hash string, keyHash string, fragments []*Fragment) error {
	dir, err := d.checkPathExists(p)
	if err != nil {
		return err
	}
	for i := 0; i < len(dir.INodes); i++ {
		if dir.INodes[i].GetName() == name {
			return errors.New("The same Name file or directory exists: " + p + name)
		}
	}
	d.lock()
	defer d.unlock()
	dir.INodes = append(dir.INodes, NewFile(name, size, hash, keyHash, fragments))
	return nil
}

// Update the data of file finding by the filename and the path of file.
func (d *Directory) UpdateFileData(p string, name string, size int64, hash string, fragments []*Fragment) error {
	file, err := d.checkFileExists(p, name)
	if err != nil {
		return err
	}
	d.lock()
	defer d.unlock()
	file.Size = size
	file.Hash = hash
	file.Fragments = fragments
	return nil
}

// Update the Key of file
func (d *Directory) UpdateFileKey(p string, name string, keyHash string, hash string, fragments []*Fragment) (operations map[string]int, err error) {
	file, err := d.checkFileExists(p, name)
	if err != nil {
		return operations, err
	}
	file.lock()
	defer file.unlock()
	operations[file.KeyIndex]--
	file.KeyIndex = keyHash
	file.Hash = hash
	file.Fragments = fragments
	operations[keyHash]++
	return operations, nil
}

// Delete the file finding by the Name under the path.
func (d *Directory) DeleteFile(p string, name string) (string, error) {
	dir, err := d.checkPathExists(p)
	if err != nil {
		return "", err
	}
	d.lock()
	defer d.unlock()
	for i := 0; i < len(dir.INodes); i++ {
		switch dir.INodes[i].(type) {
		case *File:
			file := dir.INodes[i].(*File)
			if file.GetName() == name {
				dir.INodes = append(dir.INodes[:i], dir.INodes[i+1:]...)
				return file.KeyIndex, nil
			}
		}
	}
	return "", errors.New("File doesn't exists: " + p + name)
}

func (d Directory) AddSea(p string, name string, hash string, sea *FragmentSea) error {
	file, err := d.checkFileExists(p, name)
	if err != nil {
		return err
	}
	file.lock()
	defer file.unlock()
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
func (d *Directory) List(p string) ([]INodeInfo, error) {
	dir, err := d.checkPathExists(p)
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

func (f *File) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(f)
	return buf.Bytes()
}

func FileFromBytes(data []byte) (*File, error) {
	f := &File{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(f)
	return f, err
}

func (d *Directory) ToJson() string {
	data, _ := json.MarshalIndent(d, "", "\t")
	return string(data)
}

func (f *File) ToJson() string {
	data, _ := json.MarshalIndent(f, "", "\t")
	return string(data)
}
