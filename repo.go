package main

import "database/sql"

type (
	repo struct {
		db *sql.DB
	}
)

func initRepo() (*repo, error) {
	var r repo

	return &r, nil
}
