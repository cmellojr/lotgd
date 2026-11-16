# storage

## Visão Geral
Camada de persistência responsável pelos esquemas SQLite, migrações automáticas de banco de dados, credenciais de jogadores e histórico do estado do vilarejo utilizando o driver CGO-free `modernc.org/sqlite`.

## Interface
- Conexão e Banco (`db.go`):
  - `type DB`: Invólucro de conexão com o banco de dados.
  - `func OpenDB(path string) (*DB, error)`: Abre a conexão SQLite, ativa modo WAL e chaves estrangeiras, e executa migrações automáticas.
- Entidades de Dados (`models.go`):
  - `type Player`: Registro persistido dos atributos do jogador, equipamentos, saldo bancário, turnos diários e histórico de vitórias contra o Dragão.
  - `type VillageState`: Registro persistido do estado diário do vilarejo (dragão vivo/derrotado, mestre do dia e data).
  - `type NewsEntry`: Registro histórico de eventos públicos do vilarejo (vitórias contra o dragão, mortes de heróis).
- Repositório de Jogadores (`player_repo.go`):
  - `type PlayerRepository`: Interface de repositório para contas de jogadores.
  - `func NewPlayerRepository(db *sql.DB) *PlayerRepository`
  - `func (r *PlayerRepository) Register(username, password string) (*Player, error)`: Aplica hash seguro bcrypt na senha e registra o jogador.
  - `func (r *PlayerRepository) Authenticate(username, password string) (*Player, error)`: Valida credenciais via comparação segura de hash bcrypt.
  - `func (r *PlayerRepository) GetByID(id int64) (*Player, error)`
  - `func (r *PlayerRepository) GetByUsername(username string) (*Player, error)`
  - `func (r *PlayerRepository) Save(p *Player) error`
  - `func (r *PlayerRepository) ListRankings() ([]*Player, error)`
- Repositório do Vilarejo (`village_repo.go`):
  - `type VillageRepository`: Repositório do estado comunitário e notícias do vilarejo.
  - `func NewVillageRepository(db *sql.DB, opts ...VillageRepositoryOption) *VillageRepository`
  - `func (r *VillageRepository) GetOrCreateTodayState(dateStr string) (*VillageState, error)`
  - `func (r *VillageRepository) RecordDragonSlayed(dateStr string, slayerID int64, slayerName string) error`
  - `func (r *VillageRepository) AddNews(message string) error`
  - `func (r *VillageRepository) GetLatestNews(limit int) ([]NewsEntry, error)`

## Invariantes
- Senhas nunca são salvas em texto puro; o uso de hash bcrypt é estritamente obrigatório.
- Modo WAL (Write-Ahead Logging) e timeout de espera configurados para garantir concorrência segura entre múltiplos processos/conexões.
- As operações no banco mantêm integridade referencial e transacional.

## Efeitos Colaterais
- Leituras e escritas nos arquivos de banco SQLite no disco local (`.db`).
- Modifica as tabelas: `players`, `village_states`, `news`.

## Verificação
- Executar: `go test ./internal/storage/...`
