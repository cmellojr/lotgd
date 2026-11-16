# engine

## Visão Geral
Motor central das regras de jogo, implementando a resolução de combates por turno, regras de economia bancária, ciclos de renovação de dias, progressão e teste de mestres, catálogo de itens e mecânica de prestígio do Dragão.

## Interface
- Modelos de Domínio (`models.go`):
  - `type Player`: Estado em memória do jogador com atributos, equipamentos, ouro na mão e no banco, poções, turnos e abates de dragão.
  - `func (p *Player) TotalAttack() int`
  - `func (p *Player) TotalDefense() int`
  - `func (p *Player) IsAlive() bool`
  - `func NewPlayerFromStorage(sp *storage.Player) *Player`
  - `func (p *Player) ToStorage() *storage.Player`
  - `type Monster`: Estado do monstro em combate com ataque, defesa, pontos de vida, ouro e recompensa de XP.
  - `func (m *Monster) IsAlive() bool`
- Motor de Combate (`combat.go`):
  - `type CombatEngine`: Motor sem estado para resolução de batalhas.
  - `func NewCombatEngine() *CombatEngine`
  - `type TurnResult`: Snapshot do resultado de uma rodada de combate (dano, fuga, acerto crítico).
  - `func (ce *CombatEngine) CalculateDamage(attack, defense int) int`
  - `func (ce *CombatEngine) Attack(player *Player, monster *Monster) TurnResult`
  - `func (ce *CombatEngine) AttemptFlee(player *Player, monster *Monster) TurnResult`
  - `func (ce *CombatEngine) UsePotion(player *Player) TurnResult`
- Serviço de Economia (`economy.go`):
  - `type EconomyService`: Regras do banco e penalidade de morte.
  - `func NewEconomyService() *EconomyService`
  - `func (e *EconomyService) Deposit(p *Player, amount int) error`
  - `func (e *EconomyService) Withdraw(p *Player, amount int) error`
  - `func (e *EconomyService) ProcessDeathPenalty(p *Player) (lostGold int, lostExp int)`
- Gerenciamento de Turnos (`turns.go`):
  - `const DailyForestFights = 10`
  - `type TurnManager`: Controle do ciclo do "Novo Dia" e consumo de lutas na floresta.
  - `func NewTurnManager() *TurnManager`
  - `func CurrentDateString() string`
  - `func (tm *TurnManager) CheckAndApplyNewDay(p *Player, currentDate string) bool`
  - `func (tm *TurnManager) ConsumeFight(p *Player) error`
- Progressão e Treinamento (`progression.go`):
  - `type Master`: Mestres da guilda para o desafio de passagem de nível.
  - `var LevelTable`: Tabela de requisitos de XP e ganhos de atributos por nível.
  - `var MasterCatalog`: Lista dos mestres de cada faixa de nível.
  - `func CanLevelUp(p *Player) bool`
  - `func LevelUp(p *Player)`
- Itens e Equipamentos (`items.go`):
  - `type Item`: Definições de armas, armaduras e consumíveis.
  - `var WeaponsCatalog`, `var ArmorsCatalog`, `var PotionsCatalog`
  - `func FindWeapon(id string) (Item, bool)`
  - `func FindArmor(id string) (Item, bool)`
  - `func FindPotion(id string) (Item, bool)`
  - `func CalculateTradeInQuote(currentVal, newVal int) int`

## Invariantes
- Ouro na mão e ouro no banco nunca podem assumir valores negativos.
- Turnos de floresta não podem ser consumidos quando `ForestFightsLeft <= 0`.
- A morte do personagem zera o ouro em mãos e reduz uma fração do XP atual, mas nunca afeta o ouro guardado no banco nem regride o nível do jogador.
- Derrotar o Dragão ativa a mecânica de prestígio, resetando o nível para 1 enquanto acumula bônus permanentes por dragão abatido.

## Efeitos Colaterais
- Altera os atributos e o inventário dos modelos em memória.

## Verificação
- Executar: `go test ./internal/engine/...`
