package payload

import (
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"gitlab.com/SeaStorage/SeaStorage-TP/user"
	"testing"
)

func TestSeaStoragePayload_ToBytes(t *testing.T) {
	pl := NewSeaStoragePayload(CreateUser, "Test", "", "", "", "", storage.FileInfo{}, "", user.OperationSignature{})
	t.Log(pl.ToBytes())
}

func TestSeaStoragePayloadFromBytes(t *testing.T) {
	pl := NewSeaStoragePayload(CreateUser, "Test", "", "", "", "", storage.FileInfo{}, "", user.OperationSignature{})
	data := pl.ToBytes()
	t.Log(pl)
	t.Log(SeaStoragePayloadFromBytes(data))
}
