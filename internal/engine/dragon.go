package engine

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

// DragonTitleIndex retorna um índice determinístico para selecionar o título
// temático do Dragão do Dia.
//
// Didática Go: Usamos um hash SHA-256 da string da data concatenada com um salt
// para obter um uint64, que é usado como seed no math/rand. Dessa forma, todos
// os aventureiros conectados no mesmo dia recebem o mesmo índice de título,
// renovado a cada amanhecer de forma reproduzível.
//
// A lógica é idêntica à used bestiary.GenerateDragonOfDay, garantindo coerência
// entre os atributos do monstro e o título exibido na tela.
func DragonTitleIndex(dayDate string, titleCount int) int {
	if titleCount <= 0 {
		return 0
	}

	hash := sha256.Sum256([]byte(dayDate + "-go-dragon-salt"))
	seed := binary.BigEndian.Uint64(hash[:8])
	rng := rand.New(rand.NewSource(int64(seed)))

	return rng.Intn(titleCount)
}
