package sea

import "time"

type Sea struct {
	name       string
	totalSpace uint
	freeSpace  uint
}

type Fragment struct {
	timestamp time.Time
	data      []byte
}

func NewSea(name string) *Sea {
	return &Sea{name: name, totalSpace: 0, freeSpace: 0}
}

func NewFragment(data []byte) *Fragment {
	return &Fragment{timestamp: time.Now(), data: data}
}
