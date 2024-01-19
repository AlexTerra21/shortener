package storagers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	d.createTable()
	return nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Set(index string, value string) {

}

func (d *DB) Get(url string) (string, error) {
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
	query := `CREATE TABLE IF NOT EXISTS urls (uuid VARCHAR (100) UNIQUE NOT NULL, short_url VARCHAR (100) NOT NULL, original_url VARCHAR (100) NOT NULL)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		logger.Log().Sugar().Debugf("Error %s when creating product table", err)
		return err
	}
	return nil
}
