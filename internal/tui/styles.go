package tui

import "github.com/charmbracelet/lipgloss"

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
