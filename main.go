package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/whyakari/rinha-de-backend-v2/api/handlers"
	db "github.com/whyakari/rinha-de-backend-v2/database"
)

func main() {
    if err := db.InitDB(); err != nil {
        log.Fatal("Error initializing the database:", err)
    }

	if err := executeSchemaSQL("schema.sql"); err != nil {
		log.Fatal("Error when executing schema.sql file:", err)
	}

	router := gin.Default()
	router.POST("/clientes/:id/transacoes", handlers.HandleTransacoes)
	router.GET("/clientes/:id/extrato", handlers.HandleExtrato)

	if err := router.Run(":3000"); err != nil {
		log.Fatal("Error starting Gin server:", err)
	}
}

func executeSchemaSQL(filePath string) error {
    schemaSQL, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("error reading SQL file: %v", err)
    }

    statements := strings.Split(string(schemaSQL), ";")

    for _, statement := range statements {
        trimmedStatement := strings.TrimSpace(statement)
        if trimmedStatement == "" {
            continue
        }

        _, err := db.DB.Exec(trimmedStatement)
        if err != nil {
            return fmt.Errorf("error executing instruction SQL: %v", err)
        }
    }

    return nil
}
