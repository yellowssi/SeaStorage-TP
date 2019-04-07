package seaStorageState

import (
	"github.com/mitchellh/copystructure"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
	"testing"
)

func TestSeaStorageState_GetSea(t *testing.T) {
	testMap := map[string][]byte{
		"a": []byte("apple"),
		"b": []byte("banana"),
		"c": []byte("core"),
		"d": nil,
	}
	test, ok := testMap["d"]
	if ok {
		println(test)
	} else {
		println("doesn't exists")
	}
}

func TestCopy(t *testing.T) {
	dir := storage.NewDirectory("test")
	_, err := dir.CreateDirectory("/test1/")
	if err != nil {
		t.Error(err)
	}
	test, err := copystructure.Copy(dir)
	if err != nil {
		t.Error(err)
	}
	t.Log(dir)
	t.Log(test.(*storage.Directory))
}
