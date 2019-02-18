package user

type Address string

type User struct {
	Name   string
	groups []Address
}
