package storagers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
)

type DB struct {
	db *sql.DB
}

func (d *DB) New(dbConf string) error {
	db, err := sql.Open("pgx", dbConf)
	if err != nil {
		return errors.New("error open database")
	}
	d.db = db
	err = d.createTable()
	return err
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Set(ctx context.Context, index string, value string) error {
	newURL := ShortenedURL{
		UUID:        uuid.New().String(),
		ShortURL:    index,
		OriginalURL: value,
	}
	if err := d.insertURL(ctx, newURL); err != nil {
		return err
	}
	logger.Log().Debug("Storage_Set_DB", zap.Any("new_url", newURL))
	return nil
}

func (d *DB) Get(ctx context.Context, url string) (string, error) {
	return "", nil
}

func (d *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (d *DB) createTable() error {
	query := "CREATE TABLE IF NOT EXISTS urls" +
		"(uuid VARCHAR (100) UNIQUE NOT NULL," +
		"short_url VARCHAR (100) NOT NULL," +
		"original_url VARCHAR (100) UNIQUE NOT NULL)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		logger.Log().Debug("error when creating product table", zap.Error(err))
		return err
	}
	return nil
}

func (d *DB) insertURL(ctx context.Context, url ShortenedURL) error {
	query := "INSERT INTO urls(uuid, short_url, original_url) VALUES ($1, $2, $3)"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Log().Debug("error when preparing SQL statement", zap.Error(err))
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, url.UUID, url.ShortURL, url.OriginalURL)
	if err != nil {
		logger.Log().Debug("error when inserting row into urls table", zap.Error(err))
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		logger.Log().Debug("error when finding rows affected", zap.Error(err))
		return err
	}
	logger.Log().Debug("Inserted", zap.Int64("Rows", rows))

	return nil
}
