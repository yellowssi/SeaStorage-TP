package payload

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.com/SeaStorage/SeaStorage/storage"
	"gitlab.com/SeaStorage/SeaStorage/user"
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
	UserUpdateName      uint = 12
	UserUpdateFileData  uint = 13
	UserUpdateFileKey   uint = 14
	UserPublicKey       uint = 15
)

// Group action
var (
	GroupCreateFile      uint = 20
	GroupCreateDirectory uint = 21
	GroupUpdateFileName  uint = 22
	GroupUpdateFileData  uint = 23
	GroupUpdateFileKey   uint = 24
	GroupPublicKey       uint = 25
)

// Sea Action
var (
	SeaStoreFile  uint = 30
	SeaDeleteFile uint = 31
)

type SeaStoragePayload struct {
	Action    uint                    `default:"Unset(0)"`
	Name      string                  `default:""`
	PWD       string                  `default:"/"`
	Target    string                  `default:""`
	Target2   string                  `default:""`
	Key       string                  `default:""`
	FileInfo  storage.FileInfo        `default:"FileInfo{}"`
	Hash      string                  `default:""`
	Signature user.OperationSignature `default:"OperationSignature{}"`
}

func NewSeaStoragePayload(action uint, name string, PWD string, target string, target2 string, key string, fileInfo storage.FileInfo, hash string, signature user.OperationSignature) *SeaStoragePayload {
	return &SeaStoragePayload{
		Action:    action,
		Name:      name,
		PWD:       PWD,
		Target:    target,
		Target2:   target2,
		Key:       key,
		FileInfo:  fileInfo,
		Hash:      hash,
		Signature: signature,
	}
}

func SeaStoragePayloadFromBytes(payloadData []byte) (payload *SeaStoragePayload, err error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}
	buf := bytes.NewBuffer(payloadData)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(payload)
	return payload, err
}

func (ssp *SeaStoragePayload) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(ssp)
	return buf.Bytes()
}
