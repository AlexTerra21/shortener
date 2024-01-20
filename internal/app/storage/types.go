package storage

import "context"

type Storager interface {
	New(string) error
	Close()
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
}

type Storage struct {
	S       Storager
	confStr string
}
