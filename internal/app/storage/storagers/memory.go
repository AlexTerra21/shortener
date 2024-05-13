package storagers

import (
	"context"
	"errors"
	"slices"
	"sync"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
)

// Структура для хранения данных в памяти
type Memory struct {
	data  []ShortenedURL
	mutex sync.RWMutex
}

// Инициализация хранилища
func (m *Memory) New(string) error {
	m.data = make([]ShortenedURL, 0)
	m.mutex = sync.RWMutex{}
	return nil
}

// Закрытие хранилища
func (m *Memory) Close() {
}

// Добавление данных в хранилище
func (m *Memory) Set(_ context.Context, index string, value string, userID int) error {
	newURL := ShortenedURL{
		UUID:        userID,
		IdxShortURL: index,
		OriginalURL: value,
		DeletedFlag: false,
	}
	m.data = append(m.data, newURL)
	logger.Log().Debug("Storage_Set_Memory", zap.Any("new_url", newURL))
	return nil
}

// Добавление пакета данных в хранилище
func (m *Memory) BatchSet(_ context.Context, batchValues *[]models.BatchStore, userID int) error {
	for _, url := range *batchValues {
		newURL := ShortenedURL{
			UUID:        userID,
			IdxShortURL: url.IdxShortURL,
			OriginalURL: url.OriginalURL,
			DeletedFlag: false,
		}
		m.data = append(m.data, newURL)
		logger.Log().Debug("Storage_Set_Memory", zap.Any("new_url", newURL))
	}
	return nil
}

// Получение данных из хранилища
func (m *Memory) Get(_ context.Context, idxURL string) (string, bool, error) {
	idx := slices.IndexFunc(m.data, func(c ShortenedURL) bool { return c.IdxShortURL == idxURL })
	if idx == -1 {
		return "", false, errors.New("URL not found")
	}
	return m.data[idx].OriginalURL, m.data[idx].DeletedFlag, nil
}

// Получение статистики по хранилищу
func (m *Memory) Stats(ctx context.Context) (models.StatsResp, error) {
	uniqUsersIDs := make(map[int]bool)

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, shortURL := range m.data {
		uniqUsersIDs[shortURL.UUID] = true
	}
	return models.StatsResp{
		UrlsCount: len(m.data),
		UserCount: len(uniqUsersIDs),
	}, nil
}
