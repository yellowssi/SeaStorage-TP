package sea

import "time"

type Sea struct {
	TotalSpace uint
	FreeSpace  uint
}

type Fragment struct {
	Timestamp time.Time
	Data []byte
}

type Block struct {
}
