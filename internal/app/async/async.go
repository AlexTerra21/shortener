package async

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
)

type Async struct {
	delChan chan storagers.UsersURL
	storage *storage.Storage
}

func NewAsync(s *storage.Storage) *Async {
	instance := &Async{
		delChan: make(chan storagers.UsersURL, 1024),
		storage: s,
	}

	go instance.delURLs()

	return instance
}

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

func (a *Async) Push(del storagers.UsersURL) {
	a.delChan <- del
}
