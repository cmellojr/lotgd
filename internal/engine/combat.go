package engine

import (
	"fmt"
	"math/rand"

	"lotgd/internal/i18n"
)

// CombatEngine gerencia o fluxo de combates por turnos entre o jogador e monstros.
//
// Didática Go: Receber um gerador `*rand.Rand` via struct permite injetar sementes (seeds)
// previsíveis nos testes unitários, garantindo testes 100% determinísticos sem flaky tests.
type CombatEngine struct {
	rng *rand.Rand
}

// NewCombatEngine cria uma nova instância do motor de combate.
// Se rng for nil, um gerador padrão é inicializado.
func NewCombatEngine(rng *rand.Rand) *CombatEngine {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &CombatEngine{rng: rng}
}

// TurnResult detalha o desfecho de uma rodada de combate (ação do jogador e contra-ataque).
type TurnResult struct {
	PlayerDamageDealt  int    `json:"player_damage_dealt"`
	PlayerCritical     bool   `json:"player_critical"`
	MonsterDamageDealt int    `json:"monster_damage_dealt"`
	MonsterCritical    bool   `json:"monster_critical"`
	MonsterDefeated    bool   `json:"monster_defeated"`
	PlayerDefeated     bool   `json:"player_defeated"`
	FledSuccessfully   bool   `json:"fled_successfully"`
	XPBonus            int    `json:"xp_bonus"`
	GoldBonus          int    `json:"gold_bonus"`
	Message            string `json:"message"`
}

// CalculateDamage aplica a fórmula canônica do GDD:
// Dano = max(1, (ATK_atacante + Rnd(1, 4)) - DEF_defensor)
// Há 10% de chance de Acerto Crítico (dano base aumentado em 50%).
func (ce *CombatEngine) CalculateDamage(atk, def int) (int, bool) {
	roll := ce.rng.Intn(4) + 1 // Rnd(1, 4)
	isCritical := ce.rng.Float64() < 0.10

	effectiveATK := atk + roll
	if isCritical {
		effectiveATK = int(float64(effectiveATK) * 1.5)
	}

	damage := effectiveATK - def
	if damage < 1 {
		damage = 1
	}

	return damage, isCritical
}

// Attack executa um turno completo de troca de golpes entre jogador e monstro.
//
// 1. O jogador ataca primeiro.
// 2. Se o monstro morrer, a batalha encerra imediatamente com vitória e recompensas.
// 3. Se o monstro sobreviver, ele desfere seu contra-ataque.
// 4. Se o jogador morrer, o status é registrado para processamento de derrota.
func (ce *CombatEngine) Attack(player *Player, monster *Monster) TurnResult {
	res := TurnResult{}

	// --- 1. Ataque do Jogador ---
	pDmg, pCrit := ce.CalculateDamage(player.TotalAttack(), monster.Defense)
	res.PlayerDamageDealt = pDmg
	res.PlayerCritical = pCrit

	monster.Health -= pDmg
	if monster.Health <= 0 {
		monster.Health = 0
		res.MonsterDefeated = true
		res.XPBonus = monster.XPReward
		res.GoldBonus = monster.GoldReward

		// Aplica recompensas ao jogador
		player.Experience += monster.XPReward
		player.Gold += monster.GoldReward

		if monster.IsDragon {
			player.DragonKills++
			res.Message = fmt.Sprintf("VITÓRIA LENDÁRIA! Você desferiu o golpe fatal e derrotou %s!", monster.Name)
		} else {
			res.Message = fmt.Sprintf("Você derrotou %s e ganhou %d XP e %d moedas de ouro!", monster.Name, monster.XPReward, monster.GoldReward)
		}
		return res
	}

	// --- 2. Contra-ataque do Monstro ---
	mDmg, mCrit := ce.CalculateDamage(monster.Attack, player.TotalDefense())
	res.MonsterDamageDealt = mDmg
	res.MonsterCritical = mCrit

	player.Health -= mDmg
	if player.Health <= 0 {
		player.Health = 0
		res.PlayerDefeated = true
		res.Message = fmt.Sprintf("%s desferiu um golpe mortal! Você sucumbiu na escuridão...", monster.Name)
		return res
	}

	critMsg := ""
	if pCrit {
		critMsg = " [GOLPE CRÍTICO!]"
	}
	res.Message = fmt.Sprintf("Você causou %d de dano%s. %s contra-atacou causando %d de dano.", pDmg, critMsg, monster.Name, mDmg)
	return res
}

// AttemptFlee tenta escapar da batalha.
//
// Chance base de fuga = 50%.
// Se falhar na fuga, o monstro ganha um ataque livre de oportunidade!
func (ce *CombatEngine) AttemptFlee(player *Player, monster *Monster) TurnResult {
	res := TurnResult{}

	// Monstros com afixo "Covarde" aumentam a chance de fuga para 80%
	fleeChance := 0.50
	if monster.Prefix == "Covarde" {
		fleeChance = 0.80
	}

	// Dragão não permite fuga fácil (20% de chance)
	if monster.IsDragon {
		fleeChance = 0.20
	}

	if ce.rng.Float64() < fleeChance {
		res.FledSuccessfully = true
		res.Message = fmt.Sprintf("Você conseguiu recuar estrategicamente para as sombras e escapar de %s!", monster.Name)
		return res
	}

	// Falha na fuga: monstro contra-ataca livremente
	mDmg, mCrit := ce.CalculateDamage(monster.Attack, player.TotalDefense())
	res.MonsterDamageDealt = mDmg
	res.MonsterCritical = mCrit
	player.Health -= mDmg

	if player.Health <= 0 {
		player.Health = 0
		res.PlayerDefeated = true
		res.Message = fmt.Sprintf("Você tropeçou ao tentar fugir! %s aproveitou e desferiu um golpe letal!", monster.Name)
		return res
	}

	res.Message = fmt.Sprintf("Falha ao fugir! %s bloqueou seu caminho e te atingiu causando %d de dano.", monster.Name, mDmg)
	return res
}

// UsePotion consome uma poção do inventário do jogador, restaurando vida sem gastar turno de combate livre.
func (ce *CombatEngine) UsePotion(player *Player) (int, error) {
	if player.PotionsCount <= 0 {
		return 0, fmt.Errorf("você não possui poções na bolsa")
	}

	if player.Health >= player.MaxHealth {
		return 0, fmt.Errorf("sua vida já está cheia")
	}

	potion, found := FindPotion(i18n.PotionHealth)
	heal := 30
	if found {
		heal = potion.HealAmount
	}

	oldHP := player.Health
	player.Health += heal
	if player.Health > player.MaxHealth {
		player.Health = player.MaxHealth
	}

	actualHealed := player.Health - oldHP
	player.PotionsCount--

	return actualHealed, nil
}
