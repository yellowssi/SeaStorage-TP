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

package state

import (
	"bytes"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/yellowssi/SeaStorage-TP/crypto"
	"github.com/yellowssi/SeaStorage-TP/sea"
	"github.com/yellowssi/SeaStorage-TP/storage"
	"github.com/yellowssi/SeaStorage-TP/user"
	"time"
)

type AddressType uint8

var (
	AddressTypeUser  AddressType = 0
	AddressTypeGroup AddressType = 1
	AddressTypeSea   AddressType = 2
)

var (
	Namespace      = crypto.SHA512HexFromBytes([]byte("SeaStorage"))[:6]
	UserNamespace  = crypto.SHA256HexFromBytes([]byte("User"))[:4]
	GroupNamespace = crypto.SHA256HexFromBytes([]byte("Group"))[:4]
	SeaNamespace   = crypto.SHA256HexFromBytes([]byte("Sea"))[:4]
)

type SeaStorageState struct {
	context    *processor.Context
	userCache  map[string][]byte
	groupCache map[string][]byte
	seaCache   map[string][]byte
}

func NewSeaStorageState(context *processor.Context) *SeaStorageState {
	return &SeaStorageState{
		context:    context,
		userCache:  make(map[string][]byte),
		groupCache: make(map[string][]byte),
		seaCache:   make(map[string][]byte),
	}
}

func (sss *SeaStorageState) GetUser(address string) (*user.User, error) {
	userBytes, ok := sss.userCache[address]
	if ok {
		return user.UserFromBytes(userBytes)
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	if len(results[address]) > 0 {
		sss.userCache[address] = results[address]
		return user.UserFromBytes(results[address])
	}
	return nil, &processor.InvalidTransactionError{Msg: "user doesn't exists"}
}

func (sss *SeaStorageState) CreateUser(username string, publicKey string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	_, ok := sss.userCache[address]
	if ok {
		return &processor.InvalidTransactionError{Msg: "user exists"}
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return err
	}
	if len(results[address]) > 0 {
		return &processor.InvalidTransactionError{Msg: "user exists"}
	}
	return sss.saveUser(user.GenerateUser(publicKey), address)
}

func (sss *SeaStorageState) saveUser(u *user.User, address string) error {
	uBytes := u.ToBytes()
	addresses, err := sss.context.SetState(map[string][]byte{
		address: uBytes,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}
	sss.userCache[address] = uBytes
	return nil
}

func (sss *SeaStorageState) GetGroup(address string) (*user.Group, error) {
	groupBytes, ok := sss.groupCache[address]
	if ok {
		return user.GroupFromBytes(groupBytes)
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	if len(results[address]) > 0 {
		sss.seaCache[address] = results[address]
		return user.GroupFromBytes(results[address])
	}
	return nil, &processor.InvalidTransactionError{Msg: "group doesn't exists"}
}

func (sss *SeaStorageState) CreateGroup(groupName, leader, key string) error {
	address := MakeAddress(AddressTypeGroup, groupName, "")
	_, ok := sss.groupCache[address]
	if ok {
		return &processor.InvalidTransactionError{Msg: "group exists"}
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return err
	}
	if len(results[address]) > 0 {
		return &processor.InvalidTransactionError{Msg: "group exists"}
	}
	return sss.saveGroup(user.GenerateGroup(groupName, leader), address)
}

func (sss *SeaStorageState) saveGroup(g *user.Group, address string) error {
	gBytes := g.ToBytes()
	addresses, err := sss.context.SetState(map[string][]byte{
		address: gBytes,
	})
	if err != nil {
		return err
	}
	if len(addresses) > 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}
	sss.groupCache[address] = gBytes
	return nil
}

func (sss *SeaStorageState) GetSea(address string) (*sea.Sea, error) {
	seaBytes, ok := sss.seaCache[address]
	if ok {
		return sea.SeaFromBytes(seaBytes)
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	if len(results[address]) > 0 {
		sss.seaCache[address] = results[address]
		return sea.SeaFromBytes(results[address])
	}
	return nil, &processor.InvalidTransactionError{Msg: "sea doesn't exists"}
}

func (sss *SeaStorageState) CreateSea(seaName, publicKey string) error {
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	_, ok := sss.seaCache[address]
	if ok {
		return &processor.InvalidTransactionError{Msg: "sea exists"}
	}
	results, err := sss.context.GetState([]string{address})
	if err != nil {
		return err
	}
	if len(results[address]) > 0 {
		return &processor.InvalidTransactionError{Msg: "sea exists"}
	}
	return sss.saveSea(sea.NewSea(publicKey), address)
}

func (sss *SeaStorageState) saveSea(s *sea.Sea, address string) error {
	sBytes := s.ToBytes()
	addresses, err := sss.context.SetState(map[string][]byte{
		address: sBytes,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response. "}
	}
	sss.seaCache[address] = sBytes
	return nil
}

func (sss *SeaStorageState) saveSeaOperations(address, publicKey string, data []byte, seaOperations map[string][]*sea.Operation) error {
	var err error
	seaCache := make(map[string]*sea.Sea)
	for seaAddr, operations := range seaOperations {
		s, ok := seaCache[seaAddr]
		if !ok {
			s, err = sss.GetSea(seaAddr)
			if err != nil {
				return err
			}
			seaCache[seaAddr] = s
		}
		for _, operation := range operations {
			operation.Owner = publicKey
		}
		s.AddOperation(operations)
	}
	cache := map[string][]byte{address: data}
	for addr, s := range seaCache {
		cache[addr] = s.ToBytes()
	}
	addresses, err := sss.context.SetState(cache)
	if err != nil {
		return err
	}
	if len(addresses) != len(cache) {
		return &processor.InternalError{Msg: "failed to store info"}
	}
	for addr, s := range seaCache {
		sss.seaCache[addr] = s.ToBytes()
	}
	sss.userCache[address] = data
	return nil
}

func (sss *SeaStorageState) UserShareFiles(username, publicKey, p, target, dst string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	seaOperations, _, err := u.Root.ShareFiles(p, target, dst, true)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveSeaOperations(address, u.PublicKey, u.ToBytes(), seaOperations)
}

func (sss *SeaStorageState) UserCreateDirectory(username, publicKey, p string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	err = u.Root.CreateDirectory(p)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserCreateFile(username, publicKey, p string, info storage.FileInfo) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	err = u.Root.CreateFile(p, info)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserDeleteDirectory(username, publicKey, p, target string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	seaOperations, err := u.Root.DeleteDirectory(p, target, true)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveSeaOperations(address, u.PublicKey, u.ToBytes(), seaOperations)
}

func (sss *SeaStorageState) UserDeleteFile(username, publicKey, p, target string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	seaOperations, err := u.Root.DeleteFile(p, target, true)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveSeaOperations(address, u.PublicKey, u.ToBytes(), seaOperations)
}

func (sss *SeaStorageState) UserMove(username, publicKey, p, name, newPath string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	err = u.Root.Move(p, name, newPath)
	if err != nil {
		return err
	}
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateName(username, publicKey, p, name, newName string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	err = u.Root.UpdateName(p, name, newName)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileData(username, publicKey, p string, info storage.FileInfo) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	seaOperations, err := u.Root.UpdateFileData(p, info, true)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveSeaOperations(address, u.PublicKey, u.ToBytes(), seaOperations)
}

func (sss *SeaStorageState) UserUpdateFileKey(username, publicKey, p string, info storage.FileInfo) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	seaOperations, err := u.Root.UpdateFileKey(p, info, true)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveSeaOperations(address, u.PublicKey, u.ToBytes(), seaOperations)
}

func (sss *SeaStorageState) UserPublishKey(username, publicKey, key string) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	u, err := sss.GetUser(address)
	if err != nil {
		return err
	}
	err = u.Root.PublishKey(publicKey, key)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) SeaStoreFile(seaName, publicKey string, operations []user.Operation) error {
	seaAddress := MakeAddress(AddressTypeSea, seaName, publicKey)
	s, err := sss.GetSea(seaAddress)
	if err != nil {
		return err
	}
	userCache := make(map[string]*user.User)
	for _, operation := range operations {
		if operation.Sea != publicKey {
			return &processor.InvalidTransactionError{Msg: "invalid operation"}
		}
		timestamp := time.Unix(operation.Timestamp, 0)
		if !operation.Verify() || timestamp.Before(time.Now()) {
			return &processor.InvalidTransactionError{Msg: "invalid operation"}
		}
		u, ok := userCache[operation.Address]
		if !ok {
			u, err = sss.GetUser(operation.Address)
			if err != nil {
				return err
			}
			userCache[operation.Address] = u
		}
		if !u.VerifyPublicKey(operation.PublicKey) {
			return &processor.InvalidTransactionError{Msg: "signature is invalid"}
		}
		err = u.Root.AddSea(operation.Path, operation.Name, operation.Hash, storage.NewFragmentSea(seaAddress, publicKey, timestamp))
		if err != nil {
			return &processor.InvalidTransactionError{Msg: err.Error()}
		}
		s.Handles++
	}
	cache := make(map[string][]byte)
	cache[seaAddress] = s.ToBytes()
	for address, u := range userCache {
		cache[address] = u.ToBytes()
	}
	addresses, err := sss.context.SetState(cache)
	if err != nil {
		return err
	}
	if len(addresses) != len(cache) {
		return &processor.InternalError{Msg: "failed to save data"}
	}
	for address, u := range userCache {
		sss.userCache[address] = u.ToBytes()
	}
	sss.seaCache[seaAddress] = s.ToBytes()
	return nil
}

func (sss *SeaStorageState) SeaConfirmOperations(seaName, publicKey string, operations []sea.Operation) error {
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	s, err := sss.GetSea(address)
	if err != nil {
		return err
	}
	s.RemoveOperations(operations)
	return sss.saveSea(s, address)
}

func MakeAddress(addressType AddressType, name, publicKey string) string {
	switch addressType {
	case AddressTypeUser:
		return Namespace + UserNamespace + crypto.SHA512HexFromBytes(bytes.Join([][]byte{[]byte(name), crypto.HexToBytes(publicKey)}, []byte{}))[:60]
	case AddressTypeGroup:
		return Namespace + GroupNamespace + crypto.SHA512HexFromBytes([]byte(name))[:60]
	case AddressTypeSea:
		return Namespace + SeaNamespace + crypto.SHA512HexFromBytes(bytes.Join([][]byte{[]byte(name), crypto.HexToBytes(publicKey)}, []byte{}))[:60]
	default:
		return ""
	}
}
