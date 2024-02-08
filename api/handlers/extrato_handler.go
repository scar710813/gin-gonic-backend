package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/whyakari/rinha-de-backend-v2/database"
	"github.com/whyakari/rinha-de-backend-v2/models"
)

var limiteDoCliente = 100000

func HandleExtrato(c *gin.Context) {
    clienteID := c.Param("id")

    rows, err := db.DB.Query("SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE id_cliente = $1 ORDER BY realizada_em DESC LIMIT 10", clienteID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao consultar transações no banco de dados"})
        return
    }
    defer rows.Close()

    var saldoAtual int
    err = db.DB.QueryRow("SELECT saldo FROM clientes WHERE id = $1", clienteID).Scan(&saldoAtual)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao obter saldo do cliente"})
        return
    }

    var ultimasTransacoes []models.Transacao
    for rows.Next() {
        var transacao models.Transacao
        err := rows.Scan(
			&transacao.Valor,
			&transacao.Tipo,
			&transacao.Descricao,
			&transacao.RealizadaEm)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao processar transação"})
			fmt.Println(err)
            return
        }
        ultimasTransacoes = append(ultimasTransacoes, transacao)
    }

    // Verifica se não há transações para o cliente
    if len(ultimasTransacoes) == 0 {
        ultimasTransacoes = []models.Transacao{}
    }

    // Retorna a resposta conforme especificado
    c.JSON(200, gin.H{
        "saldo": gin.H{
            "total":        saldoAtual,
            "data_extrato": time.Now().UTC().Format(time.RFC3339Nano),
            "limite":       limiteDoCliente,
        },
        "ultimas_transacoes": ultimasTransacoes,
    })
}

