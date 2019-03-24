package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hash string
type Address string
type Key string

func SHA256(data []byte) Hash {
	hashHandler := sha256.New()
	hashHandler.Write(data)
	hashBytes := hashHandler.Sum(nil)
	return Hash(hex.EncodeToString(hashBytes))
}
