package sea

import (
	"bytes"
	"encoding/gob"
)

var (
	ActionUserDelete  uint = 1
	ActionUserShared  uint = 2
	ActionGroupDelete uint = 3
	ActionGroupShared uint = 4
)

type Operation struct {
	Action uint   // delete or shared
	Owner  string // owner public key
	Hash   string // the hash of file or fragment
	Shared bool   // whether target is shared file or owner file
}

type Sea struct {
	PublicKey  string
	Handles    int
	Operations []Operation
}

func NewOperation(action uint, owner string, hash string, shared bool) *Operation {
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
		Operations: make([]Operation, 0),
	}
}

func (s *Sea) AddOperation(operations []*Operation) {
	for _, operation := range operations {
		s.Operations = append(s.Operations, *operation)
	}
}

func (s *Sea) RemoveOperations(operations []Operation) {
	for _, operation := range operations {
		for i, seaOperation := range s.Operations {
			if operation == seaOperation {
				s.Operations = append(s.Operations[:i], s.Operations[i+1:]...)
				break
			}
		}
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
