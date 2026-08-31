package bestiary

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

// GenerateDragonOfDay cria a instância lendária do Dragão com atributos
// determinísticos derivados estritamente da data do dia (YYYY-MM-DD).
//
// Didática Go: Usamos um hash SHA-256 da string da data para obter um uint64,
// que é usado como seed no math/rand. Dessa forma, todos os aventureiros conectados
// no mesmo dia enfrentam um Dragão idêntico, renovado a cada amanhecer.
func GenerateDragonOfDay(dayDate string) engine.Monster {
	if dayDate == "" {
		dayDate = engine.CurrentDateString()
	}

	hash := sha256.Sum256([]byte(dayDate + "-go-dragon-salt"))
	seed := binary.BigEndian.Uint64(hash[:8])
	rng := rand.New(rand.NewSource(int64(seed)))

	// Variações diárias do Dragão:
	// HP Base: 300 a 450
	// ATK Base: 45 a 60
	// DEF Base: 25 a 35
	// Recompensas: 5000 XP, 3000 Gold
	hp := 300 + rng.Intn(151)
	atk := 45 + rng.Intn(16)
	def := 25 + rng.Intn(11)

	// Subtítulos temáticos de acordo com o humor do dia
	titles := []string{
		"O Devorador de Goroutines",
		"O Terror dos Ponteiros",
		"A Fúria Ancestral de Gopher",
		"O Destruidor de Compiladores",
		"A Chama Vermelha do Abismo",
	}
	title := titles[rng.Intn(len(titles))]
	fullName := fmt.Sprintf("%s, %s", i18n.GetMonsterName(i18n.MonsterDragon), title)

	return engine.Monster{
		ID:         i18n.MonsterDragon,
		Name:       fullName,
		Tier:       5, // Tier lendário de Chefe Final
		Health:     hp,
		MaxHealth:  hp,
		Attack:     atk,
		Defense:    def,
		XPReward:   5000,
		GoldReward: 3000,
		Prefix:     "",
		IsDragon:   true,
	}
}
