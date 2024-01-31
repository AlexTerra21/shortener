package storage

import (
	"context"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
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

func (stor *Storage) Set(ctx context.Context, index string, value string, userID int) error {
	err := stor.S.Set(ctx, index, value, userID)
	return err
}

func (stor *Storage) BatchSet(ctx context.Context, data *[]models.BatchStore, userID int) error {
	err := stor.S.BatchSet(ctx, data, userID)
	return err
}

func (stor *Storage) Get(ctx context.Context, idxURL string, userID int) (originalURL string, err error) {
	originalURL, err = stor.S.Get(ctx, idxURL, userID)
	if err == nil {
		logger.Log().Sugar().Debugf("Founded URL %s", originalURL)
	}
	return
}
