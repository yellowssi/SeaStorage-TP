package user

import (
	"bytes"
	"encoding/gob"
	"gitlab.com/SeaStorage/SeaStorage-TP/storage"
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
	Leader  string
	Members map[string]Role
	Root    *storage.Root
}

func NewGroup(name string, leader string, members map[string]Role, root *storage.Root) *Group {
	return &Group{
		Name:    name,
		Leader:  leader,
		Members: members,
		Root:    root,
	}
}

func GenerateGroup(name string, leader string) *Group {
	return NewGroup(name, leader, map[string]Role{leader: RoleOwner}, storage.GenerateRoot())
}

func (g *Group) UpdateLeader(user string, newLeader string) bool {
	if user == g.Leader {
		g.Leader = newLeader
		return true
	}
	return false
}

func (g *Group) UpdateMemberRole(user string, member string, role Role) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	g.Members[member] = role
	return true
}

func (g *Group) RemoveMember(user string, member string) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	delete(g.Members, member)
	return true
}

func (g *Group) ToBytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(g)
	return buf.Bytes()
}

func GroupFromBytes(data []byte) (*Group, error) {
	g := &Group{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(g)
	return g, err
}
