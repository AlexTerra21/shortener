package storage

import (
	"context"

	"github.com/AlexTerra21/shortener/internal/app/models"
)

type Storager interface {
	New(string) error
	Close()
	Set(context.Context, string, string) error
	BatchSet(context.Context, *[]models.BatchStore) error
	Get(context.Context, string) (string, error)
}

type Storage struct {
	S       Storager
	confStr string
}
