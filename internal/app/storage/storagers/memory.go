package storagers

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
)

type Memory struct {
	data []ShortenedURL
}

func (m *Memory) New(string) error {
	m.data = make([]ShortenedURL, 0)
	return nil
}

func (m *Memory) Close() {
}

func (m *Memory) Set(_ context.Context, index string, value string, userID int) error {
	newURL := ShortenedURL{
		UUID:        uuid.New().String(),
		IdxShortURL: index,
		OriginalURL: value,
		UserID:      userID,
	}
	m.data = append(m.data, newURL)
	logger.Log().Debug("Storage_Set_Memory", zap.Any("new_url", newURL))
	return nil
}

func (m *Memory) BatchSet(_ context.Context, batchValues *[]models.BatchStore, userID int) error {
	for _, url := range *batchValues {
		newURL := ShortenedURL{
			UUID:        uuid.New().String(),
			IdxShortURL: url.IdxShortURL,
			OriginalURL: url.OriginalURL,
			UserID:      userID,
		}
		m.data = append(m.data, newURL)
		logger.Log().Debug("Storage_Set_Memory", zap.Any("new_url", newURL))
	}
	return nil
}

func (m *Memory) Get(_ context.Context, idxURL string, userID int) (string, error) {
	idx := slices.IndexFunc(m.data, func(c ShortenedURL) bool { return c.IdxShortURL == idxURL && c.UserID == userID })
	if idx == -1 {
		return "", errors.New("URL not found")
	}
	return m.data[idx].OriginalURL, nil
}
