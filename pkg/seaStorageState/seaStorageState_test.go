package seaStorageState

import "testing"

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
