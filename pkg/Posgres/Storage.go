package Postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type Storage struct {
	Db *sql.DB
}

func NewStorage(DBurl string) *Storage {
	conn, err := sql.Open("pgx", DBurl)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Storage{conn}
}
