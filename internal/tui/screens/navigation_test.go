package screens

import (
	"path/filepath"
	"testing"

	"lotgd/internal/engine"
	"lotgd/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func createTestDBAndPlayer(t *testing.T) (*storage.DB, *engine.Player) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_nav.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	stPlayer, err := db.CreatePlayer("TestHero", "pass123")
	if err != nil {
		db.Close()
		t.Fatalf("failed to create test player: %v", err)
	}

	player := engine.NewPlayerFromStorage(stPlayer)
	return db, player
}

func TestSmithScreen_TabShortcutsAndNavigation(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	smith := NewSmithScreen(db, player)

	// Default tab is smithTabWeapons
	if smith.tab != smithTabWeapons {
		t.Errorf("expected initial tab to be smithTabWeapons, got %v", smith.tab)
	}

	// Press "2" to switch to Armors
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if smith.tab != smithTabArmors {
		t.Errorf("expected tab smithTabArmors after '2', got %v", smith.tab)
	}

	// Press "3" to switch to Potions
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	if smith.tab != smithTabPotions {
		t.Errorf("expected tab smithTabPotions after '3', got %v", smith.tab)
	}

	// Press "1" to switch back to Weapons
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	if smith.tab != smithTabWeapons {
		t.Errorf("expected tab smithTabWeapons after '1', got %v", smith.tab)
	}

	// Navigation with j / k
	if smith.cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", smith.cursor)
	}

	// Press "j" -> cursor down
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if smith.cursor != 1 {
		t.Errorf("expected cursor 1 after 'j', got %d", smith.cursor)
	}

	// Press "k" -> cursor up
	smith.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if smith.cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", smith.cursor)
	}
}

func TestTavernScreen_VimNavigationAndNoKConflict(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	tavern := NewTavernScreen(db, player)

	if tavern.cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", tavern.cursor)
	}

	// Press "j" -> cursor moves down to 1
	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if tavern.cursor != 1 {
		t.Errorf("expected cursor 1 after 'j', got %d", tavern.cursor)
	}

	// Press "k" -> cursor moves up to 0, and infoMsg is NOT updated to duel message
	initialInfoMsg := tavern.infoMsg
	tavern.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if tavern.cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", tavern.cursor)
	}
	if tavern.infoMsg != initialInfoMsg {
		t.Errorf("expected infoMsg unchanged on 'k' navigation, got: %q", tavern.infoMsg)
	}
}

func TestGuildScreen_VimNavigation(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	guild := NewGuildScreen(db, player)

	if guild.cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", guild.cursor)
	}

	// Press "j" -> cursor moves to 1
	guild.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if guild.cursor != 1 {
		t.Errorf("expected cursor 1 after 'j', got %d", guild.cursor)
	}

	// Press "k" -> cursor moves to 0
	guild.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if guild.cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", guild.cursor)
	}
}

func TestChapelScreen_VimNavigation(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	chapel := NewChapelScreen(db, player)

	if chapel.cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", chapel.cursor)
	}

	// Press "j" -> cursor moves to 1
	chapel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if chapel.cursor != 1 {
		t.Errorf("expected cursor 1 after 'j', got %d", chapel.cursor)
	}

	// Press "k" -> cursor moves to 0
	chapel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if chapel.cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", chapel.cursor)
	}
}

func TestDragonScreen_VimNavigation(t *testing.T) {
	db, player := createTestDBAndPlayer(t)
	defer db.Close()

	dragon := NewDragonScreen(db, player, nil)
	dragon.SetPlayer(player) // Loads dragon state

	if dragon.cursor != 0 {
		t.Errorf("expected initial cursor 0, got %d", dragon.cursor)
	}

	// Press "j" -> cursor moves to 1
	dragon.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if dragon.cursor != 1 {
		t.Errorf("expected cursor 1 after 'j', got %d", dragon.cursor)
	}

	// Press "k" -> cursor moves to 0
	dragon.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if dragon.cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", dragon.cursor)
	}
}
