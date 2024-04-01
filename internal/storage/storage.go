package storage

import "service-auth/pkg/database/postgresdb"

type Storage struct {
	db *postgresdb.Postgres
}

func New(db *postgresdb.Postgres) Repository {
	return &Storage{
		db: db,
	}
}
