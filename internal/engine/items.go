package engine

import "lotgd/internal/i18n"

// ItemType categoriza os tipos de equipamentos e consumíveis no RPG.
//
// Didática Go: Criamos um tipo alias customizado `type ItemType string` baseado no tipo primitivo
// string para prover segurança de tipo em tempo de compilação (type safety), impedindo a atribuição
// inadvertida de strings arbitrárias em campos de tipo de item.
type ItemType string

const (
	// ItemTypeWeapon define armas equipáveis na mão principal que aumentam o poder de ataque.
	ItemTypeWeapon ItemType = "weapon"

	// ItemTypeArmor define armaduras e vestimentas de proteção que aumentam a defesa total.
	ItemTypeArmor ItemType = "armor"

	// ItemTypePotion define consumíveis da bolsa de utilidades (ex: cura de vida).
	ItemTypePotion ItemType = "potion"
)

// Item representa a entidade de equipamento ou consumível no modelo de domínio do jogo.
//
// Didática Go: A struct agrupa atributos relevantes para compras na ferraria e combate,
// mantendo-se puramente desacoplada de operações de banco de dados ou renderização gráfica.
type Item struct {
	ID          i18n.ItemID `json:"id"`
	Type        ItemType    `json:"type"`
	NameKey     i18n.ItemID `json:"name_key"`
	Value       int         `json:"value"`       // Preço de compra na ferraria
	PowerBonus  int         `json:"power_bonus"` // Bônus de Ataque (se arma) ou Defesa (se armadura)
	HealAmount  int         `json:"heal_amount"` // Quantidade de HP restaurado (se poção)
	Description string      `json:"description"`
}

// WeaponsCatalog contém todas as armas disponíveis no jogo, ordenadas por progressão de poder e preço.
//
// Didática Go: O uso de um slice estático de structs imutáveis evita alocações
// dinâmicas desnecessárias em tempo de execução e garante consultas de catálogo extremamente rápidas.
var WeaponsCatalog = []Item{
	{
		ID:          i18n.WeaponStick,
		Type:        ItemTypeWeapon,
		NameKey:     i18n.WeaponStick,
		Value:       0,
		PowerBonus:  1,
		Description: "Um galho seco encontrado no chão. Melhor que lutar desarmado.",
	},
	{
		ID:          i18n.WeaponDagger,
		Type:        ItemTypeWeapon,
		NameKey:     i18n.WeaponDagger,
		Value:       50,
		PowerBonus:  3,
		Description: "Pequena e afiada, ótima para perfurar pontos fracos de monstros menores.",
	},
	{
		ID:          i18n.WeaponShortSword,
		Type:        ItemTypeWeapon,
		NameKey:     i18n.WeaponShortSword,
		Value:       200,
		PowerBonus:  6,
		Description: "Espada confiável de ferro forjada na oficina do Mestre Torin.",
	},
	{
		ID:          i18n.WeaponBroadsword,
		Type:        ItemTypeWeapon,
		NameKey:     i18n.WeaponBroadsword,
		Value:       600,
		PowerBonus:  12,
		Description: "Lâmina pesada de aço temperado capaz de cortar escamas grossas.",
	},
	{
		ID:          i18n.WeaponDragonSlayer,
		Type:        ItemTypeWeapon,
		NameKey:     i18n.WeaponDragonSlayer,
		Value:       2000,
		PowerBonus:  25,
		Description: "Arma lendária banhada em sangue de monstros antigos, feita para abater dragões.",
	},
}

// ArmorsCatalog contém todas as armaduras e trajes de proteção para compra e equipamento.
var ArmorsCatalog = []Item{
	{
		ID:          i18n.ArmorClothes,
		Type:        ItemTypeArmor,
		NameKey:     i18n.ArmorClothes,
		Value:       0,
		PowerBonus:  0,
		Description: "Roupas rotas de plebeu que não oferecem proteção real contra presas e garras.",
	},
	{
		ID:          i18n.ArmorLeather,
		Type:        ItemTypeArmor,
		NameKey:     i18n.ArmorLeather,
		Value:       40,
		PowerBonus:  2,
		Description: "Colete de couro curtido, leve e flexível para iniciantes na floresta.",
	},
	{
		ID:          i18n.ArmorChainmail,
		Type:        ItemTypeArmor,
		NameKey:     i18n.ArmorChainmail,
		Value:       180,
		PowerBonus:  5,
		Description: "Anéis entrelaçados de ferro que absorvem impactos de cortes e flechadas.",
	},
	{
		ID:          i18n.ArmorPlate,
		Type:        ItemTypeArmor,
		NameKey:     i18n.ArmorPlate,
		Value:       550,
		PowerBonus:  10,
		Description: "Placas maciças de aço polido que protegem o tórax contra golpes devastadores.",
	},
	{
		ID:          i18n.ArmorDragonScale,
		Type:        ItemTypeArmor,
		NameKey:     i18n.ArmorDragonScale,
		Value:       1800,
		PowerBonus:  20,
		Description: "Armadura forjada com as escamas impenetráveis do próprio Dragão Ancestral.",
	},
}

// PotionsCatalog lista os consumíveis e poções de restauração disponíveis na ferraria/loja.
var PotionsCatalog = []Item{
	{
		ID:          i18n.PotionHealth,
		Type:        ItemTypePotion,
		NameKey:     i18n.PotionHealth,
		Value:       25,
		HealAmount:  30,
		Description: "Frasco com líquido vermelho brilhante que cura ferimentos imediatamente.",
	},
	{
		ID:          i18n.PotionGarbageCollector,
		Type:        ItemTypePotion,
		NameKey:     i18n.PotionGarbageCollector,
		Value:       100,
		HealAmount:  100,
		Description: "Poção suprema que limpa todas as imperfeições e restaura a saúde completamente.",
	},
}

// FindWeapon busca uma arma no catálogo pelo seu identificador único.
//
// Didática Go: Retorna a struct do item e um booleano de confirmação `(Item, bool)`,
// seguindo o padrão clássico `comma-ok` idiomático de Go.
func FindWeapon(id i18n.ItemID) (Item, bool) {
	for _, w := range WeaponsCatalog {
		if w.ID == id {
			return w, true
		}
	}
	return Item{}, false
}

// FindArmor busca uma armadura no catálogo pelo seu identificador único.
func FindArmor(id i18n.ItemID) (Item, bool) {
	for _, a := range ArmorsCatalog {
		if a.ID == id {
			return a, true
		}
	}
	return Item{}, false
}

// FindPotion busca uma poção no catálogo pelo seu identificador único.
func FindPotion(id i18n.ItemID) (Item, bool) {
	for _, p := range PotionsCatalog {
		if p.ID == id {
			return p, true
		}
	}
	return Item{}, false
}
