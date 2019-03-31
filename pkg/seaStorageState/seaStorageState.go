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

var Namespace = crypto.SHA512([]byte("SeaStorage"))[:6]

type SeaStorageState struct {
	userContext  *processor.Context
	groupContext *processor.Context
	seaContext   *processor.Context
	userCache    map[crypto.Address][]byte
	groupCache   map[crypto.Address][]byte
	seaCache     map[crypto.Address][]byte
}

func NewSeaStorageState(userContext *processor.Context, groupContext *processor.Context, seaContext *processor.Context) *SeaStorageState {
	return &SeaStorageState{
		userContext:  userContext,
		groupContext: groupContext,
		seaContext:   seaContext,
		userCache:    make(map[crypto.Address][]byte),
		groupCache:   make(map[crypto.Address][]byte),
		seaCache:     make(map[crypto.Address][]byte),
	}
}

func (sss *SeaStorageState) GetUser(name string, publicKey string) (*user.User, error) {
	address := makeAddress(name, publicKey)
	userBytes, ok := sss.userCache[address]
	if ok {
		return deserializeUser(userBytes)
	}
	results, err := sss.userContext.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.userCache[address] = results[string(address)]
		return deserializeUser(results[string(address)])
	}
	return nil, errors.New("User doesn't exists. ")
}

func (sss *SeaStorageState) CreateUser(name string, publicKey string) error {
	address := makeAddress(name, publicKey)
	_, ok := sss.userCache[address]
	if ok {
		return errors.New("User exists. ")
	}
	results, err := sss.userContext.GetState([]string{string(address)})
	if err != nil {
		return err
	}
	if len(string(results[string(address)])) > 0 {
		return errors.New("User exists. ")
	}
	return sss.saveUser(user.NewUser(name), address)
}

func (sss *SeaStorageState) saveUser(u *user.User, address crypto.Address) error {
	uBytes, err := serialize(u)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to serialize account: ", err)}
	}
	sss.userCache[address] = uBytes
	addresses, err := sss.userContext.SetState(map[string][]byte{
		string(address): uBytes,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}
}

func (sss *SeaStorageState) GetGroup(name string, publicKey string) (*user.Group, error) {
	address := makeAddress(name, publicKey)
	groupBytes, ok := sss.groupCache[address]
	if ok {
		return deserializeGroup(groupBytes)
	}
	results, err := sss.groupContext.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.seaCache[address] = results[string(address)]
		return deserializeGroup(results[string(address)])
	}
	return nil, errors.New("Group doesn't exists. ")
}

func (sss *SeaStorageState) GetSea(name string, publicKey string) (*sea.Sea, error) {
	address := makeAddress(name, publicKey)
	seaBytes, ok := sss.seaCache[address]
	if ok {
		return deserializeSea(seaBytes)
	}
	results, err := sss.seaContext.GetState([]string{string(address)})
	if err != nil {
		return nil, err
	}
	if len(string(results[string(address)])) > 0 {
		sss.seaCache[address] = results[string(address)]
		return deserializeSea(results[string(address)])
	}
	return nil, errors.New("Sea doesn't exists. ")
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

func makeAddress(name string, publicKey string) crypto.Address {
	return crypto.Address(Namespace + crypto.SHA512(bytes.Join([][]byte{[]byte(name), []byte(publicKey)}, []byte{}))[:64])
}
