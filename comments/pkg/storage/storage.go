package storage

type Comment struct {
	ID      int
	IdNews  int
	Content string
	PubTime int64
}

type Interface interface {
	CreateTable() error
	RemoveTable() error
	GetComments(IdNews int) ([]Comment, error)
	AddComment(comment Comment) error
	RemoveComment(comment Comment) error
}
