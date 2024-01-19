package storagers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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
