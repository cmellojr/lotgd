package ui

import (
	"fmt"
	"strings"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"

	"github.com/charmbracelet/lipgloss"
)

// Constantes de paleta de cores ANSI para o tema de terminal retro BBS.
//
// Didática Go: Usamos `lipgloss.Color` com códigos hexadecimais para criar um visual rico em
// terminais ANSI de 256 cores ou TrueColor, mantendo consistência temática em toda a interface TUI.
const (
	ColorGoldDark   = lipgloss.Color("#D4AF37")
	ColorGoldBright = lipgloss.Color("#FFD700")
	ColorRedDark    = lipgloss.Color("#8B0000")
	ColorRedBright  = lipgloss.Color("#FF4500")
	ColorGreenDark  = lipgloss.Color("#006400")
	ColorGreenLight = lipgloss.Color("#32CD32")
	ColorCyanDark   = lipgloss.Color("#008B8B")
	ColorCyanBright = lipgloss.Color("#00FFFF")
	ColorPurple     = lipgloss.Color("#9370DB")
	ColorGrayDark   = lipgloss.Color("#2E3440")
	ColorGrayMid    = lipgloss.Color("#4C566A")
	ColorGrayLight  = lipgloss.Color("#D8DEE9")
	ColorWhite      = lipgloss.Color("#ECEFF4")
)

// Estilos de UI reutilizáveis instanciados com a biblioteca Lip Gloss.
//
// Didática Go: Estilos no Lip Gloss são imutáveis e encadeáveis (builder pattern).
// Definir e reutilizar variáveis globais de estilo evita alocações redundantes a cada quadro renderizado.
var (
	// Caixas base da aplicação
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(ColorWhite)

	// Banners e Títulos
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGoldBright).
			Padding(0, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorGoldDark)

	SubtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(ColorCyanBright)

	// Barra de status no topo da tela
	StatusBarContainer = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(ColorGrayMid).
				Padding(0, 0, 1, 0).
				MarginBottom(1)

	StatusLabel = lipgloss.NewStyle().
			Foreground(ColorCyanDark).
			Bold(true)

	StatusValue = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true)

	StatusGold = lipgloss.NewStyle().
			Foreground(ColorGoldBright).
			Bold(true)

	StatusHP = lipgloss.NewStyle().
			Foreground(ColorGreenLight).
			Bold(true)

	StatusFights = lipgloss.NewStyle().
			Foreground(ColorRedBright).
			Bold(true)

	// Caixa de Conteúdo e Diálogos
	ContentBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCyanDark).
			Padding(1, 2).
			MarginBottom(1)

	// Estilos de Itens de Menu
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorGrayLight)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorGoldBright).
				Bold(true)

	KeyShortcutStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorGoldBright)

	// Estilos de Combate e Log
	CombatBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorRedDark).
			Padding(1, 2).
			MarginBottom(1)

	LogPlayerStyle = lipgloss.NewStyle().
			Foreground(ColorGreenLight)

	LogMonsterStyle = lipgloss.NewStyle().
			Foreground(ColorRedBright)

	LogSystemStyle = lipgloss.NewStyle().
			Foreground(ColorCyanBright).
			Italic(true)

	LogCriticalStyle = lipgloss.NewStyle().
				Foreground(ColorGoldBright).
				Bold(true)

	// Notificações e Mensagens de Erro
	ErrorNoticeStyle = lipgloss.NewStyle().
				Foreground(ColorRedBright).
				Bold(true)

	SuccessNoticeStyle = lipgloss.NewStyle().
				Foreground(ColorGreenLight).
				Bold(true)

	HelpFooterStyle = lipgloss.NewStyle().
			Foreground(ColorGrayMid).
			Italic(true).
			MarginTop(1)
)

// ScreenID representa o identificador único de cada tela na máquina de estados da TUI.
type ScreenID string

const (
	ScreenLogin    ScreenID = "login"
	ScreenTown     ScreenID = "town"
	ScreenForest   ScreenID = "forest"
	ScreenTavern   ScreenID = "tavern"
	ScreenChapel   ScreenID = "chapel"
	ScreenSmith    ScreenID = "smith"
	ScreenGuild    ScreenID = "guild"
	ScreenDragon   ScreenID = "dragon"
	ScreenGameOver ScreenID = "game_over"
)

// ChangeScreenMsg é uma mensagem do Bubble Tea solicitando a transição de tela ativa.
//
// Didática Go: No Bubble Tea (The Elm Architecture), as telas filhas disparam comandos (`tea.Cmd`)
// que retornam mensagens customizadas como `ChangeScreenMsg`. O modelo raiz (`MainModel`) intercepta
// essa mensagem no seu método `Update` e altera a tela visível.
type ChangeScreenMsg struct {
	Screen ScreenID
}

// PlayerUpdatedMsg notifica o modelo raiz e as sub-telas de que o estado do jogador em memória mudou.
type PlayerUpdatedMsg struct {
	Player *engine.Player
}

// RenderStatusBar renderiza o cabeçalho superior padrão com atributos do herói, barra de vida e equipamentos.
func RenderStatusBar(p *engine.Player, width int) string {
	if p == nil {
		return ""
	}

	hpPercent := float64(p.Health) / float64(p.MaxHealth)
	hpBar := renderHPBar(hpPercent, 10)

	heroInfo := fmt.Sprintf("%s %s | %s %s",
		StatusLabel.Render("Herói:"),
		StatusValue.Render(p.Username),
		StatusLabel.Render("Nível:"),
		StatusValue.Render(fmt.Sprintf("%d", p.Level)),
	)

	healthInfo := fmt.Sprintf("%s %s %s/%s",
		StatusLabel.Render("HP:"),
		hpBar,
		StatusHP.Render(fmt.Sprintf("%d", p.Health)),
		StatusValue.Render(fmt.Sprintf("%d", p.MaxHealth)),
	)

	goldInfo := fmt.Sprintf("%s %s %s",
		StatusLabel.Render("Ouro:"),
		StatusGold.Render(fmt.Sprintf("%d", p.Gold)),
		StatusValue.Render(fmt.Sprintf("(Banco: %d)", p.BankGold)),
	)

	fightsInfo := fmt.Sprintf("%s %s",
		StatusLabel.Render("Lutas Diárias:"),
		StatusFights.Render(fmt.Sprintf("%d", p.ForestFights)),
	)

	weaponName := i18n.GetItemName(p.Weapon.ID)
	armorName := i18n.GetItemName(p.Armor.ID)
	equipInfo := fmt.Sprintf("%s %s | %s %s | %s %s",
		StatusLabel.Render("Arma:"),
		StatusValue.Render(weaponName),
		StatusLabel.Render("Armadura:"),
		StatusValue.Render(armorName),
		StatusLabel.Render("Poções:"),
		StatusValue.Render(fmt.Sprintf("%d", p.PotionsCount)),
	)

	// Atributos de combate e progresso de experiência.
	//
	// Sem isto o jogador vê os atributos do monstro na tela de combate e nunca os
	// seus: TotalAttack e TotalDefense eram calculados e nunca exibidos, e a
	// experiência só aparecia na Guilda.
	xpText := fmt.Sprintf("%d (máx)", p.Experience)
	if req, ok := engine.NextLevelRequirement(p.Level); ok {
		xpText = fmt.Sprintf("%d/%d", p.Experience, req.RequiredXP)
	}

	combatInfo := fmt.Sprintf("%s %s | %s %s | %s %s",
		StatusLabel.Render("ATK:"),
		StatusValue.Render(fmt.Sprintf("%d", p.TotalAttack())),
		StatusLabel.Render("DEF:"),
		StatusValue.Render(fmt.Sprintf("%d", p.TotalDefense())),
		StatusLabel.Render("XP:"),
		StatusValue.Render(xpText),
	)

	line1 := fmt.Sprintf("%s    %s    %s    %s", heroInfo, healthInfo, goldInfo, fightsInfo)
	line2 := fmt.Sprintf("%s | %s", equipInfo, combatInfo)

	return StatusBarContainer.Width(width).Render(line1 + "\n" + line2)
}

// renderHPBar constrói a representação gráfica visual em blocos ANSI da barra de vida do jogador.
func renderHPBar(ratio float64, totalBlocks int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * float64(totalBlocks))
	empty := totalBlocks - filled

	fillStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)

	return fmt.Sprintf("[%s%s]", StatusHP.Render(fillStr), StatusLabel.Render(emptyStr))
}
