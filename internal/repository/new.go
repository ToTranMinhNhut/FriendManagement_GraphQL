package repository

import "database/sql"

type DBRepo struct {
	Db *sql.DB
}

func NewDBRepo(db *sql.DB) DBRepo {
	return DBRepo{
		Db: db,
	}
}
