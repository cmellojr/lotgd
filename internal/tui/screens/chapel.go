package screens

import (
	"fmt"
	"strings"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// ChapelScreen represents the sanctuary of healing and blessings.
type ChapelScreen struct {
	db        *storage.DB
	player    *engine.Player
	cursor    int
	menuItems []string
	infoMsg   string
	width     int
	height    int
}

// NewChapelScreen initializes the chapel screen.
func NewChapelScreen(db *storage.DB, player *engine.Player) *ChapelScreen {
	return &ChapelScreen{
		db:     db,
		player: player,
		cursor: 0,
		menuItems: []string{
			"Pedir Cura aos Céus (Custo proporcional ao dano)",
			"Fazer Doação aos Pobres (10 Ouro)",
			"Meditar em Silêncio diante do Altar",
			"Voltar para a Praça Central",
		},
		infoMsg: "Velas aromáticas iluminam os vitrais sagrados da capela.",
	}
}

// Init starts the chapel screen.
func (s *ChapelScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates player state.
func (s *ChapelScreen) SetPlayer(p *engine.Player) {
	s.player = p
}

// SetSize updates dimensions.
func (s *ChapelScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update handles chapel interactions.
func (s *ChapelScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			} else {
				s.cursor = len(s.menuItems) - 1
			}
			return s, nil
		case "down", "j":
			if s.cursor < len(s.menuItems)-1 {
				s.cursor++
			} else {
				s.cursor = 0
			}
			return s, nil
		}

		k := strings.ToUpper(msg.String())

		switch k {
		case "V", "ESC":
			return s.backToTown()
		case "C": // Curar
			s.cursor = 0
			return s.selectCurrent()
		case "D": // Doação
			s.cursor = 1
			return s.selectCurrent()
		case "M": // Meditar
			s.cursor = 2
			return s.selectCurrent()
		case "ENTER":
			return s.selectCurrent()
		}
	}

	return s, nil
}

func (s *ChapelScreen) selectCurrent() (tea.Model, tea.Cmd) {
	switch s.cursor {
	case 0: // Cura completa
		if s.player.Health >= s.player.MaxHealth {
			s.infoMsg = "Frei Anselmo sorri: 'Você já goza de plena saúde, meu filho. Guarde seu ouro.'"
			return s, nil
		}

		missingHP := s.player.MaxHealth - s.player.Health
		cost := (missingHP * 2) / 3
		if cost < 1 {
			cost = 1
		}

		if s.player.Gold < cost {
			s.infoMsg = fmt.Sprintf("Você não tem ouro suficiente para o bálsamo sagrado. Custo: %d moedas (Você tem: %d)", cost, s.player.Gold)
			return s, nil
		}

		s.player.Gold -= cost
		s.player.Health = s.player.MaxHealth
		SavePlayer(s.db, s.player)
		s.infoMsg = fmt.Sprintf("Frei Anselmo unge seus ferimentos com óleos sagrados. Vida restaurada completamente por %d moedas de ouro!", cost)

	case 1: // Doação
		if s.player.Gold < 10 {
			s.infoMsg = "Você não possui nem 10 moedas de ouro para ofertar no confessionário."
			return s, nil
		}

		s.player.Gold -= 10
		SavePlayer(s.db, s.player)
		s.infoMsg = "Você coloca 10 moedas na caixa de esmolas. Uma sensação de paz e leveza espiritual aquece sua alma."

	case 2: // Meditação
		s.infoMsg = "Você ajoelha diante do altar e contempla as provações que o aguardam contra o Dragão."

	case 3: // Voltar
		return s.backToTown()
	}

	return s, nil
}

func (s *ChapelScreen) backToTown() (tea.Model, tea.Cmd) {
	SavePlayer(s.db, s.player)
	return s, func() tea.Msg {
		return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
	}
}

// View renders the chapel UI.
func (s *ChapelScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("⛪  " + i18n.GetLocationName(i18n.LocationChapel) + "  ⛪")
	b.WriteString(title + "\n\n")

	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s recebe você com um gesto de bênção sereno.\n\n", i18n.GetNPCName(i18n.NPCAnselmo)))

	for i, item := range s.menuItems {
		if i == s.cursor {
			content.WriteString(ui.SelectedMenuItemStyle.Render("> "+item) + "\n")
		} else {
			content.WriteString(ui.MenuItemStyle.Render("  "+item) + "\n")
		}
	}

	if s.infoMsg != "" {
		content.WriteString("\n" + ui.SuccessNoticeStyle.Render(s.infoMsg))
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[↑/↓] Selecionar • [Enter] Confirmar • [V] Voltar à Praça"))

	return ui.AppStyle.Render(b.String())
}
