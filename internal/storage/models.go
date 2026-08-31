package storage

import (
	"time"
)

// Player represents a registered adventurer in the database.
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
	ForestFights int       `json:"forest_fights"`
	DragonKills  int       `json:"dragon_kills"`
	LastLoginDay string    `json:"last_login_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// VillageState holds global daily data such as the Dragon of the Day and server day.
type VillageState struct {
	DayDate     string `json:"day_date"`
	DragonAlive bool   `json:"dragon_alive"`
	DragonHP    int    `json:"dragon_hp"`
	DragonMaxHP int    `json:"dragon_max_hp"`
	DragonATK   int    `json:"dragon_atk"`
	DragonDEF   int    `json:"dragon_def"`
	SlayerName  string `json:"slayer_name"`
}

// NewsEntry represents an announcement or gossip on the town board.
type NewsEntry struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
