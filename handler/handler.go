package handler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"gitlab.com/SeaStorage/SeaStorage-TP/payload"
	"gitlab.com/SeaStorage/SeaStorage-TP/state"
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
	pl, err := payload.SeaStoragePayloadFromBytes(request.GetPayload())
	if err != nil {
		return err
	}
	st := state.NewSeaStorageState(context)

	logger.Debugf("SeaStorage txn %v: user %v: payload: Name='%v', Action='%v'", request.Signature, user, pl.Name, pl.Action)

	switch pl.Action {
	// Base Action
	case payload.CreateUser:
		return st.CreateUser(pl.Target, user)
	case payload.CreateGroup:
		return st.CreateGroup(pl.Target, state.MakeAddress(state.AddressTypeUser, pl.Name, user), pl.Key)
	case payload.CreateSea:
		return st.CreateSea(pl.Target, user)

	// User Action
	case payload.UserCreateDirectory:
		return st.UserCreateDirectory(pl.Name, user, pl.PWD, pl.Target)
	case payload.UserCreateFile:
		return st.UserCreateFile(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserUpdateName:
		return st.UserUpdateName(pl.Name, user, pl.PWD, pl.Target, pl.Target2)
	case payload.UserUpdateFileData:
		return st.UserUpdateFileData(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserUpdateFileKey:
		return st.UserUpdateFileKey(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserPublicKey:
		return st.UserPublicKey(pl.Name, user, pl.Key)
	// TODO: User Join Group & Search Group

	// Group Action
	//case payload.GroupCreateDirectory:
	//case payload.GroupCreateFile:
	//case payload.GroupUpdateFileName:
	//case payload.GroupUpdateFileData:
	//case payload.GroupUpdateFileKey:
	//case payload.GroupPublicKey:
	// TODO: Invite User & Access User Join Group & Leave Member

	// Sea Action
	case payload.SeaStoreFile:
		return st.SeaStoreFile(pl.Name, user, pl.Hash, pl.Signature)

	default:
		return &processor.InvalidTransactionError{Msg: fmt.Sprint("Invalid Action: ", pl.Action)}
	}
}
