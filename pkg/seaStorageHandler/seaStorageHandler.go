package seaStorageHandler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/seaStoragePayload"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/seaStorageState"
)

type SeaStorageHandler struct {
	Name    string
	Version []string
}

func (h *SeaStorageHandler) FamilyName() string {
	return h.Name
}

func (h *SeaStorageHandler) FamilyVersion() []string {
	return h.Version
}

func (h *SeaStorageHandler) FamilyNamespaces() []string {
	return []string{string(seaStorageState.Namespace)}
}

func (h *SeaStorageHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	header := request.GetHeader()
	user := crypto.Address(header.GetSignerPublicKey())
	payload, err := seaStoragePayload.FromBytes(request.GetPayload())
	if err != nil {
		return err
	}
	state := seaStorageState.NewSeaStorageState(context)
	switch payload.Action {
	// Base Action
	case seaStoragePayload.PayloadTypeCreateUser:
		return state.CreateUser(payload.Target, user)
	case seaStoragePayload.PayloadTypeCreateGroup:
		return state.CreateGroup(payload.Target, seaStorageState.MakeAddress(seaStorageState.AddressTypeUser, payload.Name, user), payload.Key)
	case seaStoragePayload.PayloadTypeCreateSea:
		return state.CreateSea(payload.Target, user)
	case seaStoragePayload.PayloadTypeSearchSharedFile:
	case seaStoragePayload.PayloadTypeGetSharedFileInfo:

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
	case seaStoragePayload.PayloadTypeUserListDirectory:
	case seaStoragePayload.PayloadTypeUserGetFileInfo:
	// TODO: User Join Group & Search Group

	// Group Action
	//case seaStoragePayload.PayloadTypeGroupCreateDirectory:
	//case seaStoragePayload.PayloadTypeGroupCreateFile:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileName:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileData:
	//case seaStoragePayload.PayloadTypeGroupUpdateFileKey:
	//case seaStoragePayload.PayloadTypeGroupPublicKey:
	//case seaStoragePayload.PayloadTypeGroupListDirectory:
	//case seaStoragePayload.PayloadTypeGroupGetFileInfo:
	// TODO: Invite User & Access User Join Group & Leave Member

	// Sea Action
	case seaStoragePayload.PayloadTypeSeaCheckStatus:
	case seaStoragePayload.PayloadTypeSeaSetStatus:
	case seaStoragePayload.PayloadTypeSeaSetSpace:
	case seaStoragePayload.PayloadTypeSeaStoreFile:
	case seaStoragePayload.PayloadTypeSeaUpdateFile:

	default:
		return &processor.InvalidTransactionError{Msg: fmt.Sprint("Invalid Action: ", payload.Action)}
	}
}
