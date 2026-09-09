package storage

import (
	"context"
	"fmt"
	"time"
)

// DragonGenerator define o contrato da closure geradora dos atributos do Dragão do Dia.
//
// Didática Go: O pacote `storage` não importa o pacote `bestiary` diretamente para evitar ciclos
// de dependência circular (`storage` -> `bestiary` -> `engine` -> `storage`). Em seu lugar, usamos
// injeção de dependência via closure/tipo função (`DragonGenerator`), que é instanciado na `main`.
//
// Retorna os atributos brutos: `(hp, atk, def, goldReward)`.
type DragonGenerator func(dayDate string) (hp, atk, def, goldReward int)

// VillageRepositoryOption define o tipo para aplicação do padrão Functional Options no repositório.
type VillageRepositoryOption func(*VillageRepository)

// WithDragonGenerator injeta a função geradora dos atributos do Dragão do Dia no repositório.
func WithDragonGenerator(gen DragonGenerator) VillageRepositoryOption {
	return func(r *VillageRepository) { r.dragonGen = gen }
}

// VillageRepository gerencia o estado global do vilarejo, o status diário do Dragão e os murais de notícias.
type VillageRepository struct {
	db        *DB
	dragonGen DragonGenerator
}

// NewVillageRepository cria e inicializa um novo repositório do vilarejo aplicando opções funcionais.
func NewVillageRepository(db *DB, opts ...VillageRepositoryOption) *VillageRepository {
	r := &VillageRepository{db: db}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// GetOrCreateTodayState busca ou inicializa atomicamente o estado do vilarejo e do Dragão do Dia no SQLite.
//
// Didática Go:
// 1. Abrimos uma transação explícita com `r.db.BeginTx(ctx, nil)` para garantir isolamento ACID.
// 2. Usamos `defer tx.Rollback()` para garantir o rollback automático se a função retornar antes de `tx.Commit()`.
// 3. Executamos `INSERT OR IGNORE` para criar a linha do dia corrente sem sobrescrever se já existir.
func (r *VillageRepository) GetOrCreateTodayState(ctx context.Context) (*VillageState, error) {
	today := time.Now().UTC().Format("2006-01-02")

	var maxHP, atk, def, goldReward int
	if r.dragonGen != nil {
		maxHP, atk, def, goldReward = r.dragonGen(today)
	} else {
		maxHP, atk, def, goldReward = 300, 45, 25, 3000
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	insertQuery := `
	INSERT OR IGNORE INTO village_state (day_date, dragon_alive, dragon_hp, dragon_max_hp, dragon_atk, dragon_def, dragon_gold_reward, slayer_name)
	VALUES (?, 1, ?, ?, ?, ?, ?, '')
	`
	if _, err := tx.ExecContext(ctx, insertQuery, today, maxHP, maxHP, atk, def, goldReward); err != nil {
		return nil, fmt.Errorf("falha ao assegurar estado do dragão do dia: %w", err)
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
		return nil, fmt.Errorf("falha ao ler estado do vilarejo: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("falha ao efetuar commit da transação: %w", err)
	}

	state.DragonAlive = dragonAliveInt == 1
	return &state, nil
}

// RecordDragonSlayed registra o abate do Dragão do Dia por um herói, atualizando o estado e inserindo um anúncio no jornal.
//
// Didática Go: A transação utiliza controle de concorrência otimista (`WHERE day_date = ? AND dragon_alive = 1`).
// Se `res.RowsAffected()` retornar 0, significa que outro jogador abateu o Dragão no mesmo dia milissegundos antes,
// abortando a operação de forma segura.
func (r *VillageRepository) RecordDragonSlayed(ctx context.Context, slayerName string) error {
	today := time.Now().UTC().Format("2006-01-02")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE village_state SET dragon_alive = 0, dragon_hp = 0, slayer_name = ?
	WHERE day_date = ? AND dragon_alive = 1
	`
	res, err := tx.ExecContext(ctx, query, slayerName, today)
	if err != nil {
		return fmt.Errorf("falha ao atualizar status do dragão: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("o dragão já foi derrotado hoje por outro herói")
	}

	news := fmt.Sprintf("GLÓRIA AO HERÓI! %s derrotou o Dragão e salvou o Vilarejo hoje!", slayerName)
	newsQuery := `INSERT INTO news (message, created_at) VALUES (?, CURRENT_TIMESTAMP)`
	if _, err := tx.ExecContext(ctx, newsQuery, news); err != nil {
		return fmt.Errorf("falha ao registrar notícia do feito: %w", err)
	}

	return tx.Commit()
}

// AddNews publica um novo anúncio ou comunicado no mural de notícias do vilarejo.
func (r *VillageRepository) AddNews(ctx context.Context, message string) error {
	query := `INSERT INTO news (message, created_at) VALUES (?, CURRENT_TIMESTAMP)`
	_, err := r.db.ExecContext(ctx, query, message)
	return err
}

// GetLatestNews recupera os comunicados mais recentes do mural de notícias ordenados do mais novo ao mais antigo.
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
