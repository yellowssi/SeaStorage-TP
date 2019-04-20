package state

import (
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"testing"
)

func TestSeaStorageState_GetSea(t *testing.T) {
	testMap := map[string][]byte{
		"a": []byte("apple"),
		"b": []byte("banana"),
		"c": []byte("core"),
		"d": nil,
	}
	test, ok := testMap["d"]
	if ok {
		println(test)
	} else {
		println("doesn't exists")
	}
}

func TestMakeAddress(t *testing.T) {
	cont := signing.NewSecp256k1Context()
	priv := cont.NewRandomPrivateKey()
	pub := cont.GetPublicKey(priv)
	address := MakeAddress(AddressTypeUser, "Test", pub.AsHex())
	t.Log(address)
	t.Log(len(address))
}
