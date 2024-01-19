package storage

import (
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
)

func NewStorage(fname string, dbstr string) (*Storage, error) {
	var stor Storage
	if dbstr != "" {
		stor = Storage{
			S:       &storagers.DB{},
			confStr: dbstr,
		}
		logger.Log().Debug("Database mode")
	} else if fname != "" {
		stor = Storage{
			S:       &storagers.File{},
			confStr: fname,
		}
		logger.Log().Debug("File mode")
	} else {
		stor = Storage{
			S:       &storagers.Memory{},
			confStr: "",
		}
		logger.Log().Debug("Memory mode")
	}

	if err := stor.S.New(stor.confStr); err != nil {
		return nil, err
	}

	return &stor, nil
}

func (stor *Storage) Close() {
	stor.S.Close()
}

func (stor *Storage) Set(index string, value string) {
	stor.S.Set(index, value)
}

func (stor *Storage) Get(url string) (string, error) {
	data, err := stor.S.Get(url)
	return data, err
}
