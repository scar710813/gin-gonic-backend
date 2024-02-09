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
        log.Fatal("Erro ao inicializar o banco de dados:", err)
    }

	if err := executeSchemaSQL("database/schema.sql"); err != nil {
		log.Fatal("Erro ao executar o arquivo schema.sql:", err)
	}

	router := gin.Default()
	router.POST("/clientes/:id/transacoes", handlers.HandleTransacoes)
	router.GET("/clientes/:id/extrato", handlers.HandleExtrato)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor Gin:", err)
	}
}

func executeSchemaSQL(filePath string) error {
    schemaSQL, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("erro ao ler o arquivo SQL: %v", err)
    }

    statements := strings.Split(string(schemaSQL), ";")

    for _, statement := range statements {
        trimmedStatement := strings.TrimSpace(statement)
        if trimmedStatement == "" {
            continue
        }

        _, err := db.DB.Exec(trimmedStatement)
        if err != nil {
            return fmt.Errorf("erro ao executar a instrução SQL: %v", err)
        }
    }

    return nil
}
