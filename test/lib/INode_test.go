package lib

import (
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/internal/lib"
	"testing"
)

func Test_ValidName(t *testing.T) {
	var err error
	err = lib.ValidName("abcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()-_=+[]{}|;:'\",.<>?")
	if err != nil {
		t.Fail()
	}
	err = lib.ValidName("te/st")
	if err == nil {
		t.Fail()
	}
}

func Test_ValidPath(t *testing.T) {
	var err error
	err = lib.ValidPath("/etc/lib/SeaStorage/test/")
	if err != nil {
		t.Fail()
	}
	err = lib.ValidPath("/")
	if err != nil {
		t.Fail()
	}
	err = lib.ValidPath("/test")
	if err == nil {
		t.Fail()
	}
	err = lib.ValidPath("test/")
	if err == nil {
		t.Fail()
	}
	err = lib.ValidPath("/test//test/")
	if err == nil {
		t.Fail()
	}
}

func Test_ValidFile(t *testing.T) {
	var err error
	err, _ = lib.NewFile("/test/", "test", 1, make([]byte, 512), make([]*lib.Fragment, 3))
	if err != nil {
		t.Error(err)
	}
}
