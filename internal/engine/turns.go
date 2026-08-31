package engine

import (
	"fmt"
	"time"
)

// DailyForestFights define a cota canônica de 15 lutas por dia na floresta.
const DailyForestFights = 15

// TurnManager controla a cota diária de ações do herói e a mecânica do "Novo Dia".
type TurnManager struct{}

// NewTurnManager instancia o gerenciador de turnos.
func NewTurnManager() *TurnManager {
	return &TurnManager{}
}

// CurrentDateString retorna a data do servidor formatada em YYYY-MM-DD.
func CurrentDateString() string {
	return time.Now().Format("2006-01-02")
}

// CheckAndApplyNewDay verifica se o aventureiro está logando em um novo dia do calendário.
// Em caso afirmativo, restaura os 15 turnos de combate e atualiza a data de login.
func (tm *TurnManager) CheckAndApplyNewDay(p *Player, today string) bool {
	if today == "" {
		today = CurrentDateString()
	}

	if p.LastLoginDay != today {
		p.LastLoginDay = today
		p.ForestFights = DailyForestFights
		return true // Novo dia aplicado!
	}
	return false
}

// ConsumeFight consome 1 turno de combate na floresta.
func (tm *TurnManager) ConsumeFight(p *Player) error {
	if p.ForestFights <= 0 {
		return fmt.Errorf("você já gastou todos os seus turnos de exploração por hoje. Descanse na taverna até o Novo Dia")
	}
	p.ForestFights--
	return nil
}
