package storage

type Stop struct {
	ID          int
	BannedWords string
}

type Interface interface {
	GetList() ([]Stop, error)
	AddList(c Stop) error
	CreateStopTable() error
	DropStopTable() error
}
