package async

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
)

// Структура для системы асинхронного удаления
type Async struct {
	delChan chan storagers.UsersURL // Канал, куда помещаются записи подлежащие удалению
	storage *storage.Storage        // Ссылка на хранилище
}

// Инициализация системы асинхронного удаления
func NewAsync(s *storage.Storage) *Async {
	instance := &Async{
		delChan: make(chan storagers.UsersURL, 1024),
		storage: s,
	}

	go instance.delURLs()

	return instance
}

// Асинхронный метод удаления URL из базы
func (a *Async) delURLs() {
	ticker := time.NewTicker(10 * time.Second)

	var del []storagers.UsersURL

	for {
		select {
		case urlID := <-a.delChan:
			del = append(del, urlID)
		case <-ticker.C:
			if len(del) == 0 {
				continue
			}
			db, ok := a.storage.S.(*storagers.DB)
			if !ok {
				del = nil
				continue
			}
			err := db.Delete(context.Background(), del)
			if err != nil {
				logger.Log().Debug("cannot delete messages", zap.Error(err))
				continue
			}
			del = nil
		}
	}
}

// Метод помещающий объекты для удаления в канал
func (a *Async) Push(del storagers.UsersURL) {
	a.delChan <- del
}
