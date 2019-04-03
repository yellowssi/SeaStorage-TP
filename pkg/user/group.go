package user

import (
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
)

type Role uint8

var (
	RoleGuest      Role = 1
	RoleDeveloper  Role = 2
	RoleMaintainer Role = 3
	RoleOwner      Role = 4
)

type Group struct {
	Name    string
	Leader  crypto.Address
	Members map[crypto.Address]Role
	Root    *storage.Root
}

func NewGroup(name string, leader crypto.Address, members map[crypto.Address]Role, root *storage.Root) *Group {
	return &Group{
		Name:    name,
		Leader:  leader,
		Members: members,
		Root:    root,
	}
}

func GenerateGroup(name string, leader crypto.Address) *Group {
	return NewGroup(name, leader, map[crypto.Address]Role{leader: RoleOwner}, storage.GenerateRoot())
}

func (g *Group) Rename(user crypto.Address, name string) bool {
	if g.Leader != user {
		return false
	}
	g.Name = name
	return false
}

func (g *Group) UpdateLeader(user crypto.Address, newLeader crypto.Address) bool {
	if user == g.Leader {
		g.Leader = newLeader
		return true
	}
	return false
}

func (g *Group) UpdateMemberRole(user crypto.Address, member crypto.Address, role Role) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	g.Members[member] = role
	return true
}

func (g *Group) RemoveMember(user crypto.Address, member crypto.Address) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	delete(g.Members, member)
	return true
}
