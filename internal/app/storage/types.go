package storage

type Storager interface {
	New(string) error
	Close()
	Set(string, string)
	Get(string) (string, error)
}

type Storage struct {
	S       Storager
	confStr string
}
