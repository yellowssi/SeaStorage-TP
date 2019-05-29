package sea

import "testing"

var s *Sea

func init() {
	s = NewSea("public key")
}

func TestSea_AddOperation(t *testing.T) {
	s.AddOperation([]Operation{{
		Action: ActionUserDelete,
		Owner:  "test",
		Hash:   "hash",
		Shared: false,
	}})
	t.Log(s)
}

func TestSea_RemoveOperations(t *testing.T) {
	s.RemoveOperations([]Operation{{
		Action: ActionUserDelete,
		Owner:  "test",
		Hash:   "hash",
		Shared: false,
	}})
	t.Log(s)
}
