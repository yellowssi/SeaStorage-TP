// Copyright © 2019 yellowsea <hh1271941291@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/mitchellh/copystructure"
	"github.com/yellowssi/SeaStorage-TP/crypto"
	"github.com/yellowssi/SeaStorage-TP/sea"
	"strings"
)

func init() {
	gob.Register(&File{})
	gob.Register(&Directory{})
}

// Root store information of files and Keys used to encryption.
// Store the information of private files in 'Home' directory.
// Store the information of shared files in 'Shared' directory.
type Root struct {
	Home   *Directory
	Shared *Directory
	Keys   *FileKeyMap
}

// FileInfo is the information of files for usage.
type FileInfo struct {
	Name      string
	Size      int64
	Hash      string
	Key       string
	Fragments []*Fragment
}

// NewRoot is the construct for Root.
func NewRoot(home, share *Directory, keyMap *FileKeyMap) *Root {
	return &Root{
		Home:   home,
		Shared: share,
		Keys:   keyMap,
	}
}

// NewFileInfo is the construct for FileInfo.
func NewFileInfo(name string, size int64, hash string, key string, fragments []*Fragment) *FileInfo {
	return &FileInfo{
		Name:      name,
		Size:      size,
		Hash:      hash,
		Key:       key,
		Fragments: fragments,
	}
}

// GenerateRoot generate new root for usage.
func GenerateRoot() *Root {
	return NewRoot(NewDirectory("home"), NewDirectory("shared"), NewFileKeyMap())
}

// Check the path whether valid.
// Valid Path
// (1) start and end with '/'
// (2) not contain '//'
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
// Valid Name shouldn't contain '/'
func validName(name string) error {
	if len(name) == 0 {
		return errors.New("Name shouldn't be nil: " + name)
	}
	if strings.Contains(name, "/") {
		return errors.New("Name shouldn't contain '/': " + name)
	}
	return nil
}

// Check both path and name valid.
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

// CreateFile generate file in the path and store its information.
func (root *Root) CreateFile(p string, info FileInfo) error {
	err := validInfo(p, info.Name)
	if err != nil {
		return err
	}
	fileKeyIndex := root.Keys.AddKey(info.Key, true)
	err = root.Home.CreateFile(p, info.Name, info.Hash, fileKeyIndex, info.Size, info.Fragments)
	if err != nil {
		return err
	}
	root.Home.updateDirectorySize(p)
	return nil
}

// UpdateName change the target iNode's name to new name.
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

// UpdateFileData change the information of file.
func (root *Root) UpdateFileData(p string, info FileInfo, userOrGroup bool) (map[string][]*sea.Operation, error) {
	err := validInfo(p, info.Name)
	if err != nil {
		return nil, err
	}
	return root.Home.UpdateFileData(p, info.Name, info.Hash, info.Size, info.Fragments, userOrGroup, false)
}

// UpdateFileKey change the encryption key of file and its information.
func (root *Root) UpdateFileKey(p string, info FileInfo, userOrGroup bool) (map[string][]*sea.Operation, error) {
	err := validInfo(p, info.Name)
	if err != nil {
		return nil, err
	}
	root.Keys.AddKey(info.Key, false)
	keyUsed, seaOperations, err := root.Home.UpdateFileKey(p, info.Name, crypto.SHA512HexFromHex(info.Key), info.Hash, info.Size, info.Fragments, userOrGroup, false)
	if err != nil {
		return nil, err
	}
	root.Keys.UpdateKeyUsed(keyUsed)
	return seaOperations, nil
}

// PublishKey publish the key encrypted by public key.
func (root *Root) PublishKey(publicKey, key string) error {
	return root.Keys.PublishKey(publicKey, key)
}

// DeleteFile delete file in the path.
func (root *Root) DeleteFile(p, name string, userOrGroup bool) (map[string][]*sea.Operation, error) {
	err := validInfo(p, name)
	if err != nil {
		return nil, err
	}
	seaOperations, keyIndex, err := root.Home.DeleteFile(p, name, userOrGroup, false)
	if err != nil {
		return nil, err
	}
	root.Keys.UpdateKeyUsed(map[string]int{keyIndex: -1})
	root.Home.updateDirectorySize(p)
	return seaOperations, nil
}

// CreateDirectory create directory in the path.
func (root *Root) CreateDirectory(p string) error {
	err := validPath(p)
	if err != nil {
		return err
	}
	_, err = root.Home.CreateDirectory(p)
	return err
}

// DeleteDirectory delete directory and files in it.
func (root *Root) DeleteDirectory(p, name string, userOrGroup bool) (map[string][]*sea.Operation, error) {
	err := validInfo(p, name)
	if err != nil {
		return nil, err
	}
	seaOperations, keyUsed, err := root.Home.DeleteDirectory(p, name, userOrGroup, false)
	if err != nil {
		return nil, err
	}
	root.Keys.UpdateKeyUsed(keyUsed)
	root.Home.updateDirectorySize(p)
	return seaOperations, nil
}

// Move change the iNode parent path to new path.
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

// GetFile returns the information of file.
func (root *Root) GetFile(p, name string) (file FileInfo, err error) {
	err = validInfo(p, name)
	if err != nil {
		return
	}
	f, err := root.Home.checkFileExists(p, name)
	if err != nil {
		return
	}
	key := root.Keys.GetKey(f.KeyIndex)
	return *NewFileInfo(f.Name, f.Size, f.Hash, key.Key, f.Fragments), nil
}

// GetDirectory returns the information of directory.
func (root *Root) GetDirectory(p string) (dir *Directory, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Home.checkPathExists(p)
}

// GetINode returns the information of iNode.
func (root *Root) GetINode(p, name string) (INode, error) {
	return root.Home.checkINodeExists(p, name)
}

// ListDirectory list the information of iNodes in the directory.
func (root *Root) ListDirectory(p string) (iNodes []INodeInfo, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Home.List(p)
}

// GetSharedFile returns the information of file in the 'shared' directory.
func (root *Root) GetSharedFile(p, name string) (file FileInfo, err error) {
	err = validInfo(p, name)
	if err != nil {
		return
	}
	f, err := root.Shared.checkFileExists(p, name)
	if err != nil {
		return
	}
	key := root.Keys.GetKey(f.KeyIndex)
	return *NewFileInfo(f.Name, f.Size, f.Hash, key.Key, f.Fragments), nil
}

// GetSharedDirectory returns the information of directory in the 'shared' directory.
func (root *Root) GetSharedDirectory(p string) (dir *Directory, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Shared.checkPathExists(p)
}

// GetSharedINode returns the information of iNode in the 'shared' directory.
func (root *Root) GetSharedINode(p, name string) (INode, error) {
	return root.Shared.checkINodeExists(p, name)
}

// ListSharedDirectory list the information of iNodes in the 'shared' directory.
func (root *Root) ListSharedDirectory(p string) (iNodes []INodeInfo, err error) {
	err = validPath(p)
	if err != nil {
		return
	}
	return root.Shared.List(p)
}

// AddSea add fragment stored sea's information to its file.
func (root *Root) AddSea(p, name, hash string, sea *FragmentSea) error {
	return root.Home.AddSea(p, name, hash, sea)
}

// ShareFiles copy the information of file to 'shared' directory.
func (root *Root) ShareFiles(p, name, dst string, userOrGroup bool) (map[string][]*sea.Operation, []string, error) {
	iNode, err := root.GetINode(p, name)
	if err != nil {
		return nil, nil, err
	}
	target, err := copystructure.Copy(iNode)
	if err != nil {
		return nil, nil, err
	}
	var seaOperations map[string][]*sea.Operation
	if userOrGroup {
		seaOperations = iNode.GenerateSeaOperations(sea.ActionUserShared, true)
	} else {
		seaOperations = iNode.GenerateSeaOperations(sea.ActionGroupShared, true)
	}
	destination, _ := root.Shared.CreateDirectory(p)
	destination.INodes = append(destination.INodes, target.(INode))
	var keys = make([]string, 0)
	var keyUsed = make(map[string]int)
	keyIndexes := iNode.GetKeys()
	for _, keyIndex := range keyIndexes {
		fileKey := root.Keys.GetKey(keyIndex)
		keyUsed[keyIndex]--
		keys = append(keys, fileKey.Key)
	}
	root.Keys.UpdateKeyUsed(keyUsed)
	return seaOperations, keys, nil
}

// ToBytes convert root to byte slice.
func (root *Root) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(root)
	return buf.Bytes()
}

// RootFromBytes convert root from byte slice.
func RootFromBytes(data []byte) (*Root, error) {
	root := &Root{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(root)
	return root, err
}
