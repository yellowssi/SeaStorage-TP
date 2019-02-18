package payload

type SeaStoragePayload struct {
	Name   string
	Action string
	Space  int
}

func FromBytes(payloadData []byte) (*SeaStoragePayload, error) {

}
