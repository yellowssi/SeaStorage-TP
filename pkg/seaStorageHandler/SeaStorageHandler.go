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
		return state.CreateUser(payload.Create, user)
	case seaStoragePayload.PayloadTypeCreateGroup:
		return state.CreateGroup(payload.Create, seaStorageState.MakeAddress(seaStorageState.AddressTypeUser, payload.Name, user), payload.Key)
	case seaStoragePayload.PayloadTypeCreateSea:
		return state.CreateSea(payload.Create, user)
	case seaStoragePayload.PayloadTypeSearchSharedFile:

	// User Action
	case seaStoragePayload.PayloadTypeUserCreateDirectory:
		u, err := state.GetUser(payload.Name, user)
		if err != nil {
			return err
		}
		return u.Root.CreateDirectory(payload.PWD + "/" + payload.Create)
	case seaStoragePayload.PayloadTypeUserCreateFile:
		u, err := state.GetUser(payload.Name, user)
		if err != nil {
			return err
		}
		return u.Root.CreateFile(payload.PWD, payload.FileInfo)
	case seaStoragePayload.PayloadTypeUserUpdateFile:
	case seaStoragePayload.PayloadTypeUserShareFiles:
	case seaStoragePayload.PayloadTypeUserPublicKey:
	case seaStoragePayload.PayloadTypeUserListDirectory:
	case seaStoragePayload.PayloadTypeUserGetFileInfo:
	case seaStoragePayload.PayloadTypeUserListSharedDirectory:
	case seaStoragePayload.PayloadTypeUserGetSharedFileInfo:
	// TODO: User Join Group & Search Group

	// Group Action
	case seaStoragePayload.PayloadTypeGroupCreateDirectory:
	case seaStoragePayload.PayloadTypeGroupCreateFile:
	case seaStoragePayload.PayloadTypeGroupUpdateFile:
	case seaStoragePayload.PayloadTypeGroupShareFiles:
	case seaStoragePayload.PayloadTypeGroupPublicKey:
	case seaStoragePayload.PayloadTypeGroupListDirectory:
	case seaStoragePayload.PayloadTypeGroupGetFileInfo:
	case seaStoragePayload.PayloadTypeGroupListSharedDirectory:
	case seaStoragePayload.PayloadTypeGroupGetSharedFileInfo:
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
