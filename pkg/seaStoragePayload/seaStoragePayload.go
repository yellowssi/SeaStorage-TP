package seaStoragePayload

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

const _ = proto.ProtoPackageIsVersion3

type PayloadType uint16

var (
	PayloadTypeUnset       PayloadType = 0
	PayloadTypeCreateUser  PayloadType = 1
	PayloadTypeCreateGroup PayloadType = 2
	PayloadTypeCreateSea   PayloadType = 3
)

var (
	PayloadTypeUserCreateFile      PayloadType = 100
	PayloadTypeUserCreateDirectory PayloadType = 101
	PayloadTypeUserUpdateFileName  PayloadType = 102
	PayloadTypeUserUpdateFileData  PayloadType = 103
	PayloadTypeUserUpdateFileKey   PayloadType = 104
	PayloadTypeUserPublicKey       PayloadType = 105
)

var (
	PayloadTypeGroupCreateFile      PayloadType = 200
	PayloadTypeGroupCreateDirectory PayloadType = 201
	PayloadTypeGroupUpdateFileName  PayloadType = 202
	PayloadTypeGroupUpdateFileData  PayloadType = 203
	PayloadTypeGroupUpdateFileKey   PayloadType = 204
	PayloadTypeGroupPublicKey       PayloadType = 205
)

var (
	PayloadTypeSeaSetStatus   PayloadType = 300
	PayloadTypeSeaSetSpace    PayloadType = 301
	PayloadTypeSeaStoreFile   PayloadType = 302
	PayloadTypeSeaUpdateFile  PayloadType = 303
	PayloadTypeSeaCheckStatus PayloadType = 304
)

type SeaStoragePayload struct {
	Action   PayloadType
	Name     string // default: ""
	PWD      string // default: "/"
	Target   string // default: ""
	Target2  string // default: ""
	Key      crypto.Key
	FileInfo storage.FileInfo
}

func NewSeaStoragePayload(action PayloadType, name string, pwd string, target string, target2 string, info storage.FileInfo) *SeaStoragePayload {
	return &SeaStoragePayload{
		Action:   action,
		Name:     name,
		PWD:      pwd,
		Target:   target,
		Target2:  target2,
		FileInfo: info,
	}
}

func FromBytes(payloadData []byte) (payload *SeaStoragePayload, err error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}
	buf := bytes.NewBuffer(payloadData)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(payload)
	return payload, err
}

func (ssp *SeaStoragePayload) ToBytes() (data []byte, err error) {
	buf := bytes.NewBuffer(data)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(ssp)
	return data, err
}
