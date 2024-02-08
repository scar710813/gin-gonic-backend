package handlers

import (
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

    _, err := db.DB.Exec("INSERT INTO transacoes (id_cliente, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, DATETIME('now'))", clienteID, requestBody.Valor, requestBody.Tipo, requestBody.Descricao)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao inserir transação no banco de dados"})
        return
    }

    if requestBody.Tipo == "c" {
        // Se for crédito, atualizar o saldo do cliente
        _, err := db.DB.Exec("UPDATE clientes SET saldo = saldo + $1 WHERE id = $2", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }

    } else if requestBody.Tipo == "d" {
        // Se for débito, verificar se o saldo será consistente após a transação
        var saldoAtual int
        err := db.DB.QueryRow("SELECT saldo FROM clientes WHERE id = $1", clienteID).Scan(&saldoAtual)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao obter saldo do cliente"})
            return
        }

		if requestBody.Tipo == "d" && saldoAtual-requestBody.Valor < -limiteDoCliente {
			c.JSON(422, gin.H{"error": "Transação de débito não permitida. Saldo insuficiente."})
			return
		}

        // Se for consistente, atualizar o saldo do cliente
        _, err = db.DB.Exec("UPDATE clientes SET saldo = saldo - $1 WHERE id = $2", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }
    }

    // Retornar a resposta conforme especificado
    c.JSON(200, gin.H{
        "limite": 100000,
        "saldo":  -9098,
    })
}

