package sea

import (
	"bytes"
	"encoding/gob"
	"gitlab.com/SeaStorage/SeaStorage/pkg/crypto"
	"time"
)

var (
	OperationActionUpdate int8 = 1
	OperationActionDelete int8 = 2
	OperationActionShared int8 = 3
)

type Operation struct {
	Action int8
	Owner  crypto.Address
	Hash   crypto.Hash
	Shared bool
}

type Sea struct {
	Handles    int
	Operations []Operation
}

type Fragment struct {
	Timestamp time.Time
	Shared    bool
	Data      []byte
}

type Status struct {
	Name       string
	TotalSpace uint
	FreeSpace  uint
	Operations []Operation
	BasePath   string
}

func NewOperation(action int8, owner crypto.Address, hash crypto.Hash, shared bool) *Operation {
	return &Operation{
		Action: action,
		Owner:  owner,
		Hash:   hash,
		Shared: shared,
	}
}

func NewSea() *Sea {
	return &Sea{
		Handles:    0,
		Operations: make([]Operation, 0),
	}
}

func NewFragment(shared bool, data []byte) *Fragment {
	return &Fragment{
		Timestamp: time.Now(),
		Shared:    shared,
		Data:      data,
	}
}

func (f Fragment) ToBytes() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(f)
	return buf.Bytes(), err
}

func FragmentFromBytes(data []byte) (fragment Fragment, err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(fragment)
	return fragment, err
}
