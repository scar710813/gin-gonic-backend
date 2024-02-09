package handlers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	db "github.com/whyakari/rinha-de-backend-v2/database"
)

func HandleTransacoes(c *gin.Context) {
    if err := db.InitDB(); err != nil {
        log.Fatal("Erro ao inicializar o banco de dados:", err)
    }
    
    clienteID := c.Param("id")

    var requestBody struct {
        Valor     int    `json:"valor"`
        Tipo      string `json:"tipo"`
        Descricao string `json:"descricao"`
    }

    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(400, gin.H{"error": "Dados inválidos"})
        return
    }

    _, err := db.DB.Exec("INSERT INTO transacoes (id_cliente, valor, tipo, descricao, realizada_em) VALUES (?, ?, ?, ?, NOW())", clienteID, requestBody.Valor, requestBody.Tipo, requestBody.Descricao)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao inserir transação no banco de dados"})
		fmt.Println(err)
        return
    }

    if requestBody.Tipo == "c" {
        _, err := db.DB.Exec("UPDATE clientes SET saldo = saldo + ? WHERE id = ?", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }
    } else if requestBody.Tipo == "d" {
        var saldoAtual int
        err := db.DB.QueryRow("SELECT saldo FROM clientes WHERE id = ?", clienteID).Scan(&saldoAtual)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao obter saldo do cliente"})
            return
        }

        if requestBody.Tipo == "d" && saldoAtual-requestBody.Valor < -limiteDoCliente {
            c.JSON(422, gin.H{"error": "Transação de débito não permitida. Saldo insuficiente."})
            return
        }

        _, err = db.DB.Exec("UPDATE clientes SET saldo = saldo - ? WHERE id = ?", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }
    }

    c.JSON(200, gin.H{
        "limite": 100000,
        "saldo":  -9098,
    })
}
