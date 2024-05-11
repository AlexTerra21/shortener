package storage

import (
	"context"

	"github.com/AlexTerra21/shortener/internal/app/models"
)

// Интерфейс, который должны реализовать "хранители" (память, файл, база)
type Storager interface {
	New(string) error
	Close()
	Set(context.Context, string, string, int) error
	BatchSet(context.Context, *[]models.BatchStore, int) error
	Get(context.Context, string) (string, bool, error)
	Stats(ctx context.Context) (models.StatsResp, error)
}

// Описание "хранителя"
type Storage struct {
	S       Storager // ссылка на "хранителя"
	confStr string   // парамеры настройки
}
