package state

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"strings"
)

var Namespace = hexdigest("SeaStorage")[:6]

type SeaStorageState struct {
	context   *processor.Context
	userCache map[string][]byte
	seaCache  map[string][]byte
}

func hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}
