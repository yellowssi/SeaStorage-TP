package sea

import "time"

type Sea struct {
	Name       string
	TotalSpace uint
	FreeSpace  uint
}

type Fragment struct {
	Timestamp time.Time
	Data      []byte
}

func NewSea(name string) *Sea {
	return &Sea{Name: name, TotalSpace: 0, FreeSpace: 0}
}

func NewFragment(data []byte) *Fragment {
	return &Fragment{Timestamp: time.Now(), Data: data}
}
