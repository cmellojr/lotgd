package storage

import (
	"context"
	"database/sql"
	"errors"
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
	query := `
	SELECT day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, dragon_gold_reward, slayer_name
	FROM village_state WHERE day_date = ?
	`
	row := r.db.QueryRowContext(ctx, query, today)
	var state VillageState
	var dragonAliveInt int

	err := row.Scan(
		&state.DayDate, &dragonAliveInt, &state.DragonHP, &state.DragonMaxHP,
		&state.DragonATK, &state.DragonDEF, &state.DragonGoldReward, &state.SlayerName,
	)
	if err == nil {
		state.DragonAlive = dragonAliveInt == 1
		return &state, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query village state: %w", err)
	}

	// Gera o Dragão do Dia via generator injetado (bestiary.GenerateDragonOfDay).
	// Se o generator não foi injetado (ex.: testes), usa fallback seguro.
	var maxHP, atk, def, goldReward int
	if r.dragonGen != nil {
		maxHP, atk, def, goldReward = r.dragonGen(today)
	} else {
		maxHP, atk, def, goldReward = 300, 45, 25, 3000
	}

	insertQuery := `
	INSERT INTO village_state (day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, dragon_gold_reward, slayer_name)
	VALUES (?, 1, ?, ?, ?, ?, ?, '')
	`
	_, err = r.db.ExecContext(ctx, insertQuery, today, maxHP, maxHP, atk, def, goldReward)
	if err != nil {
		return nil, fmt.Errorf("failed to generate daily dragon: %w", err)
	}

	return &VillageState{
		DayDate:          today,
		DragonAlive:      true,
		DragonHP:         maxHP,
		DragonMaxHP:      maxHP,
		DragonATK:        atk,
		DragonDEF:        def,
		DragonGoldReward: goldReward,
		SlayerName:       "",
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
