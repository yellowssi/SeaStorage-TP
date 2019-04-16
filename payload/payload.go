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

type PayloadType uint8

var (
	PayloadTypeUnset       PayloadType = 0
	PayloadTypeCreateUser  PayloadType = 1
	PayloadTypeCreateGroup PayloadType = 2
	PayloadTypeCreateSea   PayloadType = 3
)

var (
	PayloadTypeUserCreateFile      PayloadType = 10
	PayloadTypeUserCreateDirectory PayloadType = 11
	PayloadTypeUserUpdateName      PayloadType = 12
	PayloadTypeUserUpdateFileData  PayloadType = 13
	PayloadTypeUserUpdateFileKey   PayloadType = 14
	PayloadTypeUserPublicKey       PayloadType = 15
)

var (
	PayloadTypeGroupCreateFile      PayloadType = 20
	PayloadTypeGroupCreateDirectory PayloadType = 21
	PayloadTypeGroupUpdateFileName  PayloadType = 22
	PayloadTypeGroupUpdateFileData  PayloadType = 23
	PayloadTypeGroupUpdateFileKey   PayloadType = 24
	PayloadTypeGroupPublicKey       PayloadType = 25
)

var (
	PayloadTypeSeaStoreFile  PayloadType = 30
	PayloadTypeSeaDeleteFile PayloadType = 31
)

type SeaStoragePayload struct {
	Action    PayloadType
	Name      string // default: ""
	PWD       string // default: "/"
	Target    string // default: ""
	Target2   string // default: ""
	Key       string
	FileInfo  storage.FileInfo
	Hash      string
	Signature user.OperationSignature
}

func NewSeaStoragePayload(action PayloadType, name string, PWD string, target string, target2 string, key string, fileInfo storage.FileInfo, hash string, signature user.OperationSignature) *SeaStoragePayload {
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

func PayloadFromBytes(payloadData []byte) (payload *SeaStoragePayload, err error) {
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
