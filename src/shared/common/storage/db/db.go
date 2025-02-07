package db

import (
	"go-mma/shared/common/storage/db/transactor"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type closeDB func() error

func New(dsn string) (transactor.Transactor, transactor.DBContext, closeDB, error) {
	db, err := establishConn(dsn)
	if err != nil {
		return nil, nil, nil, err
	}

	tx, dbCtx := transactor.New(db, transactor.NestedTransactionsSavepoints)

	return tx, dbCtx, func() error { return close(db) }, err
}

func close(db *sqlx.DB) error {
	return db.Close()
}

func establishConn(dsn string) (*sqlx.DB, error) {
	// this Pings the database trying to connect
	return sqlx.Connect("postgres", dsn)
}
