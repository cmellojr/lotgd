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

type townMenuItem struct {
	key         string
	label       string
	target      ui.ScreenID
	description string
}

// TownScreen represents the main village hub.
type TownScreen struct {
	db       *storage.DB
	player   *engine.Player
	cursor   int
	items    []townMenuItem
	infoMsg  string
	bankMode bool
	width    int
	height   int
}

// NewTownScreen initializes the town hub.
func NewTownScreen(db *storage.DB, player *engine.Player) *TownScreen {
	items := []townMenuItem{
		{key: "F", label: "Floresta Sombria", target: ui.ScreenForest, description: "Procure monstros, lute por ouro e experiência."},
		{key: "T", label: "Taverna da Dona Rosalinda", target: ui.ScreenTavern, description: "Ouça fofocas, converse com aventureiros e veja notícias."},
		{key: "C", label: "Capela do Frei Anselmo", target: ui.ScreenChapel, description: "Cure seus ferimentos com o curandeiro do vilarejo."},
		{key: "M", label: "Ferraria do Mestre Torin", target: ui.ScreenSmith, description: "Compre armas melhores, armaduras e poções de cura."},
		{key: "G", label: "Guilda dos Aventureiros", target: ui.ScreenGuild, description: "Treine com o Mestre Tobias para subir de nível."},
		{key: "D", label: "Covil do Dragão Ancestral", target: ui.ScreenDragon, description: "O confronto final! Requer nível máximo e coragem."},
		{key: "B", label: "Banco do Vilarejo", target: "", description: "Deposite seu ouro para não perder ao morrer na floresta."},
		{key: "S", label: "Salvar e Sair (Logout)", target: ui.ScreenLogin, description: "Encerra a sessão e guarda o progresso no templo."},
	}

	return &TownScreen{
		db:      db,
		player:  player,
		cursor:  0,
		items:   items,
		infoMsg: "Bem-vindo à Praça Central. Escolha para onde deseja ir.",
	}
}

// Init initializes the town screen.
func (s *TownScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the active player reference.
func (s *TownScreen) SetPlayer(p *engine.Player) {
	s.player = p
}

// SetSize updates screen dimensions.
func (s *TownScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processes navigation and actions in the town square.
func (s *TownScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToUpper(msg.String())

		// Modo Banco Interativo
		if s.bankMode {
			switch k {
			case "D": // Depositar tudo
				if s.player.Gold <= 0 {
					s.infoMsg = "Você não possui moedas de ouro na bolsa para depositar."
				} else {
					econ := engine.NewEconomyService()
					deposited := s.player.Gold
					_ = econ.Deposit(s.player, deposited)
					_ = s.db.SavePlayer(s.player.ToStorage())
					s.infoMsg = fmt.Sprintf("Você depositou %d moedas de ouro no cofre com segurança!", deposited)
				}
				s.bankMode = false
				return s, nil
			case "R", "W": // Retirar tudo
				if s.player.BankGold <= 0 {
					s.infoMsg = "Seu cofre no banco está vazio."
				} else {
					econ := engine.NewEconomyService()
					withdrawn := s.player.BankGold
					_ = econ.Withdraw(s.player, withdrawn)
					_ = s.db.SavePlayer(s.player.ToStorage())
					s.infoMsg = fmt.Sprintf("Você retirou %d moedas de ouro do seu cofre.", withdrawn)
				}
				s.bankMode = false
				return s, nil
			case "ESC", "V", "B":
				s.bankMode = false
				s.infoMsg = "Você retornou para o centro da praça."
				return s, nil
			}
			return s, nil
		}

		// Navegação de Menu Padrão
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			} else {
				s.cursor = len(s.items) - 1
			}
			return s, nil
		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			} else {
				s.cursor = 0
			}
			return s, nil
		case "enter":
			return s.selectItem(s.items[s.cursor])
		}

		// Atalhos BBS diretos por tecla rápida
		for _, item := range s.items {
			if k == item.key {
				return s.selectItem(item)
			}
		}
	}

	return s, nil
}

func (s *TownScreen) selectItem(item townMenuItem) (tea.Model, tea.Cmd) {
	if item.key == "B" {
		s.bankMode = true
		s.infoMsg = "Bem-vindo ao Banco do Vilarejo. [D]epositar todo o ouro | [R]etirar todo o ouro | [Esc] Voltar"
		return s, nil
	}

	if item.target == ui.ScreenLogin {
		_ = s.db.SavePlayer(s.player.ToStorage())
		return s, func() tea.Msg {
			return ui.ChangeScreenMsg{Screen: ui.ScreenLogin}
		}
	}

	if item.target != "" {
		return s, func() tea.Msg {
			return ui.ChangeScreenMsg{Screen: item.target}
		}
	}

	return s, nil
}

// View renders the village square.
func (s *TownScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("🏛  " + i18n.GetLocationName(i18n.LocationTown) + "  🏛")
	b.WriteString(title + "\n\n")

	var menuContent strings.Builder
	for i, item := range s.items {
		shortcut := ui.KeyShortcutStyle.Render("[" + item.key + "]")
		line := fmt.Sprintf("%s %s — %s", shortcut, item.label, ui.HelpFooterStyle.Render(item.description))

		if i == s.cursor {
			menuContent.WriteString(ui.SelectedMenuItemStyle.Render("> " + line) + "\n")
		} else {
			menuContent.WriteString(ui.MenuItemStyle.Render("  " + line) + "\n")
		}
	}

	if s.infoMsg != "" {
		menuContent.WriteString("\n" + ui.LogSystemStyle.Render("ℹ "+s.infoMsg))
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(menuContent.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[↑/↓] Selecionar • [Enter] Entrar • [Letras] Atalho Rápido • [Ctrl+C] Sair"))

	return ui.AppStyle.Render(b.String())
}
