package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	ellcurv "github.com/btcsuite/btcd/btcec"
)

type Address string
type Key string

func (address Address) GetBytes() []byte {
	addressBytes, _ := hex.DecodeString(string(address))
	return addressBytes
}

func GenerateRandomKey(len int) []byte {
	if len != 128 && len != 256 && len != 512 {
		panic(errors.New("AES Key length should be 128 or 256 or 512"))
	}
	key := make([]byte, len)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}
	return key
}

func NewKey(len int) Key {
	key := GenerateRandomKey(len)
	return Key(hex.EncodeToString(key))
}

func (k Key) GetBytes() []byte {
	keyBytes, err := hex.DecodeString(string(k))
	if err != nil {
		panic(err.Error())
	}
	return keyBytes
}

func (k Key) Encrypt(address Address) []byte {
	pub, err := ellcurv.ParsePubKey(address.GetBytes(), ellcurv.S256())
	result, err := ellcurv.Encrypt(pub, k.GetBytes())
	if err != nil {
		panic(err)
	}
	return result
}

func (k Key) Verify(address Address, key Key) bool {
	pub, err := ellcurv.ParsePubKey(address.GetBytes(), ellcurv.S256())
	if err != nil {
		panic(err.Error())
	}
	result, err := ellcurv.Encrypt(pub, key.GetBytes())
	if err != nil {
		panic(err.Error())
	}
	return bytes.Equal(result, k.GetBytes())
}
