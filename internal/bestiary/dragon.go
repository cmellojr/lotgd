package bestiary

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

// GenerateDragonOfDay instancia a criatura lendária do Dragão Ancestral com atributos
// determinísticos derivados estritamente da data do calendário (YYYY-MM-DD).
//
// Didática Go: Usamos um hash de alta dispersão SHA-256 (`sha256.Sum256`) a partir da string da data
// concatenada com um salt para derivar um `uint64` utilizado como seed em `math/rand`. Dessa forma,
// todos os aventureiros conectados no mesmo dia enfrentam um Dragão idêntico com o mesmo título e poder,
// renovado de forma 100% determinística e reproduzível a cada novo amanhecer.
func GenerateDragonOfDay(dayDate string) engine.Monster {
	if dayDate == "" {
		dayDate = engine.CurrentDateString()
	}

	hash := sha256.Sum256([]byte(dayDate + "-go-dragon-salt"))
	seed := binary.BigEndian.Uint64(hash[:8])
	rng := rand.New(rand.NewSource(int64(seed)))

	// Variações diárias dos atributos do Dragão:
	// - HP Base: 300 a 450
	// - Ataque Base: 45 a 60
	// - Defesa Base: 25 a 35
	// - Recompensas Fixas: 5000 XP e 3000 moedas de Ouro
	hp := 300 + rng.Intn(151)
	atk := 45 + rng.Intn(16)
	def := 25 + rng.Intn(11)

	// Seleção determinística do subtítulo temático do dia via DragonTitleIndex
	titleIdx := engine.DragonTitleIndex(dayDate, len(i18n.DragonTitlesPTBR))
	fullName := fmt.Sprintf("%s, %s", i18n.GetMonsterName(i18n.MonsterDragon), i18n.DragonTitlesPTBR[titleIdx])

	return engine.Monster{
		ID:         i18n.MonsterDragon,
		Name:       fullName,
		Tier:       5, // Tier lendário especial reservado ao Chefe Final
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
