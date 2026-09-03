package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB instance for LOTGD operations.
type DB struct {
	*sql.DB
}

// OpenDB opens a connection to SQLite database and applies schemas.
func OpenDB(dsn string) (*DB, error) {
	pragmaParams := "_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)"
	formattedDSN := dsn
	if strings.Contains(dsn, "?") {
		formattedDSN += "&" + pragmaParams
	} else {
		formattedDSN += "?" + pragmaParams
	}

	db, err := sql.Open("sqlite", formattedDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// BBS terminal load profile: limit connection pool to 1
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Enable WAL mode (persisted in DB header)
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL journal mode: %w", err)
	}

	lotgdDB := &DB{DB: db}
	if err := lotgdDB.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return lotgdDB, nil
}

// SavePlayer persists player struct changes.
func (d *DB) SavePlayer(p Player) error {
	repo := NewPlayerRepository(d)
	return repo.Save(context.Background(), &p)
}

// AuthenticatePlayer verifies user credentials.
func (d *DB) AuthenticatePlayer(username, password string) (Player, error) {
	repo := NewPlayerRepository(d)
	p, err := repo.Authenticate(context.Background(), username, password)
	if err != nil {
		return Player{}, err
	}
	return *p, nil
}

// CreatePlayer registers a new player.
func (d *DB) CreatePlayer(username, password string) (Player, error) {
	repo := NewPlayerRepository(d)
	p, err := repo.Register(context.Background(), username, password)
	if err != nil {
		return Player{}, err
	}
	return *p, nil
}

func (d *DB) migrate() error {
	var userVersion int
	if err := d.QueryRow("PRAGMA user_version;").Scan(&userVersion); err != nil {
		return fmt.Errorf("failed to read user_version: %w", err)
	}

	if userVersion < 2 {
		var tableExists int
		err := d.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='players'").Scan(&tableExists)
		if err != nil {
			return fmt.Errorf("failed to check players table existence: %w", err)
		}

		if tableExists > 0 {
			var columnExists int
			err := d.QueryRow("SELECT COUNT(*) FROM pragma_table_info('players') WHERE name='potions_count'").Scan(&columnExists)
			if err != nil {
				return fmt.Errorf("failed to check potions_count column: %w", err)
			}
			if columnExists == 0 {
				if _, err := d.Exec("ALTER TABLE players ADD COLUMN potions_count INTEGER NOT NULL DEFAULT 0;"); err != nil {
					return fmt.Errorf("failed to add potions_count column: %w", err)
				}
			}
		}

		if _, err := d.Exec("PRAGMA user_version = 2;"); err != nil {
			return fmt.Errorf("failed to set user_version: %w", err)
		}
	}

	schema := `
	CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL COLLATE NOCASE,
		password_hash TEXT NOT NULL,
		level INTEGER NOT NULL DEFAULT 1,
		experience INTEGER NOT NULL DEFAULT 0,
		gold INTEGER NOT NULL DEFAULT 50,
		bank_gold INTEGER NOT NULL DEFAULT 0,
		health INTEGER NOT NULL DEFAULT 20,
		max_health INTEGER NOT NULL DEFAULT 20,
		attack INTEGER NOT NULL DEFAULT 5,
		defense INTEGER NOT NULL DEFAULT 2,
		weapon_id TEXT NOT NULL DEFAULT 'stick',
		armor_id TEXT NOT NULL DEFAULT 'clothes',
		potions_count INTEGER NOT NULL DEFAULT 0,
		forest_fights INTEGER NOT NULL DEFAULT 15,
		dragon_kills INTEGER NOT NULL DEFAULT 0,
		last_login_day TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS village_state (
		day_date TEXT PRIMARY KEY,
		dragon_alive INTEGER NOT NULL DEFAULT 1,
		dragon_hp INTEGER NOT NULL DEFAULT 250,
		dragon_max_hp INTEGER NOT NULL DEFAULT 250,
		dragon_atk INTEGER NOT NULL DEFAULT 35,
		dragon_def INTEGER NOT NULL DEFAULT 20,
		dragon_gold_reward INTEGER NOT NULL DEFAULT 3000,
		slayer_name TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS graveyard (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		victim_name TEXT NOT NULL,
		killer_name TEXT NOT NULL,
		killed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.Exec(schema)
	return err
}
