package models

import (
	"database/sql"
	"os"
)

func OpenDB(dsn, dbtype, path string) (*sql.DB, error) {
	db, err := sql.Open(dbtype, dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(f))
	if err != nil {
		return nil, err
	}

	return db, nil
}
