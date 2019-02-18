package user

type Role bool

var (
	Role_Read  = true
	Role_Write = false
)

type Group struct {
	Name    string
	Leader  Address
	Members map[Address]Role
}
