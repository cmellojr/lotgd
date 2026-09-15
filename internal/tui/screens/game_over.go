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

// GameOverScreen exibe a tela de derrota moratória, informando as perdas de ouro/XP e ressurreição na capela.
//
// Didática TEA: O `GameOverScreen` é acionado quando a vida do herói chega a zero.
// Ele exibe os valores exatos de ouro perdido da bolsa e experiência reduzida calculados por `engine.EconomyService`,
// lembrando o jogador de que os fundos guardados no cofre do banco permanecem 100% seguros.
type GameOverScreen struct {
	db       *storage.DB
	player   *engine.Player
	lostGold int
	lostXP   int
	width    int
	height   int
}

// NewGameOverScreen inicializa a tela de derrota com os valores das penalidades moratórias calculadas.
func NewGameOverScreen(db *storage.DB, player *engine.Player, lostGold, lostXP int) *GameOverScreen {
	return &GameOverScreen{
		db:       db,
		player:   player,
		lostGold: lostGold,
		lostXP:   lostXP,
	}
}

// Init inicializa a tela de game over.
func (s *GameOverScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer atualiza a referência ao herói ativo em memória.
func (s *GameOverScreen) SetPlayer(p *engine.Player) {
	s.player = p
	s.lostGold = 0
	s.lostXP = 0
}

// SetSize atualiza as dimensões de largura e altura da tela.
func (s *GameOverScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update aguarda a confirmação do jogador (`Enter`) para redirecioná-lo ressuscitado para a Capela do Frei Anselmo.
func (s *GameOverScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "v", "c":
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenChapel}
			}
		}
	}
	return s, nil
}

// View renderiza o painel de derrota e o resumo das perdas moratórias no terminal.
func (s *GameOverScreen) View() string {
	var b strings.Builder

	title := ui.TitleStyle.Render("☠  " + i18n.GetUIText(i18n.UIGameOverTitle) + "  ☠")
	b.WriteString(title + "\n\n")

	var content strings.Builder
	content.WriteString("Sua visão escurece enquanto seu corpo cai inerte no chão frio.\n")
	content.WriteString("Salteadores e criaturas da floresta vasculharam seus pertences...\n\n")

	content.WriteString("Penalidades da Derrota:\n")
	content.WriteString(" • " + ui.ErrorNoticeStyle.Render(i18n.GetMessage(i18n.MsgGameOverLostGold, s.lostGold)) + "\n")
	content.WriteString(" • " + ui.ErrorNoticeStyle.Render(i18n.GetMessage(i18n.MsgGameOverLostXP, s.lostXP)) + "\n")
	if s.player != nil {
		content.WriteString(fmt.Sprintf(" • Ouro protegido no cofre do banco: %s moedas\n\n", ui.StatusGold.Render(fmt.Sprintf("%d", s.player.BankGold))))
	}

	content.WriteString(i18n.GetUIText(i18n.UIGameOverRespawnInfo) + "\n\n")
	content.WriteString(ui.SelectedMenuItemStyle.Render("> Pressione [Enter] para despertar na Capela..."))

	b.WriteString(ui.CombatBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render(i18n.GetUIText(i18n.UIFooterGameOver)))

	return ui.AppStyle.Render(b.String())
}
