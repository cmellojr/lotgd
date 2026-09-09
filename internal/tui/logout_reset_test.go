package tui

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lotgd/internal/storage"
)

func typeRunes(m tea.Model, text string) tea.Model {
	for _, r := range text {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return m
}

// Sair da sessão ("Salvar e Sair" na Praça) tem de devolver a tela de login
// vazia: sem nome, sem senha e sem jogador em memória.
func TestLogout_ResetsLoginScreenAndPlayer(t *testing.T) {
	db, err := storage.OpenDB(filepath.Join(t.TempDir(), "logout.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer db.Close()

	var m tea.Model = NewMainModel(db, nil)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	// O rodapé do login usa "•" como separador; a senha ecoa "•" por caractere.
	// A referência é a contagem na tela vazia.
	bulletsEmpty := strings.Count(m.View(), "•")

	// Login: nome, Tab, senha, Enter (cria a conta), depois a PlayerUpdatedMsg.
	m = typeRunes(m, "Mozozinho")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = typeRunes(m, "segredo")
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter no login não produziu comando")
	}
	m, _ = m.Update(cmd())
	mm := m.(*MainModel)
	if mm.currentScreen != ScreenTown || mm.player == nil {
		t.Fatalf("após login: tela=%q jogador=%v", mm.currentScreen, mm.player != nil)
	}

	// Sair: atalho [S] na Praça emite ChangeScreenMsg{ScreenLogin}.
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if cmd == nil {
		t.Fatal("[S] na Praça não produziu comando")
	}
	m, _ = m.Update(cmd())
	mm = m.(*MainModel)

	if mm.currentScreen != ScreenLogin {
		t.Fatalf("após sair: tela=%q, esperado %q", mm.currentScreen, ScreenLogin)
	}
	if mm.player != nil {
		t.Errorf("após sair, o jogador continua em memória: %s", mm.player.Username)
	}
	view := mm.View()
	if strings.Contains(view, "Mozozinho") {
		t.Errorf("após sair, o nome continua preenchido na tela de login")
	}
	if got := strings.Count(view, "•"); got != bulletsEmpty {
		t.Errorf("após sair, a senha continua preenchida na tela de login (%d caracteres ecoados)", got-bulletsEmpty)
	}
}
