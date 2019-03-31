package sea

import "time"

type Sea struct {
	totalSpace uint
	freeSpace  uint
}

type Fragment struct {
	timestamp time.Time
	data      []byte
}

type Block struct {
}
