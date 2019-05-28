package payload

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.com/SeaStorage/SeaStorage-TP/sea"
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"gitlab.com/SeaStorage/SeaStorage-TP/user"
)

const _ = proto.ProtoPackageIsVersion3

// Common action
var (
	Unset       uint = 0
	CreateUser  uint = 1
	CreateGroup uint = 2
	CreateSea   uint = 3
)

// User action
var (
	UserCreateFile      uint = 10
	UserCreateDirectory uint = 11
	UserDeleteFile      uint = 12
	UserDeleteDirectory uint = 13
	UserUpdateName      uint = 14
	UserUpdateFileData  uint = 15
	UserUpdateFileKey   uint = 16
	UserPublicKey       uint = 17
	UserMove            uint = 18
	UserShare           uint = 19
)

// Group action
var (
	GroupCreateFile      uint = 20
	GroupCreateDirectory uint = 21
	GroupDeleteFile      uint = 22
	GroupDeleteDirectory uint = 23
	GroupUpdateFileName  uint = 24
	GroupUpdateFileData  uint = 25
	GroupUpdateFileKey   uint = 26
	GroupPublicKey       uint = 27
)

// Sea Action
var (
	SeaStoreFile         uint = 30
	SeaConfirmOperations uint = 31
)

type SeaStoragePayload struct {
	Action         uint             `default:"Unset(0)"`
	Name           string           `default:""`
	PWD            string           `default:"/"`
	Target         []string         `default:"nil"`
	Key            string           `default:""`
	FileInfo       storage.FileInfo `default:"FileInfo{}"`
	UserOperations []user.Operation `default:"nil"`
	SeaOperations  []sea.Operation  `default:"nil"`
}

func NewSeaStoragePayload(action uint, name string, PWD string, target []string, key string, fileInfo storage.FileInfo, userOperations []user.Operation, seaOperations []sea.Operation) *SeaStoragePayload {
	return &SeaStoragePayload{
		Action:         action,
		Name:           name,
		PWD:            PWD,
		Target:         target,
		Key:            key,
		FileInfo:       fileInfo,
		UserOperations: userOperations,
		SeaOperations:  seaOperations,
	}
}

func SeaStoragePayloadFromBytes(payloadData []byte) (*SeaStoragePayload, error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}
	pl := &SeaStoragePayload{}
	buf := bytes.NewBuffer(payloadData)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(pl)
	return pl, err
}

func (ssp *SeaStoragePayload) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(ssp)
	return buf.Bytes()
}
