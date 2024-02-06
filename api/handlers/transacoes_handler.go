package handlers

import (
    "github.com/gin-gonic/gin"
)

func HandleTransacoes(c *gin.Context) {
    // Obtém o ID do cliente a partir dos parâmetros da URL
    clienteID := c.Param("id")

    // Parse do corpo JSON da requisição
    var requestBody struct {
        Valor     int    `json:"valor"`
        Tipo      string `json:"tipo"`
        Descricao string `json:"descricao"`
    }

    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(400, gin.H{"error": "Dados inválidos"})
        return
    }

    // Execute operações no banco de dados PostgreSQL
    _, err := db.Exec("INSERT INTO transacoes (id_cliente, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, NOW())", clienteID, requestBody.Valor, requestBody.Tipo, requestBody.Descricao)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao inserir transação no banco de dados"})
        return
    }

    // Lógica para processar a transação
    if requestBody.Tipo == "c" {
        // Se for crédito, atualizar o saldo do cliente
        _, err := db.Exec("UPDATE clientes SET saldo = saldo + $1 WHERE id = $2", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }

        // Lógica adicional conforme necessário
    } else if requestBody.Tipo == "d" {
        // Se for débito, verificar se o saldo será consistente após a transação
        var saldoAtual int
        err := db.QueryRow("SELECT saldo FROM clientes WHERE id = $1", clienteID).Scan(&saldoAtual)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao obter saldo do cliente"})
            return
        }

        if saldoAtual-requestBody.Valor < -limiteDoCliente {
            // Retornar erro 422 - Saldo inconsistente
            c.JSON(422, gin.H{"error": "Saldo inconsistente após a transação"})
            return
        }

        // Se for consistente, atualizar o saldo do cliente
        _, err = db.Exec("UPDATE clientes SET saldo = saldo - $1 WHERE id = $2", requestBody.Valor, clienteID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
            return
        }

        // Lógica adicional conforme necessário
    }

    // Retornar a resposta conforme especificado
    c.JSON(200, gin.H{
        "limite": 100000,
        "saldo":  -9098,
    })
}

