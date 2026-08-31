package screens

import (
	"fmt"
	"strings"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// GameOverScreen displays the death penalty and resurrection screen.
type GameOverScreen struct {
	db       *storage.DB
	player   *engine.Player
	lostGold int
	lostXP   int
	width    int
	height   int
}

// NewGameOverScreen initializes the death screen.
func NewGameOverScreen(db *storage.DB, player *engine.Player) *GameOverScreen {
	econ := engine.NewEconomyService()
	var lGold, lXP int
	if player != nil {
		lGold, lXP = econ.ProcessDeathPenalty(player)
		if db != nil {
			_ = db.SavePlayer(player.ToStorage())
		}
	}

	return &GameOverScreen{
		db:       db,
		player:   player,
		lostGold: lGold,
		lostXP:   lXP,
	}
}

// Init starts the game over screen.
func (s *GameOverScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer sets the player and applies death penalties.
func (s *GameOverScreen) SetPlayer(p *engine.Player) {
	s.player = p
	econ := engine.NewEconomyService()
	if p != nil {
		s.lostGold, s.lostXP = econ.ProcessDeathPenalty(p)
		if s.db != nil {
			_ = s.db.SavePlayer(p.ToStorage())
		}
	}
}

// SetSize updates screen size.
func (s *GameOverScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update waits for confirmation to return to the chapel/town.
func (s *GameOverScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "v", "c":
			if s.player != nil && s.db != nil {
				_ = s.db.SavePlayer(s.player.ToStorage())
			}
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenChapel}
			}
		}
	}
	return s, nil
}

// View renders the death screen.
func (s *GameOverScreen) View() string {
	var b strings.Builder

	title := ui.TitleStyle.Render("☠  O DESTINO DE UM HERÓI CAÍDO  ☠")
	b.WriteString(title + "\n\n")

	var content strings.Builder
	content.WriteString("Sua visão escurece enquanto seu corpo cai inerte no chão frio.\n")
	content.WriteString("Salteadores e criaturas da floresta vasculharam seus pertences...\n\n")

	content.WriteString("Penalidades da Derrota:\n")
	content.WriteString(fmt.Sprintf(" • Ouro perdido na bolsa: %s moedas\n", ui.ErrorNoticeStyle.Render(fmt.Sprintf("%d", s.lostGold))))
	content.WriteString(fmt.Sprintf(" • Experiência perdida: %s XP\n", ui.ErrorNoticeStyle.Render(fmt.Sprintf("%d", s.lostXP))))
	if s.player != nil {
		content.WriteString(fmt.Sprintf(" • Ouro protegido no cofre do banco: %s moedas\n\n", ui.StatusGold.Render(fmt.Sprintf("%d", s.player.BankGold))))
	}

	content.WriteString("O Frei Anselmo encontrou seu corpo à beira da estrada e o ressuscitou na Capela com 1 HP.\n\n")
	content.WriteString(ui.SelectedMenuItemStyle.Render("> Pressione [Enter] para despertar na Capela..."))

	b.WriteString(ui.CombatBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[Enter] Despertar na Capela"))

	return ui.AppStyle.Render(b.String())
}
