package storagers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/errs"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
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
		IdxShortURL: index,
		OriginalURL: value,
	}
	if err := d.insertURL(ctx, newURL); err != nil {
		return err
	}
	logger.Log().Debug("Storage_Set_DB", zap.Any("new_url", newURL))
	return nil
}

func (d *DB) BatchSet(ctx context.Context, batchValues *[]models.BatchStore) error {
	newURLs := make([]ShortenedURL, 0)
	for _, url := range *batchValues {
		newURL := ShortenedURL{
			UUID:        uuid.New().String(),
			IdxShortURL: url.IdxShortURL,
			OriginalURL: url.OriginalURL,
		}
		newURLs = append(newURLs, newURL)
		logger.Log().Debug("Storage_Set_File", zap.Any("new_url", newURL))
	}
	err := d.insertURLs(ctx, &newURLs)
	if err != nil {
		logger.Log().Error("Error write URL to file", zap.Error(err))
		return err
	}
	return nil
}

func (d *DB) Get(ctx context.Context, idxURL string) (originalURL string, err error) {
	row := d.db.QueryRowContext(ctx, `SELECT original_url FROM urls WHERE short_url = $1`, idxURL)
	err = row.Scan(&originalURL)
	return
}

func (d *DB) GetShortURL(ctx context.Context, originalURL string) (idxURL string, err error) {
	row := d.db.QueryRowContext(ctx, `SELECT short_url FROM urls WHERE original_url = $1`, originalURL)
	err = row.Scan(&idxURL)
	return
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
	// запускаем транзакцию
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Debug("error when creating transaction", zap.Error(err))
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			uuid VARCHAR (100) UNIQUE NOT NULL,
			short_url VARCHAR (100) UNIQUE NOT NULL,
			original_url VARCHAR (100) UNIQUE NOT NULL
		)
	`)

	if err != nil {
		logger.Log().Debug("error when creating urls table", zap.Error(err))
		return tx.Rollback()
	}

	return tx.Commit()
}

func (d *DB) insertURL(ctx context.Context, url ShortenedURL) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	// можно вызвать Rollback в defer,
	// если Commit будет раньше, то откат проигнорируется
	defer tx.Rollback()

	query := `INSERT INTO urls(uuid, short_url, original_url) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		logger.Log().Debug("error when preparing SQL statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, url.UUID, url.IdxShortURL, url.OriginalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgerrcode.UniqueViolation == pgErr.Code) {
			err = errs.ErrConflict
		}
		logger.Log().Debug("error when inserting row into urls table", zap.Error(err))
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		logger.Log().Debug("error when finding rows affected", zap.Error(err))
		return err
	}
	logger.Log().Debug("Inserted", zap.Int64("Rows", rows))

	return tx.Commit()
}

func (d *DB) insertURLs(ctx context.Context, urls *[]ShortenedURL) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	// можно вызвать Rollback в defer,
	// если Commit будет раньше, то откат проигнорируется
	defer tx.Rollback()
	query := `
		INSERT INTO urls(uuid, short_url, original_url) 
			VALUES ($1, $2, $3)
	`
	ctxLocal, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	stmt, err := tx.PrepareContext(ctxLocal, query)
	if err != nil {
		logger.Log().Debug("error when preparing SQL statement", zap.Error(err))
		return err
	}
	defer stmt.Close()
	var allRows int64
	for _, url := range *urls {
		res, err := stmt.ExecContext(ctxLocal, url.UUID, url.IdxShortURL, url.OriginalURL)
		if err != nil {
			logger.Log().Debug("error when inserting row into urls table", zap.Error(err))
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			logger.Log().Debug("error when finding rows affected", zap.Error(err))
			return err
		}
		allRows += rows
	}
	logger.Log().Debug("Inserted", zap.Int64("Rows", allRows))

	return tx.Commit()
}

func (d *DB) GetAll(ctx context.Context, shortURLPrefix string) ([]models.BatchStore, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT short_url, original_url  FROM urls`)
	if err != nil {
		return nil, err
	}
	// не забываем закрыть курсор после завершения работы с данными
	defer rows.Close()
	var allURLs []models.BatchStore
	for rows.Next() {
		var u models.BatchStore
		if err := rows.Scan(&u.IdxShortURL, &u.OriginalURL); err != nil {
			return nil, err
		}
		u.IdxShortURL = shortURLPrefix + "/" + u.IdxShortURL
		allURLs = append(allURLs, u)
	}
	// необходимо проверить ошибки уровня курсора
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return allURLs, nil
}
