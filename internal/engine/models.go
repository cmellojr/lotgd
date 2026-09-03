package engine

import (
	"time"

	"lotgd/internal/i18n"
	"lotgd/internal/storage"
)

// Player representa o modelo de domínio do herói em memória durante o jogo.
//
// Separamos o modelo de domínio (engine.Player) do modelo de banco (storage.Player)
// para manter a lógica de negócio pura, isolada de detalhes de persistência e serialização.
type Player struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Level        int       `json:"level"`
	Experience   int       `json:"experience"`
	Gold         int       `json:"gold"`
	BankGold     int       `json:"bank_gold"`
	Health       int       `json:"health"`
	MaxHealth    int       `json:"max_health"`
	BaseAttack   int       `json:"base_attack"`
	BaseDefense  int       `json:"base_defense"`
	Weapon       Item      `json:"weapon"`
	Armor        Item      `json:"armor"`
	PotionsCount int       `json:"potions_count"`
	ForestFights int       `json:"forest_fights"`
	DragonKills  int       `json:"dragon_kills"`
	LastLoginDay string    `json:"last_login_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TotalAttack calcula o poder de ataque total do jogador (Base + Bônus da Arma).
func (p *Player) TotalAttack() int {
	return p.BaseAttack + p.Weapon.PowerBonus
}

// TotalDefense calcula a capacidade defensiva total do jogador (Base + Bônus da Armadura).
func (p *Player) TotalDefense() int {
	return p.BaseDefense + p.Armor.PowerBonus
}

// IsAlive verifica se o jogador ainda possui pontos de vida restantes.
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// Monster representa um inimigo encontrado na floresta ou no covil.
type Monster struct {
	ID         i18n.MonsterID `json:"id"`
	Name       string         `json:"name"` // Nome completo formatado com prefixo (ex: "Feroz Rato-do-Esgoto")
	Tier       int            `json:"tier"` // Nível de ameaça (1 a 4, ou 5 para o Dragão)
	Health     int            `json:"health"`
	MaxHealth  int            `json:"max_health"`
	Attack     int            `json:"attack"`
	Defense    int            `json:"defense"`
	XPReward   int            `json:"xp_reward"`
	GoldReward int            `json:"gold_reward"`
	Prefix     string         `json:"prefix,omitempty"` // Afixo procedural aplicado
	IsDragon   bool           `json:"is_dragon"`
}

// IsAlive indica se o monstro ainda está em combate.
func (m *Monster) IsAlive() bool {
	return m.Health > 0
}

// NewPlayerFromStorage converte a struct persistida do SQLite para o modelo de domínio.
//
// Didática Go: Esta função fábrica hidrata o modelo de domínio buscando referências
// de itens no catálogo de forma segura.
func NewPlayerFromStorage(sp storage.Player) *Player {
	weapon, found := FindWeapon(i18n.ItemID(sp.WeaponID))
	if !found {
		weapon = WeaponsCatalog[0] // Fallback para pedaço de pau
	}

	armor, found := FindArmor(i18n.ItemID(sp.ArmorID))
	if !found {
		armor = ArmorsCatalog[0] // Fallback para roupas simples
	}

	return &Player{
		ID:           sp.ID,
		Username:     sp.Username,
		Level:        sp.Level,
		Experience:   sp.Experience,
		Gold:         sp.Gold,
		BankGold:     sp.BankGold,
		Health:       sp.Health,
		MaxHealth:    sp.MaxHealth,
		BaseAttack:   sp.Attack,
		BaseDefense:  sp.Defense,
		Weapon:       weapon,
		Armor:        armor,
		PotionsCount: sp.PotionsCount,
		ForestFights: sp.ForestFights,
		DragonKills:  sp.DragonKills,
		LastLoginDay: sp.LastLoginDay,
		CreatedAt:    sp.CreatedAt,
		UpdatedAt:    sp.UpdatedAt,
	}
}

// ToStorage converte o modelo de domínio de volta para a struct de persistência.
func (p *Player) ToStorage() storage.Player {
	return storage.Player{
		ID:           p.ID,
		Username:     p.Username,
		Level:        p.Level,
		Experience:   p.Experience,
		Gold:         p.Gold,
		BankGold:     p.BankGold,
		Health:       p.Health,
		MaxHealth:    p.MaxHealth,
		Attack:       p.BaseAttack,
		Defense:      p.BaseDefense,
		WeaponID:     string(p.Weapon.ID),
		ArmorID:      string(p.Armor.ID),
		PotionsCount: p.PotionsCount,
		ForestFights: p.ForestFights,
		DragonKills:  p.DragonKills,
		LastLoginDay: p.LastLoginDay,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
