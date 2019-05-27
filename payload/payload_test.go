package payload

import (
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"testing"
)

func TestSeaStoragePayload_ToBytes(t *testing.T) {
	pl := NewSeaStoragePayload(CreateUser, "Test", "", "", "", "", storage.FileInfo{}, nil)
	t.Log(pl.ToBytes())
}

func TestSeaStoragePayloadFromBytes(t *testing.T) {
	pl := NewSeaStoragePayload(CreateUser, "Test", "", "", "", "", storage.FileInfo{}, nil)
	data := pl.ToBytes()
	t.Log(pl)
	t.Log(SeaStoragePayloadFromBytes(data))
}
