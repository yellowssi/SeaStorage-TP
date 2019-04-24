package state

import (
	"bytes"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/mitchellh/copystructure"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
	"gitlab.com/SeaStorage/SeaStorage-TP/sea"
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"gitlab.com/SeaStorage/SeaStorage-TP/user"
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
	Namespace       = crypto.SHA512HexFromBytes([]byte("SeaStorage"))[:6]
	UserNamespace   = crypto.SHA256HexFromBytes([]byte("User"))[:4]
	GroupNamespace  = crypto.SHA256HexFromBytes([]byte("Group"))[:4]
	SeaNamespace    = crypto.SHA256HexFromBytes([]byte("Sea"))[:4]
	SharedNamespace = crypto.SHA256HexFromBytes([]byte("Shared"))[:4]
)

type SeaStorageState struct {
	context     *processor.Context
	userCache   map[string][]byte
	groupCache  map[string][]byte
	seaCache    map[string][]byte
	sharedCache map[string][]byte
}

func NewSeaStorageState(context *processor.Context) *SeaStorageState {
	return &SeaStorageState{
		context:    context,
		userCache:  make(map[string][]byte),
		groupCache: make(map[string][]byte),
		seaCache:   make(map[string][]byte),
	}
}

func (sss *SeaStorageState) GetUser(username string, publicKey string) (*user.User, error) {
	address := MakeAddress(AddressTypeUser, username, publicKey)
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
	return sss.saveUser(user.GenerateUser(), address)
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

func (sss *SeaStorageState) GetGroup(groupName string) (*user.Group, error) {
	address := MakeAddress(AddressTypeGroup, groupName, "")
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

func (sss *SeaStorageState) CreateGroup(groupName string, leader string, key string) error {
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

func (sss *SeaStorageState) GetSea(seaName string, publicKey string) (*sea.Sea, error) {
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
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

func (sss *SeaStorageState) CreateSea(seaName string, publicKey string) error {
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
	return sss.saveSea(sea.NewSea(), address)
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

func (sss *SeaStorageState) UserShareFile(username string, publicKey string, path string, target string) error {
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

func (sss *SeaStorageState) saveSharedFiles(node storage.INode, address string) error {
	// TODO: Judge Shared Files Exists (Target / Update)
	nBytes := node.ToBytes()
	addresses, err := sss.context.SetState(map[string][]byte{
		address: nBytes,
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response. "}
	}
	sss.sharedCache[address] = nBytes
	// TODO: Add Event
	return nil
}

func (sss *SeaStorageState) UserCreateDirectory(username string, publicKey string, path string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.CreateDirectory(path)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserCreateFile(username string, publicKey string, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.CreateFile(path, info)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateName(username string, publicKey string, path string, name string, newName string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateName(path, name, newName)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileData(username string, publicKey string, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateFileData(path, info)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserUpdateFileKey(username string, publicKey string, path string, info storage.FileInfo) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.UpdateFileKey(path, info)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) UserPublicKey(username string, publicKey string, key string) error {
	u, err := sss.GetUser(username, publicKey)
	if err != nil {
		return err
	}
	err = u.Root.PublicKey(publicKey, key)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	address := MakeAddress(AddressTypeUser, username, publicKey)
	return sss.saveUser(u, address)
}

func (sss *SeaStorageState) SeaStoreFile(seaName string, publicKey string, hash string, sign user.OperationSignature) error {
	s, err := sss.GetSea(seaName, publicKey)
	if err != nil {
		return err
	}
	if !sign.Verify() || sign.Operation.Timestamp.Add(deadlineTime).Before(time.Now()) {
		return &processor.InvalidTransactionError{Msg: "signature is invalid"}
	}
	u, err := sss.GetUser(sign.Operation.Owner, sign.Operation.PublicKey)
	if err != nil {
		return err
	}
	err = u.Root.AddSea(sign.Operation.Path, sign.Operation.Name, hash, storage.NewFragmentSea(publicKey))
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	s.Handles++
	address := MakeAddress(AddressTypeSea, seaName, publicKey)
	return sss.saveSea(s, address)
}

func MakeAddress(addressType AddressType, name string, publicKey string) string {
	switch addressType {
	case AddressTypeUser:
		return Namespace + UserNamespace + crypto.SHA512HexFromBytes(bytes.Join([][]byte{[]byte(name), crypto.HexToBytes(publicKey)}, []byte{}))[:60]
	case AddressTypeGroup:
		return Namespace + GroupNamespace + crypto.SHA512HexFromBytes([]byte(name))[:60]
	case AddressTypeSea:
		return Namespace + SeaNamespace + crypto.SHA512HexFromBytes(bytes.Join([][]byte{[]byte(name), crypto.HexToBytes(publicKey)}, []byte{}))[:60]
	case AddressTypeShared:
		return Namespace + SharedNamespace + crypto.SHA512HexFromBytes(bytes.Join([][]byte{[]byte(name), crypto.HexToBytes(publicKey)}, []byte{}))[:60]
	default:
		return ""
	}
}
