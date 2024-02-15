package db

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHostname := os.Getenv("DB_HOSTNAME")

    if dbHostname == "" {
        return fmt.Errorf("DB_HOSTNAME environment variable is not set")
    }

    if dbUser == "" {
        return fmt.Errorf("DB_USER environment variable is not set")
    }

    if dbName == "" {
        return fmt.Errorf("DB_NAME environment variable is not set")
    }

    if dbPort == "" {
        dbPort = "3306" // port default to MySQL
    }

    dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHostname, dbPort, dbName)

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

