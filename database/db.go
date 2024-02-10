package db

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    var err error
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return fmt.Errorf("DATABASE_URL não definida")
    }

    DB, err = sql.Open("mysql", dbURL)
    if err != nil {
        return fmt.Errorf("erro ao conectar-se ao banco de dados: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("erro ao verificar a conexão com o banco de dados: %v", err)
    }

    return nil
}

