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

type dragonState int

const (
	dragonStateApproach dragonState = iota
	dragonStateCombat
	dragonStateVictory
	dragonStateDefeat
)

// DragonScreen handles the legendary confrontation against the Daily Dragon boss.
type DragonScreen struct {
	db        *storage.DB
	player    *engine.Player
	ce        *engine.CombatEngine
	state     dragonState
	dragon    *engine.Monster
	slayer    string
	isAlive   bool
	combatLog []string
	width     int
	height    int
	dragonGen storage.DragonGenerator
}

// NewDragonScreen initializes the Dragon's Lair.
func NewDragonScreen(db *storage.DB, player *engine.Player, dragonGen storage.DragonGenerator) *DragonScreen {
	return &DragonScreen{
		db:        db,
		player:    player,
		ce:        engine.NewCombatEngine(nil),
		state:     dragonStateApproach,
		combatLog: make([]string, 0),
		dragonGen: dragonGen,
	}
}

// Init starts the dragon screen.
func (s *DragonScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the player state and queries daily dragon status.
func (s *DragonScreen) SetPlayer(p *engine.Player) {
	s.player = p
	s.state = dragonStateApproach
	s.combatLog = nil
	s.loadDragonState()
}

// SetSize updates screen dimensions.
func (s *DragonScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

func (s *DragonScreen) loadDragonState() {
	vRepo := storage.NewVillageRepository(s.db, storage.WithDragonGenerator(s.dragonGen))
	st, err := vRepo.GetOrCreateTodayState(context.Background())
	if err == nil {
		s.isAlive = st.DragonAlive
		s.slayer = st.SlayerName

		// Título determinístico baseado na data (mesmo seed SHA-256 que o bestiary)
		titleIdx := engine.DragonTitleIndex(st.DayDate, len(i18n.DragonTitlesPTBR))
		fullName := fmt.Sprintf("%s, %s", i18n.GetMonsterName(i18n.MonsterDragon), i18n.DragonTitlesPTBR[titleIdx])

		s.dragon = &engine.Monster{
			ID:         i18n.MonsterDragon,
			Name:       fullName,
			Tier:       5,
			Health:     st.DragonHP,
			MaxHealth:  st.DragonMaxHP,
			Attack:     st.DragonATK,
			Defense:    st.DragonDEF,
			XPReward:   5000,
			GoldReward: st.DragonGoldReward,
			IsDragon:   true,
		}
	}
}

// Update processes dragon combat and lair interactions.
func (s *DragonScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToUpper(msg.String())

		if k == "V" || k == "ESC" {
			if s.state == dragonStateCombat {
				s.appendLog("Não há como recuar agora! O calor sufocante e a fúria do Dragão bloqueiam a saída!")
				return s, nil
			}
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
			}
		}

		switch s.state {
		case dragonStateApproach:
			switch k {
			case "D", "ENTER", "L":
				return s.startDragonFight()
			}

		case dragonStateCombat:
			switch k {
			case "A", "ENTER":
				return s.handleAttack()
			case "P", "U":
				return s.handlePotion()
			case "F":
				return s.handleFlee()
			}

		case dragonStateVictory:
			switch k {
			case "ENTER", "V", "C":
				return s.backToTown()
			}

		case dragonStateDefeat:
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenGameOver}
			}
		}
	}

	return s, nil
}

func (s *DragonScreen) startDragonFight() (tea.Model, tea.Cmd) {
	if !s.isAlive {
		s.appendLog(fmt.Sprintf("O Dragão já foi derrotado hoje pelo lendário herói %s! Aguarde o próximo alvorecer.", s.slayer))
		return s, nil
	}

	if s.player.Level < 5 {
		s.appendLog("Seu nível é muito baixo! Treine na guilda até pelo menos o Nível 5 antes de desafiar o Dragão.")
		return s, nil
	}

	s.state = dragonStateCombat
	s.combatLog = []string{
		"Você empunha sua arma e adentra a câmara vulcânica! O Dragão ruge chamas ancestrais!",
	}
	return s, nil
}

func (s *DragonScreen) handleAttack() (tea.Model, tea.Cmd) {
	if s.dragon == nil {
		return s, nil
	}

	res := s.ce.Attack(s.player, s.dragon)
	s.appendLog(res.Message)

	if res.MonsterDefeated {
		s.state = dragonStateVictory
		vRepo := storage.NewVillageRepository(s.db, storage.WithDragonGenerator(s.dragonGen))
		_ = vRepo.RecordDragonSlayed(context.Background(), s.player.Username)
		SavePlayer(s.db, s.player)
		s.appendLog("🔥 O DRAGÃO CAIU! Seus restos viraram lenda e você salvou todo o Vilarejo! 🔥")
		s.appendLog("Pressione [Enter] para retornar triunfante à Praça do Vilarejo!")
	} else if res.PlayerDefeated {
		s.state = dragonStateDefeat
		s.appendLog("Pressione [Enter] para sucumbir...")
		SavePlayer(s.db, s.player)
	}

	return s, nil
}

func (s *DragonScreen) handlePotion() (tea.Model, tea.Cmd) {
	healed, err := s.ce.UsePotion(s.player)
	if err != nil {
		s.appendLog(fmt.Sprintf("⚠ %v", err))
		return s, nil
	}

	s.appendLog(fmt.Sprintf("Você usou uma Poção de Vida (+%d HP)! Poções restantes: %d", healed, s.player.PotionsCount))
	SavePlayer(s.db, s.player)
	return s, nil
}

func (s *DragonScreen) handleFlee() (tea.Model, tea.Cmd) {
	res := s.ce.AttemptFlee(s.player, s.dragon)
	s.appendLog(res.Message)

	if res.FledSuccessfully {
		s.state = dragonStateApproach
		s.appendLog("Você escapou milagrosamente do sopro abrasador do Dragão e voltou à entrada do covil.")
		SavePlayer(s.db, s.player)
	} else if res.PlayerDefeated {
		s.state = dragonStateDefeat
		s.appendLog("Pressione [Enter] para sucumbir...")
		SavePlayer(s.db, s.player)
	}

	return s, nil
}

func (s *DragonScreen) appendLog(msg string) {
	s.combatLog = append(s.combatLog, msg)
	if len(s.combatLog) > 8 {
		s.combatLog = s.combatLog[len(s.combatLog)-8:]
	}
}

func (s *DragonScreen) backToTown() (tea.Model, tea.Cmd) {
	SavePlayer(s.db, s.player)
	return s, func() tea.Msg {
		return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
	}
}

// View renders the dragon's lair.
func (s *DragonScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("🌋  " + i18n.GetLocationName(i18n.LocationDragon) + "  🌋")
	b.WriteString(title + "\n\n")

	var content strings.Builder

	if s.state == dragonStateApproach {
		if !s.isAlive {
			content.WriteString("As cavernas profundas estão em silêncio. A carcaça do Dragão jaz no abismo.\n")
			content.WriteString(fmt.Sprintf("Ele foi derrotado hoje por %s!\n\n", ui.StatusGold.Render(s.slayer)))
			content.WriteString(ui.MenuItemStyle.Render("  [V]oltar em paz para a Praça Central") + "\n")
		} else {
			content.WriteString("Rios de lava iluminam a colossal câmara rochosa. O ar queima os pulmões.\n")
			content.WriteString(fmt.Sprintf("O Dragão do Dia aguarda: %s (%d/%d HP, ATK %d, DEF %d)\n\n",
				ui.LogMonsterStyle.Render(s.dragon.Name), s.dragon.Health, s.dragon.MaxHealth, s.dragon.Attack, s.dragon.Defense))
			content.WriteString(ui.SelectedMenuItemStyle.Render("> [D]esafiar o Dragão para o Combate Final") + "\n")
			content.WriteString(ui.MenuItemStyle.Render("  [V]oltar para a Praça Central") + "\n")
		}
	} else {
		dHPPercent := float64(s.dragon.Health) / float64(s.dragon.MaxHealth)
		dHPBar := renderSimpleBar(dHPPercent, 12)

		content.WriteString(fmt.Sprintf("Chefe Lendário: %s\n", ui.LogMonsterStyle.Render(s.dragon.Name)))
		content.WriteString(fmt.Sprintf("Vida do Dragão: %s %d/%d  |  ATK: %d  |  DEF: %d\n\n",
			dHPBar, s.dragon.Health, s.dragon.MaxHealth, s.dragon.Attack, s.dragon.Defense))

		content.WriteString("Registro do Confronto:\n")
		for _, line := range s.combatLog {
			content.WriteString(" • " + line + "\n")
		}
		content.WriteString("\n")

		if s.state == dragonStateCombat {
			content.WriteString(fmt.Sprintf("%s   %s   %s\n",
				ui.KeyShortcutStyle.Render("[A]tacar com Fúria"),
				ui.KeyShortcutStyle.Render(fmt.Sprintf("[P]oção (%d)", s.player.PotionsCount)),
				ui.KeyShortcutStyle.Render("[F]ugir Desesperadamente"),
			))
		} else if s.state == dragonStateVictory {
			content.WriteString(ui.SuccessNoticeStyle.Render("🏆 VITÓRIA LENDÁRIA! Pressione [Enter] para comemorar no Vilarejo!"))
		} else if s.state == dragonStateDefeat {
			content.WriteString(ui.ErrorNoticeStyle.Render("💀 O DRAGÃO TE DESTRUIU! Pressione [Enter] para continuar..."))
		}
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[A] Atacar • [P] Poção • [F] Fugir • [V] Voltar"))

	return ui.AppStyle.Render(b.String())
}
