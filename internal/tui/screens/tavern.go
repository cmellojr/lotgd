package screens

import (
	"context"
	"fmt"
	"strings"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// TavernScreen represents the bustling social hub of the village.
type TavernScreen struct {
	db        *storage.DB
	player    *engine.Player
	cursor    int
	menuItems []string
	infoMsg   string
	news      []*storage.NewsEntry
	width     int
	height    int
}

// NewTavernScreen initializes the tavern screen.
func NewTavernScreen(db *storage.DB, player *engine.Player) *TavernScreen {
	return &TavernScreen{
		db:     db,
		player: player,
		cursor: 0,
		menuItems: []string{
			"Ouvir Fofocas com Dona Rosalinda",
			"Conversar e Flerte com Cassandra",
			"Desafiar o Cavaleiro Vermelho (Duelo)",
			"Ler o Mural de Notícias do Vilarejo",
			"Voltar para a Praça Central",
		},
		infoMsg: "O aroma de hidromel e ensopado preenche o salão aquecido pela lareira.",
	}
}

// Init starts the tavern screen.
func (s *TavernScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the player reference.
func (s *TavernScreen) SetPlayer(p *engine.Player) {
	s.player = p
	s.loadNews()
}

// SetSize updates screen dimensions.
func (s *TavernScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *TavernScreen) loadNews() {
	vRepo := storage.NewVillageRepository(s.db)
	entries, err := vRepo.GetLatestNews(context.Background(), 5)
	if err == nil {
		s.news = entries
	}
}

// Update processes player interactions in the tavern.
func (s *TavernScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "R", "F": // Fofoca / Rosalinda
			s.cursor = 0
			return s.selectCurrent()
		case "C": // Cassandra
			s.cursor = 1
			return s.selectCurrent()
		case "D": // Duelo
			s.cursor = 2
			return s.selectCurrent()
		case "N", "M": // Notícias
			s.cursor = 3
			return s.selectCurrent()
		case "ENTER":
			return s.selectCurrent()
		}
	}

	return s, nil
}

func (s *TavernScreen) selectCurrent() (tea.Model, tea.Cmd) {
	switch s.cursor {
	case 0: // Dona Rosalinda
		s.infoMsg = "Dona Rosalinda enxuga uma caneca: 'Ouvi dizer que o Mestre Torin forjou armas novas na ferraria, e que o Dragão anda mais agitado nos picos!'"
	case 1: // Cassandra
		if s.player.Gold >= 10 {
			s.player.Gold -= 10
			s.player.ForestFights++ // Ganha 1 turno de inspiração
			SavePlayer(s.db, s.player)
			s.infoMsg = "Você oferece uma bebida a Cassandra. Ela sorri graciosamente e sua determinação é renovada! (+1 Luta na Floresta!)"
		} else {
			s.infoMsg = "Cassandra te olha com desdém: 'Volte quando tiver pelo menos 10 moedas de ouro para pagar uma rodada, aventureiro.'"
		}
	case 2: // Cavaleiro Vermelho
		s.infoMsg = "O Cavaleiro Vermelho ergue a viseira: 'Você ainda não possui a têmpera necessária para cruzar lâminas comigo, garoto. Vá caçar na floresta!'"
	case 3: // Mural de Notícias
		s.loadNews()
		if len(s.news) == 0 {
			s.infoMsg = "O mural está vazio hoje. Nenhuma grande façanha recente foi registrada."
		} else {
			var sb strings.Builder
			sb.WriteString("Mural de Notícias Recentes:\n")
			for _, entry := range s.news {
				sb.WriteString(fmt.Sprintf(" 📜 [%s] %s\n", entry.CreatedAt.Format("15:04"), entry.Message))
			}
			s.infoMsg = sb.String()
		}
	case 4: // Voltar
		return s.backToTown()
	}

	return s, nil
}

func (s *TavernScreen) backToTown() (tea.Model, tea.Cmd) {
	SavePlayer(s.db, s.player)
	return s, func() tea.Msg {
		return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
	}
}

// View renders the tavern interface.
func (s *TavernScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("🍺  " + i18n.GetLocationName(i18n.LocationTavern) + "  🍺")
	b.WriteString(title + "\n\n")

	var content strings.Builder
	content.WriteString("Música de alaúde ecoa suavemente entre as mesas de carvalho maciço.\n\n")

	for i, item := range s.menuItems {
		if i == s.cursor {
			content.WriteString(ui.SelectedMenuItemStyle.Render("> "+item) + "\n")
		} else {
			content.WriteString(ui.MenuItemStyle.Render("  "+item) + "\n")
		}
	}

	if s.infoMsg != "" {
		content.WriteString("\n" + ui.LogSystemStyle.Render(s.infoMsg))
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[↑/↓] Navegar • [Enter] Conversar/Ação • [V] Voltar à Praça"))

	return ui.AppStyle.Render(b.String())
}
