// Copyright © 2019 yellowsea <hh1271941291@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/yellowssi/SeaStorage-TP/payload"
	"github.com/yellowssi/SeaStorage-TP/state"
)

var logger = logging.Get()

type SeaStorageHandler struct {
	Name    string
	Version []string
}

func NewSeaStorageHandler(name string, version []string) *SeaStorageHandler {
	return &SeaStorageHandler{Name: name, Version: version}
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

	logger.Debugf("SeaStorage txn %v: user %v: payload: Name='%v', Action='%v', Target='%v'", request.Signature, user, pl.Name, pl.Action, pl.Target)

	switch pl.Action {
	// Base Action
	case payload.CreateUser:
		if len(pl.Target) != 1 || pl.Target[0] == "" {
			return &processor.InvalidTransactionError{Msg: "username is nil"}
		}
		return st.CreateUser(pl.Target[0], user)
	case payload.CreateGroup:
		if len(pl.Target) != 1 || pl.Target[0] == "" {
			return &processor.InvalidTransactionError{Msg: "group name is nil"}
		}
		return st.CreateGroup(pl.Target[0], state.MakeAddress(state.AddressTypeUser, pl.Name, user), pl.Key)
	case payload.CreateSea:
		if len(pl.Target) != 1 || pl.Target[0] == "" {
			return &processor.InvalidTransactionError{Msg: "sea name is nil"}
		}
		return st.CreateSea(pl.Target[0], user)

	// User Action
	case payload.UserCreateDirectory:
		return st.UserCreateDirectory(pl.Name, user, pl.PWD)
	case payload.UserCreateFile:
		return st.UserCreateFile(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserDeleteDirectory:
		if len(pl.Target) != 1 || pl.Target[0] == "" {
			return &processor.InvalidTransactionError{Msg: "directory name is nil"}
		}
		return st.UserDeleteDirectory(pl.Name, user, pl.PWD, pl.Target[0])
	case payload.UserDeleteFile:
		if len(pl.Target) != 1 || pl.Target[0] == "" {
			return &processor.InvalidTransactionError{Msg: "filename is nil"}
		}
		return st.UserDeleteFile(pl.Name, user, pl.PWD, pl.Target[0])
	case payload.UserUpdateName:
		if len(pl.Target) != 2 || pl.Target[0] == "" || pl.Target[1] == "" {
			return &processor.InvalidTransactionError{Msg: "the name of file or directory is nil"}
		}
		return st.UserUpdateName(pl.Name, user, pl.PWD, pl.Target[0], pl.Target[1])
	case payload.UserUpdateFileData:
		return st.UserUpdateFileData(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserUpdateFileKey:
		return st.UserUpdateFileKey(pl.Name, user, pl.PWD, pl.FileInfo)
	case payload.UserPublishKey:
		return st.UserPublishKey(pl.Name, user, pl.Key)
	case payload.UserMove:
		if len(pl.Target) != 2 || pl.Target[0] == "" || pl.Target[1] == "" {
			return &processor.InvalidTransactionError{Msg: "the name of file or directory is nil"}
		}
		return st.UserMove(pl.Name, user, pl.PWD, pl.Target[0], pl.Target[1])
	case payload.UserShare:
		if len(pl.Target) != 2 || pl.Target[0] == "" || pl.Target[1] == "" {
			return &processor.InvalidTransactionError{Msg: "the name of file or directory is nil"}
		}
		return st.UserShareFiles(pl.Name, user, pl.PWD, pl.Target[0], pl.Target[1])
	// TODO: User Join Group & Search Group

	// Group Action
	//case payload.GroupCreateDirectory:
	//case payload.GroupCreateFile:
	//case payload.GroupDeleteDirectory:
	//case payload.GroupDeleteFile:
	//case payload.GroupUpdateFileName:
	//case payload.GroupUpdateFileData:
	//case payload.GroupUpdateFileKey:
	//case payload.GroupPublishKey:
	// TODO: Invite User & Access User Join Group & Leave Member

	// Sea Action
	case payload.SeaStoreFile:
		return st.SeaStoreFile(pl.Name, user, pl.UserOperations)
	case payload.SeaConfirmOperations:
		return st.SeaConfirmOperations(pl.Name, user, pl.SeaOperations)

	default:
		return &processor.InvalidTransactionError{Msg: fmt.Sprint("Invalid Action: ", pl.Action)}
	}
}
