package engine

import (
	"fmt"
)

// EconomyService centraliza as operações financeiras e bancárias do jogo.
//
// Didática Go: Encapsular a lógica de negócios em uma struct de serviço sem campos (`struct{}`)
// permite organizar métodos relacionados (Deposit, Withdraw, ProcessDeathPenalty) de forma limpa,
// facilitando a extensão e garantindo que todas as regras de saldo sejam validadas em um único lugar.
type EconomyService struct{}

// NewEconomyService cria uma nova instância do serviço financeiro.
func NewEconomyService() *EconomyService {
	return &EconomyService{}
}

// Deposit transfere moedas de ouro da bolsa de mão do herói para o cofre seguro do banco do vilarejo.
//
// Regras de Negócio:
// - A quantia deve ser estritamente positiva.
// - O jogador deve possuir saldo suficiente na bolsa (`Gold >= amount`).
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

// Withdraw retira moedas de ouro do cofre do banco e as transfere para a bolsa de mão do herói.
//
// Regras de Negócio:
// - A quantia deve ser estritamente positiva.
// - O saldo bancário deve ser suficiente (`BankGold >= amount`).
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

// ProcessDeathPenalty aplica a regra moratória clássica de morte em BBS Door Games:
// - Perda de 100% do ouro que o herói carregava na bolsa (o saldo do banco fica 100% protegido).
// - Perda de 10% da experiência acumulada no nível atual (sem rebaixamento de nível).
// - Restauração da vida para 1 HP na Capela do Frei Anselmo.
//
// Didática Go: Retorna múltiplos valores `(lostGold, lostXP)` informando exatamente os valores
// deduzidos para exibição informativa na tela de Game Over.
func (e *EconomyService) ProcessDeathPenalty(p *Player) (lostGold int, lostXP int) {
	lostGold = p.Gold
	p.Gold = 0

	lostXP = int(float64(p.Experience) * 0.10)
	p.Experience -= lostXP
	if p.Experience < 0 {
		p.Experience = 0
	}

	p.Health = 1 // Ressuscita o jogador com 1 ponto de vida
	return lostGold, lostXP
}
