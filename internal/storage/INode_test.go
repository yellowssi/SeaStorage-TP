package storage

import (
	"testing"
)

func Test_ValidName(t *testing.T) {
	var err error
	err = ValidName("abcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()-_=+[]{}|;:'\",.<>?")
	if err != nil {
		t.Fail()
	}
	err = ValidName("te/st")
	if err == nil {
		t.Fail()
	}
}

func Test_ValidPath(t *testing.T) {
	var err error
	err = ValidPath("/etc/storage/SeaStorage/test/")
	if err != nil {
		t.Fail()
	}
	err = ValidPath("/")
	if err != nil {
		t.Fail()
	}
	err = ValidPath("/test")
	if err == nil {
		t.Fail()
	}
	err = ValidPath("test/")
	if err == nil {
		t.Fail()
	}
	err = ValidPath("/test//test/")
	if err == nil {
		t.Fail()
	}
}

func Test_ValidFile(t *testing.T) {
	err := ValidFile("/test/", "test", make([]*Fragment, 3))
	if err != nil {
		t.Error(err)
	}
}
