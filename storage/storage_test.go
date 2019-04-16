package storage

import (
	"testing"
)

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

func TestGob(t *testing.T) {
	root := NewDirectory("root")
	_, _ = root.CreateDirectory("/test")
	_ = root.CreateFile("/", "testFile", 10, "test", "key", make([]*Fragment, 0))
	_, _ = root.CreateDirectory("/test/testDir")
	t.Log(root)
	data, err := root.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	test, err := DirectoryFromBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(test)
}
