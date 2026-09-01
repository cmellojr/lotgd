package bestiary

import (
	"fmt"
	"math/rand"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

// MonsterGenerator produz criaturas procedurais com modificadores estocásticos.
type MonsterGenerator struct {
	rng *rand.Rand
}

// NewMonsterGenerator cria o gerador com um RNG injetável para reprodutibilidade.
func NewMonsterGenerator(rng *rand.Rand) *MonsterGenerator {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &MonsterGenerator{rng: rng}
}

// GenerateForPlayer escolhe um monstro balanceado de acordo com o nível atual do jogador.
//
// Regra de Tier:
// - Nível 1 a 2 -> Tier 1 (com 10% de chance de encontrar Tier 2 desafiador)
// - Nível 3 a 4 -> Tier 2 (com 15% de chance de Tier 3)
// - Nível 5 a 7 -> Tier 3 (com 15% de chance de Tier 4)
// - Nível 8 a 10 -> Tier 4
func (mg *MonsterGenerator) GenerateForPlayer(playerLevel int) engine.Monster {
	var tier int
	roll := mg.rng.Float64()

	switch {
	case playerLevel <= 2:
		tier = 1
		if roll < 0.10 {
			tier = 2
		}
	case playerLevel <= 4:
		tier = 2
		if roll < 0.15 {
			tier = 3
		}
	case playerLevel <= 7:
		tier = 3
		if roll < 0.15 {
			tier = 4
		}
	default:
		tier = 4
	}

	return mg.GenerateByTier(tier)
}

// GenerateByTier gera um monstro procedural pertencente ao Tier requisitado.
func (mg *MonsterGenerator) GenerateByTier(tier int) engine.Monster {
	if tier < 1 {
		tier = 1
	}
	if tier > 4 {
		tier = 4
	}

	monsterList := TierMonstersLists[tier]
	chosenID := monsterList[mg.rng.Intn(len(monsterList))]
	tpl := CanonicalTemplates[chosenID]

	baseName := i18n.GetMonsterName(chosenID)

	// 50% de chance de receber um afixo especial
	var prefix string
	hp := tpl.BaseHP
	atk := tpl.BaseATK
	def := tpl.BaseDEF
	xp := tpl.BaseXP
	gold := tpl.BaseGold

	if mg.rng.Float64() < 0.50 {
		affix := AvailableAffixes[mg.rng.Intn(len(AvailableAffixes))]
		prefix = affix.NamePTBR
		hp = int(float64(hp) * affix.HPMult)
		atk = int(float64(atk) * affix.ATKMult)
		def = int(float64(def) * affix.DEFMult)
		xp = int(float64(xp) * affix.XPMult)
		gold = int(float64(gold) * affix.GoldMult)
	}

	fullName := baseName
	if prefix != "" {
		fullName = fmt.Sprintf("%s %s", prefix, baseName)
	}

	// Garante valores mínimos viáveis
	if hp < 5 {
		hp = 5
	}
	if atk < 1 {
		atk = 1
	}
	if xp < 1 {
		xp = 1
	}
	if gold < 1 {
		gold = 1
	}

	return engine.Monster{
		ID:         chosenID,
		Name:       fullName,
		Tier:       tier,
		Health:     hp,
		MaxHealth:  hp,
		Attack:     atk,
		Defense:    def,
		XPReward:   xp,
		GoldReward: gold,
		Prefix:     prefix,
		IsDragon:   false,
	}
}
