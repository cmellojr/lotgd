package storage

import (
	"context"
	"path/filepath"
	"testing"
)

func TestStorageWorkflow(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_lotgd.db")

	db, err := OpenDB(dbPath)
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}
	defer db.Close()

	playerRepo := NewPlayerRepository(db)
	villageRepo := NewVillageRepository(db)

	// 1. Test Player Registration
	p, err := playerRepo.Register(ctx, "GopherHero", "secret123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if p.Username != "GopherHero" || p.Level != 1 || p.ForestFights != 15 {
		t.Errorf("Unexpected player defaults: %+v", p)
	}

	// 2. Test Player Authentication
	authPlayer, err := playerRepo.Authenticate(ctx, "GopherHero", "secret123")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if authPlayer.ID != p.ID {
		t.Errorf("expected ID %d, got %d", p.ID, authPlayer.ID)
	}

	// 3. Test Invalid Password
	if _, err := playerRepo.Authenticate(ctx, "GopherHero", "wrongpass"); err != ErrInvalidPass {
		t.Errorf("expected ErrInvalidPass, got %v", err)
	}

	// 4. Test Player Save & Ranking & Potions
	p.Level = 5
	p.DragonKills = 1
	p.Gold = 500
	p.PotionsCount = 3
	if err := playerRepo.Save(ctx, p); err != nil {
		t.Fatalf("Save player failed: %v", err)
	}

	savedP, err := playerRepo.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if savedP.PotionsCount != 3 {
		t.Errorf("expected PotionsCount 3, got %d", savedP.PotionsCount)
	}

	rankings, err := playerRepo.ListRankings(ctx, 10)
	if err != nil {
		t.Fatalf("ListRankings failed: %v", err)
	}
	if len(rankings) != 1 || rankings[0].DragonKills != 1 {
		t.Errorf("expected ranking with 1 dragon kill, got %+v", rankings)
	}

	// 5. Test Village State & Daily Dragon
	state, err := villageRepo.GetOrCreateTodayState(ctx)
	if err != nil {
		t.Fatalf("GetOrCreateTodayState failed: %v", err)
	}
	if !state.DragonAlive || state.DragonHP <= 0 {
		t.Errorf("expected alive dragon with HP > 0, got %+v", state)
	}

	// 6. Test Record Dragon Slayed & News
	if err := villageRepo.RecordDragonSlayed(ctx, "GopherHero"); err != nil {
		t.Fatalf("RecordDragonSlayed failed: %v", err)
	}

	news, err := villageRepo.GetLatestNews(ctx, 5)
	if err != nil {
		t.Fatalf("GetLatestNews failed: %v", err)
	}
	if len(news) == 0 {
		t.Fatalf("expected at least 1 news entry")
	}
}

func TestLegacySchemaMigration(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "legacy_lotgd.db")

	// Create legacy database without potions_count column and with user_version = 0 (or 1)
	legacyDB, err := OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open initial db: %v", err)
	}

	// Re-create table without potions_count to simulate a v0.0.1 database
	if _, err := legacyDB.Exec(`
		DROP TABLE players;
		CREATE TABLE players (
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
			forest_fights INTEGER NOT NULL DEFAULT 15,
			dragon_kills INTEGER NOT NULL DEFAULT 0,
			last_login_day TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		PRAGMA user_version = 1;
	`); err != nil {
		t.Fatalf("failed to create legacy schema: %v", err)
	}

	// Insert legacy player record
	if _, err := legacyDB.Exec(`
		INSERT INTO players (username, password_hash, last_login_day)
		VALUES ('LegacyHero', 'hash123', '2026-01-01');
	`); err != nil {
		t.Fatalf("failed to insert legacy hero: %v", err)
	}
	legacyDB.Close()

	// Now open database again with OpenDB which triggers migrate()
	db, err := OpenDB(dbPath)
	if err != nil {
		t.Fatalf("OpenDB failed on legacy database: %v", err)
	}
	defer db.Close()

	repo := NewPlayerRepository(db)
	p, err := repo.GetByUsername(ctx, "LegacyHero")
	if err != nil {
		t.Fatalf("GetByUsername failed on migrated player: %v", err)
	}

	if p.PotionsCount != 0 {
		t.Errorf("expected default PotionsCount 0, got %d", p.PotionsCount)
	}

	p.PotionsCount = 5
	if err := repo.Save(ctx, p); err != nil {
		t.Fatalf("Save failed on migrated player: %v", err)
	}

	updatedP, err := repo.GetByUsername(ctx, "LegacyHero")
	if err != nil {
		t.Fatalf("GetByUsername failed after save: %v", err)
	}
	if updatedP.PotionsCount != 5 {
		t.Errorf("expected PotionsCount 5, got %d", updatedP.PotionsCount)
	}
}
