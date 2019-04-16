package handler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"gitlab.com/SeaStorage/SeaStorage/payload"
	"gitlab.com/SeaStorage/SeaStorage/state"
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
	return []string{string(state.Namespace)}
}

func (h *SeaStorageHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	header := request.GetHeader()
	user := header.GetSignerPublicKey()
	pl, err := payload.PayloadFromBytes(request.GetPayload())
	if err != nil {
		return err
	}
	st := state.NewSeaStorageState(context)

	logger.Debugf("SeaStorage txn %v: user %v: payload: Name='%v', Action='%v'", request.Signature, user, pl.Name, pl.Action)

	switch pl.Action {
	// Base Action
	case payload.PayloadTypeCreateUser:
		return st.CreateUser(pl.Target, user)
	case payload.PayloadTypeCreateGroup:
		return st.CreateGroup(pl.Target, state.MakeAddress(state.AddressTypeUser, pl.Name, user), pl.Key)
	case payload.PayloadTypeCreateSea:
		return st.CreateSea(pl.Target, user)

	// User Action
	case payload.PayloadTypeUserCreateDirectory:
		return st.UserCreateDirectory(pl.Name, user, pl.PWD, pl.Target)
	case payload.PayloadTypeUserCreateFile:
		return st.UserCreateFile(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.PayloadTypeUserUpdateName:
		return st.UserUpdateName(pl.Name, user, pl.PWD, pl.Target, pl.Target2)
	case payload.PayloadTypeUserUpdateFileData:
		return st.UserUpdateFileData(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.PayloadTypeUserUpdateFileKey:
		return st.UserUpdateFileKey(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.PayloadTypeUserPublicKey:
		return st.UserPublicKey(pl.Name, user, pl.Key)
	// TODO: User Join Group & Search Group

	// Group Action
	//case payload.PayloadTypeGroupCreateDirectory:
	//case payload.PayloadTypeGroupCreateFile:
	//case payload.PayloadTypeGroupUpdateFileName:
	//case payload.PayloadTypeGroupUpdateFileData:
	//case payload.PayloadTypeGroupUpdateFileKey:
	//case payload.PayloadTypeGroupPublicKey:
	// TODO: Invite User & Access User Join Group & Leave Member

	// Sea Action
	case payload.PayloadTypeSeaStoreFile:
		return st.SeaStoreFile(pl.Name, user, pl.Hash, pl.Signature)

	default:
		return &processor.InvalidTransactionError{Msg: fmt.Sprint("Invalid Action: ", pl.Action)}
	}
}
