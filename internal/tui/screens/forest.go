package screens

import (
	"fmt"
	"strings"

	"lotgd/internal/bestiary"
	"lotgd/internal/engine"
	"lotgd/internal/i18n"
	"lotgd/internal/storage"
	"lotgd/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type forestState int

const (
	forestStateExploring forestState = iota
	forestStateCombat
	forestStateVictory
	forestStateDefeat
	forestStateFled
)

// ForestScreen handles forest exploration and turn-based combat.
type ForestScreen struct {
	db        *storage.DB
	player    *engine.Player
	gen       *bestiary.MonsterGenerator
	ce        *engine.CombatEngine
	tm        *engine.TurnManager
	state     forestState
	monster   *engine.Monster
	combatLog []string
	width     int
	height    int
}

// NewForestScreen initializes the forest exploration and combat screen.
func NewForestScreen(db *storage.DB, player *engine.Player) *ForestScreen {
	return &ForestScreen{
		db:        db,
		player:    player,
		gen:       bestiary.NewMonsterGenerator(nil),
		ce:        engine.NewCombatEngine(nil),
		tm:        engine.NewTurnManager(),
		state:     forestStateExploring,
		combatLog: make([]string, 0),
	}
}

// Init starts the forest screen.
func (s *ForestScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the active player.
func (s *ForestScreen) SetPlayer(p *engine.Player) {
	s.player = p
	s.state = forestStateExploring
	s.monster = nil
	s.combatLog = nil
}

// SetSize updates screen dimensions.
func (s *ForestScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processes forest actions and combat rounds.
func (s *ForestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToUpper(msg.String())

		// Retornar à cidade
		if k == "V" || k == "ESC" || (s.state == forestStateExploring && k == "C") {
			if s.state == forestStateCombat {
				s.appendLog("Você não pode fugir sem tentar uma retirada estratégica! [F]ugir")
				return s, nil
			}
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
			}
		}

		switch s.state {
		case forestStateExploring:
			switch k {
			case "P", "E", "ENTER":
				return s.startExploration()
			}

		case forestStateCombat:
			switch k {
			case "A", "ENTER":
				return s.handleAttack()
			case "F":
				return s.handleFlee()
			case "P", "U":
				return s.handlePotion()
			}

		case forestStateVictory, forestStateFled:
			switch k {
			case "P", "ENTER", "E":
				return s.startExploration()
			}

		case forestStateDefeat:
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenGameOver}
			}
		}
	}

	return s, nil
}

func (s *ForestScreen) startExploration() (tea.Model, tea.Cmd) {
	if s.player.ForestFights <= 0 {
		s.appendLog("Você está exausto e não possui mais turnos de luta hoje. Retorne à cidade!")
		return s, nil
	}

	if s.player.Health <= 0 {
		return s, func() tea.Msg {
			return ui.ChangeScreenMsg{Screen: ui.ScreenGameOver}
		}
	}

	_ = s.tm.ConsumeFight(s.player)

	m := s.gen.GenerateForPlayer(s.player.Level)
	s.monster = &m
	s.state = forestStateCombat
	s.combatLog = []string{
		fmt.Sprintf("Você adentra a mata fechada e é surpreendido por um %s!", m.Name),
	}

	SavePlayer(s.db, s.player)
	return s, nil
}

func (s *ForestScreen) handleAttack() (tea.Model, tea.Cmd) {
	if s.monster == nil {
		return s, nil
	}

	res := s.ce.Attack(s.player, s.monster)
	s.appendLog(res.Message)

	if res.MonsterDefeated {
		s.state = forestStateVictory
		s.appendLog("Pressione [P]rocurar para outra luta ou [V]oltar para a cidade.")
		SavePlayer(s.db, s.player)
	} else if res.PlayerDefeated {
		s.state = forestStateDefeat
		s.appendLog("Pressione qualquer tecla para prosseguir...")
		SavePlayer(s.db, s.player)
	}

	return s, nil
}

func (s *ForestScreen) handleFlee() (tea.Model, tea.Cmd) {
	if s.monster == nil {
		return s, nil
	}

	res := s.ce.AttemptFlee(s.player, s.monster)
	s.appendLog(res.Message)

	if res.FledSuccessfully {
		s.state = forestStateFled
		s.appendLog("Pressione [P]rocurar para outra luta ou [V]oltar para a cidade.")
		SavePlayer(s.db, s.player)
	} else if res.PlayerDefeated {
		s.state = forestStateDefeat
		s.appendLog("Pressione qualquer tecla para prosseguir...")
		SavePlayer(s.db, s.player)
	}

	return s, nil
}

func (s *ForestScreen) handlePotion() (tea.Model, tea.Cmd) {
	healed, err := s.ce.UsePotion(s.player)
	if err != nil {
		s.appendLog(fmt.Sprintf("⚠ %v", err))
		return s, nil
	}

	s.appendLog(fmt.Sprintf("Você bebeu uma Poção de Vida e recuperou %d pontos de HP! (Poções restantes: %d)", healed, s.player.PotionsCount))
	SavePlayer(s.db, s.player)
	return s, nil
}

func (s *ForestScreen) appendLog(msg string) {
	s.combatLog = append(s.combatLog, msg)
	if len(s.combatLog) > 8 {
		s.combatLog = s.combatLog[len(s.combatLog)-8:]
	}
}

// View renders the forest environment and combat screen.
func (s *ForestScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("🌲  " + i18n.GetLocationName(i18n.LocationForest) + "  🌲")
	b.WriteString(title + "\n\n")

	var content strings.Builder

	if s.state == forestStateExploring {
		content.WriteString("As árvores antigas sussurram enquanto a névoa fria cobre o chão de folhas.\n")
		content.WriteString("Criaturas perigosas espreitam nas sombras à espera de aventureiros desavisados.\n\n")
		content.WriteString(fmt.Sprintf("Lutas restantes hoje: %s\n\n", ui.StatusFights.Render(fmt.Sprintf("%d", s.player.ForestFights))))
		content.WriteString(ui.SelectedMenuItemStyle.Render("> [P]rocurar Monstros na Mata") + "\n")
		content.WriteString(ui.MenuItemStyle.Render("  [V]oltar para a Praça Central") + "\n")
	} else {
		if s.monster != nil {
			mHPPercent := float64(s.monster.Health) / float64(s.monster.MaxHealth)
			mHPBar := renderSimpleBar(mHPPercent, 10)

			content.WriteString(fmt.Sprintf("Inimigo: %s  (Tier %d)\n", ui.LogMonsterStyle.Render(s.monster.Name), s.monster.Tier))
			content.WriteString(fmt.Sprintf("Vida do Inimigo: %s %d/%d  |  ATK: %d  |  DEF: %d\n\n",
				mHPBar, s.monster.Health, s.monster.MaxHealth, s.monster.Attack, s.monster.Defense))
		}

		content.WriteString("Registro da Batalha:\n")
		for _, line := range s.combatLog {
			content.WriteString(" • " + line + "\n")
		}
		content.WriteString("\n")

		if s.state == forestStateCombat {
			content.WriteString(fmt.Sprintf("%s   %s   %s   %s\n",
				ui.KeyShortcutStyle.Render("[A]tacar"),
				ui.KeyShortcutStyle.Render("[F]ugir"),
				ui.KeyShortcutStyle.Render(fmt.Sprintf("[P]oção (%d)", s.player.PotionsCount)),
				ui.KeyShortcutStyle.Render("[V]oltar (Apenas fora de combate)"),
			))
		} else if s.state == forestStateVictory || s.state == forestStateFled {
			content.WriteString(fmt.Sprintf("%s   %s\n",
				ui.SelectedMenuItemStyle.Render("> [P]rocurar Próximo Inimigo"),
				ui.KeyShortcutStyle.Render("[V]oltar para a Cidade"),
			))
		} else if s.state == forestStateDefeat {
			content.WriteString(ui.ErrorNoticeStyle.Render("💀 VOCÊ FOI DERROTADO! Pressione [Enter] para continuar..."))
		}
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[A] Atacar • [F] Fugir • [P] Poção/Procurar • [V] Voltar à Praça"))

	return ui.AppStyle.Render(b.String())
}

func renderSimpleBar(ratio float64, blocks int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(blocks))
	empty := blocks - filled
	return fmt.Sprintf("[%s%s]", strings.Repeat("█", filled), strings.Repeat("░", empty))
}
