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

	// 4. Test Player Save & Ranking
	p.Level = 5
	p.DragonKills = 1
	p.Gold = 500
	if err := playerRepo.Save(ctx, p); err != nil {
		t.Fatalf("Save player failed: %v", err)
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
