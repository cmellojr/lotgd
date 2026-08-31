package engine_test

import (
	"math/rand"
	"testing"
	"time"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
	"lotgd/internal/storage"
)

func TestPlayerStatsCalculation(t *testing.T) {
	p := &engine.Player{
		Level:       1,
		BaseAttack:  10,
		BaseDefense: 5,
		Weapon: engine.Item{
			ID:         i18n.WeaponDagger,
			PowerBonus: 3,
		},
		Armor: engine.Item{
			ID:         i18n.ArmorLeather,
			PowerBonus: 2,
		},
	}

	if p.TotalAttack() != 13 {
		t.Fatalf("esperado TotalAttack = 13, obtido %d", p.TotalAttack())
	}

	if p.TotalDefense() != 7 {
		t.Fatalf("esperado TotalDefense = 7, obtido %d", p.TotalDefense())
	}
}

func TestCombatEngine_DeterministicDamage(t *testing.T) {
	// RNG com seed fixa para assertividade total
	rng := rand.New(rand.NewSource(42))
	ce := engine.NewCombatEngine(rng)

	// ATK 10 vs DEF 5 -> Dano esperado >= 1
	dmg, crit := ce.CalculateDamage(10, 5)
	if dmg < 1 {
		t.Fatalf("dano nunca deve ser inferior a 1, obtido: %d", dmg)
	}
	_ = crit

	// Dano mínimo garantido mesmo com defesa gigantesca
	minDmg, _ := ce.CalculateDamage(2, 999)
	if minDmg != 1 {
		t.Fatalf("com defesa massiva, dano deve ser exatamente 1, obtido: %d", minDmg)
	}
}

func TestCombatEngine_AttackVictory(t *testing.T) {
	rng := rand.New(rand.NewSource(100))
	ce := engine.NewCombatEngine(rng)

	p := &engine.Player{
		Username:     "GopherHero",
		Level:        1,
		Health:       50,
		MaxHealth:    50,
		BaseAttack:   20,
		BaseDefense:  10,
		Experience:   0,
		Gold:         10,
		Weapon:       engine.WeaponsCatalog[1],
		Armor:        engine.ArmorsCatalog[1],
		ForestFights: 10,
	}

	monster := &engine.Monster{
		ID:         i18n.MonsterSewerRat,
		Name:       "Rato-do-Esgoto",
		Health:     10,
		MaxHealth:  10,
		Attack:     2,
		Defense:    1,
		XPReward:   25,
		GoldReward: 15,
	}

	res := ce.Attack(p, monster)
	if !res.MonsterDefeated {
		t.Fatalf("esperado que o monstro fosse derrotado no primeiro turno")
	}

	if p.Experience != 25 {
		t.Fatalf("esperado XP do jogador = 25, obtido: %d", p.Experience)
	}

	if p.Gold != 25 {
		t.Fatalf("esperado Gold do jogador = 25 (10 + 15), obtido: %d", p.Gold)
	}
}

func TestCombatEngine_FleeAttempt(t *testing.T) {
	rng := rand.New(rand.NewSource(200))
	ce := engine.NewCombatEngine(rng)

	p := &engine.Player{
		Health:      50,
		MaxHealth:   50,
		BaseDefense: 5,
	}

	monster := &engine.Monster{
		Name:    "Covarde Goblin",
		Prefix:  "Covarde",
		Attack:  5,
		Defense: 2,
	}

	// Com afixo Covarde, a taxa de fuga é 80%
	res := ce.AttemptFlee(p, monster)
	if !res.FledSuccessfully && !res.PlayerDefeated && res.MonsterDamageDealt == 0 {
		t.Fatalf("resultado de fuga inconsistente")
	}
}

func TestProgression_LevelUp(t *testing.T) {
	p := &engine.Player{
		Level:       1,
		Experience:  150, // Requisito para Nível 2 é 100 XP
		Gold:        60,  // Custo é 50 Gold
		Health:      20,
		MaxHealth:   50,
		BaseAttack:  10,
		BaseDefense: 5,
	}

	can, msg := engine.CanLevelUp(p)
	if !can {
		t.Fatalf("jogador deveria poder subir de nível: %s", msg)
	}

	err := engine.LevelUp(p)
	if err != nil {
		t.Fatalf("erro inesperado ao subir de nível: %v", err)
	}

	if p.Level != 2 {
		t.Fatalf("esperado nível 2, obtido %d", p.Level)
	}
	if p.Gold != 10 {
		t.Fatalf("esperado 10 moedas de ouro restantes, obtido %d", p.Gold)
	}
	if p.MaxHealth != 65 {
		t.Fatalf("esperado MaxHealth = 65 (50 + 15), obtido %d", p.MaxHealth)
	}
	if p.Health != 65 {
		t.Fatalf("ao subir de nível a vida deve ser completamente restaurada, obtido %d", p.Health)
	}
}

func TestEconomyService(t *testing.T) {
	econ := engine.NewEconomyService()
	p := &engine.Player{
		Gold:       100,
		BankGold:   50,
		Experience: 1000,
	}

	// Depósito
	err := econ.Deposit(p, 40)
	if err != nil || p.Gold != 60 || p.BankGold != 90 {
		t.Fatalf("falha no depósito: gold=%d, bank=%d, err=%v", p.Gold, p.BankGold, err)
	}

	// Saque
	err = econ.Withdraw(p, 20)
	if err != nil || p.Gold != 80 || p.BankGold != 70 {
		t.Fatalf("falha no saque: gold=%d, bank=%d, err=%v", p.Gold, p.BankGold, err)
	}

	// Penalidade de Morte
	lostGold, lostXP := econ.ProcessDeathPenalty(p)
	if lostGold != 80 || p.Gold != 0 {
		t.Fatalf("esperado perda total do ouro na bolsa (80), obtido %d", lostGold)
	}
	if p.BankGold != 70 {
		t.Fatalf("ouro no banco deve estar 100%% seguro (70), obtido %d", p.BankGold)
	}
	if lostXP != 100 || p.Experience != 900 {
		t.Fatalf("esperado perda de 10%% do XP (100), restante 900, obtido lost=%d exp=%d", lostXP, p.Experience)
	}
}

func TestTurnManager(t *testing.T) {
	tm := engine.NewTurnManager()
	p := &engine.Player{
		ForestFights: 3,
		LastLoginDay: "2026-08-24",
	}

	// Novo dia
	applied := tm.CheckAndApplyNewDay(p, "2026-08-25")
	if !applied || p.ForestFights != 15 || p.LastLoginDay != "2026-08-25" {
		t.Fatalf("falha ao renovar turnos no novo dia")
	}

	// Consumo de turnos
	err := tm.ConsumeFight(p)
	if err != nil || p.ForestFights != 14 {
		t.Fatalf("falha ao consumir turno: %v", err)
	}
}

func TestPlayerStorageConversion(t *testing.T) {
	sp := storage.Player{
		ID:           1,
		Username:     "Valente",
		Level:        3,
		Experience:   500,
		Gold:         120,
		BankGold:     300,
		Health:       80,
		MaxHealth:    80,
		Attack:       15,
		Defense:      8,
		WeaponID:     string(i18n.WeaponShortSword),
		ArmorID:      string(i18n.ArmorLeather),
		ForestFights: 12,
		DragonKills:  0,
		LastLoginDay: "2026-08-25",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	dp := engine.NewPlayerFromStorage(sp)
	if dp.Weapon.ID != i18n.WeaponShortSword || dp.Armor.ID != i18n.ArmorLeather {
		t.Fatalf("falha ao hidratar itens do jogador a partir do storage")
	}

	spConverted := dp.ToStorage()
	if spConverted.WeaponID != string(i18n.WeaponShortSword) || spConverted.ArmorID != string(i18n.ArmorLeather) {
		t.Fatalf("falha ao converter de volta para storage")
	}
}
