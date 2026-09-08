package screens

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/ui"
)

// screenAfterKey envia uma tecla à tela e devolve o destino da ChangeScreenMsg
// resultante, ou "" se a tecla não provocou mudança de tela.
func screenAfterKey(t *testing.T, m tea.Model, key string) ui.ScreenID {
	t.Helper()
	var msg tea.KeyMsg
	switch key {
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	_, cmd := m.Update(msg)
	if cmd == nil {
		return ""
	}
	if cs, ok := cmd().(ui.ChangeScreenMsg); ok {
		return cs.Screen
	}
	return ""
}

func newDefeatedPlayer(t *testing.T) (*storage.DB, *engine.Player) {
	t.Helper()
	db, err := storage.OpenDB(filepath.Join(t.TempDir(), "defeat.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	sp, err := db.CreatePlayer("derrotado", "senha")
	if err != nil {
		t.Fatalf("CreatePlayer: %v", err)
	}
	p := engine.NewPlayerFromStorage(sp)
	p.Health = 0
	return db, p
}

// Em estado de derrota, [V], [Esc] e [C] têm de levar ao Game Over (onde a
// penalidade é aplicada), nunca de volta à Praça.
func TestForestScreen_DefeatHasNoSideExit(t *testing.T) {
	for _, key := range []string{"v", "V", "esc", "c"} {
		db, p := newDefeatedPlayer(t)
		s := NewForestScreen(db, p)
		s.state = forestStateDefeat
		if got := screenAfterKey(t, s, key); got != ui.ScreenGameOver {
			t.Errorf("tecla %q em derrota na floresta levou a %q, esperado %q", key, got, ui.ScreenGameOver)
		}
	}
}

func TestDragonScreen_DefeatHasNoSideExit(t *testing.T) {
	for _, key := range []string{"v", "V", "esc"} {
		db, p := newDefeatedPlayer(t)
		s := NewDragonScreen(db, p, nil)
		s.state = dragonStateDefeat
		if got := screenAfterKey(t, s, key); got != ui.ScreenGameOver {
			t.Errorf("tecla %q em derrota no covil levou a %q, esperado %q", key, got, ui.ScreenGameOver)
		}
	}
}

// Fora da derrota, [V] continua a voltar à Praça (regressão do comportamento normal).
func TestForestScreen_ExploringStillReturnsToTown(t *testing.T) {
	db, p := newDefeatedPlayer(t)
	p.Health = p.MaxHealth
	s := NewForestScreen(db, p)
	if got := screenAfterKey(t, s, "v"); got != ui.ScreenTown {
		t.Errorf("[V] a explorar levou a %q, esperado %q", got, ui.ScreenTown)
	}
}
