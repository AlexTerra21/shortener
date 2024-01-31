package storage

import (
	"context"

	"github.com/AlexTerra21/shortener/internal/app/models"
)

type Storager interface {
	New(string) error
	Close()
	Set(context.Context, string, string, int) error
	BatchSet(context.Context, *[]models.BatchStore, int) error
	Get(context.Context, string, int) (string, error)
}

type Storage struct {
	S       Storager
	confStr string
}
