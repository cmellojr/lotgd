package engine

import (
	"fmt"
)

// LevelRequirement define os pré-requisitos (XP e Ouro) e as recompensas de atributos ao atingir um novo nível.
type LevelRequirement struct {
	Level      int `json:"level"`
	RequiredXP int `json:"required_xp"`
	CostGold   int `json:"cost_gold"`
	HealthGain int `json:"health_gain"`
	AttackGain int `json:"attack_gain"`
	DefGain    int `json:"def_gain"`
}

// Master representa um Mestre de Treinamento na Guilda dos Aventureiros (LORD 1989 / ADR-0006).
type Master struct {
	TargetLevel int    `json:"target_level"`
	Name        string `json:"name"`
	Health      int    `json:"health"`
	MaxHealth   int    `json:"max_health"`
	Attack      int    `json:"attack"`
	Defense     int    `json:"defense"`
}

// LevelTable especifica a curva de progressão do nível 1 ao nível 12 (LORD 1989 / ADR-0006).
//
// Didática Go: A tabela estática de structs desacopla as fórmulas numéricas e permite ajustes
// rápidos de balanceamento no GDD sem alterar o fluxo do algoritmo de verificação.
var LevelTable = []LevelRequirement{
	{Level: 1, RequiredXP: 0, CostGold: 0, HealthGain: 0, AttackGain: 0, DefGain: 0},
	{Level: 2, RequiredXP: 100, CostGold: 0, HealthGain: 15, AttackGain: 2, DefGain: 2},
	{Level: 3, RequiredXP: 300, CostGold: 0, HealthGain: 20, AttackGain: 3, DefGain: 2},
	{Level: 4, RequiredXP: 700, CostGold: 0, HealthGain: 25, AttackGain: 4, DefGain: 3},
	{Level: 5, RequiredXP: 1500, CostGold: 0, HealthGain: 30, AttackGain: 5, DefGain: 4},
	{Level: 6, RequiredXP: 3000, CostGold: 0, HealthGain: 35, AttackGain: 6, DefGain: 5},
	{Level: 7, RequiredXP: 5500, CostGold: 0, HealthGain: 40, AttackGain: 7, DefGain: 6},
	{Level: 8, RequiredXP: 9000, CostGold: 0, HealthGain: 45, AttackGain: 8, DefGain: 7},
	{Level: 9, RequiredXP: 14000, CostGold: 0, HealthGain: 50, AttackGain: 10, DefGain: 8},
	{Level: 10, RequiredXP: 22000, CostGold: 0, HealthGain: 60, AttackGain: 12, DefGain: 10},
	{Level: 11, RequiredXP: 35000, CostGold: 0, HealthGain: 75, AttackGain: 15, DefGain: 12},
	{Level: 12, RequiredXP: 55000, CostGold: 0, HealthGain: 90, AttackGain: 18, DefGain: 15},
}

// MasterCatalog define a lista canônica dos 11 Mestres de Treinamento (do Nível 2 ao 12) conforme ADR-0006.
var MasterCatalog = []Master{
	{TargetLevel: 2, Name: "Mestre Halder", Health: 25, MaxHealth: 25, Attack: 7, Defense: 3},
	{TargetLevel: 3, Name: "Mestre Tobias", Health: 45, MaxHealth: 45, Attack: 12, Defense: 6},
	{TargetLevel: 4, Name: "Mestre Kaelen", Health: 75, MaxHealth: 75, Attack: 18, Defense: 10},
	{TargetLevel: 5, Name: "Mestra Vanya", Health: 110, MaxHealth: 110, Attack: 25, Defense: 15},
	{TargetLevel: 6, Name: "Mestre Roderick", Health: 150, MaxHealth: 150, Attack: 34, Defense: 21},
	{TargetLevel: 7, Name: "Mestra Elora", Health: 200, MaxHealth: 200, Attack: 44, Defense: 28},
	{TargetLevel: 8, Name: "Mestre Thorgrim", Health: 260, MaxHealth: 260, Attack: 55, Defense: 36},
	{TargetLevel: 9, Name: "Mestre Valerius", Health: 330, MaxHealth: 330, Attack: 68, Defense: 45},
	{TargetLevel: 10, Name: "Mestre Arthorian", Health: 420, MaxHealth: 420, Attack: 83, Defense: 55},
	{TargetLevel: 11, Name: "Mestra Ignis", Health: 530, MaxHealth: 530, Attack: 100, Defense: 68},
	{TargetLevel: 12, Name: "Mestre Turgon", Health: 680, MaxHealth: 680, Attack: 120, Defense: 85},
}

// GetMasterForTargetLevel obtém as estatísticas do Mestre correspondente ao nível alvo.
func GetMasterForTargetLevel(targetLevel int) (Master, bool) {
	for _, m := range MasterCatalog {
		if m.TargetLevel == targetLevel {
			return m, true
		}
	}
	return Master{}, false
}

// MaxLevel define o teto máximo de nível alcançável antes do confronto final contra o Dragão.
const MaxLevel = 12

// NextLevelRequirement obtém a estrutura de requisitos do próximo nível do jogador.
//
// Didática Go: Retorna `(LevelRequirement, bool)` no padrão `comma-ok`. Se o jogador já estiver no nível máximo,
// o booleano retornado será `false`.
func NextLevelRequirement(currentLevel int) (LevelRequirement, bool) {
	targetLevel := currentLevel + 1
	if targetLevel > MaxLevel {
		return LevelRequirement{}, false
	}

	for _, req := range LevelTable {
		if req.Level == targetLevel {
			return req, true
		}
	}
	return LevelRequirement{}, false
}

// CanLevelUp valida se o herói possui a experiência necessária e se ainda tem a tentativa diária disponível.
func CanLevelUp(p *Player) (bool, string) {
	req, ok := NextLevelRequirement(p.Level)
	if !ok {
		return false, "Você já alcançou o nível máximo de maestria!"
	}

	master, hasMaster := GetMasterForTargetLevel(req.Level)
	if !hasMaster {
		return false, "Mestre de Treinamento não encontrado para este nível."
	}

	if p.Experience < req.RequiredXP {
		return false, fmt.Sprintf("Experiência insuficiente. Necessário: %d XP (Você tem: %d XP)", req.RequiredXP, p.Experience)
	}

	if p.MasterFoughtToday {
		return false, fmt.Sprintf("Você já desafiou %s hoje! Apenas 1 tentativa por dia é permitida. Retorne no próximo alvorecer.", master.Name)
	}

	return true, fmt.Sprintf("Pronto para desafiar %s pelo Nível %d!", master.Name, req.Level)
}

// LevelUp executa a promoção de nível na Guilda dos Aventureiros após vencer o combate contra o Mestre.
//
// Didática Go: O método altera o ponteiro do jogador em memória (`*Player`), atualizando o nível, curando-o totalmente
// e marcando a tentativa diária como realizada.
func LevelUp(p *Player) error {
	req, ok := NextLevelRequirement(p.Level)
	if !ok {
		return fmt.Errorf("jogador já se encontra no nível máximo")
	}

	if p.Experience < req.RequiredXP {
		return fmt.Errorf("experiência insuficiente para o próximo nível")
	}

	// Aplica a progressão de atributos e cura completa
	p.Level = req.Level
	p.MaxHealth += req.HealthGain
	p.Health = p.MaxHealth // Cura completa e imediata ao subir de nível
	p.BaseAttack += req.AttackGain
	p.BaseDefense += req.DefGain
	p.MasterFoughtToday = true

	return nil
}
