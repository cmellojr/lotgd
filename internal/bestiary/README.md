# bestiary

## Visão Geral
Responsável pela geração procedural de monstros comuns da floresta e pelo encontro diário com o Dragão do Dia. Combina arquétipos canônicos de monstros com prefixos e sufixos de afixos estatísticos balanceados com base no nível do jogador.

## Interface
- `type MonsterTemplate`: Estrutura definindo atributos base, tier de dificuldade, recompensas de experiência e ouro (`types.go`).
- `type AffixModifier`: Estrutura definindo modificadores de status, prefixos e sufixos para variação procedural (`types.go`).
- `var CanonicalTemplates`: Lista canônica de monstros base divididos por faixas de dificuldade (`types.go`).
- `var AvailableAffixes`: Conjunto de modificadores aplicados dinamicamente aos modelos base (`types.go`).
- `type MonsterGenerator`: Gerador e fábrica procedural de monstros (`generator.go`).
- `func NewMonsterGenerator(seed int64) *MonsterGenerator`: Cria e instancia um novo gerador de monstros.
- `func (mg *MonsterGenerator) GenerateForPlayer(playerLevel int) engine.Monster`: Gera um monstro adequado ao nível do aventureiro.
- `func (mg *MonsterGenerator) GenerateByTier(tier int) engine.Monster`: Gera um monstro diretamente pelo nível de dificuldade/tier.
- `func GenerateDragonOfDay(dateStr string) engine.Monster`: Gera de forma determinística o chefe Dragão do Dia com base na data do calendário (`dragon.go`).

## Invariantes
- A geração determinística do Dragão produz exatamente os mesmos atributos e título para todos os aventureiros no mesmo dia civil.
- Os cálculos de geração de monstros por tier nunca geram valores negativos de ataque ou defesa.

## Efeitos Colaterais
- Utiliza geradores de números pseudo-aleatórios (`math/rand`) para seleção procedural de moldes e afixos.

## Verificação
- Executar: `go test ./internal/bestiary/...`
