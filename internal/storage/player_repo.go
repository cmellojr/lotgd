package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrUserExists     = errors.New("username already registered")
	ErrInvalidPass    = errors.New("invalid password")
)

// PlayerRepository gerencia todas as operações de persistência e consulta SQL relativas aos jogadores.
type PlayerRepository struct {
	db *DB
}

// NewPlayerRepository instancia o repositório de jogadores.
func NewPlayerRepository(db *DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// Register cria uma nova conta de jogador com a senha criptografada usando Bcrypt.
//
// Didática Go: Criptografamos a senha com `bcrypt.GenerateFromPassword` antes da gravação no banco,
// garantindo que senhas puras em texto simples jamais sejam armazenadas no disco.
func (r *PlayerRepository) Register(ctx context.Context, username, password string) (*Player, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("falha ao criptografar senha: %w", err)
	}

	today := time.Now().UTC().Format("2006-01-02")
	query := `
	INSERT INTO players (username, password_hash, last_login_day, created_at, updated_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	res, err := r.db.ExecContext(ctx, query, username, string(hash), today)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("falha ao registrar jogador: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

// Authenticate valida as credenciais do jogador comparando o hash Bcrypt armazenado com a senha fornecida.
func (r *PlayerRepository) Authenticate(ctx context.Context, username, password string) (*Player, error) {
	player, err := r.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPass
	}

	return player, nil
}

// GetByID busca um jogador no banco através de seu identificador numérico único (`ID`).
func (r *PlayerRepository) GetByID(ctx context.Context, id int64) (*Player, error) {
	query := `
	SELECT id, username, password_hash, level, experience, gold, bank_gold, health, max_health,
	       attack, defense, weapon_id, armor_id, potions_count, forest_fights, dragon_kills, last_login_day,
	       created_at, updated_at
	FROM players WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanPlayer(row)
}

// GetByUsername busca um jogador no banco pelo seu nome de usuário (case-insensitive via NOCASE na tabela).
func (r *PlayerRepository) GetByUsername(ctx context.Context, username string) (*Player, error) {
	query := `
	SELECT id, username, password_hash, level, experience, gold, bank_gold, health, max_health,
	       attack, defense, weapon_id, armor_id, potions_count, forest_fights, dragon_kills, last_login_day,
	       created_at, updated_at
	FROM players WHERE username = ?
	`
	row := r.db.QueryRowContext(ctx, query, username)
	return scanPlayer(row)
}

// Save atualiza os atributos e progresso do jogador no banco de dados SQLite.
func (r *PlayerRepository) Save(ctx context.Context, p *Player) error {
	query := `
	UPDATE players SET
		level = ?, experience = ?, gold = ?, bank_gold = ?, health = ?, max_health = ?,
		attack = ?, defense = ?, weapon_id = ?, armor_id = ?, potions_count = ?, forest_fights = ?,
		dragon_kills = ?, last_login_day = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		p.Level, p.Experience, p.Gold, p.BankGold, p.Health, p.MaxHealth,
		p.Attack, p.Defense, p.WeaponID, p.ArmorID, p.PotionsCount, p.ForestFights,
		p.DragonKills, p.LastLoginDay, p.ID,
	)
	return err
}

// ListRankings retorna os melhores jogadores ordenados por abates de dragão, nível e experiência acumulada.
func (r *PlayerRepository) ListRankings(ctx context.Context, limit int) ([]*Player, error) {
	query := `
	SELECT id, username, password_hash, level, experience, gold, bank_gold, health, max_health,
	       attack, defense, weapon_id, armor_id, potions_count, forest_fights, dragon_kills, last_login_day,
	       created_at, updated_at
	FROM players
	ORDER BY dragon_kills DESC, level DESC, experience DESC
	LIMIT ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*Player
	for rows.Next() {
		p, err := scanPlayer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

// rowScanner abstrai os métodos Scan do `*sql.Row` e `*sql.Rows`.
//
// Didática Go: A definição de uma interface local com um único método `Scan(dest ...any) error`
// permite reutillizar a mesma função utilitária `scanPlayer` tanto para consultas de registro único (`QueryRow`)
// quanto para iterações em conjunto de linhas (`Query`).
type rowScanner interface {
	Scan(dest ...any) error
}

func scanPlayer(scanner rowScanner) (*Player, error) {
	var p Player
	err := scanner.Scan(
		&p.ID, &p.Username, &p.PasswordHash, &p.Level, &p.Experience, &p.Gold, &p.BankGold,
		&p.Health, &p.MaxHealth, &p.Attack, &p.Defense, &p.WeaponID, &p.ArmorID, &p.PotionsCount,
		&p.ForestFights, &p.DragonKills, &p.LastLoginDay, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return &p, nil
}
