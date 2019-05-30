package user

import (
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"testing"
	"time"
)

var signer *signing.Signer

func init() {
	cont := signing.NewSecp256k1Context()
	priv := cont.NewRandomPrivateKey()
	signer = signing.NewCryptoFactory(cont).NewSigner(priv)
}

func TestOperation(t *testing.T) {
	o := NewOperation("address", signer.GetPublicKey().AsHex(), "sea", "path", "name", "hash", 10, time.Now().Unix(), *signer)
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
