// Copyright © 2019 yellowsea <hh1271941291@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	ellcurv "github.com/btcsuite/btcd/btcec"
)

// SHA224
func SHA224BytesFromBytes(data []byte) []byte {
	hashHandler := sha256.New224()
	hashHandler.Write(data)
	return hashHandler.Sum(nil)
}

func SHA224BytesFromHex(data string) []byte {
	return SHA224BytesFromBytes(HexToBytes(data))
}

func SHA224HexFromBytes(data []byte) string {
	return BytesToHex(SHA224BytesFromBytes(data))
}

func SHA224HexFromHex(data string) string {
	return BytesToHex(SHA224BytesFromHex(data))
}

// SHA256
func SHA256BytesFromBytes(data []byte) []byte {
	hashHandler := sha256.New()
	hashHandler.Write(data)
	return hashHandler.Sum(nil)
}

func SHA256BytesFromHex(data string) []byte {
	return SHA256BytesFromBytes(HexToBytes(data))
}

func SHA256HexFromBytes(data []byte) string {
	return BytesToHex(SHA256BytesFromBytes(data))
}

func SHA256HexFromHex(data string) string {
	return BytesToHex(SHA256BytesFromHex(data))
}

// SHA384
func SHA384BytesFromBytes(data []byte) []byte {
	hashHandler := sha512.New384()
	hashHandler.Write(data)
	return hashHandler.Sum(nil)
}

func SHA384BytesFromHex(data string) []byte {
	return SHA384BytesFromBytes(HexToBytes(data))
}

func SHA384FromBytes(data []byte) string {
	return BytesToHex(SHA384BytesFromBytes(data))
}

func SHA384HexFromHex(data string) string {
	return BytesToHex(SHA384BytesFromHex(data))
}

// SHA512
func SHA512BytesFromBytes(data []byte) []byte {
	hashHandler := sha512.New()
	hashHandler.Write(data)
	return hashHandler.Sum(nil)
}

func SHA512BytesFromHex(data string) []byte {
	return SHA512BytesFromBytes(HexToBytes(data))
}

func SHA512HexFromBytes(data []byte) string {
	return BytesToHex(SHA512BytesFromBytes(data))
}

func SHA512HexFromHex(data string) string {
	return BytesToHex(SHA512BytesFromHex(data))
}

// Ellcurv
func Encryption(publicKey, data string) ([]byte, error) {
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

func Decryption(privateKey, data string) ([]byte, error) {
	priv, _ := ellcurv.PrivKeyFromBytes(ellcurv.S256(), HexToBytes(privateKey))
	result, err := ellcurv.Decrypt(priv, HexToBytes(data))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AES
func GenerateRandomAESKey(len int) []byte {
	if len != 128 && len != 192 && len != 256 {
		panic(aes.KeySizeError(len))
	}
	key := make([]byte, len/8)
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

func AESKeyEncryptedByPublicKey(key, publicKey string) []byte {
	pub, err := ellcurv.ParsePubKey(HexToBytes(publicKey), ellcurv.S256())
	if err != nil {
		panic(err)
	}
	result, err := ellcurv.Encrypt(pub, HexToBytes(key))
	if err != nil {
		panic(err)
	}
	return result
}

func AESKeyEncryption(key, data string) (result []byte, err error) {
	cipher, err := aes.NewCipher(HexToBytes(key))
	if err != nil {
		return nil, err
	}
	cipher.Encrypt(result, HexToBytes(data))
	return
}

func AESKeyDecryption(key, data string) (result []byte, err error) {
	cipher, err := aes.NewCipher(HexToBytes(key))
	if err != nil {
		return nil, err
	}
	cipher.Decrypt(result, HexToBytes(data))
	return
}

// Convert between Hex and Bytes
func HexToBytes(str string) []byte {
	data, _ := hex.DecodeString(str)
	return data
}

func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}
