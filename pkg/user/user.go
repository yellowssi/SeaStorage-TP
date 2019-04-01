package user

import (
	"github.com/deckarep/golang-set"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
)

type User struct {
	Name   string
	Groups mapset.Set
}

func NewUser(name string) *User {
	var groups mapset.Set
	return &User{Name: name, Groups: groups}
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
