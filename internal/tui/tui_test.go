package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"lotgd/internal/storage"
	"lotgd/internal/ui"
)

func TestTUI_MainModel_Initialization(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_tui.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()

	model := NewMainModel(db)
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
