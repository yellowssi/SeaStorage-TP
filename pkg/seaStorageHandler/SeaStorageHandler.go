package seaStorageHandler

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
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
	return []string{seaStorageState.Namespace}
}

func (h *SeaStorageHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	header := request.GetHeader()
	user := header.GetSignerPublicKey()
}
