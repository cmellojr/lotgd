package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// VillageRepository manages global village state, the Dragon of the Day, and town news.
type VillageRepository struct {
	db *DB
}

// NewVillageRepository creates a new instance of VillageRepository.
func NewVillageRepository(db *DB) *VillageRepository {
	return &VillageRepository{db: db}
}

// GetOrCreateTodayState retrieves or generates the Dragon and state for today.
func (r *VillageRepository) GetOrCreateTodayState(ctx context.Context) (*VillageState, error) {
	today := time.Now().Format("2006-01-02")
	query := `
	SELECT day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, slayer_name
	FROM village_state WHERE day_date = ?
	`
	row := r.db.QueryRowContext(ctx, query, today)
	var state VillageState
	var dragonAliveInt int

	err := row.Scan(
		&state.DayDate, &dragonAliveInt, &state.DragonHP, &state.DragonMaxHP,
		&state.DragonATK, &state.DragonDEF, &state.SlayerName,
	)
	if err == nil {
		state.DragonAlive = dragonAliveInt == 1
		return &state, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query village state: %w", err)
	}

	// Generate a new Dragon with dynamic attributes for the new day
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxHP := 240 + rng.Intn(60) // 240 - 300 HP
	atk := 32 + rng.Intn(10)    // 32 - 42 ATK
	def := 18 + rng.Intn(8)     // 18 - 26 DEF

	insertQuery := `
	INSERT INTO village_state (day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, slayer_name)
	VALUES (?, 1, ?, ?, ?, ?, '')
	`
	_, err = r.db.ExecContext(ctx, insertQuery, today, maxHP, maxHP, atk, def)
	if err != nil {
		return nil, fmt.Errorf("failed to generate daily dragon: %w", err)
	}

	return &VillageState{
		DayDate:     today,
		DragonAlive: true,
		DragonHP:    maxHP,
		DragonMaxHP: maxHP,
		DragonATK:   atk,
		DragonDEF:   def,
		SlayerName:  "",
	}, nil
}

// RecordDragonSlayed updates the Dragon state when defeated by a hero.
func (r *VillageRepository) RecordDragonSlayed(ctx context.Context, slayerName string) error {
	today := time.Now().Format("2006-01-02")
	query := `
	UPDATE village_state SET dragon_alive = 0, dragon_hp = 0, slayer_name = ?
	WHERE day_date = ?
	`
	_, err := r.db.ExecContext(ctx, query, slayerName, today)
	if err != nil {
		return err
	}

	news := fmt.Sprintf("GLÓRIA AO HERÓI! %s derrotou o Dragão e salvou o Vilarejo hoje!", slayerName)
	return r.AddNews(ctx, news)
}

// AddNews inserts a new event to the village board.
func (r *VillageRepository) AddNews(ctx context.Context, message string) error {
	query := `INSERT INTO news (message, created_at) VALUES (?, CURRENT_TIMESTAMP)`
	_, err := r.db.ExecContext(ctx, query, message)
	return err
}

// GetLatestNews retrieves recent news entries.
func (r *VillageRepository) GetLatestNews(ctx context.Context, limit int) ([]*NewsEntry, error) {
	query := `SELECT id, message, created_at FROM news ORDER BY id DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []*NewsEntry
	for rows.Next() {
		var n NewsEntry
		if err := rows.Scan(&n.ID, &n.Message, &n.CreatedAt); err != nil {
			return nil, err
		}
		news = append(news, &n)
	}
	return news, nil
}
