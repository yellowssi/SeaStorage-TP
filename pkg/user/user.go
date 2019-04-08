package user

import (
	"github.com/deckarep/golang-set"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

type User struct {
	Groups mapset.Set
	Root   *storage.Root
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
