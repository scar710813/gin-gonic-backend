package handlers

import (
	"database/sql"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/whyakari/rinha-de-backend-v2/models"
)

var db * sql.DB
var limiteDoCliente = 100000

func HandleExtrato(c *gin.Context) {
    // Obtém o ID do cliente a partir dos parâmetros da URL
    clienteID := c.Param("id")

    // Execute operações no banco de dados PostgreSQL
    rows, err := db.Query("SELECT * FROM transacoes WHERE id_cliente = $1 ORDER BY realizada_em DESC LIMIT 10", clienteID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao consultar transações no banco de dados"})
        return
    }
    defer rows.Close()

    // Lógica para gerar o extrato
    var saldoAtual int
    err = db.QueryRow("SELECT saldo FROM clientes WHERE id = $1", clienteID).Scan(&saldoAtual)
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
            return
        }
        ultimasTransacoes = append(ultimasTransacoes, transacao)
    }

    // Retornar a resposta conforme especificado
    c.JSON(200, gin.H{
        "saldo": saldoAtual,
        "data_extrato": time.Now().UTC().Format(time.RFC3339Nano),
        "limite": limiteDoCliente,
        "ultimas_transacoes": ultimasTransacoes,
    })
}

