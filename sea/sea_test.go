package sea

import (
	"testing"
)

var s *Sea

func init() {
	s = NewSea("public key")
}

func TestSea_AddOperation(t *testing.T) {
	operations := []*Operation{{
		Action: ActionUserDelete,
		Owner:  "test1",
		Hash:   "hash",
		Shared: false,
	}, {
		Action: ActionUserDelete,
		Owner:  "test2",
		Hash:   "hash",
		Shared: false,
	}}
	for _, operation := range operations {
		s.AddOperation([]*Operation{operation})
	}
	t.Log(s)
}

func TestSea_RemoveOperations(t *testing.T) {
	s.RemoveOperations([]Operation{{
		Action: ActionUserDelete,
		Owner:  "test1",
		Hash:   "hash",
		Shared: false,
	}, {
		Action: ActionUserDelete,
		Owner:  "test2",
		Hash:   "hash",
		Shared: false,
	}})
	t.Log(s)
}
