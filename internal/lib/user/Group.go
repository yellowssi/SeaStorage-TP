package user

type Role uint8

var (
	RoleGuest      Role = 1
	RoleDeveloper  Role = 2
	RoleMaintainer Role = 3
	RoleOwner      Role = 4
)

type Group struct {
	Name    string
	Leader  Address
	Members map[Address]Role
}

func (g *Group) Rename(user Address, name string) bool {
	if g.Leader != user {
		return false
	}
	g.Name = name
	return false
}

func (g *Group) UpdateLeader(user Address, newLeader Address) bool {
	if user == g.Leader {
		g.Leader = newLeader
		return true
	}
	return false
}

func (g *Group) UpdateMemberRole(user Address, member Address, role Role) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	g.Members[member] = role
	return true
}

func (g *Group) RemoveMember(user Address, member Address) bool {
	if g.Members[user] != RoleOwner {
		return false
	} else if g.Members[member] == RoleOwner && g.Leader != user {
		return false
	}
	delete(g.Members, member)
	return true
}
