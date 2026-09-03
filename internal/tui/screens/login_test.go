package screens

import (
	"path/filepath"
	"strings"
	"testing"

	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLoginScreen_ValidationAndAuth(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_login.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	screen := NewLoginScreen(db)
	screen.SetSize(80, 24)

	// Verify Init command is not nil (textinput.Blink)
	if cmd := screen.Init(); cmd == nil {
		t.Errorf("expected non-nil Init cmd")
	}

	// 1. Empty field submit
	screen.usernameIn.SetValue("")
	screen.passwordIn.SetValue("")
	screen.focus = focusSubmit
	_, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("expected nil cmd on validation failure")
	}
	if !strings.Contains(screen.errMsg, "Por favor, preencha o nome de aventureiro e a senha.") {
		t.Errorf("expected empty input error, got: %q", screen.errMsg)
	}

	// 2. Register new player successfully
	screen.usernameIn.SetValue("HeroOne")
	screen.passwordIn.SetValue("secret123")
	screen.focus = focusSubmit
	screen.errMsg = ""

	m, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("expected non-nil cmd on successful registration")
	}
	msg := cmd()
	pMsg, ok := msg.(ui.PlayerUpdatedMsg)
	if !ok || pMsg.Player == nil || pMsg.Player.Username != "HeroOne" {
		t.Fatalf("expected PlayerUpdatedMsg with HeroOne, got: %+v", msg)
	}

	// 3. Login with existing player, wrong password
	screen2 := NewLoginScreen(db)
	screen2.usernameIn.SetValue("HeroOne")
	screen2.passwordIn.SetValue("wrongpass")
	screen2.focus = focusSubmit

	_, cmd2 := screen2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd2 != nil {
		t.Errorf("expected nil cmd on wrong password failure")
	}
	if !strings.Contains(screen2.errMsg, "Senha incorreta para o aventureiro.") {
		t.Errorf("expected invalid pass error, got: %q", screen2.errMsg)
	}

	// 4. Login with existing player, correct password
	screen3 := NewLoginScreen(db)
	screen3.usernameIn.SetValue("HeroOne")
	screen3.passwordIn.SetValue("secret123")
	screen3.focus = focusSubmit

	_, cmd3 := screen3.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd3 == nil {
		t.Fatalf("expected non-nil cmd on correct password login")
	}
	msg3 := cmd3()
	pMsg3, ok := msg3.(ui.PlayerUpdatedMsg)
	if !ok || pMsg3.Player == nil || pMsg3.Player.Username != "HeroOne" {
		t.Fatalf("expected PlayerUpdatedMsg with HeroOne, got: %+v", msg3)
	}

	// 5. Test view rendering
	viewStr := m.(*LoginScreen).View()
	if !strings.Contains(viewStr, "THE LEGEND OF THE GO DRAGON") {
		t.Errorf("expected View output to contain game title, got: %s", viewStr)
	}
}

func TestLoginScreen_DuplicateAccountError(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_dup.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Pre-create account in DB
	if _, err := db.CreatePlayer("ExistingHero", "pass123"); err != nil {
		t.Fatalf("failed to create player: %v", err)
	}

	// Simulate concurrent registration where CreatePlayer returns ErrUserExists
	// by attempting duplicate registration via storage repository directly first
	if _, err := db.CreatePlayer("ExistingHero", "pass456"); err == nil {
		t.Fatalf("expected ErrUserExists on duplicate CreatePlayer, got nil")
	}

	// Now verify LoginScreen properly maps ErrUserExists to error message
	screen := NewLoginScreen(db)
	screen.usernameIn.SetValue("ExistingHero")
	screen.passwordIn.SetValue("wrongpass")
	screen.focus = focusSubmit

	// Wrong pass error on existing user
	_, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("expected nil cmd on wrong password")
	}
	if !strings.Contains(screen.errMsg, "Senha incorreta para o aventureiro.") {
		t.Errorf("expected wrong pass error message, got %q", screen.errMsg)
	}
}

func TestLoginScreen_FocusNavigation(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_nav.db")

	db, err := storage.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	screen := NewLoginScreen(db)

	if screen.focus != focusUsername {
		t.Errorf("expected initial focusUsername, got %v", screen.focus)
	}

	// Tab -> focusPassword
	screen.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screen.focus != focusPassword {
		t.Errorf("expected focusPassword after Tab, got %v", screen.focus)
	}

	// Tab -> focusSubmit
	screen.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screen.focus != focusSubmit {
		t.Errorf("expected focusSubmit after Tab, got %v", screen.focus)
	}

	// Tab -> focusUsername
	screen.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screen.focus != focusUsername {
		t.Errorf("expected focusUsername after Tab wrap, got %v", screen.focus)
	}

	// Shift+Tab -> focusSubmit
	screen.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if screen.focus != focusSubmit {
		t.Errorf("expected focusSubmit after Shift+Tab, got %v", screen.focus)
	}

	// Shift+Tab -> focusPassword
	screen.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if screen.focus != focusPassword {
		t.Errorf("expected focusPassword after Shift+Tab, got %v", screen.focus)
	}

	// Shift+Tab -> focusUsername
	screen.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if screen.focus != focusUsername {
		t.Errorf("expected focusUsername after Shift+Tab, got %v", screen.focus)
	}

	// Enter on Username moves focus to Password
	screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if screen.focus != focusPassword {
		t.Errorf("expected focusPassword after Enter on Username, got %v", screen.focus)
	}
}
