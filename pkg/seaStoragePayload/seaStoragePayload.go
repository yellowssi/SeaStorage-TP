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
	PayloadTypeUserUploadFile          PayloadType = 102
	PayloadTypeUserUploadDirectory     PayloadType = 103
	PayloadTypeUserUpdateFile          PayloadType = 104
	PayloadTypeUserShareFiles          PayloadType = 105
	PayloadTypeUserPublicKey           PayloadType = 106
	PayloadTypeUserListDirectory       PayloadType = 107
	PayloadTypeUserGetFileInfo         PayloadType = 108
	PayloadTypeUserListSharedDirectory PayloadType = 109
	PayloadTypeUserGetSharedFileInfo   PayloadType = 110
)

var (
	PayloadTypeGroupCreateFile          PayloadType = 200
	PayloadTypeGroupCreateDirectory     PayloadType = 201
	PayloadTypeGroupUploadFile          PayloadType = 202
	PayloadTypeGroupUploadDirectory     PayloadType = 203
	PayloadTypeGroupUpdateFile          PayloadType = 204
	PayloadTypeGroupShareFiles          PayloadType = 205
	PayloadTypeGroupPublicKey           PayloadType = 206
	PayloadTypeGroupListDirectory       PayloadType = 207
	PayloadTypeGroupGetFileInfo         PayloadType = 208
	PayloadTypeGroupListSharedDirectory PayloadType = 209
	PayloadTypeGroupGetSharedFileInfo   PayloadType = 210
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
