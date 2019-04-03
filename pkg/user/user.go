package user

import (
	"github.com/deckarep/golang-set"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

type User struct {
	Name   string
	Groups mapset.Set
	Root   *storage.Root
}

type EncryptedUser struct {
	Name  string
	Group mapset.Set
	Root  *storage.EncryptedRoot
}

func NewUser(name string) *User {
	var groups mapset.Set
	return &User{Name: name, Groups: groups, Root: storage.NewRoot()}
}

func NewEncryptedUser(name string, group mapset.Set, root *storage.EncryptedRoot) *EncryptedUser {
	return &EncryptedUser{
		Name:  name,
		Group: group,
		Root:  root,
	}
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

func (u *User) Rename(name string) {
	u.Name = name
}
