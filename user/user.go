// Copyright © 2019 yellowsea <hh1271941291@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
	"strconv"
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
	Sea       string
	Path      string
	Name      string
	Size      int64
	Hash      string
	Timestamp int64
	Signature string
}

func NewOperation(address, publicKey, sea, path, name, hash string, size, timestamp int64, signer signing.Signer) *Operation {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(size))
	sign := signer.Sign(bytes.Join([][]byte{[]byte(address + publicKey + sea + path + name + hash), buf, []byte(strconv.Itoa(int(timestamp)))}, []byte{}))
	return &Operation{
		Address:   address,
		PublicKey: publicKey,
		Sea:       sea,
		Path:      path,
		Name:      name,
		Size:      size,
		Hash:      hash,
		Timestamp: timestamp,
		Signature: crypto.BytesToHex(sign),
	}
}

func (o *Operation) Verify() bool {
	pub := signing.NewSecp256k1PublicKey(crypto.HexToBytes(o.PublicKey))
	cont := signing.NewSecp256k1Context()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(o.Size))
	return cont.Verify(crypto.HexToBytes(o.Signature), bytes.Join([][]byte{[]byte(o.Address + o.PublicKey + o.Sea + o.Path + o.Name + o.Hash), buf, []byte(strconv.Itoa(int(o.Timestamp)))}, []byte{}), pub)
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
