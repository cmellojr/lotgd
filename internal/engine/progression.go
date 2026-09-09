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

// LevelTable especifica a curva de progressão do nível 1 ao nível 10 com o Mestre Tobias na Guilda.
//
// Didática Go: A tabela estática de structs desacopla as fórmulas numéricas e permite ajustes
// rápidos de balanceamento no GDD sem alterar o fluxo do algoritmo de verificação.
var LevelTable = []LevelRequirement{
	{Level: 1, RequiredXP: 0, CostGold: 0, HealthGain: 0, AttackGain: 0, DefGain: 0},
	{Level: 2, RequiredXP: 100, CostGold: 50, HealthGain: 15, AttackGain: 2, DefGain: 2},
	{Level: 3, RequiredXP: 300, CostGold: 150, HealthGain: 20, AttackGain: 3, DefGain: 2},
	{Level: 4, RequiredXP: 700, CostGold: 350, HealthGain: 25, AttackGain: 4, DefGain: 3},
	{Level: 5, RequiredXP: 1500, CostGold: 800, HealthGain: 30, AttackGain: 5, DefGain: 4},
	{Level: 6, RequiredXP: 3000, CostGold: 1600, HealthGain: 35, AttackGain: 6, DefGain: 5},
	{Level: 7, RequiredXP: 5500, CostGold: 3000, HealthGain: 40, AttackGain: 7, DefGain: 6},
	{Level: 8, RequiredXP: 9000, CostGold: 5000, HealthGain: 45, AttackGain: 8, DefGain: 7},
	{Level: 9, RequiredXP: 14000, CostGold: 8000, HealthGain: 50, AttackGain: 10, DefGain: 8},
	{Level: 10, RequiredXP: 22000, CostGold: 12000, HealthGain: 60, AttackGain: 12, DefGain: 10},
}

// MaxLevel define o teto máximo de nível alcançável antes do confronto final contra o Dragão.
const MaxLevel = 10

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

// CanLevelUp valida se o herói possui a experiência e o ouro mínimos necessários para ser promovido.
func CanLevelUp(p *Player) (bool, string) {
	req, ok := NextLevelRequirement(p.Level)
	if !ok {
		return false, "Você já alcançou o nível máximo de maestria!"
	}

	if p.Experience < req.RequiredXP {
		return false, fmt.Sprintf("Experiência insuficiente. Necessário: %d XP (Você tem: %d XP)", req.RequiredXP, p.Experience)
	}

	if p.Gold < req.CostGold {
		return false, fmt.Sprintf("Ouro insuficiente para o treinamento. Custo: %d moedas (Você tem: %d)", req.CostGold, p.Gold)
	}

	return true, fmt.Sprintf("Pronto para avançar para o Nível %d com o Mestre Tobias!", req.Level)
}

// LevelUp executa o treinamento na Guilda dos Aventureiros, deduzindo o ouro e aumentando os atributos de forma permanente.
//
// Didática Go: O método altera o ponteiro do jogador em memória (`*Player`), atualizando o nível e curando-o totalmente.
func LevelUp(p *Player) error {
	req, ok := NextLevelRequirement(p.Level)
	if !ok {
		return fmt.Errorf("jogador já se encontra no nível máximo")
	}

	if p.Experience < req.RequiredXP {
		return fmt.Errorf("experiência insuficiente para o próximo nível")
	}

	if p.Gold < req.CostGold {
		return fmt.Errorf("ouro insuficiente para o treinamento na guilda")
	}

	// Deduz o custo de treinamento em ouro e aplica a progressão de atributos
	p.Gold -= req.CostGold
	p.Level = req.Level
	p.MaxHealth += req.HealthGain
	p.Health = p.MaxHealth // Cura completa e imediata ao subir de nível
	p.BaseAttack += req.AttackGain
	p.BaseDefense += req.DefGain

	return nil
}
