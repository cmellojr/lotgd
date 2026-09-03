package screens

import (
	"errors"
	"fmt"
	"strings"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginFocus int

const (
	focusUsername loginFocus = iota
	focusPassword
	focusSubmit
)

// LoginScreen handles authentication and player account creation.
type LoginScreen struct {
	db         *storage.DB
	usernameIn textinput.Model
	passwordIn textinput.Model
	focus      loginFocus
	errMsg     string
	infoMsg    string
	width      int
	height     int
}

// NewLoginScreen initializes the login interface.
func NewLoginScreen(db *storage.DB) *LoginScreen {
	u := textinput.New()
	u.Placeholder = "Digite seu nome de aventureiro..."
	u.Focus()
	u.CharLimit = 32
	u.Width = 35

	p := textinput.New()
	p.Placeholder = "Digite sua senha secreta..."
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.CharLimit = 32
	p.Width = 35

	return &LoginScreen{
		db:         db,
		usernameIn: u,
		passwordIn: p,
		focus:      focusUsername,
		infoMsg:    "Se a conta não existir, ela será criada automaticamente.",
	}
}

// Init returns the initial command for the login screen.
func (s *LoginScreen) Init() tea.Cmd {
	return textinput.Blink
}

// SetSize updates the screen dimensions.
func (s *LoginScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processes input events on the login form.
func (s *LoginScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			s.nextFocus()
			return s, nil
		case "shift+tab", "up":
			s.prevFocus()
			return s, nil
		case "enter":
			if s.focus == focusSubmit || s.focus == focusPassword {
				return s.handleLogin()
			}
			s.nextFocus()
			return s, nil
		}
	}

	var cmdU, cmdP tea.Cmd
	s.usernameIn, cmdU = s.usernameIn.Update(msg)
	s.passwordIn, cmdP = s.passwordIn.Update(msg)

	return s, tea.Batch(cmdU, cmdP)
}

func (s *LoginScreen) nextFocus() {
	if s.focus == focusUsername {
		s.focus = focusPassword
		s.usernameIn.Blur()
		s.passwordIn.Focus()
	} else if s.focus == focusPassword {
		s.focus = focusSubmit
		s.passwordIn.Blur()
	} else {
		s.focus = focusUsername
		s.passwordIn.Blur()
		s.usernameIn.Focus()
	}
}

func (s *LoginScreen) prevFocus() {
	if s.focus == focusSubmit {
		s.focus = focusPassword
		s.passwordIn.Focus()
	} else if s.focus == focusPassword {
		s.focus = focusUsername
		s.passwordIn.Blur()
		s.usernameIn.Focus()
	} else {
		s.focus = focusSubmit
		s.usernameIn.Blur()
	}
}

func (s *LoginScreen) handleLogin() (tea.Model, tea.Cmd) {
	user := strings.TrimSpace(s.usernameIn.Value())
	pass := strings.TrimSpace(s.passwordIn.Value())

	if user == "" || pass == "" {
		s.errMsg = "Por favor, preencha o nome de aventureiro e a senha."
		return s, nil
	}

	// Tenta autenticar
	sp, err := s.db.AuthenticatePlayer(user, pass)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidPass):
			s.errMsg = "Senha incorreta para o aventureiro."
			return s, nil
		case errors.Is(err, storage.ErrPlayerNotFound):
			// Se não existe, cria
			newSP, createErr := s.db.CreatePlayer(user, pass)
			if createErr != nil {
				if errors.Is(createErr, storage.ErrUserExists) {
					s.errMsg = "Este nome de aventureiro já está registrado."
				} else {
					s.errMsg = fmt.Sprintf("Erro ao criar conta: %v", createErr)
				}
				return s, nil
			}
			sp = newSP
		default:
			s.errMsg = fmt.Sprintf("Erro ao autenticar: %v", err)
			return s, nil
		}
	}

	// Converte para modelo de domínio
	player := engine.NewPlayerFromStorage(sp)

	// Verifica e processa virada do dia
	tm := engine.NewTurnManager()
	today := engine.CurrentDateString()
	if tm.CheckAndApplyNewDay(player, today) {
		SavePlayer(s.db, player)
	}

	return s, func() tea.Msg {
		return ui.PlayerUpdatedMsg{Player: player}
	}
}

// View renders the login screen.
func (s *LoginScreen) View() string {
	title := ui.TitleStyle.Render("🏰 THE LEGEND OF THE GO DRAGON 🐉")
	subtitle := ui.SubtitleStyle.Render("Uma aventura épica no terminal (BBS RPG Clássico)")

	var b strings.Builder
	b.WriteString(title + "\n")
	b.WriteString(subtitle + "\n\n")

	content := fmt.Sprintf(
		"Nome de Aventureiro:\n%s\n\nSenha Secreta:\n%s\n\n",
		s.usernameIn.View(),
		s.passwordIn.View(),
	)

	submitBtn := "[ Entrar no Reino / Criar Personagem ]"
	if s.focus == focusSubmit {
		submitBtn = ui.SelectedMenuItemStyle.Render("> " + submitBtn)
	} else {
		submitBtn = ui.MenuItemStyle.Render(submitBtn)
	}

	content += submitBtn + "\n\n"

	if s.errMsg != "" {
		content += ui.ErrorNoticeStyle.Render("⚠ "+s.errMsg) + "\n"
	}
	if s.infoMsg != "" {
		content += ui.LogSystemStyle.Render("ℹ "+s.infoMsg) + "\n"
	}

	b.WriteString(ui.ContentBoxStyle.Width(60).Render(content))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[Tab/Shift+Tab] Navegar • [Enter] Confirmar • [Ctrl+C] Sair"))

	return ui.AppStyle.Render(b.String())
}
