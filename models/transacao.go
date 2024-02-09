package models

type Transacao struct {
    Valor       int       `json:"valor"`
    Tipo        string    `json:"tipo"`
    Descricao   string    `json:"descricao"`
    RealizadaEm string	  `json:"realizada_em"`
}
