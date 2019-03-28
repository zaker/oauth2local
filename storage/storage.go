package storage

type Storage interface {
	SetCode(string) error
	GetCode() (string, error)
}
