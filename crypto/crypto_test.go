package crypto

import (
	"testing"
)

func TestSHA512(t *testing.T) {
	hash := SHA512Hex("SeaStorage")
	println(hash)
	println(hash[:64])
}
