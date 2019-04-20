package user

import (
	"bytes"
	"encoding/gob"
	"github.com/deckarep/golang-set"
	"gitlab.com/SeaStorage/SeaStorage/crypto"
	"gitlab.com/SeaStorage/SeaStorage/storage"
	"time"
)

type User struct {
	Groups mapset.Set
	Root   *storage.Root
}

type Operation struct {
	Owner     string
	PublicKey string
	Path      string
	Name      string
	Timestamp time.Time
}

type OperationSignature struct {
	Operation Operation
	Signature string
}

func NewUser(groups mapset.Set, root *storage.Root) *User {
	return &User{
		Groups: groups,
		Root:   root,
	}
}

func GenerateUser() *User {
	return NewUser(mapset.NewSet(), storage.GenerateRoot())
}

func NewOperation(owner string, publicKey string, path string, name string, timestamp time.Time) *Operation {
	return &Operation{
		Owner:     owner,
		PublicKey: publicKey,
		Path:      path,
		Name:      name,
		Timestamp: timestamp,
	}
}

func NewOperationSignature(operation Operation, privateKey string) (*OperationSignature, error) {
	operationBytes := operation.ToBytes()
	signature, err := crypto.Sign(privateKey, crypto.BytesToHex(operationBytes))
	if err != nil {
		return nil, err
	}
	return &OperationSignature{Operation: operation, Signature: crypto.BytesToHex(signature)}, nil
}

func (u *User) JoinGroup(group string) bool {
	if u.Groups.Contains(group) {
		return false
	} else {
		u.Groups.Add(group)
		return true
	}
}

func (u *User) LeaveGroup(group string) bool {
	if u.Groups.Contains(group) {
		u.Groups.Remove(group)
		return true
	}
	return false
}

func (u *User) IsInGroup(group string) bool {
	return u.Groups.Contains(group)
}

func (u *User) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(u)
	return buf.Bytes()
}

func UserFromBytes(data []byte) (*User, error) {
	buf := bytes.NewBuffer(data)
	u := &User{}
	dec := gob.NewDecoder(buf)
	err := dec.Decode(u)
	return u, err
}

func (o Operation) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(o)
	return buf.Bytes()
}

func (o Operation) ToHex() string {
	return crypto.BytesToHex(o.ToBytes())
}

func (ops OperationSignature) Verify() bool {
	return crypto.Verify(ops.Operation.PublicKey, ops.Signature, ops.Operation.ToHex())
}
