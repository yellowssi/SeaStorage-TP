package crypto

import "testing"

func TestHash(t *testing.T) {
	test := "SeaStorage"
	hash := Hash(test)
	println(string(hash))
}

func TestAddress(t *testing.T) {
	test := "SeaStorage"
	hash := SHA512([]byte(test))
	address := Address(hash)
	println(hash)
	println(address)
}
