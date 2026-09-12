package engine

import (
	"time"

	"lotgd/internal/i18n"
	"lotgd/internal/storage"
)

// Player representa o modelo de domínio do herói mantido em memória durante a sessão de jogo.
//
// Didática Go: Separamos o modelo de domínio (`engine.Player`) do modelo de banco de dados (`storage.Player`).
// Essa separação arquitetural isola as regras de negócio de detalhes de persistência e serialização SQL,
// permitindo que o modelo de domínio possua referências ricas (como as structs `Item` em vez de apenas IDs em string).
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
	ForestFights      int       `json:"forest_fights"`
	DragonKills       int       `json:"dragon_kills"`
	MasterFoughtToday bool      `json:"master_fought_today"`
	LastLoginDay      string    `json:"last_login_day"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TotalAttack calcula o poder de ataque total do jogador (Ataque Base + Bônus da Arma equipada).
func (p *Player) TotalAttack() int {
	return p.BaseAttack + p.Weapon.PowerBonus
}

// TotalDefense calcula a capacidade defensiva total do jogador (Defesa Base + Bônus da Armadura equipada).
func (p *Player) TotalDefense() int {
	return p.BaseDefense + p.Armor.PowerBonus
}

// IsAlive verifica se o herói ainda possui pontos de vida positivos restantes.
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// Monster representa a entidade de um inimigo encontrado na floresta ou no covil do chefe.
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

// IsAlive indica se o monstro ainda está ativo no combate.
func (m *Monster) IsAlive() bool {
	return m.Health > 0
}

// NewPlayerFromStorage converte e hidrata a struct de persistência (`storage.Player`) para o modelo de domínio (`engine.Player`).
//
// Didática Go: Esta função fábrica busca e vincula as instâncias completas dos itens (`Weapon` e `Armor`)
// a partir dos catálogos em memória usando seus IDs em string, provendo fallbacks seguros caso um ID seja inválido.
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
		ID:                sp.ID,
		Username:          sp.Username,
		Level:             sp.Level,
		Experience:        sp.Experience,
		Gold:              sp.Gold,
		BankGold:          sp.BankGold,
		Health:            sp.Health,
		MaxHealth:         sp.MaxHealth,
		BaseAttack:        sp.Attack,
		BaseDefense:       sp.Defense,
		Weapon:            weapon,
		Armor:             armor,
		PotionsCount:      sp.PotionsCount,
		ForestFights:      sp.ForestFights,
		DragonKills:       sp.DragonKills,
		MasterFoughtToday: sp.MasterFoughtToday,
		LastLoginDay:      sp.LastLoginDay,
		CreatedAt:         sp.CreatedAt,
		UpdatedAt:         sp.UpdatedAt,
	}
}

// ToStorage converte o modelo de domínio de volta para a struct plana de persistência do SQLite.
//
// Didática Go: Mapeia as referências ricas de volta para IDs em string e atualiza o timestamp `UpdatedAt`.
func (p *Player) ToStorage() storage.Player {
	now := time.Now().UTC()
	p.UpdatedAt = now
	return storage.Player{
		ID:                p.ID,
		Username:          p.Username,
		Level:             p.Level,
		Experience:        p.Experience,
		Gold:              p.Gold,
		BankGold:          p.BankGold,
		Health:            p.Health,
		MaxHealth:         p.MaxHealth,
		Attack:            p.BaseAttack,
		Defense:           p.BaseDefense,
		WeaponID:          string(p.Weapon.ID),
		ArmorID:           string(p.Armor.ID),
		PotionsCount:      p.PotionsCount,
		ForestFights:      p.ForestFights,
		DragonKills:       p.DragonKills,
		MasterFoughtToday: p.MasterFoughtToday,
		LastLoginDay:      p.LastLoginDay,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         now,
	}
}
