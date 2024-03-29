package storage

import (
	"github.com/yellowssi/SeaStorage-TP/sea"
	"testing"
	"time"
)

var root = GenerateRoot()

func TestValidName(t *testing.T) {
	var err error
	err = validName("abcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()-_=+[]{}|;:'\",.<>?")
	if err != nil {
		t.Fail()
	}
	err = validName("te/st")
	if err == nil {
		t.Fail()
	}
}

func TestValidPath(t *testing.T) {
	var err error
	err = validPath("/etc/storage/SeaStorage/test/")
	if err != nil {
		t.Fail()
	}
	err = validPath("/")
	if err != nil {
		t.Fail()
	}
	err = validPath("/test")
	if err == nil {
		t.Fail()
	}
	err = validPath("test/")
	if err == nil {
		t.Fail()
	}
	err = validPath("/test//test/")
	if err == nil {
		t.Fail()
	}
}

func TestRoot_CreateDirectory(t *testing.T) {
	err := root.CreateDirectory("/home/SeaStorage/")
	if err != nil {
		t.Error(err)
	}
	err = root.CreateDirectory("/lib/")
	if err != nil {
		t.Error(err)
	}
	t.Log(root.ListDirectory("/"))
}

func TestRoot_GetDirectory(t *testing.T) {
	dir, err := root.GetDirectory("/home/SeaStorage/")
	if err != nil {
		t.Error(err)
	}
	t.Log(dir.ToJson())
}

func TestRoot_CreateFile(t *testing.T) {
	err := root.CreateFile("/home/SeaStorage/", *NewFileInfo("test", 256, "hash", "key", []*Fragment{{Hash: "test", Size: 1}}))
	if err != nil {
		t.Error(err)
	}
	t.Log(root.Home.ToJson())
}

func TestRoot_AddSea(t *testing.T) {
	err := root.AddSea("/home/SeaStorage/", "test", "test", NewFragmentSea("address", "publicKey", time.Now()))
	if err != nil {
		t.Error(err)
	}
	t.Log(root.Home.ToJson())
}

func TestRoot_GenerateSeaOperations(t *testing.T) {
	t.Log(root.Home.GenerateSeaOperations(sea.ActionUserDelete, false))
}

func TestRoot_ShareFiles(t *testing.T) {
	operations, keys, err := root.ShareFiles("/home/SeaStorage/", "test", "/", true)
	if err != nil {
		t.Error(err)
	}
	t.Log(operations)
	t.Log(keys)
	t.Log(root.Shared.ToJson())
}

func TestRoot_DeleteFile(t *testing.T) {
	seaOperations, err := root.DeleteFile("/home/SeaStorage/", "test", true)
	if err != nil {
		t.Error(err)
	}
	t.Log(seaOperations)
	t.Log(root.Home.ToJson())
}

func TestToBytesAndFromBytes(t *testing.T) {
	data := root.Home.ToBytes()
	t.Log(data)
	test, err := DirectoryFromBytes(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(test.ToJson())
}

func TestRoot_Move(t *testing.T) {
	root.CreateDirectory("/test1/")
	root.CreateDirectory("/test2/")
	err := root.Move("/", "test1", "/test2/")
	if err != nil {
		t.Error(err)
	}
}

var data []byte

func TestRoot_ToBytes(t *testing.T) {
	data = root.ToBytes()
}

func TestRootFromBytes(t *testing.T) {
	test, err := RootFromBytes(data)
	if err != nil {
		t.Errorf("failed to unmarshal root: %v", err)
	} else {
		t.Log(test)
	}
}
