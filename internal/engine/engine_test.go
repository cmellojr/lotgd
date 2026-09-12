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

func TestCombatEngine_UsePotionEdgeCases(t *testing.T) {
	ce := engine.NewCombatEngine(nil)

	tests := []struct {
		name        string
		player      engine.Player
		wantHealth  int
		wantPotions int
		wantHealed  int
		wantError   bool
	}{
		{
			name:        "rejects when inventory is empty",
			player:      engine.Player{Health: 10, MaxHealth: 50},
			wantHealth:  10,
			wantPotions: 0,
			wantError:   true,
		},
		{
			name:        "rejects when health is full",
			player:      engine.Player{Health: 50, MaxHealth: 50, PotionsCount: 1},
			wantHealth:  50,
			wantPotions: 1,
			wantError:   true,
		},
		{
			name:        "caps healing at maximum health",
			player:      engine.Player{Health: 40, MaxHealth: 50, PotionsCount: 1},
			wantHealth:  50,
			wantPotions: 0,
			wantHealed:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := tt.player
			healed, err := ce.UsePotion(&player)

			if (err != nil) != tt.wantError {
				t.Fatalf("UsePotion() error = %v, want error: %v", err, tt.wantError)
			}
			if healed != tt.wantHealed || player.Health != tt.wantHealth || player.PotionsCount != tt.wantPotions {
				t.Fatalf("unexpected potion result: healed=%d health=%d potions=%d", healed, player.Health, player.PotionsCount)
			}
		})
	}
}

func TestEconomyService_RejectsInvalidTransfers(t *testing.T) {
	econ := engine.NewEconomyService()

	tests := []struct {
		name string
		call func(*engine.Player) error
		gold int
		bank int
	}{
		{
			name: "rejects zero deposit",
			call: func(p *engine.Player) error { return econ.Deposit(p, 0) },
			gold: 100,
			bank: 50,
		},
		{
			name: "rejects deposit above wallet balance",
			call: func(p *engine.Player) error { return econ.Deposit(p, 101) },
			gold: 100,
			bank: 50,
		},
		{
			name: "rejects negative withdrawal",
			call: func(p *engine.Player) error { return econ.Withdraw(p, -1) },
			gold: 100,
			bank: 50,
		},
		{
			name: "rejects withdrawal above bank balance",
			call: func(p *engine.Player) error { return econ.Withdraw(p, 51) },
			gold: 100,
			bank: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &engine.Player{Gold: tt.gold, BankGold: tt.bank}
			if err := tt.call(player); err == nil {
				t.Fatal("expected transfer to be rejected")
			}
			if player.Gold != tt.gold || player.BankGold != tt.bank {
				t.Fatalf("rejected transfer mutated balances: gold=%d bank=%d", player.Gold, player.BankGold)
			}
		})
	}
}

func TestProgression_LevelUp(t *testing.T) {
	p := &engine.Player{
		Level:       1,
		Experience:  150, // Requisito para Nível 2 é 100 XP
		Gold:        60,  // Ouro não é mais consumido no treino (ADR-0006)
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
	if p.Gold != 60 {
		t.Fatalf("esperado 60 moedas de ouro (sem alteração), obtido %d", p.Gold)
	}
	if p.MaxHealth != 65 {
		t.Fatalf("esperado MaxHealth = 65 (50 + 15), obtido %d", p.MaxHealth)
	}
	if p.Health != 65 {
		t.Fatalf("ao subir de nível a vida deve ser completamente restaurada, obtido %d", p.Health)
	}
	if !p.MasterFoughtToday {
		t.Fatalf("ao subir de nível a tentativa diária deve ser marcada como utilizada")
	}
}

func TestProgression_Boundaries(t *testing.T) {
	if req, ok := engine.NextLevelRequirement(12); ok || req != (engine.LevelRequirement{}) {
		t.Fatalf("expected no requirement beyond max level, got %+v, ok=%v", req, ok)
	}

	player := &engine.Player{Level: 1, Experience: 99, Gold: 50}
	if can, _ := engine.CanLevelUp(player); can {
		t.Fatal("player should not level up below the exact XP threshold")
	}

	player.Experience = 100
	if can, _ := engine.CanLevelUp(player); !can {
		t.Fatal("player should level up at the exact XP threshold")
	}
}

// TestMasterCombatWinRateProgression simula 1.000 combates entre um jogador adequadamente
// equipado para o seu nível e o Mestre de Treinamento correspondente (ADR-0006).
//
// Critério de Aceite:
// - Jogador com equipamento recomendado do nível: Taxa de vitória entre 65% e 85%.
// - Jogador sem equipamentos (atributos base apenas): Taxa de vitória < 35%.
func TestMasterCombatWinRateProgression(t *testing.T) {
	type gearBonus struct {
		weaponAtk int
		armorDef  int
	}

	// Bônus de equipamento recomendados por nível alvo (2 a 12)
	recommendedGear := map[int]gearBonus{
		2:  {weaponAtk: 3, armorDef: 2},
		3:  {weaponAtk: 5, armorDef: 4},
		4:  {weaponAtk: 10, armorDef: 6},
		5:  {weaponAtk: 14, armorDef: 9},
		6:  {weaponAtk: 18, armorDef: 11},
		7:  {weaponAtk: 23, armorDef: 14},
		8:  {weaponAtk: 27, armorDef: 17},
		9:  {weaponAtk: 34, armorDef: 21},
		10: {weaponAtk: 40, armorDef: 26},
		11: {weaponAtk: 47, armorDef: 30},
		12: {weaponAtk: 55, armorDef: 36},
	}

	const iterations = 1000

	// Simula para cada nível alvo de 2 a 12
	for targetLevel := 2; targetLevel <= 12; targetLevel++ {
		master, ok := engine.GetMasterForTargetLevel(targetLevel)
		if !ok {
			t.Fatalf("Mestre para nível alvo %d não encontrado", targetLevel)
		}

		// Calcula os atributos base acumulados do jogador no nível anterior
		baseHP := 20
		baseAtk := 5
		baseDef := 2
		for lvl := 2; lvl <= targetLevel-1; lvl++ {
			for _, req := range engine.LevelTable {
				if req.Level == lvl {
					baseHP += req.HealthGain
					baseAtk += req.AttackGain
					baseDef += req.DefGain
				}
			}
		}

		gear := recommendedGear[targetLevel]

		// 1. Simulação com jogador equipado
		equippedWins := 0
		for i := 0; i < iterations; i++ {
			rng := rand.New(rand.NewSource(int64(targetLevel*10000 + i)))
			ce := engine.NewCombatEngine(rng)

			player := &engine.Player{
				Username:    "HeroEquipped",
				Health:      baseHP,
				MaxHealth:   baseHP,
				BaseAttack:  baseAtk,
				BaseDefense: baseDef,
				Weapon:      engine.Item{PowerBonus: gear.weaponAtk},
				Armor:       engine.Item{PowerBonus: gear.armorDef},
			}

			masterMonster := &engine.Monster{
				Name:      master.Name,
				Health:    master.Health,
				MaxHealth: master.MaxHealth,
				Attack:    master.Attack,
				Defense:   master.Defense,
			}

			for player.IsAlive() && masterMonster.IsAlive() {
				res := ce.Attack(player, masterMonster)
				if res.MonsterDefeated {
					equippedWins++
					break
				}
				if res.PlayerDefeated {
					break
				}
			}
		}

		equippedWinRate := float64(equippedWins) / float64(iterations)

		// 2. Simulação com jogador desarmado/sem equipamento
		unequippedWins := 0
		for i := 0; i < iterations; i++ {
			rng := rand.New(rand.NewSource(int64(targetLevel*20000 + i)))
			ce := engine.NewCombatEngine(rng)

			player := &engine.Player{
				Username:    "HeroBare",
				Health:      baseHP,
				MaxHealth:   baseHP,
				BaseAttack:  baseAtk,
				BaseDefense: baseDef,
				Weapon:      engine.Item{PowerBonus: 0},
				Armor:       engine.Item{PowerBonus: 0},
			}

			masterMonster := &engine.Monster{
				Name:      master.Name,
				Health:    master.Health,
				MaxHealth: master.MaxHealth,
				Attack:    master.Attack,
				Defense:   master.Defense,
			}

			for player.IsAlive() && masterMonster.IsAlive() {
				res := ce.Attack(player, masterMonster)
				if res.MonsterDefeated {
					unequippedWins++
					break
				}
				if res.PlayerDefeated {
					break
				}
			}
		}

		unequippedWinRate := float64(unequippedWins) / float64(iterations)

		if equippedWinRate < 0.65 || equippedWinRate > 0.85 {
			t.Errorf("Nível %d (%s): taxa de vitória equipado fora do intervalo [0.65, 0.85]: %.3f",
				targetLevel, master.Name, equippedWinRate)
		}

		if unequippedWinRate >= 0.35 {
			t.Errorf("Nível %d (%s): taxa de vitória desarmado deve ser < 0.35, obtido: %.3f",
				targetLevel, master.Name, unequippedWinRate)
		}
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

func TestTurnManager_RejectsExhaustedFights(t *testing.T) {
	tm := engine.NewTurnManager()
	player := &engine.Player{ForestFights: 0}

	if err := tm.ConsumeFight(player); err == nil {
		t.Fatal("expected exhausted fights to be rejected")
	}
	if player.ForestFights != 0 {
		t.Fatalf("rejected fight changed remaining fights to %d", player.ForestFights)
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
