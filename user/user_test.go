package user

import (
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"testing"
)

var singer *signing.Signer

func init() {
	cont := signing.NewSecp256k1Context()
	priv := cont.NewRandomPrivateKey()
	singer = signing.NewCryptoFactory(cont).NewSigner(priv)
}

func TestOperation(t *testing.T) {
	o := NewOperation("address", "PublicKey", "path", "name", "hash", 10, *singer)
	t.Log(o)
	result := o.Verify()
	t.Log("Verify result:", result)
	data := o.ToBytes()
	testOperation, err := OperationFromBytes(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(testOperation)
}
