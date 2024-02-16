package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/whyakari/rinha-de-backend-v2/database"
	"github.com/whyakari/rinha-de-backend-v2/models"
)

var limiteDoCliente = 100000

func clienteExists(clienteID string) bool {
    var exists int
    err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE id = ?)", clienteID).Scan(&exists)
    if err != nil {
        return false
    }
    return exists == 1
}

func HandleExtrato(c *gin.Context) {
    clienteID := c.Param("id")

    if !clienteExists(clienteID) {
        c.AbortWithStatus(404)
        return
    }

    rows, err := db.DB.Query("SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE id_cliente = ? ORDER BY realizada_em DESC LIMIT 10", clienteID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Error when querying transactions in the database"})
        return
    }
    defer rows.Close()

    var saldoAtual int
    err = db.DB.QueryRow("SELECT saldo FROM clientes WHERE id = ?", clienteID).Scan(&saldoAtual)
    if err != nil {
        c.JSON(500, gin.H{"error": "Error getting customer balance"})
        return
    }

    var ultimasTransacoes []models.Transacao

    for rows.Next() {
        var transacao models.Transacao
        err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
        if err != nil {
            c.JSON(500, gin.H{"error": "Error processing transaction"})
            fmt.Println(err)
            return
        }
        ultimasTransacoes = append(ultimasTransacoes, transacao)
    }

    if len(ultimasTransacoes) == 0 {
        ultimasTransacoes = []models.Transacao{}
    }

    c.JSON(200, gin.H{
        "saldo": gin.H{
            "total":        saldoAtual,
            "data_extrato": time.Now().UTC().Format(time.RFC3339Nano),
            "limite":       limiteDoCliente,
        },
        "ultimas_transacoes": ultimasTransacoes,
    })
}
