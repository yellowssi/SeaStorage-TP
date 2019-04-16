package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	ellcurv "github.com/btcsuite/btcd/btcec"
)

func SHA256Bytes(data string) []byte {
	hashHandler := sha256.New()
	hashHandler.Write(HexToBytes(data))
	return hashHandler.Sum(nil)
}

func SHA256Hex(data string) string {
	return BytesToHex(SHA256Bytes(data))
}

func SHA384Bytes(data string) []byte {
	hashHandler := sha512.New384()
	hashHandler.Write(HexToBytes(data))
	return hashHandler.Sum(nil)
}

func SHA384Hex(data string) string {
	return BytesToHex(SHA384Bytes(data))
}

func SHA512Bytes(data string) []byte {
	hashHandler := sha512.New()
	hashHandler.Write(HexToBytes(data))
	return hashHandler.Sum(nil)
}

func SHA512Hex(data string) string {
	return BytesToHex(SHA512Bytes(data))
}

func Encryption(publicKey string, data string) ([]byte, error) {
	pub, err := ellcurv.ParsePubKey(HexToBytes(publicKey), ellcurv.S256())
	if err != nil {
		return nil, err
	}
	result, err := ellcurv.Encrypt(pub, HexToBytes(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Verify(publicKey string, sign string, data string) bool {
	pub, err := ellcurv.ParsePubKey(HexToBytes(publicKey), ellcurv.S256())
	if err != nil {
		return false
	}
	signature, err := ellcurv.ParseSignature(HexToBytes(sign), ellcurv.S256())
	if err != nil {
		return false
	}
	hash := SHA512Bytes(data)
	return signature.Verify(hash, pub)
}

func Sign(privateKey string, data string) ([]byte, error) {
	priv, _ := ellcurv.PrivKeyFromBytes(ellcurv.S256(), HexToBytes(privateKey))
	hash := SHA512Bytes(data)
	signature, err := priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return signature.Serialize(), nil
}

func Decryption(privateKey string, data string) ([]byte, error) {
	priv, _ := ellcurv.PrivKeyFromBytes(ellcurv.S256(), HexToBytes(privateKey))
	result, err := ellcurv.Decrypt(priv, HexToBytes(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GenerateRandomAESKey(len int) []byte {
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

func NewAESKey(len int) string {
	keyBytes := GenerateRandomAESKey(len)
	return BytesToHex(keyBytes)
}

func AESKeyEncryptedByPublicKey(key string, publicKey string) []byte {
	pub, err := ellcurv.ParsePubKey(HexToBytes(publicKey), ellcurv.S256())
	result, err := ellcurv.Encrypt(pub, HexToBytes(key))
	if err != nil {
		panic(err)
	}
	return result
}

func AESKeyVerify(publicKey string, key string, encryptedKey string) bool {
	pub, err := ellcurv.ParsePubKey(HexToBytes(publicKey), ellcurv.S256())
	if err != nil {
		panic(err.Error())
	}
	result, err := ellcurv.Encrypt(pub, HexToBytes(key))
	if err != nil {
		panic(err.Error())
	}
	return bytes.Equal(result, HexToBytes(encryptedKey))
}

func AESKeyEncryption(key string, data string) (result []byte, err error) {
	cipher, err := aes.NewCipher(HexToBytes(key))
	if err != nil {
		return nil, err
	}
	cipher.Encrypt(result, HexToBytes(data))
	return
}

func AESKeyDecryption(key string, data string) (result []byte, err error) {
	cipher, err := aes.NewCipher(HexToBytes(key))
	if err != nil {
		return nil, err
	}
	cipher.Decrypt(result, HexToBytes(data))
	return
}

func HexToBytes(str string) []byte {
	data, _ := hex.DecodeString(str)
	return data
}

func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}
