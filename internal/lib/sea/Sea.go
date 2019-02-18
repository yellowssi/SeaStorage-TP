package sea

type Hash [512]byte

type Sea struct {
	TotalSpace uint
	FreeSpace  uint
}

type SeaFragment struct {
	hash Hash
}

type SeaBlock struct {
}
