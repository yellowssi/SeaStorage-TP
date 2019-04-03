package main

import (
	"encoding/gob"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

func init() {
	gob.Register(&storage.File{})
	gob.Register(&storage.Directory{})
}

func main() {
}
