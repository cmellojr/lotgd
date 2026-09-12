package engine

import (
	"fmt"
	"time"
)

// DailyForestFights define a cota canônica de 15 lutas de exploração por dia do calendário real.
const DailyForestFights = 15

// TurnManager controla a cota diária de combates do herói e gerencia o reset do "Novo Dia".
//
// Didática Go: O TurnManager atua como um serviço utilitário para garantir que nenhum jogador
// exceda o limite diário de combates, reabastecendo o contador a cada virada de dia no calendário.
type TurnManager struct{}

// NewTurnManager instancia o gerenciador de turnos de combate.
func NewTurnManager() *TurnManager {
	return &TurnManager{}
}

// CurrentDateString retorna a data atual do servidor formatada no padrão ISO 8601 (YYYY-MM-DD).
func CurrentDateString() string {
	return time.Now().Format("2006-01-02")
}

// CheckAndApplyNewDay verifica se o herói está realizando login em uma nova data do calendário.
//
// Em caso afirmativo, restaura os 15 turnos de combate da cota diária, cura completamente a saúde
// do personagem e atualiza o registro `LastLoginDay`. Esta função é a ÚNICA fonte de verdade (Single Source of Truth)
// para o ciclo do Novo Dia no jogo.
func (tm *TurnManager) CheckAndApplyNewDay(p *Player, today string) bool {
	if today == "" {
		today = CurrentDateString()
	}

	if p.LastLoginDay != today {
		p.LastLoginDay = today
		p.ForestFights = DailyForestFights
		p.Health = p.MaxHealth
		p.MasterFoughtToday = false
		return true // Novo dia aplicado com sucesso!
	}
	return false
}

// ConsumeFight debita 1 turno de combate da cota diária do jogador.
//
// Retorna um erro amigável se a cota já tiver sido completamente esgotada no dia corrente.
func (tm *TurnManager) ConsumeFight(p *Player) error {
	if p.ForestFights <= 0 {
		return fmt.Errorf("você já gastou todos os seus turnos de exploração por hoje. Descanse na taverna até o Novo Dia")
	}
	p.ForestFights--
	return nil
}
