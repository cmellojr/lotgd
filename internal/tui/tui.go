package tui

import (
	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// MainModel is the root Bubble Tea model managing screen routing and global state.
type MainModel struct {
	db            *storage.DB
	player        *engine.Player
	currentScreen ScreenID
	width         int
	height        int

	// Sub-screens
	loginScreen    *screens.LoginScreen
	townScreen     *screens.TownScreen
	forestScreen   *screens.ForestScreen
	tavernScreen   *screens.TavernScreen
	chapelScreen   *screens.ChapelScreen
	smithScreen    *screens.SmithScreen
	guildScreen    *screens.GuildScreen
	dragonScreen   *screens.DragonScreen
	gameOverScreen *screens.GameOverScreen
}

// NewMainModel instantiates the root TUI model.
func NewMainModel(db *storage.DB) *MainModel {
	return &MainModel{
		db:             db,
		currentScreen:  ScreenLogin,
		loginScreen:    screens.NewLoginScreen(db),
		townScreen:     screens.NewTownScreen(db, nil),
		forestScreen:   screens.NewForestScreen(db, nil),
		tavernScreen:   screens.NewTavernScreen(db, nil),
		chapelScreen:   screens.NewChapelScreen(db, nil),
		smithScreen:    screens.NewSmithScreen(db, nil),
		guildScreen:    screens.NewGuildScreen(db, nil),
		dragonScreen:   screens.NewDragonScreen(db, nil),
		gameOverScreen: screens.NewGameOverScreen(db, nil),
	}
}

// Init starts the root model.
func (m *MainModel) Init() tea.Cmd {
	return m.loginScreen.Init()
}

// Update delegates messages to the active screen and handles global transitions.
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncSizes()
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			if m.player != nil && m.db != nil {
				_ = m.db.SavePlayer(m.player.ToStorage())
			}
			return m, tea.Quit
		}

	case PlayerUpdatedMsg:
		m.player = msg.Player
		m.syncPlayerState()
		m.currentScreen = ScreenTown
		return m, nil

	case ChangeScreenMsg:
		m.currentScreen = msg.Screen
		m.syncPlayerState()
		return m, nil
	}

	var cmd tea.Cmd
	switch m.currentScreen {
	case ScreenLogin:
		var sub tea.Model
		sub, cmd = m.loginScreen.Update(msg)
		m.loginScreen = sub.(*screens.LoginScreen)

	case ScreenTown:
		var sub tea.Model
		sub, cmd = m.townScreen.Update(msg)
		m.townScreen = sub.(*screens.TownScreen)

	case ScreenForest:
		var sub tea.Model
		sub, cmd = m.forestScreen.Update(msg)
		m.forestScreen = sub.(*screens.ForestScreen)

	case ScreenTavern:
		var sub tea.Model
		sub, cmd = m.tavernScreen.Update(msg)
		m.tavernScreen = sub.(*screens.TavernScreen)

	case ScreenChapel:
		var sub tea.Model
		sub, cmd = m.chapelScreen.Update(msg)
		m.chapelScreen = sub.(*screens.ChapelScreen)

	case ScreenSmith:
		var sub tea.Model
		sub, cmd = m.smithScreen.Update(msg)
		m.smithScreen = sub.(*screens.SmithScreen)

	case ScreenGuild:
		var sub tea.Model
		sub, cmd = m.guildScreen.Update(msg)
		m.guildScreen = sub.(*screens.GuildScreen)

	case ScreenDragon:
		var sub tea.Model
		sub, cmd = m.dragonScreen.Update(msg)
		m.dragonScreen = sub.(*screens.DragonScreen)

	case ScreenGameOver:
		var sub tea.Model
		sub, cmd = m.gameOverScreen.Update(msg)
		m.gameOverScreen = sub.(*screens.GameOverScreen)
	}

	return m, cmd
}

func (m *MainModel) syncSizes() {
	m.loginScreen.SetSize(m.width, m.height)
	m.townScreen.SetSize(m.width, m.height)
	m.forestScreen.SetSize(m.width, m.height)
	m.tavernScreen.SetSize(m.width, m.height)
	m.chapelScreen.SetSize(m.width, m.height)
	m.smithScreen.SetSize(m.width, m.height)
	m.guildScreen.SetSize(m.width, m.height)
	m.dragonScreen.SetSize(m.width, m.height)
	m.gameOverScreen.SetSize(m.width, m.height)
}

func (m *MainModel) syncPlayerState() {
	if m.player == nil {
		return
	}
	m.townScreen.SetPlayer(m.player)
	m.forestScreen.SetPlayer(m.player)
	m.tavernScreen.SetPlayer(m.player)
	m.chapelScreen.SetPlayer(m.player)
	m.smithScreen.SetPlayer(m.player)
	m.guildScreen.SetPlayer(m.player)
	m.dragonScreen.SetPlayer(m.player)
	m.gameOverScreen.SetPlayer(m.player)
}

// View delegates rendering to the active screen.
func (m *MainModel) View() string {
	switch m.currentScreen {
	case ScreenLogin:
		return m.loginScreen.View()
	case ScreenTown:
		return m.townScreen.View()
	case ScreenForest:
		return m.forestScreen.View()
	case ScreenTavern:
		return m.tavernScreen.View()
	case ScreenChapel:
		return m.chapelScreen.View()
	case ScreenSmith:
		return m.smithScreen.View()
	case ScreenGuild:
		return m.guildScreen.View()
	case ScreenDragon:
		return m.dragonScreen.View()
	case ScreenGameOver:
		return m.gameOverScreen.View()
	default:
		return "Tela desconhecida."
	}
}
