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

// GuildScreen handles player level-up promotions and lore with Master Tobias.
type GuildScreen struct {
	db        *storage.DB
	player    *engine.Player
	cursor    int
	menuItems []string
	infoMsg   string
	width     int
	height    int
}

// NewGuildScreen initializes the adventurers guild.
func NewGuildScreen(db *storage.DB, player *engine.Player) *GuildScreen {
	return &GuildScreen{
		db:     db,
		player: player,
		cursor: 0,
		menuItems: []string{
			"Avançar de Nível (Treinamento com Mestre Tobias)",
			"Consultar Requisitos de Níveis e Maestria",
			"Ler Pergaminhos Antigos da Lenda do Dragão",
			"Voltar para a Praça Central",
		},
		infoMsg: "Prateleiras infinitas de tomos antigos cobrem as paredes de pedra da guilda.",
	}
}

// Init starts the guild screen.
func (s *GuildScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the player reference.
func (s *GuildScreen) SetPlayer(p *engine.Player) {
	s.player = p
}

// SetSize updates screen dimensions.
func (s *GuildScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processes guild interactions.
func (s *GuildScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToUpper(msg.String())

		switch k {
		case "V", "ESC":
			return s.backToTown()
		case "A", "T":
			s.cursor = 0
			return s.selectCurrent()
		case "R", "C":
			s.cursor = 1
			return s.selectCurrent()
		case "P", "L":
			s.cursor = 2
			return s.selectCurrent()
		case "UP":
			if s.cursor > 0 {
				s.cursor--
			} else {
				s.cursor = len(s.menuItems) - 1
			}
			return s, nil
		case "DOWN":
			if s.cursor < len(s.menuItems)-1 {
				s.cursor++
			} else {
				s.cursor = 0
			}
			return s, nil
		case "ENTER":
			return s.selectCurrent()
		}
	}

	return s, nil
}

func (s *GuildScreen) selectCurrent() (tea.Model, tea.Cmd) {
	switch s.cursor {
	case 0: // Treinar e subir de nível
		canUp, reason := engine.CanLevelUp(s.player)
		if !canUp {
			s.infoMsg = reason
			return s, nil
		}

		err := engine.LevelUp(s.player)
		if err != nil {
			s.infoMsg = fmt.Sprintf("Erro ao realizar treinamento: %v", err)
			return s, nil
		}

		SavePlayer(s.db, s.player)
		s.infoMsg = fmt.Sprintf("PARABÉNS! Você foi promovido com sucesso para o NÍVEL %d! Seus atributos e vida foram aprimorados!", s.player.Level)

	case 1: // Consultar tabela de níveis
		req, ok := engine.NextLevelRequirement(s.player.Level)
		if !ok {
			s.infoMsg = "Você já atingiu o Nível 10 (Grau Máximo de Mestre da Guilda)!"
		} else {
			s.infoMsg = fmt.Sprintf("Próximo Nível (%d): Requer %d XP e %d Ouro. Ganhos: +%d HP, +%d ATK, +%d DEF.",
				req.Level, req.RequiredXP, req.CostGold, req.HealthGain, req.AttackGain, req.DefGain)
		}

	case 2: // Pergaminhos do Dragão
		s.infoMsg = "Mestre Tobias ajusta seus óculos de leitura: 'Reza a lenda que o Dragão Ancestral desperta a cada alvorecer com forças renovadas. Apenas um herói armado até os dentes e com nível 5 ou superior ousará desafiar seu covil!'"

	case 3: // Voltar
		return s.backToTown()
	}

	return s, nil
}

func (s *GuildScreen) backToTown() (tea.Model, tea.Cmd) {
	SavePlayer(s.db, s.player)
	return s, func() tea.Msg {
		return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
	}
}

// View renders the guild interface.
func (s *GuildScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("📜  " + i18n.GetLocationName(i18n.LocationGuild) + "  📜")
	b.WriteString(title + "\n\n")

	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s ergue a cabeça por cima de uma pilha de tomos encadernados.\n", i18n.GetNPCName(i18n.NPCTobias)))
	content.WriteString(fmt.Sprintf("Sua Experiência Atual: %s XP  |  Ouro na Bolsa: %s moedas\n\n",
		ui.StatusValue.Render(fmt.Sprintf("%d", s.player.Experience)),
		ui.StatusGold.Render(fmt.Sprintf("%d", s.player.Gold)),
	))

	for i, item := range s.menuItems {
		if i == s.cursor {
			content.WriteString(ui.SelectedMenuItemStyle.Render("> "+item) + "\n")
		} else {
			content.WriteString(ui.MenuItemStyle.Render("  "+item) + "\n")
		}
	}

	if s.infoMsg != "" {
		content.WriteString("\n" + ui.LogSystemStyle.Render("ℹ "+s.infoMsg))
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[↑/↓] Selecionar • [Enter] Confirmar • [V] Voltar à Praça"))

	return ui.AppStyle.Render(b.String())
}
