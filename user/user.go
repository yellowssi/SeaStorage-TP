package user

import (
	"bytes"
	"encoding/gob"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"strconv"
	"time"
)

type User struct {
	PublicKey string
	Groups    []string
	Root      *storage.Root
}

func NewUser(publicKey string, groups []string, root *storage.Root) *User {
	return &User{
		PublicKey: publicKey,
		Groups:    groups,
		Root:      root,
	}
}

func GenerateUser(publicKey string) *User {
	return NewUser(publicKey, make([]string, 0), storage.GenerateRoot())
}

func (u *User) VerifyPublicKey(publicKey string) bool {
	return publicKey == u.PublicKey
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

type Operation struct {
	Address   string
	PublicKey string
	Path      string
	Name      string
	Size      int64
	Hash      string
	Timestamp int64
	Signature []byte
}

func NewOperation(address, publicKey, path, name, hash string, size int64, signer signing.Signer) *Operation {
	timestamp := time.Now().Unix()
	sign := signer.Sign(bytes.Join([][]byte{[]byte(address + publicKey + path + name + hash), []byte(strconv.Itoa(int(timestamp)))}, []byte{}))
	return &Operation{
		Address:   address,
		PublicKey: publicKey,
		Path:      path,
		Name:      name,
		Size:      size,
		Hash:      hash,
		Timestamp: timestamp,
		Signature: sign,
	}
}

func (o *Operation) Verify() bool {
	pub := signing.NewSecp256k1PublicKey(crypto.HexToBytes(o.PublicKey))
	cont := signing.NewSecp256k1Context()
	return cont.Verify(o.Signature, bytes.Join([][]byte{[]byte(o.Address + o.PublicKey + o.Path + o.Name + o.Hash), []byte(strconv.Itoa(int(o.Timestamp)))}, []byte{}), pub)
}

func (o *Operation) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(o)
	return buf.Bytes()
}

func OperationFromBytes(data []byte) (*Operation, error) {
	buf := bytes.NewBuffer(data)
	o := &Operation{}
	dec := gob.NewDecoder(buf)
	err := dec.Decode(o)
	return o, err
}
