package storage

import (
	"testing"
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
	t.Log(dir)
}
