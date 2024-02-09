package db

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("mysql", "root:love@tcp(localhost:3306)/rinha")
    if err != nil {
        return fmt.Errorf("erro ao conectar-se ao banco de dados: %v", err)
    }
 
    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("erro ao verificar a conexão com o banco de dados: %v", err)
    }
    
    return nil
}

