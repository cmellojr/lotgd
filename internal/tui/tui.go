package tui

import (
	"log"

	"lotgd/internal/engine"
	"lotgd/internal/storage"
	"lotgd/internal/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// MainModel é o modelo raiz da arquitetura TEA (The Elm Architecture) do Bubble Tea no LOTGD.
//
// Didática Go / Arquitetura TUI:
// O `MainModel` funciona como o orquestrador central e roteador de estado do jogo.
// Ele mantém a referência ao banco de dados, o estado do herói ativo (`*engine.Player`),
// a tela visível atual (`currentScreen`) e as instâncias de todos os modelos de sub-telas.
type MainModel struct {
	db            *storage.DB
	player        *engine.Player
	currentScreen ScreenID
	width         int
	height        int

	// Sub-telas da aplicação
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

// NewMainModel instancia o modelo TUI raiz e inicializa todas as sub-telas.
func NewMainModel(db *storage.DB, dragonGen storage.DragonGenerator) *MainModel {
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
		dragonScreen:   screens.NewDragonScreen(db, nil, dragonGen),
		gameOverScreen: screens.NewGameOverScreen(db, nil, 0, 0),
	}
}

// Init inicializa a aplicação executando o comando de inicialização da tela de login.
//
// Didática TEA: O método `Init()` é o primeiro estágio do ciclo de vida (Init -> Update -> View).
// Ele retorna comandos assíncronos (`tea.Cmd`) que devem ser executados ao iniciar o programa.
func (m *MainModel) Init() tea.Cmd {
	return m.loginScreen.Init()
}

// Update intercepta e roteia mensagens (`tea.Msg`) para a tela ativa ou processa eventos globais.
//
// Fluxo de Execução do Roteador:
// 1. Redimensionamento de Terminal (`tea.WindowSizeMsg`): Atualiza dimensões globais e sincroniza com sub-telas.
// 2. Atalhos Globais (`Ctrl+C`): Retorna `tea.Quit` para encerrar o programa graciosamente.
// 3. Atualização de Jogador (`PlayerUpdatedMsg`): Registra o herói logado e navega para a Praça Central.
// 4. Mudança de Tela (`ChangeScreenMsg`): Atualiza `currentScreen`, processa penalidades moratórias se for Game Over e sincroniza estados.
// 5. Delegação: Encaminha a mensagem não capturada para o método `Update()` da sub-tela visível no momento.
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncSizes()
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case PlayerUpdatedMsg:
		m.player = msg.Player
		m.syncPlayerState()
		m.currentScreen = ScreenTown
		return m, nil

	case ChangeScreenMsg:
		if msg.Screen == ScreenGameOver && m.player != nil {
			econ := engine.NewEconomyService()
			lostGold, lostXP := econ.ProcessDeathPenalty(m.player)
			if err := m.Save(); err != nil {
				log.Printf("WARN: failed to save player state during Game Over: %v", err)
			}
			m.gameOverScreen = screens.NewGameOverScreen(m.db, m.player, lostGold, lostXP)
			m.gameOverScreen.SetSize(m.width, m.height)
		}
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
}

// View renderiza a interface textual chamando o método `View()` da tela atualmente ativa.
//
// Didática TEA: O método `View()` é a função pura de renderização no ciclo TEA (Init -> Update -> View).
// Ele converte o estado atual do modelo em uma string formatada em ANSI para exibição no terminal.
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

// Save persiste o estado atual do herói no banco de dados SQLite de forma segura, idônea e idempotente.
func (m *MainModel) Save() error {
	if m == nil || m.db == nil || m.player == nil {
		return nil
	}
	return m.db.SavePlayer(m.player.ToStorage())
}
