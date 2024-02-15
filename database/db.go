package db

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    dbHostname := os.Getenv("DB_HOSTNAME")
    if dbHostname == "" {
        return fmt.Errorf("DB_HOSTNAME environment variable is not set")
    }

    dbConnString := fmt.Sprintf("root:love@tcp(%s:3306)/rinha", dbHostname)
    var err error
    DB, err = sql.Open("mysql", dbConnString)
    if err != nil {
        return fmt.Errorf("error connecting to database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("error checking database connection: %v", err)
    }

    return nil
}
