package storage

import (
	"context"
	"fmt"
	"time"
)

// DragonGenerator é um contrato de geração do Dragão do Dia.
// O pacote storage não importa bestiary diretamente (evita ciclo de dependência);
// a implementação é injetada pelo chamador (cmd/lotgd, cmd/server) via closure.
//
// Retorna: hp, atk, def, goldReward — os atributos brutos do Dragão para persistência.
type DragonGenerator func(dayDate string) (hp, atk, def, goldReward int)

// VillageRepositoryOption configura opções opcionais do repositório.
type VillageRepositoryOption func(*VillageRepository)

// WithDragonGenerator injeta a função de geração do Dragão do Dia.
// Se não fornecido, GetOrCreateTodayState usa um fallback seguro (para testes).
func WithDragonGenerator(gen DragonGenerator) VillageRepositoryOption {
	return func(r *VillageRepository) { r.dragonGen = gen }
}

// VillageRepository manages global village state, the Dragon of the Day, and town news.
type VillageRepository struct {
	db        *DB
	dragonGen DragonGenerator
}

// NewVillageRepository creates a new instance of VillageRepository.
func NewVillageRepository(db *DB, opts ...VillageRepositoryOption) *VillageRepository {
	r := &VillageRepository{db: db}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// GetOrCreateTodayState retrieves or generates the Dragon and state for today.
func (r *VillageRepository) GetOrCreateTodayState(ctx context.Context) (*VillageState, error) {
	today := time.Now().Format("2006-01-02")

	var maxHP, atk, def, goldReward int
	if r.dragonGen != nil {
		maxHP, atk, def, goldReward = r.dragonGen(today)
	} else {
		maxHP, atk, def, goldReward = 300, 45, 25, 3000
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertQuery := `
	INSERT OR IGNORE INTO village_state (day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, dragon_gold_reward, slayer_name)
	VALUES (?, 1, ?, ?, ?, ?, ?, '')
	`
	if _, err := tx.ExecContext(ctx, insertQuery, today, maxHP, maxHP, atk, def, goldReward); err != nil {
		return nil, fmt.Errorf("failed to ensure daily dragon state: %w", err)
	}

	selectQuery := `
	SELECT day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, dragon_gold_reward, slayer_name
	FROM village_state WHERE day_date = ?
	`
	row := tx.QueryRowContext(ctx, selectQuery, today)
	var state VillageState
	var dragonAliveInt int

	err = row.Scan(
		&state.DayDate, &dragonAliveInt, &state.DragonHP, &state.DragonMaxHP,
		&state.DragonATK, &state.DragonDEF, &state.DragonGoldReward, &state.SlayerName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan village state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	state.DragonAlive = dragonAliveInt == 1
	return &state, nil
}

// RecordDragonSlayed updates the Dragon state when defeated by a hero.
func (r *VillageRepository) RecordDragonSlayed(ctx context.Context, slayerName string) error {
	today := time.Now().Format("2006-01-02")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE village_state SET dragon_alive = 0, dragon_hp = 0, slayer_name = ?
	WHERE day_date = ? AND dragon_alive = 1
	`
	res, err := tx.ExecContext(ctx, query, slayerName, today)
	if err != nil {
		return fmt.Errorf("failed to update dragon status: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("dragon already slain or not found for %s", today)
	}

	news := fmt.Sprintf("GLÓRIA AO HERÓI! %s derrotou o Dragão e salvou o Vilarejo hoje!", slayerName)
	newsQuery := `INSERT INTO news (message, created_at) VALUES (?, CURRENT_TIMESTAMP)`
	if _, err := tx.ExecContext(ctx, newsQuery, news); err != nil {
		return fmt.Errorf("failed to record news: %w", err)
	}

	return tx.Commit()
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
