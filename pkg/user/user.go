package user

import (
	"bytes"
	"encoding/gob"
	"github.com/deckarep/golang-set"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
	"time"
)

type User struct {
	Groups mapset.Set
	Root   *storage.Root
}

type Operation struct {
	Owner     string
	PublicKey crypto.Address
	Path      string
	Name      string
	Timestamp time.Time
}

type OperationSignature struct {
	Operation Operation
	Signature []byte
}

func NewUser(groups mapset.Set, root *storage.Root) *User {
	return &User{
		Groups: groups,
		Root:   root,
	}
}

func GenerateUser() *User {
	var group mapset.Set
	return NewUser(group, storage.GenerateRoot())
}

func NewOperation(owner string, publicKey crypto.Address, path string, name string, timestamp time.Time) *Operation {
	return &Operation{
		Owner:     owner,
		PublicKey: publicKey,
		Path:      path,
		Name:      name,
		Timestamp: timestamp,
	}
}

func NewOperationSignature(operation Operation, privateKey []byte) (*OperationSignature, error) {
	operationBytes, err := operation.ToBytes()
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(privateKey, operationBytes)
	if err != nil {
		return nil, err
	}
	return &OperationSignature{Operation: operation, Signature: signature}, nil
}

func (u *User) JoinGroup(group crypto.Address) bool {
	if u.Groups.Contains(group) {
		return false
	} else {
		u.Groups.Add(group)
		return true
	}
}

func (u *User) LeaveGroup(group crypto.Address) bool {
	if u.Groups.Contains(group) {
		u.Groups.Remove(group)
		return true
	}
	return false
}

func (u *User) IsInGroup(group crypto.Address) bool {
	return u.Groups.Contains(group)
}

func (o Operation) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	return buf.Bytes(), err
}

func (ops OperationSignature) Verify() bool {
	operationBytes, err := ops.Operation.ToBytes()
	if err != nil {
		return false
	}
	return ops.Operation.PublicKey.Verify(ops.Signature, operationBytes)
}
