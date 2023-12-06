package storage

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(index string, value string) {
	s.data[index] = value
}

func (s *Storage) Get(index string) string {
	return s.data[index]
}
