package engine

import (
	"fmt"
)

// EconomyService centraliza as operações financeiras entre a bolsa de mão e o cofre do banco.
type EconomyService struct{}

// NewEconomyService cria uma nova instância do serviço de economia.
func NewEconomyService() *EconomyService {
	return &EconomyService{}
}

// Deposit transfere moedas da bolsa do jogador para o cofre seguro do banco.
func (e *EconomyService) Deposit(p *Player, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("a quantia de depósito deve ser positiva")
	}

	if p.Gold < amount {
		return fmt.Errorf("ouro insuficiente na bolsa para depósito (disponível: %d)", p.Gold)
	}

	p.Gold -= amount
	p.BankGold += amount
	return nil
}

// Withdraw retira moedas do banco para a bolsa do jogador.
func (e *EconomyService) Withdraw(p *Player, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("a quantia de saque deve ser positiva")
	}

	if p.BankGold < amount {
		return fmt.Errorf("saldo bancário insuficiente para saque (disponível: %d)", p.BankGold)
	}

	p.BankGold -= amount
	p.Gold += amount
	return nil
}

// ProcessDeathPenalty aplica a regra clássica de morte em BBS Door Games:
// - O jogador perde 100% do ouro que carregava na bolsa (o banco permanece 100% seguro).
// - Perde 10% da experiência atual (sem rebaixamento de nível).
// - Vida é restaurada para 1 HP na capela do Frei Anselmo.
func (e *EconomyService) ProcessDeathPenalty(p *Player) (lostGold int, lostXP int) {
	lostGold = p.Gold
	p.Gold = 0

	lostXP = int(float64(p.Experience) * 0.10)
	p.Experience -= lostXP
	if p.Experience < 0 {
		p.Experience = 0
	}

	p.Health = 1 // Ressuscita com 1 ponto de vida
	return lostGold, lostXP
}
