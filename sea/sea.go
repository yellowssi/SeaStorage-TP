package sea

import (
	"bytes"
	"encoding/gob"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
	"time"
)

var (
	OperationActionUpdate int8 = 1
	OperationActionDelete int8 = 2
	OperationActionShared int8 = 3
)

type Operation struct {
	Action int8
	Owner  string
	Hash   string
	Shared bool
}

type Sea struct {
	PublicKey  string
	Handles    int
	Operations map[string]Operation
}

type Fragment struct {
	Timestamp time.Time
	Shared    bool
	Data      []byte
}

type Status struct {
	Name       string
	TotalSpace int
	FreeSpace  int
	Operations []Operation
	BasePath   string
}

func NewOperation(action int8, owner string, hash string, shared bool) *Operation {
	return &Operation{
		Action: action,
		Owner:  owner,
		Hash:   hash,
		Shared: shared,
	}
}

func NewSea(publicKey string) *Sea {
	return &Sea{
		PublicKey:  publicKey,
		Handles:    0,
		Operations: make(map[string]Operation),
	}
}

func (s *Sea) AddOperation(operations []Operation) {
	for _, operation := range operations {
		hash := crypto.SHA256HexFromBytes(operation.ToBytes())
		s.Operations[hash] = operation
	}
}

func (s *Sea) RemoveOperations(operations []Operation) {
	for _, operation := range operations {
		hash := crypto.SHA256HexFromBytes(operation.ToBytes())
		delete(s.Operations, hash)
	}
}

func (s *Sea) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(s)
	return buf.Bytes()
}

func SeaFromBytes(data []byte) (*Sea, error) {
	s := &Sea{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(s)
	return s, err
}

func NewFragment(shared bool, data []byte) *Fragment {
	return &Fragment{
		Timestamp: time.Now(),
		Shared:    shared,
		Data:      data,
	}
}

func (f Fragment) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(f)
	return buf.Bytes()
}

func FragmentFromBytes(data []byte) (Fragment, error) {
	fragment := Fragment{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&fragment)
	return fragment, err
}

func (o Operation) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(o)
	return buf.Bytes()
}

func OperationFromBytes(data []byte) (Operation, error) {
	operation := Operation{}
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&operation)
	return operation, err
}
