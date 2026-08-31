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

// SmithScreen handles weapons, armors, and potions trading with Master Torin.
type SmithScreen struct {
	db       *storage.DB
	player   *engine.Player
	tab      smithTab
	cursor   int
	infoMsg  string
	width    int
	height   int
}

// NewSmithScreen initializes the blacksmith forge.
func NewSmithScreen(db *storage.DB, player *engine.Player) *SmithScreen {
	return &SmithScreen{
		db:      db,
		player:  player,
		tab:     smithTabWeapons,
		cursor:  0,
		infoMsg: "Mestre Torin martela uma lâmina incandescente: 'Procurando aço de qualidade, forasteiro?'",
	}
}

// Init starts the smith screen.
func (s *SmithScreen) Init() tea.Cmd {
	return nil
}

// SetPlayer updates the player reference.
func (s *SmithScreen) SetPlayer(p *engine.Player) {
	s.player = p
}

// SetSize updates screen dimensions.
func (s *SmithScreen) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Update processes blacksmith shop navigation and purchases.
func (s *SmithScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToUpper(msg.String())

		switch k {
		case "V", "ESC":
			_ = s.db.SavePlayer(s.player.ToStorage())
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
		case "UP":
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil
		case "DOWN":
			maxLen := s.getCurrentCatalogLen()
			if s.cursor < maxLen-1 {
				s.cursor++
			}
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
		_ = s.db.SavePlayer(s.player.ToStorage())
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
		_ = s.db.SavePlayer(s.player.ToStorage())
		s.infoMsg = fmt.Sprintf("Você comprou e equipou: %s (+%d DEF)!", i18n.GetItemName(armor.ID), armor.PowerBonus)

	case smithTabPotions:
		potion := engine.PotionsCatalog[s.cursor]
		if s.player.Gold < potion.Value {
			s.infoMsg = fmt.Sprintf("Ouro insuficiente! %s custa %d moedas.", i18n.GetItemName(potion.ID), potion.Value)
			return s, nil
		}

		s.player.Gold -= potion.Value
		s.player.PotionsCount++
		_ = s.db.SavePlayer(s.player.ToStorage())
		s.infoMsg = fmt.Sprintf("Você comprou uma %s! (Total na bolsa: %d)", i18n.GetItemName(potion.ID), s.player.PotionsCount)
	}

	return s, nil
}

// View renders the blacksmith forge.
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
	b.WriteString("\n" + ui.HelpFooterStyle.Render("[Tab/Setas] Trocar Categoria • [Enter] Comprar Item • [V] Voltar"))

	return ui.AppStyle.Render(b.String())
}
