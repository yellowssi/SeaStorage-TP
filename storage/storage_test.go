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

func TestRoot_CreateDirectory(t *testing.T) {
	root := NewDirectory("root")
	_, err := root.CreateDirectory("/home/seastorage/")
	if err != nil {
		t.Error(err)
	}
	_, err = root.CreateDirectory("/lib/")
	if err != nil {
		t.Error(err)
	}
	t.Log(root)
	t.Log(root.List("/"))
}
