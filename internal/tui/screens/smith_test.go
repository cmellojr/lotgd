package screens

import (
	"math/rand"
	"strings"
	"testing"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSmithScreen_TradeInCreditOnPurchase(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	smith := NewSmithScreen(db, player)
	smith.SetPlayer(player)
	rng := rand.New(rand.NewSource(42))
	smith.SetRNG(rng)

	// Jogador possui a Adaga Afiada equipada (Value 50) e 150 de ouro
	player.Weapon = engine.WeaponsCatalog[1] // Adaga Afiada (Value: 50)
	player.Gold = 150

	// Selecionar Espada Curta (Index 2, Value: 200) e comprar (ENTER)
	smith.cursor = 2
	smith.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Com seed 42, percent = 76%, cotação da Adaga Afiada = (50 * 76 / 100) = 38 moedas de crédito
	// Custo líquido = 200 - 38 = 162 moedas.
	// Ouro do jogador (150) é menor que 162 -> deve recusar por ouro insuficiente informando crédito
	if player.Weapon.ID != i18n.WeaponDagger {
		t.Fatalf("arma não deveria ter sido alterada por falta de ouro líquido")
	}
	if !strings.Contains(smith.infoMsg, "Ouro insuficiente!") || !strings.Contains(smith.infoMsg, "crédito") {
		t.Fatalf("mensagem esperada sobre ouro insuficiente com crédito, obtido: %s", smith.infoMsg)
	}

	// Agora dando 200 moedas ao jogador
	player.Gold = 200
	smith.SetRNG(rand.New(rand.NewSource(42))) // mesma seed -> crédito 38
	smith.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if player.Weapon.ID != i18n.WeaponShortSword {
		t.Fatalf("esperado equipar Espada Curta, atual: %s", player.Weapon.ID)
	}
	// Ouro final: 200 - (200 - 38) = 38
	if player.Gold != 38 {
		t.Fatalf("esperado 38 moedas de ouro restantes, obtido: %d", player.Gold)
	}
	if !strings.Contains(smith.infoMsg, "Mestre Torin deu 38 moedas de crédito") {
		t.Fatalf("mensagem de sucesso deve mencionar o crédito concedido: %s", smith.infoMsg)
	}
}

func TestSmithScreen_SellEquippedItem(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	smith := NewSmithScreen(db, player)
	smith.SetPlayer(player)
	rng := rand.New(rand.NewSource(100))
	smith.SetRNG(rng)

	// 1. Tentar vender arma inicial (Pedaço de Pau, Value 0) -> deve negar
	player.Weapon = engine.WeaponsCatalog[0]
	smith.tab = smithTabWeapons
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

	if !strings.Contains(smith.infoMsg, "não possui nenhum equipamento") {
		t.Fatalf("esperado aviso de que não possui item para vender, obtido: %s", smith.infoMsg)
	}

	// 2. Equipar Montante de Aço (Value 600) e vender com 'V'
	player.Weapon = engine.WeaponsCatalog[3]
	player.Gold = 100
	goldBefore := player.Gold

	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

	if player.Weapon.ID != engine.WeaponsCatalog[0].ID {
		t.Fatalf("após a venda a arma deve retornar ao item inicial (Stick), obtido: %s", player.Weapon.ID)
	}
	if player.Gold <= goldBefore {
		t.Fatalf("ouro deveria ter aumentado após a venda, obtido: %d (antes: %d)", player.Gold, goldBefore)
	}
	if !strings.Contains(smith.infoMsg, "avaliou e comprou") {
		t.Fatalf("mensagem de confirmação de venda esperada, obtido: %s", smith.infoMsg)
	}
}
