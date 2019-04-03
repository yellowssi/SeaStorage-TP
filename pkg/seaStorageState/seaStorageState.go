package seaStorageState

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/sea"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/user"
)

type AddressType uint8

var (
	AddressTypeUser  AddressType = 0
	AddressTypeGroup AddressType = 1
	AddressTypeSea   AddressType = 2
)

var (
	Namespace      = crypto.SHA512([]byte("SeaStorage"))[:6]
	UserNamespace  = crypto.SHA512([]byte("user"))[:16]
	GroupNamespace = crypto.SHA512([]byte("group"))[:16]
	SeaNamespace   = crypto.SHA512([]byte("sea"))[:16]
)

type SeaStorageState struct {
	context    *processor.Context
	userCache  map[crypto.Address][]byte
	groupCache map[crypto.Address][]byte
	seaCache   map[crypto.Address][]byte
}

func NewSeaStorageState(context *processor.Context) *SeaStorageState {
	return &SeaStorageState{
		context:    context,
		userCache:  make(map[crypto.Address][]byte),
		groupCache: make(map[crypto.Address][]byte),
		seaCache:   make(map[crypto.Address][]byte),
	}
}

func (sss *SeaStorageState) GetUser(name string, publicKey crypto.Address) (*user.User, error) {
	address := MakeAddress(AddressTypeUser, name, publicKey)
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
	return nil, errors.New("User doesn't exists. ")
}

func (sss *SeaStorageState) CreateUser(name string, publicKey crypto.Address) error {
	address := MakeAddress(AddressTypeUser, name, publicKey)
	_, ok := sss.userCache[address]
	if ok {
		return errors.New("User exists. ")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(string(results[string(address)])) > 0 {
		return errors.New("User exists. ")
	}
	return sss.saveUser(user.GenerateUser(name), address)
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

func (sss *SeaStorageState) GetGroup(name string) (*user.Group, error) {
	address := MakeAddress(AddressTypeGroup, name, "")
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
	return nil, errors.New("Group doesn't exists. ")
}

func (sss *SeaStorageState) CreateGroup(name string, leader crypto.Address, key crypto.Key) error {
	address := MakeAddress(AddressTypeGroup, name, "")
	_, ok := sss.groupCache[address]
	if ok {
		return errors.New("Group exists. ")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(results[string(address)]) > 0 {
		return errors.New("Group exists. ")
	}
	return sss.saveGroup(user.GenerateGroup(name, leader), address)
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

func (sss *SeaStorageState) GetSea(name string, publicKey crypto.Address) (*sea.Sea, error) {
	address := MakeAddress(AddressTypeSea, name, publicKey)
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
	return nil, errors.New("Sea doesn't exists. ")
}

func (sss *SeaStorageState) CreateSea(name string, publicKey crypto.Address) error {
	address := MakeAddress(AddressTypeSea, name, publicKey)
	_, ok := sss.seaCache[address]
	if ok {
		return errors.New("Sea exists. ")
	}
	results, err := sss.context.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(string(results[string(address)])) > 0 {
		return errors.New("Sea exists. ")
	}
	return sss.saveSea(sea.NewSea(name), address)
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
		return crypto.Address(Namespace + UserNamespace + crypto.SHA384(bytes.Join([][]byte{[]byte(name), publicKey.ToBytes()}, []byte{})))
	case AddressTypeGroup:
		return crypto.Address(Namespace + GroupNamespace + crypto.SHA384([]byte(name)))
	case AddressTypeSea:
		return crypto.Address(Namespace + SeaNamespace + crypto.SHA384(bytes.Join([][]byte{[]byte(name), publicKey.ToBytes()}, []byte{})))
	}
	return crypto.Address("")
}
