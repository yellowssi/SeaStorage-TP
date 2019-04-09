package seaStorageState

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/mitchellh/copystructure"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/sea"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/user"
	"time"
)

type AddressType uint8

var deadlineTime, _ = time.ParseDuration("3h")

var (
	AddressTypeUser   AddressType = 0
	AddressTypeGroup  AddressType = 1
	AddressTypeSea    AddressType = 2
	AddressTypeShared AddressType = 3
)

var (
	Namespace       = crypto.SHA512([]byte("SeaStorage"))[:6]
	UserNamespace   = crypto.Hash(AddressTypeUser)
	GroupNamespace  = crypto.Hash(AddressTypeGroup)
	SeaNamespace    = crypto.Hash(AddressTypeSea)
	SharedNamespace = crypto.Hash(AddressTypeShared)
)

type SeaStorageState struct {
	context     *processor.Context
	userCache   map[crypto.Address][]byte
	groupCache  map[crypto.Address][]byte
	seaCache    map[crypto.Address][]byte
	sharedCache map[crypto.Address][]byte
}

func NewSeaStorageState(context *processor.Context) *SeaStorageState {
	return &SeaStorageState{
		context:    context,
		userCache:  make(map[crypto.Address][]byte),
		groupCache: make(map[crypto.Address][]byte),
		seaCache:   make(map[crypto.Address][]byte),
	}
}

func (sss *SeaStorageState) GetUser(username string, publicKey crypto.Address) (*user.User, error) {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	userBytes, ok := sss.userCache[address]
	if ok {
		return deserializeUser(userBytes)
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.userCache[address] = results[string(address)]
		return deserializeUser(results[string(address)])
	}
	return nil, errors.New("user doesn't exists")
}

func (sss *SeaStorageState) CreateUser(username string, publicKey crypto.Address) error {
	address := MakeAddress(AddressTypeUser, username, publicKey)
	_, ok := sss.userCache[address]
	if ok {
		return errors.New("user exists")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(string(results[string(address)])) > 0 {
		return errors.New("user exists")
	}
	return sss.saveUser(user.GenerateUser(), address)
}

func (sss *SeaStorageState) saveUser(u *user.User, address crypto.Address) error {
	uBytes, err := serialize(u)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to serialize account: ", err)}
	}
	addresses, err := sss.context.SetState(map[string][]byte{
		string(address): uBytes,
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

func (sss *SeaStorageState) GetGroup(groupName string) (*user.Group, error) {
	address := MakeAddress(AddressTypeGroup, groupName, "")
	groupBytes, ok := sss.groupCache[address]
	if ok {
		return deserializeGroup(groupBytes)
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.seaCache[address] = results[string(address)]
		return deserializeGroup(results[string(address)])
	}
	return nil, errors.New("group doesn't exists")
}

func (sss *SeaStorageState) CreateGroup(groupName string, leader crypto.Address, key crypto.Key) error {
	address := MakeAddress(AddressTypeGroup, groupName, "")
	_, ok := sss.groupCache[address]
	if ok {
		return errors.New("group exists")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(results[string(address)]) > 0 {
		return errors.New("group exists")
	}
	return sss.saveGroup(user.GenerateGroup(groupName, leader), address)
}

func (sss *SeaStorageState) saveGroup(g *user.Group, address crypto.Address) error {
	gBytes, err := serialize(g)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to serialize group: ", err)}
	}
	addresses, err := sss.context.SetState(map[string][]byte{
		string(address): gBytes,
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

func (sss *SeaStorageState) GetSea(seaName string, publicKey crypto.Address) (*sea.Sea, error) {
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	seaBytes, ok := sss.seaCache[address]
	if ok {
		return deserializeSea(seaBytes)
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.seaCache[address] = results[string(address)]
		return deserializeSea(results[string(address)])
	}
	return nil, errors.New("sea doesn't exists")
}

func (sss *SeaStorageState) CreateSea(seaName string, publicKey crypto.Address) error {
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	_, ok := sss.seaCache[address]
	if ok {
		return errors.New("sea exists")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(string(results[string(address)])) > 0 {
		return errors.New("sea exists")
	}
	return sss.saveSea(sea.NewSea(), address)
}

func (sss *SeaStorageState) saveSea(s *sea.Sea, address crypto.Address) error {
	sBytes, err := serialize(s)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to serialize sea: ", err)}
	}
	addresses, err := sss.context.SetState(map[string][]byte{
		string(address): sBytes,
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

func (sss *SeaStorageState) UserShareFile(username string, publicKey crypto.Address, path string, target string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	iNode, err := u.Root.GetINode(path, target)
	if err != nil {
		return err
	}
	dst, err := copystructure.Copy(iNode)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeShared, username, publicKey)
	return sss.saveSharedFiles(dst.(storage.INode), address)
}

func (sss *SeaStorageState) saveSharedFiles(node storage.INode, address crypto.Address) error {
	// TODO: Judge Shared Files Exists (Target / Update)
	nBytes, err := serialize(node)
	if err != nil {
		return err
	}
	addresses, err := sss.context.SetState(map[string][]byte{
		string(address): nBytes,
	})
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response. "}
	}
	sss.sharedCache[address] = nBytes
	return nil
}

func (sss *SeaStorageState) UserCreateDirectory(username string, publicKey crypto.Address, path string, name string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.CreateDirectory(path + "/" + name)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserCreateFile(username string, publicKey crypto.Address, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.CreateFile(path, info)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileName(username string, publicKey crypto.Address, path string, name string, newName string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateFileName(path, name, newName)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileData(username string, publicKey crypto.Address, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateFileData(path, info)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileKey(username string, publicKey crypto.Address, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateFileKey(path, info)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserPublicKey(username string, publicKey crypto.Address, key crypto.Key) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.PublicKey(publicKey, key)
	if err != nil {
		return err
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) SeaStoreFile(seaName string, publicKey crypto.Address, hash crypto.Hash, sign user.OperationSignature) error {
	s, err := sss.GetSea(seaName, publicKey)
	if err != nil {
		return err
	}
	if !sign.Verify() || sign.Operation.Timestamp.Add(deadlineTime).Before(time.Now()) {
		return errors.New("signature is invalid")
	}
	u, err := sss.GetUser(sign.Operation.Owner, sign.Operation.PublicKey)
	if err != nil {
		return err
	}
	err = u.Root.AddSea(sign.Operation.Path, sign.Operation.Name, hash, storage.NewFragmentSea(publicKey))
	if err != nil {
		return err
	}
	s.Handles++
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	return sss.saveSea(s, address)
}

func serialize(i interface{}) (data []byte, err error) {
	buf := bytes.NewBuffer(data)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(i)
	return buf.Bytes(), err
}

func deserializeUser(data []byte) (user *user.User, err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(user)
	return user, err
}

func deserializeGroup(data []byte) (group *user.Group, err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(group)
	return group, err
}

func deserializeSea(data []byte) (sea *sea.Sea, err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(sea)
	return sea, err
}

func MakeAddress(addressType AddressType, name string, publicKey crypto.Address) crypto.Address {
	switch addressType {
	case AddressTypeUser:
		return crypto.Address(Namespace + UserNamespace + crypto.SHA512(bytes.Join([][]byte{[]byte(name), publicKey.ToBytes()}, []byte{}))[:63])
	case AddressTypeGroup:
		return crypto.Address(Namespace + GroupNamespace + crypto.SHA512([]byte(name))[:63])
	case AddressTypeSea:
		return crypto.Address(Namespace + SeaNamespace + crypto.SHA512(bytes.Join([][]byte{[]byte(name), publicKey.ToBytes()}, []byte{}))[:63])
	case AddressTypeShared:
		return crypto.Address(Namespace + SharedNamespace + crypto.SHA512(bytes.Join([][]byte{[]byte(name), publicKey.ToBytes()}, []byte{}))[:63])
	default:
		return crypto.Address("")
	}
}
