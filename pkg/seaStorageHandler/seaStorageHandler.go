package seaStorageHandler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"gitlab.com/SeaStorage/SeaStorage/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage/pkg/seaStoragePayload"
	"gitlab.com/SeaStorage/SeaStorage/pkg/seaStorageState"
)

var logger = logging.Get()

type SeaStorageHandler struct {
	Name    string
	Version []string
}

func NewSeaStorageHandler(version []string) *SeaStorageHandler {
	return &SeaStorageHandler{Name: "SeaStorage", Version: version}
}

func (h *SeaStorageHandler) FamilyName() string {
	return h.Name
}

func (h *SeaStorageHandler) FamilyVersions() []string {
	return h.Version
}

func (h *SeaStorageHandler) Namespaces() []string {
	return []string{string(seaStorageState.Namespace)}
}

func (h *SeaStorageHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	header := request.GetHeader()
	user := crypto.Address(header.GetSignerPublicKey())
	payload, err := seaStoragePayload.PayloadFromBytes(request.GetPayload())
	if err != nil {
		return err
	}
	state := seaStorageState.NewSeaStorageState(context)

	logger.Debugf("SeaStorage txn %v: user %v: payload: Name='%v', Action='%v'", request.Signature, user, payload.Name, payload.Action)

	switch payload.Action {
	// Base Action
	case seaStoragePayload.PayloadTypeCreateUser:
		return state.CreateUser(payload.Target, user)
	case seaStoragePayload.PayloadTypeCreateGroup:
		return state.CreateGroup(payload.Target, seaStorageState.MakeAddress(seaStorageState.AddressTypeUser, payload.Name, user), payload.Key)
	case seaStoragePayload.PayloadTypeCreateSea:
		return state.CreateSea(payload.Target, user)

	// User Action
	case seaStoragePayload.PayloadTypeUserCreateDirectory:
		return state.UserCreateDirectory(payload.Name, user, payload.PWD, payload.Target)
	case seaStoragePayload.PayloadTypeUserCreateFile:
		return state.UserCreateFile(payload.Name, user, payload.PWD, payload.FileInfo)
	case seaStoragePayload.PayloadTypeUserUpdateFileName:
		return state.UserUpdateFileName(payload.Name, user, payload.PWD, payload.Target, payload.Target2)
	case seaStoragePayload.PayloadTypeUserUpdateFileData:
		return state.UserUpdateFileData(payload.Name, user, payload.PWD, payload.FileInfo)
	case seaStoragePayload.PayloadTypeUserUpdateFileKey:
		return state.UserUpdateFileKey(payload.Name, user, payload.PWD, payload.FileInfo)
	case seaStoragePayload.PayloadTypeUserPublicKey:
		return state.UserPublicKey(payload.Name, user, payload.Key)
	// TODO: User Join Group & Search Group

	// Group Action
	//case seaStoragePayload.PayloadTypeGroupCreateDirectory:
	//case seaStoragePayload.PayloadTypeGroupCreateFile:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileName:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileData:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileKey:
	//case seaStoragePayload.PayloadTypeGroupPublicKey:
	// TODO: Invite User & Access User Join Group & Leave Member

	// Sea Action
	case seaStoragePayload.PayloadTypeSeaStoreFile:
		return state.SeaStoreFile(payload.Name, user, payload.Hash, payload.Signature)

	default:
		return &processor.InvalidTransactionError{Msg: fmt.Sprint("Invalid Action: ", payload.Action)}
	}
}
