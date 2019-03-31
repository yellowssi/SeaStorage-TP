package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	ellcurv "github.com/btcsuite/btcd/btcec"
)

type Hash string
type Address string
type Key string

func SHA256(data []byte) Hash {
	hashHandler := sha256.New()
	hashHandler.Write(data)
	hashBytes := hashHandler.Sum(nil)
	return HashFromBytes(hashBytes)
}

func SHA512(data []byte) Hash {
	hashHandler := sha512.New()
	hashHandler.Write(data)
	hashBytes := hashHandler.Sum(nil)
	return HashFromBytes(hashBytes)
}

func AddressFromBytes(addressBytes []byte) Address {
	return Address(hex.EncodeToString(addressBytes))
}

func (address Address) ToBytes() []byte {
	addressBytes, _ := hex.DecodeString(string(address))
	return addressBytes
}

func GenerateRandomKey(len int) []byte {
	if len != 128 && len != 256 && len != 512 {
		panic(errors.New("AES key length should be 128 or 256 or 512"))
	}
	key := make([]byte, len)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}
	return key
}

func NewKey(len int) Key {
	keyBytes := GenerateRandomKey(len)
	return KeyFromBytes(keyBytes)
}

func KeyFromBytes(keyBytes []byte) Key {
	return Key(hex.EncodeToString(keyBytes))
}

func (k Key) ToBytes() []byte {
	keyBytes, err := hex.DecodeString(string(k))
	if err != nil {
		panic(err.Error())
	}
	return keyBytes
}

func (k Key) Encrypted(address Address) []byte {
	pub, err := ellcurv.ParsePubKey(address.ToBytes(), ellcurv.S256())
	result, err := ellcurv.Encrypt(pub, k.ToBytes())
	if err != nil {
		panic(err)
	}
	return result
}

func (k Key) Verify(address Address, key Key) bool {
	pub, err := ellcurv.ParsePubKey(address.ToBytes(), ellcurv.S256())
	if err != nil {
		panic(err.Error())
	}
	result, err := ellcurv.Encrypt(pub, key.ToBytes())
	if err != nil {
		panic(err.Error())
	}
	return bytes.Equal(result, k.ToBytes())
}

func HashFromBytes(hashBytes []byte) Hash {
	return Hash(hex.EncodeToString(hashBytes))
}

func (h Hash) ToBytes() []byte {
	hashBytes, err := hex.DecodeString(string(h))
	if err != nil {
		panic(err)
	}
	return hashBytes
}
