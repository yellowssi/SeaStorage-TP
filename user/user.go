package user

import (
	"bytes"
	"encoding/gob"
	"gitlab.com/SeaStorage/SeaStorage/crypto"
	"gitlab.com/SeaStorage/SeaStorage/storage"
	"time"
)

type User struct {
	Groups []string
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

func NewUser(groups []string, root *storage.Root) *User {
	return &User{
		Groups: groups,
		Root:   root,
	}
}

func GenerateUser() *User {
	return NewUser(make([]string, 0), storage.GenerateRoot())
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
	for _, g := range u.Groups {
		if g == group {
			return false
		}
	}
	u.Groups = append(u.Groups, group)
	return true
}

func (u *User) LeaveGroup(group string) bool {
	for i, g := range u.Groups {
		if g == group {
			u.Groups = append(u.Groups[:i], u.Groups[i+1:]...)
			return true
		}
	}
	return false
}

func (u *User) IsInGroup(group string) bool {
	for _, g := range u.Groups {
		if g == group {
			return true
		}
	}
	return false
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
