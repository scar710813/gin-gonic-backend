package main

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/whyakari/rinha-de-backend-v2/api/handlers"
)

var db *sql.DB

func main() {
    // Configurar a conexão com o PostgreSQL
	connectionString := "postgresql://admin:123@rinha/todos?sslmode=disable"

    var err error
    db, err = sql.Open("postgres", connectionString)
    if err != nil {
        fmt.Println("Erro ao conectar ao PostgreSQL:", err)
        return
    }
    defer db.Close()

    // Inicializar o esquema do banco de dados
    if err := initializeDatabaseSchema(); err != nil {
        fmt.Println("Erro ao inicializar o esquema do banco de dados:", err)
        return
    }

    // Resto do seu código para iniciar o servidor...
	router := gin.Default()

	// Rota para transações
	router.POST("/clientes/:id/transacoes", handlers.HandleTransacoes)

	// Rota para extrato
	router.GET("/clientes/:id/extrato", handlers.HandleExtrato)

	// Inicializar o servidor na porta 8080
	router.Run(":8080")
}

func initializeDatabaseSchema() error {
    // Ler o conteúdo do arquivo schema.sql
    content, err := os.ReadFile("database/schema.sql")
    if err != nil {
        return err
    }

    // Executar o conteúdo do arquivo SQL para criar as tabelas
    _, err = db.Exec(string(content))
    return err
}

