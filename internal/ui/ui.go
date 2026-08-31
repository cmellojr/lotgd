package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

// Color palette constants for the retro BBS ANSI theme.
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

// UI styles using Lip Gloss.
var (
	// Base application box style
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(ColorWhite)

	// Banner & Titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGoldBright).
			Padding(0, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorGoldDark)

	SubtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(ColorCyanBright)

	// Status bar at top of screen
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

	// Dialog & Content Box
	ContentBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCyanDark).
			Padding(1, 2).
			MarginBottom(1)

	// Menu Item Styles
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

	// Combat & Log Styles
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

	// Error & Notice Messages
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

// ScreenID represents unique screen identifiers in the TUI state machine.
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

// ChangeScreenMsg requests the root model to navigate to a new screen.
type ChangeScreenMsg struct {
	Screen ScreenID
}

// PlayerUpdatedMsg notifies the root model and sub-models that the active player state changed.
type PlayerUpdatedMsg struct {
	Player *engine.Player
}

// RenderStatusBar renders the standard top status bar for the active player.
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

	line1 := fmt.Sprintf("%s    %s    %s    %s", heroInfo, healthInfo, goldInfo, fightsInfo)
	line2 := equipInfo

	return StatusBarContainer.Width(width).Render(line1 + "\n" + line2)
}

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
