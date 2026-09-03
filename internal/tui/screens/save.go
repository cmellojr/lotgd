package screens

import (
	"log"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
)

// SavePlayer persists player state to database and logs any error at WARN level.
func SavePlayer(db *storage.DB, p *engine.Player) {
	if db == nil || p == nil {
		return
	}
	if err := db.SavePlayer(p.ToStorage()); err != nil {
		log.Printf("WARN: failed to save player %s: %v", p.Username, err)
	}
}
