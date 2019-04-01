package user

import "gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/crypto"

type Role uint8

var (
	RoleGuest      Role = 1
	RoleDeveloper  Role = 2
	RoleMaintainer Role = 3
	RoleOwner      Role = 4
)

type Group struct {
	name    string
	leader  crypto.Address
	members map[crypto.Address]Role
}

func NewGroup(name string, leader crypto.Address) *Group {
	return &Group{name: name, leader: leader, members: map[crypto.Address]Role{leader: RoleOwner}}
}

func (g *Group) Rename(user crypto.Address, name string) bool {
	if g.leader != user {
		return false
	}
	g.name = name
	return false
}

func (g *Group) UpdateLeader(user crypto.Address, newLeader crypto.Address) bool {
	if user == g.leader {
		g.leader = newLeader
		return true
	}
	return false
}

func (g *Group) UpdateMemberRole(user crypto.Address, member crypto.Address, role Role) bool {
	if g.members[user] != RoleOwner {
		return false
	} else if g.members[member] == RoleOwner && g.leader != user {
		return false
	}
	g.members[member] = role
	return true
}

func (g *Group) RemoveMember(user crypto.Address, member crypto.Address) bool {
	if g.members[user] != RoleOwner {
		return false
	} else if g.members[member] == RoleOwner && g.leader != user {
		return false
	}
	delete(g.members, member)
	return true
}
