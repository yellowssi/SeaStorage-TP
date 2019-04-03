package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"testing"
)

func TestGobFile(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	gob.Register(&File{})
	gob.Register(&Directory{})

	fragments := make([]*Fragment, 0)
	file := NewFile("test", 100, crypto.Hash("test"), crypto.Hash("test"), fragments)

	err := enc.Encode(file)
	if err != nil {
		t.Error(err)
	}
	t.Log(buf)

	var f File
	err = dec.Decode(&f)
	if err != nil {
		t.Error(err)
	}
	t.Log(*file)
}

func TestGobDirectory(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	gob.Register(&File{})
	gob.Register(&Directory{})

	fragments := make([]*Fragment, 0)
	file := NewFile("testFile", 100, crypto.Hash("test"), crypto.Hash("test"), fragments)
	dir1 := NewDirectory("testDir1")
	dir2 := NewDirectory("testDir2")
	dir1.INodes = append(dir1.INodes, dir2, file)
	err := enc.Encode(dir1)
	if err != nil {
		t.Error(err)
	}
	t.Log(buf)

	//buf = *bytes.NewBuffer([]byte{73, 255, 139, 3, 1, 1, 9, 68, 105, 114, 101, 99, 116, 111, 114, 121, 1, 255, 140, 0, 1, 5, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 6, 0, 1, 4, 72, 97, 115, 104, 1, 12, 0, 1, 6, 73, 78, 111, 100, 101, 115, 1, 255, 142, 0, 1, 6, 83, 104, 97, 114, 101, 100, 1, 2, 0, 0, 0, 29, 255, 141, 2, 1, 1, 15, 91, 93, 115, 116, 111, 114, 97, 103, 101, 46, 73, 78, 111, 100, 101, 1, 255, 142, 0, 1, 16, 0, 0, 255, 145, 255, 140, 1, 8, 116, 101, 115, 116, 68, 105, 114, 49, 3, 2, 18, 42, 115, 116, 111, 114, 97, 103, 101, 46, 68, 105, 114, 101, 99, 116, 111, 114, 121, 255, 140, 11, 1, 8, 116, 101, 115, 116, 68, 105, 114, 50, 0, 13, 42, 115, 116, 111, 114, 97, 103, 101, 46, 70, 105, 108, 101, 255, 129, 3, 1, 1, 4, 70, 105, 108, 101, 1, 255, 130, 0, 1, 6, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 6, 0, 1, 4, 72, 97, 115, 104, 1, 12, 0, 1, 8, 75, 101, 121, 73, 110, 100, 101, 120, 1, 12, 0, 1, 9, 70, 114, 97, 103, 109, 101, 110, 116, 115, 1, 255, 138, 0, 1, 6, 83, 104, 97, 114, 101, 100, 1, 2, 0, 0, 0, 34, 255, 137, 2, 1, 1, 19, 91, 93, 42, 115, 116, 111, 114, 97, 103, 101, 46, 70, 114, 97, 103, 109, 101, 110, 116, 1, 255, 138, 0, 1, 255, 132, 0, 0, 31, 255, 131, 3, 1, 2, 255, 132, 0, 1, 2, 1, 4, 72, 97, 115, 104, 1, 12, 0, 1, 4, 83, 101, 97, 115, 1, 255, 136, 0, 0, 0, 37, 255, 135, 2, 1, 1, 22, 91, 93, 42, 115, 116, 111, 114, 97, 103, 101, 46, 70, 114, 97, 103, 109, 101, 110, 116, 83, 101, 97, 1, 255, 136, 0, 1, 255, 134, 0, 0, 35, 255, 133, 3, 1, 2, 255, 134, 0, 1, 2, 1, 7, 65, 100, 100, 114, 101, 115, 115, 1, 12, 0, 1, 6, 87, 101, 105, 103, 104, 116, 1, 4, 0, 0, 0, 29, 255, 130, 25, 1, 8, 116, 101, 115, 116, 70, 105, 108, 101, 1, 100, 1, 4, 116, 101, 115, 116, 1, 4, 116, 101, 115, 116, 0, 0})
	var d Directory
	err = dec.Decode(&d)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", d)
	t.Log(d)
	t.Log(d.INodes[1])
}
