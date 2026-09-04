package screens

import (
	"strings"
	"testing"

	"lotgd/internal/engine"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTavernScreen_CassandraFlirtLimit(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	tavern := NewTavernScreen(db, player)
	tavern.SetPlayer(player)

	// Scenario 1: Player has 10 gold, 15 forest fights (<= engine.DailyForestFights).
	// Buying flirt gives 1 extra fight (+1) -> 16 fights, costs 10 gold.
	player.Gold = 100
	player.ForestFights = engine.DailyForestFights // 15

	// Select Cassandra (shortcut 'C' or cursor = 1 + ENTER)
	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if player.Gold != 90 {
		t.Errorf("expected gold to be 90, got %d", player.Gold)
	}
	if player.ForestFights != engine.DailyForestFights+1 {
		t.Errorf("expected ForestFights to be %d, got %d", engine.DailyForestFights+1, player.ForestFights)
	}
	if !strings.Contains(tavern.infoMsg, "+1 Luta na Floresta!") {
		t.Errorf("expected success msg containing '+1 Luta na Floresta!', got: %q", tavern.infoMsg)
	}

	// Scenario 2: Player attempts to flirt AGAIN on the same day when ForestFights > engine.DailyForestFights (now 16 > 15).
	// Flirt must be blocked.
	goldBefore := player.Gold
	fightsBefore := player.ForestFights

	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if player.Gold != goldBefore {
		t.Errorf("expected gold unchanged (%d), got %d", goldBefore, player.Gold)
	}
	if player.ForestFights != fightsBefore {
		t.Errorf("expected ForestFights unchanged (%d), got %d", fightsBefore, player.ForestFights)
	}
	expectedCapMsg := "Cassandra já te inspirou hoje. Volte amanhã para nova dose de coragem!"
	if tavern.infoMsg != expectedCapMsg {
		t.Errorf("expected infoMsg %q, got %q", expectedCapMsg, tavern.infoMsg)
	}
}

func TestTavernScreen_CassandraInsufficientGold(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	tavern := NewTavernScreen(db, player)
	tavern.SetPlayer(player)

	player.Gold = 5
	player.ForestFights = 10

	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if player.Gold != 5 {
		t.Errorf("expected gold unchanged (5), got %d", player.Gold)
	}
	if player.ForestFights != 10 {
		t.Errorf("expected ForestFights unchanged (10), got %d", player.ForestFights)
	}
	if !strings.Contains(tavern.infoMsg, "Volte quando tiver pelo menos 10 moedas de ouro") {
		t.Errorf("expected insufficient gold message, got %q", tavern.infoMsg)
	}
}

func TestTavernScreen_RedKnightMessage(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	tavern := NewTavernScreen(db, player)
	tavern.SetPlayer(player)

	// Press 'D' to challenge Red Knight
	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	expectedMsg := "O Cavaleiro Vermelho ergue a viseira: 'Você ainda não possui a tempera necessária para cruzar lâminas comigo. Volte quando estiver mais experiente!'"
	if tavern.infoMsg != expectedMsg {
		t.Errorf("expected infoMsg %q, got %q", expectedMsg, tavern.infoMsg)
	}
}
