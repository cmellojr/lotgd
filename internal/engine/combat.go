package engine

import (
	"fmt"
	"math/rand"

	"lotgd/internal/i18n"
)

// CombatEngine gerencia o fluxo de combates por turnos entre o herói e os monstros.
//
// Didática Go: Receber um gerador `*rand.Rand` via struct permite injetar sementes (seeds)
// previsíveis nos testes unitários, garantindo testes 100% determinísticos sem depender
// do estado global de aleatoriedade do sistema.
type CombatEngine struct {
	rng *rand.Rand
}

// NewCombatEngine cria e inicializa uma nova instância do motor de combate.
// Se `rng` for nil, um novo gerador pseudo-aleatório é instanciado automaticamente.
func NewCombatEngine(rng *rand.Rand) *CombatEngine {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &CombatEngine{rng: rng}
}

// TurnResult detalha o desfecho completo de uma rodada de combate (ação do jogador,
// acertos críticos, dano causado/sofrido, derrotas, fuga e recompensas acumuladas).
//
// Didática Go: As tags `json:"..."` permitem serializar e deserializar os resultados do turno
// caso seja necessário registrar logs de combate ou salvar estados temporários em JSON.
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

// CalculateDamage aplica a fórmula canônica de cálculo de dano do Game Design Document (GDD):
// Dano Base = max(1, (Ataque_Atacante + Rnd(1, 4)) - Defesa_Defensor).
// Há uma chance estocástica de 10% de Acerto Crítico (dano base amplificado em +50%).
//
// Didática Go: O método retorna múltiplos valores `(int, bool)` para informar simultaneamente
// o valor numérico do dano e se houve um evento de acerto crítico.
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

// Attack executa um turno completo de troca de golpes em combate direto.
//
// Fluxo de Execução:
// 1. O herói ataca primeiro com base no seu ataque total (Ataque Base + Arma).
// 2. Se a vida do monstro chegar a zero, o combate se encerra com vitória, concedendo XP e Ouro.
// 3. Caso o monstro sobreviva, ele desfere seu contra-ataque contra a defesa total do jogador.
// 4. Se a vida do herói chegar a zero, a derrota é sinalizada no TurnResult para processamento moratório.
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

		// Aplica recompensas diretamente ao estado do jogador em memória
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

// AttemptFlee processa a tentativa de recuo tático do combate.
//
// Regras de Fuga:
// - Chance base de fuga: 50%.
// - Monstros com o afixo "Covarde": chance aumentada para 80%.
// - O Dragão Ancestral dificulta a fuga: apenas 20% de chance.
// - Falha na fuga concede um contra-ataque livre de oportunidade ao inimigo.
func (ce *CombatEngine) AttemptFlee(player *Player, monster *Monster) TurnResult {
	res := TurnResult{}

	// Modificadores de chance de fuga de acordo com o tipo/afixo do monstro
	fleeChance := 0.50
	if monster.Prefix == "Covarde" {
		fleeChance = 0.80
	}

	if monster.IsDragon {
		fleeChance = 0.20
	}

	if ce.rng.Float64() < fleeChance {
		res.FledSuccessfully = true
		res.Message = fmt.Sprintf("Você conseguiu recuar estrategicamente para as sombras e escapar de %s!", monster.Name)
		return res
	}

	// Falha na fuga: o monstro ataca de graça
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

// UsePotion consome uma poção de cura da bolsa do jogador sem gastar turno de combate livre.
//
// Didática Go: Retorna um erro explícito `(int, error)` para indicar falha em pré-condições
// (ex: sem poções na bolsa ou vida máxima já atingida), seguindo o padrão idiomático de Go.
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
