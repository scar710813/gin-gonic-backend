package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "/root/src/github.com/whyakari/rinha-de-backend-v2/rinha.db")
    if err != nil {
        return err
    }
    return nil
}

