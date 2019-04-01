package seaStoragePayload

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

const _ = proto.ProtoPackageIsVersion3

type PayloadType uint16

var (
	PayloadTypeUnset            PayloadType = 0
	PayloadTypeCreateUser       PayloadType = 1
	PayloadTypeCreateGroup      PayloadType = 2
	PayloadTypeCreateSea        PayloadType = 3
	PayloadTypeSearchSharedFile PayloadType = 4
)

var (
	PayloadTypeUserCreateFile          PayloadType = 100
	PayloadTypeUserCreateDirectory     PayloadType = 101
	PayloadTypeUserUpdateFile          PayloadType = 102
	PayloadTypeUserShareFiles          PayloadType = 103
	PayloadTypeUserPublicKey           PayloadType = 104
	PayloadTypeUserListDirectory       PayloadType = 105
	PayloadTypeUserGetFileInfo         PayloadType = 106
	PayloadTypeUserListSharedDirectory PayloadType = 107
	PayloadTypeUserGetSharedFileInfo   PayloadType = 108
)

var (
	PayloadTypeGroupCreateFile          PayloadType = 200
	PayloadTypeGroupCreateDirectory     PayloadType = 201
	PayloadTypeGroupUpdateFile          PayloadType = 202
	PayloadTypeGroupShareFiles          PayloadType = 203
	PayloadTypeGroupPublicKey           PayloadType = 204
	PayloadTypeGroupListDirectory       PayloadType = 205
	PayloadTypeGroupGetFileInfo         PayloadType = 206
	PayloadTypeGroupListSharedDirectory PayloadType = 207
	PayloadTypeGroupGetSharedFileInfo   PayloadType = 208
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
	Create   string // default: ""
	FileInfo storage.FileInfo
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
