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

type guildState int

const (
	guildStateMenu guildState = iota
	guildStateCombat
	guildStateVictory
	guildStateDefeat
)

// GuildScreen gerencia o treinamento de promoção de nível e avanço de atributos com os Mestres na Guilda.
//
// Didática TEA: O `GuildScreen` orquestra o combate de promoção por turnos contra o Mestre de Nível
// correspondente em Turgon's Warrior Training (ADR-0006 / LORD 1989).
type GuildScreen struct {
	db            *storage.DB
	player        *engine.Player
	ce            *engine.CombatEngine
	state         guildState
	masterMonster *engine.Monster
	cursor        int
	menuItems     []string
	combatLog     []string
	infoMsg       string
	width         int
	height        int
}

// NewGuildScreen inicializa a interface da guilda dos aventureiros.
func NewGuildScreen(db *storage.DB, player *engine.Player) *GuildScreen {
	return &GuildScreen{
		db:        db,
		player:    player,
		ce:        engine.NewCombatEngine(nil),
		state:     guildStateMenu,
		cursor:    0,
		combatLog: make([]string, 0),
		menuItems: []string{
			"Desafiar Mestre de Nível (Turgon's Warrior Training)",
			"Consultar Requisitos de Níveis e Maestria",
			"Ler Pergaminhos Antigos da Lenda do Dragão",
			"Voltar para a Praça Central",
		},
		infoMsg: "Prateleiras infinitas de tomos antigos cobrem as paredes de pedra da guilda.",
	}
}

// Init inicializa a tela da guilda.
func (s *GuildScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer atualiza a referência ao herói ativo em memória.
func (s *GuildScreen) SetPlayer(p *engine.Player) {
	s.player = p
	s.state = guildStateMenu
	s.cursor = 0
	s.combatLog = nil
}

// SetSize atualiza as dimensões de largura e altura da tela.
func (s *GuildScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processa interações de promoção e consulta na guilda dos aventureiros.
func (s *GuildScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.state == guildStateMenu {
			switch msg.String() {
			case "up", "k":
				if s.cursor > 0 {
					s.cursor--
				} else {
					s.cursor = len(s.menuItems) - 1
				}
				return s, nil
			case "down", "j":
				if s.cursor < len(s.menuItems)-1 {
					s.cursor++
				} else {
					s.cursor = 0
				}
				return s, nil
			}

			k := strings.ToUpper(msg.String())

			switch k {
			case "V", "ESC":
				return s.backToTown()
			case "D", "T", "A":
				s.cursor = 0
				return s.selectCurrent()
			case "R", "C":
				s.cursor = 1
				return s.selectCurrent()
			case "P", "L":
				s.cursor = 2
				return s.selectCurrent()
			case "ENTER":
				return s.selectCurrent()
			}
		} else if s.state == guildStateCombat {
			k := strings.ToUpper(msg.String())
			switch k {
			case "A", "ENTER":
				return s.handleAttack()
			case "P", "U":
				return s.handlePotion()
			case "F":
				return s.handleFlee()
			}
		} else if s.state == guildStateVictory || s.state == guildStateDefeat {
			k := strings.ToUpper(msg.String())
			switch k {
			case "ENTER", "V", "C":
				s.state = guildStateMenu
				return s, nil
			}
		}
	}

	return s, nil
}

func (s *GuildScreen) selectCurrent() (tea.Model, tea.Cmd) {
	switch s.cursor {
	case 0: // Desafiar Mestre de Nível
		req, ok := engine.NextLevelRequirement(s.player.Level)
		if !ok {
			s.infoMsg = "Você já alcançou o nível máximo de maestria!"
			return s, nil
		}

		master, hasMaster := engine.GetMasterForTargetLevel(req.Level)
		if !hasMaster {
			s.infoMsg = "Mestre de Treinamento não encontrado para este nível."
			return s, nil
		}

		if s.player.Experience < req.RequiredXP {
			s.infoMsg = fmt.Sprintf("Experiência insuficiente para desafiar %s. Necessário: %d XP (Você tem: %d XP)",
				master.Name, req.RequiredXP, s.player.Experience)
			return s, nil
		}

		if s.player.MasterFoughtToday {
			s.infoMsg = "Você já desafiou seu Mestre hoje! Apenas 1 tentativa por dia é permitida. Retorne amanhã para tentar novamente."
			return s, nil
		}

		// Instancia o combate contra o Mestre
		s.masterMonster = &engine.Monster{
			ID:        i18n.MonsterID(fmt.Sprintf("master_%d", req.Level)),
			Name:      master.Name,
			Tier:      s.player.Level,
			Health:    master.Health,
			MaxHealth: master.MaxHealth,
			Attack:    master.Attack,
			Defense:   master.Defense,
		}

		s.state = guildStateCombat
		s.combatLog = []string{
			fmt.Sprintf("Você adentra a arena da guilda e desafia %s pelo Nível %d!", master.Name, req.Level),
		}
		return s, nil

	case 1: // Consultar tabela de níveis
		req, ok := engine.NextLevelRequirement(s.player.Level)
		if !ok {
			s.infoMsg = "Você já atingiu o Nível 12 (Grau Máximo de Mestre da Guilda)!"
		} else {
			master, _ := engine.GetMasterForTargetLevel(req.Level)
			s.infoMsg = fmt.Sprintf("Próximo Nível (%d - %s): Requer %d XP. Recompensas: +%d HP, +%d ATK, +%d DEF.",
				req.Level, master.Name, req.RequiredXP, req.HealthGain, req.AttackGain, req.DefGain)
		}

	case 2: // Pergaminhos do Dragão
		s.infoMsg = "Mestre Tobias ajusta seus óculos de leitura: 'Reza a lenda que o Dragão Ancestral desperta a cada alvorecer. Apenas um herói que tenha atingido o Nível 12 e provado sua maestria contra o Mestre Turgon poderá desafiar seu covil!'"

	case 3: // Voltar
		return s.backToTown()
	}

	return s, nil
}

func (s *GuildScreen) handleAttack() (tea.Model, tea.Cmd) {
	if s.masterMonster == nil {
		return s, nil
	}

	res := s.ce.Attack(s.player, s.masterMonster)
	s.appendLog(res.Message)

	if res.MonsterDefeated {
		s.state = guildStateVictory
		_ = engine.LevelUp(s.player)
		SavePlayer(s.db, s.player)
		s.appendLog(fmt.Sprintf("🏆 VITÓRIA! Você derrotou %s e foi promovido ao NÍVEL %d!", s.masterMonster.Name, s.player.Level))
		s.appendLog("Pressione [Enter] para retornar ao hall da Guilda!")
	} else if res.PlayerDefeated {
		s.state = guildStateDefeat
		// Regra LORD 1989: Sem morte e sem perda de HP ao ser derrotado pelo Mestre!
		s.player.Health = s.player.MaxHealth
		s.player.MasterFoughtToday = true
		SavePlayer(s.db, s.player)
		s.appendLog(fmt.Sprintf("🛡️ %s venceu este duelo, mas elogiou sua coragem. Sem perdas graves, retorne amanhã!", s.masterMonster.Name))
		s.appendLog("Pressione [Enter] para continuar...")
	}

	return s, nil
}

func (s *GuildScreen) handlePotion() (tea.Model, tea.Cmd) {
	healed, err := s.ce.UsePotion(s.player)
	if err != nil {
		s.appendLog(fmt.Sprintf("⚠ %v", err))
		return s, nil
	}

	s.appendLog(fmt.Sprintf("Você usou uma Poção de Vida (+%d HP)! Poções restantes: %d", healed, s.player.PotionsCount))
	SavePlayer(s.db, s.player)
	return s, nil
}

func (s *GuildScreen) handleFlee() (tea.Model, tea.Cmd) {
	res := s.ce.AttemptFlee(s.player, s.masterMonster)
	s.appendLog(res.Message)

	if res.FledSuccessfully {
		s.state = guildStateMenu
		s.infoMsg = "Você recuou do duelo e preservou sua tentativa diária."
		SavePlayer(s.db, s.player)
	} else if res.PlayerDefeated {
		s.state = guildStateDefeat
		s.player.Health = s.player.MaxHealth
		s.player.MasterFoughtToday = true
		SavePlayer(s.db, s.player)
		s.appendLog("Pressione [Enter] para continuar...")
	}

	return s, nil
}

func (s *GuildScreen) appendLog(msg string) {
	s.combatLog = append(s.combatLog, msg)
	if len(s.combatLog) > 8 {
		s.combatLog = s.combatLog[len(s.combatLog)-8:]
	}
}

func (s *GuildScreen) backToTown() (tea.Model, tea.Cmd) {
	SavePlayer(s.db, s.player)
	return s, func() tea.Msg {
		return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
	}
}

// View renderiza a interface da guilda dos aventureiros no terminal.
func (s *GuildScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("📜  " + i18n.GetLocationName(i18n.LocationGuild) + "  📜")
	b.WriteString(title + "\n\n")

	var content strings.Builder

	if s.state == guildStateMenu {
		content.WriteString(fmt.Sprintf("%s ergue a cabeça por cima de uma pilha de tomos encadernados.\n", i18n.GetNPCName(i18n.NPCTobias)))
		content.WriteString(fmt.Sprintf("Sua Experiência Atual: %s XP  |  Nível Atual: %s\n\n",
			ui.StatusValue.Render(fmt.Sprintf("%d", s.player.Experience)),
			ui.StatusValue.Render(fmt.Sprintf("%d", s.player.Level)),
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
	} else {
		mHPPercent := float64(s.masterMonster.Health) / float64(s.masterMonster.MaxHealth)
		mHPBar := renderSimpleBar(mHPPercent, 12)

		content.WriteString(fmt.Sprintf("Mestre de Treinamento: %s\n", ui.LogMonsterStyle.Render(s.masterMonster.Name)))
		content.WriteString(fmt.Sprintf("Vida do Mestre: %s %d/%d  |  ATK: %d  |  DEF: %d\n\n",
			mHPBar, s.masterMonster.Health, s.masterMonster.MaxHealth, s.masterMonster.Attack, s.masterMonster.Defense))

		content.WriteString("Registro do Duelo:\n")
		for _, line := range s.combatLog {
			content.WriteString(" • " + line + "\n")
		}
		content.WriteString("\n")

		if s.state == guildStateCombat {
			content.WriteString(fmt.Sprintf("%s   %s   %s\n",
				ui.KeyShortcutStyle.Render("[A]tacar"),
				ui.KeyShortcutStyle.Render(fmt.Sprintf("[P]oção (%d)", s.player.PotionsCount)),
				ui.KeyShortcutStyle.Render("[F]ugir"),
			))
		} else if s.state == guildStateVictory {
			content.WriteString(ui.SuccessNoticeStyle.Render("🏆 PROMOÇÃO CONQUISTADA! Pressione [Enter] para continuar..."))
		} else if s.state == guildStateDefeat {
			content.WriteString(ui.ErrorNoticeStyle.Render("🛡️ DUELO ENCERRADO SEM PENALIDADES! Pressione [Enter] para continuar..."))
		}
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[↑/↓] Selecionar • [Enter] Confirmar • [V] Voltar"))

	return ui.AppStyle.Render(b.String())
}
