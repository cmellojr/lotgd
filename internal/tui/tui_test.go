package tui

import (
	"path/filepath"
	"testing"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUI_MainModel_Initialization(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_tui.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()

	model := NewMainModel(db, nil)
	if model == nil {
		t.Fatalf("expected non-nil MainModel")
	}

	// Verifica tela inicial
	if model.currentScreen != ScreenLogin {
		t.Errorf("expected initial screen %s, got %s", ScreenLogin, model.currentScreen)
	}

	// Testa redimensionamento
	resizedModel, cmd := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if cmd != nil {
		t.Errorf("expected nil cmd on window resize")
	}
	m := resizedModel.(*MainModel)
	if m.width != 80 || m.height != 24 {
		t.Errorf("expected dimensions 80x24, got %dx%d", m.width, m.height)
	}

	// Testa renderização de View da tela de login
	viewOutput := m.View()
	if viewOutput == "" {
		t.Errorf("expected non-empty view output")
	}

	// Testa transição de telas via ChangeScreenMsg
	navModel, _ := m.Update(ui.ChangeScreenMsg{Screen: ui.ScreenTown})
	nm := navModel.(*MainModel)
	if nm.currentScreen != ScreenTown {
		t.Errorf("expected current screen %s, got %s", ScreenTown, nm.currentScreen)
	}
}

func TestTUI_DeathPenalty_AppliedOnce(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_death.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()

	model := NewMainModel(db, nil)

	sp, err := db.CreatePlayer("hero", "pass123")
	if err != nil {
		t.Fatalf("failed to create player: %v", err)
	}

	// Login / update player
	p := engine.NewPlayerFromStorage(sp)
	p.Gold = 100
	p.Experience = 100
	p.MaxHealth = 20
	p.Health = 0

	updatedModel, _ := model.Update(PlayerUpdatedMsg{Player: p})
	m := updatedModel.(*MainModel)

	goldBefore := m.player.Gold
	xpBefore := m.player.Experience

	if goldBefore != 100 || xpBefore != 100 {
		t.Fatalf("expected 100 gold and 100 xp, got gold %d, xp %d", goldBefore, xpBefore)
	}

	// Transition to ScreenGameOver
	overModel, _ := m.Update(ui.ChangeScreenMsg{Screen: ui.ScreenGameOver})
	om := overModel.(*MainModel)

	goldAfterDeath := om.player.Gold
	xpAfterDeath := om.player.Experience

	// EconomyService ProcessDeathPenalty loses all pocket gold (100 -> 0) and 10% XP (100 -> 90)
	if goldAfterDeath != 0 {
		t.Errorf("expected 0 gold after death penalty, got %d", goldAfterDeath)
	}
	if xpAfterDeath != 90 {
		t.Errorf("expected 90 xp after death penalty, got %d", xpAfterDeath)
	}

	// Simulate screen transitions / syncPlayerState to ensure death penalty is NOT re-applied
	navModel, _ := om.Update(ui.ChangeScreenMsg{Screen: ui.ScreenTown})
	nm := navModel.(*MainModel)

	if nm.player.Gold != goldAfterDeath {
		t.Errorf("gold changed after screen transition: expected %d, got %d", goldAfterDeath, nm.player.Gold)
	}
	if nm.player.Experience != xpAfterDeath {
		t.Errorf("xp changed after screen transition: expected %d, got %d", xpAfterDeath, nm.player.Experience)
	}

	// Re-trigger syncPlayerState directly via PlayerUpdatedMsg
	syncModel, _ := nm.Update(PlayerUpdatedMsg{Player: nm.player})
	sm := syncModel.(*MainModel)

	if sm.player.Gold != goldAfterDeath {
		t.Errorf("gold changed after syncPlayerState: expected %d, got %d", goldAfterDeath, sm.player.Gold)
	}
	if sm.player.Experience != xpAfterDeath {
		t.Errorf("xp changed after syncPlayerState: expected %d, got %d", xpAfterDeath, sm.player.Experience)
	}
}
