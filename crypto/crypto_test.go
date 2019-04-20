package crypto

import (
	"testing"
)

func TestSHA512(t *testing.T) {
	hash := SHA512HexFromBytes([]byte("SeaStorage"))
	println(hash)
	println(hash[:64])
}
