package storage

import (
	"time"
)

// Player representa a entidade plana de persistência relacional do herói no banco de dados SQLite.
//
// Didática Go: A tag `json:"-"` no campo `PasswordHash` impede que o hash Bcrypt da senha do usuário
// seja exposto acidentalmente ao serializar o jogador em logs ou JSON.
type Player struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Level        int       `json:"level"`
	Experience   int       `json:"experience"`
	Gold         int       `json:"gold"`
	BankGold     int       `json:"bank_gold"`
	Health       int       `json:"health"`
	MaxHealth    int       `json:"max_health"`
	Attack       int       `json:"attack"`
	Defense      int       `json:"defense"`
	WeaponID     string    `json:"weapon_id"`
	ArmorID      string    `json:"armor_id"`
	PotionsCount int       `json:"potions_count"`
	ForestFights      int       `json:"forest_fights"`
	DragonKills       int       `json:"dragon_kills"`
	MasterFoughtToday bool      `json:"master_fought_today"`
	LastLoginDay      string    `json:"last_login_day"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// VillageState armazena os dados globais diários do vilarejo, como o status do Dragão do Dia e o herói vitorioso.
type VillageState struct {
	DayDate          string `json:"day_date"`
	DragonAlive      bool   `json:"dragon_alive"`
	DragonHP         int    `json:"dragon_hp"`
	DragonMaxHP      int    `json:"dragon_max_hp"`
	DragonATK        int    `json:"dragon_atk"`
	DragonDEF        int    `json:"dragon_def"`
	DragonGoldReward int    `json:"dragon_gold_reward"`
	SlayerName       string `json:"slayer_name"`
}

// NewsEntry representa um comunicado, notícia ou fofoca publicado no mural da taverna do vilarejo.
type NewsEntry struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
