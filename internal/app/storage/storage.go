package storage

import (
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
)

// Инициализация хранилища
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
