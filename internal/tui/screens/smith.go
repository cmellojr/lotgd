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

type smithTab int

const (
	smithTabWeapons smithTab = iota
	smithTabArmors
	smithTabPotions
)

// SmithScreen gerencia o comércio de armas, armaduras e consumíveis na Ferraria do Mestre Torin.
//
// Didática TEA: O `SmithScreen` alterna entre abas (`smithTabWeapons`, `smithTabArmors`, `smithTabPotions`),
// exibindo os catálogos estáticos do pacote `engine`, validando moedas do jogador e equipando novos itens.
type SmithScreen struct {
	db      *storage.DB
	player  *engine.Player
	tab     smithTab
	cursor  int
	infoMsg string
	width   int
	height  int
}

// NewSmithScreen inicializa a loja da ferraria com mensagens e aba padrão.
func NewSmithScreen(db *storage.DB, player *engine.Player) *SmithScreen {
	return &SmithScreen{
		db:      db,
		player:  player,
		tab:     smithTabWeapons,
		cursor:  0,
		infoMsg: "Mestre Torin martela uma lâmina incandescente: 'Procurando aço de qualidade, forasteiro?'",
	}
}

// Init inicializa a tela do ferreiro.
func (s *SmithScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer atualiza a referência ao herói ativo em memória.
func (s *SmithScreen) SetPlayer(p *engine.Player) {
	s.player = p
}

// SetSize atualiza as dimensões de largura e altura da tela.
func (s *SmithScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processa a navegação por abas (`1`, `2`, `3`, `Tab`), movimentação do cursor e compras (`Enter`).
func (s *SmithScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			s.tab = smithTabWeapons
			s.cursor = 0
			return s, nil
		case "2":
			s.tab = smithTabArmors
			s.cursor = 0
			return s, nil
		case "3":
			s.tab = smithTabPotions
			s.cursor = 0
			return s, nil
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil
		case "down", "j":
			maxLen := s.getCurrentCatalogLen()
			if s.cursor < maxLen-1 {
				s.cursor++
			}
			return s, nil
		}

		k := strings.ToUpper(msg.String())

		switch k {
		case "V", "ESC":
			SavePlayer(s.db, s.player)
			return s, func() tea.Msg {
				return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
			}
		case "TAB", "RIGHT":
			s.tab = (s.tab + 1) % 3
			s.cursor = 0
			return s, nil
		case "LEFT":
			if s.tab > 0 {
				s.tab--
			} else {
				s.tab = smithTabPotions
			}
			s.cursor = 0
			return s, nil
		case "ENTER", "C":
			return s.handlePurchase()
		}
	}

	return s, nil
}

func (s *SmithScreen) getCurrentCatalogLen() int {
	switch s.tab {
	case smithTabWeapons:
		return len(engine.WeaponsCatalog)
	case smithTabArmors:
		return len(engine.ArmorsCatalog)
	case smithTabPotions:
		return len(engine.PotionsCatalog)
	}
	return 0
}

func (s *SmithScreen) handlePurchase() (tea.Model, tea.Cmd) {
	switch s.tab {
	case smithTabWeapons:
		weapon := engine.WeaponsCatalog[s.cursor]
		if weapon.ID == s.player.Weapon.ID {
			s.infoMsg = "Você já está empunhando esta arma."
			return s, nil
		}
		if s.player.Gold < weapon.Value {
			s.infoMsg = fmt.Sprintf("Ouro insuficiente! %s custa %d moedas.", i18n.GetItemName(weapon.ID), weapon.Value)
			return s, nil
		}

		s.player.Gold -= weapon.Value
		s.player.Weapon = weapon
		SavePlayer(s.db, s.player)
		s.infoMsg = fmt.Sprintf("Você comprou e equipou: %s (+%d ATK)!", i18n.GetItemName(weapon.ID), weapon.PowerBonus)

	case smithTabArmors:
		armor := engine.ArmorsCatalog[s.cursor]
		if armor.ID == s.player.Armor.ID {
			s.infoMsg = "Você já está vestindo esta armadura."
			return s, nil
		}
		if s.player.Gold < armor.Value {
			s.infoMsg = fmt.Sprintf("Ouro insuficiente! %s custa %d moedas.", i18n.GetItemName(armor.ID), armor.Value)
			return s, nil
		}

		s.player.Gold -= armor.Value
		s.player.Armor = armor
		SavePlayer(s.db, s.player)
		s.infoMsg = fmt.Sprintf("Você comprou e equipou: %s (+%d DEF)!", i18n.GetItemName(armor.ID), armor.PowerBonus)

	case smithTabPotions:
		potion := engine.PotionsCatalog[s.cursor]
		if s.player.Gold < potion.Value {
			s.infoMsg = fmt.Sprintf("Ouro insuficiente! %s custa %d moedas.", i18n.GetItemName(potion.ID), potion.Value)
			return s, nil
		}

		s.player.Gold -= potion.Value
		s.player.PotionsCount++
		SavePlayer(s.db, s.player)
		s.infoMsg = fmt.Sprintf("Você comprou uma %s! (Total na bolsa: %d)", i18n.GetItemName(potion.ID), s.player.PotionsCount)
	}

	return s, nil
}

// View renderiza os catálogos da ferraria organizados em abas de navegação.
func (s *SmithScreen) View() string {
	var b strings.Builder

	b.WriteString(ui.RenderStatusBar(s.player, s.width) + "\n")

	title := ui.TitleStyle.Render("⚒  " + i18n.GetLocationName(i18n.LocationSmith) + "  ⚒")
	b.WriteString(title + "\n\n")

	tabWeapons := "[1] Armas de Combate"
	tabArmors := "[2] Armaduras de Proteção"
	tabPotions := "[3] Poções & Elixires"

	if s.tab == smithTabWeapons {
		tabWeapons = ui.SelectedMenuItemStyle.Render("> " + tabWeapons + " <")
	}
	if s.tab == smithTabArmors {
		tabArmors = ui.SelectedMenuItemStyle.Render("> " + tabArmors + " <")
	}
	if s.tab == smithTabPotions {
		tabPotions = ui.SelectedMenuItemStyle.Render("> " + tabPotions + " <")
	}

	tabsLine := fmt.Sprintf("%s    %s    %s", tabWeapons, tabArmors, tabPotions)
	b.WriteString(tabsLine + "\n\n")

	var content strings.Builder

	switch s.tab {
	case smithTabWeapons:
		for i, w := range engine.WeaponsCatalog {
			equipped := ""
			if w.ID == s.player.Weapon.ID {
				equipped = " [EQUIPADO]"
			}
			name := i18n.GetItemName(w.ID)
			line := fmt.Sprintf("%-28s | Custo: %4d Ouro | Poder: +%2d ATK%s", name, w.Value, w.PowerBonus, equipped)

			if i == s.cursor {
				content.WriteString(ui.SelectedMenuItemStyle.Render("> "+line) + "\n")
				content.WriteString(ui.HelpFooterStyle.Render("    └ "+w.Description) + "\n")
			} else {
				content.WriteString(ui.MenuItemStyle.Render("  "+line) + "\n")
			}
		}

	case smithTabArmors:
		for i, a := range engine.ArmorsCatalog {
			equipped := ""
			if a.ID == s.player.Armor.ID {
				equipped = " [EQUIPADO]"
			}
			name := i18n.GetItemName(a.ID)
			line := fmt.Sprintf("%-28s | Custo: %4d Ouro | Defesa: +%2d DEF%s", name, a.Value, a.PowerBonus, equipped)

			if i == s.cursor {
				content.WriteString(ui.SelectedMenuItemStyle.Render("> "+line) + "\n")
				content.WriteString(ui.HelpFooterStyle.Render("    └ "+a.Description) + "\n")
			} else {
				content.WriteString(ui.MenuItemStyle.Render("  "+line) + "\n")
			}
		}

	case smithTabPotions:
		for i, p := range engine.PotionsCatalog {
			name := i18n.GetItemName(p.ID)
			line := fmt.Sprintf("%-28s | Custo: %4d Ouro | Cura: +%2d HP", name, p.Value, p.HealAmount)

			if i == s.cursor {
				content.WriteString(ui.SelectedMenuItemStyle.Render("> "+line) + "\n")
				content.WriteString(ui.HelpFooterStyle.Render("    └ "+p.Description) + "\n")
			} else {
				content.WriteString(ui.MenuItemStyle.Render("  "+line) + "\n")
			}
		}
	}

	if s.infoMsg != "" {
		content.WriteString("\n" + ui.LogSystemStyle.Render("ℹ "+s.infoMsg))
	}

	b.WriteString(ui.ContentBoxStyle.Width(76).Render(content.String()))
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[1-3/Tab] Categorias • [↑/↓] Selecionar • [Enter] Comprar • [V] Voltar"))

	return ui.AppStyle.Render(b.String())
}
