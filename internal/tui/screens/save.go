package screens

import (
	"log"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
)

// SavePlayer persiste o estado do herói no banco de dados SQLite e registra qualquer erro em nível WARN.
//
// Didática Go: Esta função auxiliar centraliza a conversão do modelo de domínio (`engine.Player`)
// para o modelo de persistência (`storage.Player`) e executa o salvamento de forma segura com tratamento de ponteiros nulos.
func SavePlayer(db *storage.DB, p *engine.Player) {
	if db == nil || p == nil {
		return
	}
	if err := db.SavePlayer(p.ToStorage()); err != nil {
		log.Printf("WARN: failed to save player %s: %v", p.Username, err)
	}
}
