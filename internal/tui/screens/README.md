# tui/screens

## Visão Geral
Componentes modulares de tela do Bubble Tea que implementam as localidades e interações do vilarejo, como a Praça do Vilarejo, Floresta Sombria, Taverna, Ferraria, Guilda de Treinamento, Capela, Covil do Dragão e Telas de Login/Game Over.

## Interface
- Telas do Vilarejo:
  - `type LoginScreen`: Autenticação e cadastro de contas de jogadores (`login.go`).
  - `type TownScreen`: Hub central de navegação no vilarejo (`town.go`).
  - `type ForestScreen`: Busca de monstros, combate por turnos, fuga e coleta de recompensas (`forest.go`).
  - `type TavernScreen`: Conversas com Dona Rosalinda, fofocas e operações bancárias (`tavern.go`).
  - `type SmithScreen`: Compra de armas, armaduras e sistema de troca com desconto (`smith.go`).
  - `type GuildScreen`: Desafio de combate com mestres para promoção de nível (`guild.go`).
  - `type ChapelScreen`: Serviços de cura e bênção de ressurreição com Frei Anselmo (`chapel.go`).
  - `type DragonScreen`: O grande confronto final contra o Dragão e ciclo de prestígio (`dragon.go`).
  - `type GameOverScreen`: Notificação de derrota, epitáfio e retorno ao vilarejo (`game_over.go`).
- Métodos do Ciclo de Vida:
  - `func New[Screen]Screen(...) *[Screen]Screen`
  - `func (s *[Screen]Screen) Init() tea.Cmd`
  - `func (s *[Screen]Screen) SetPlayer(p *engine.Player)`
  - `func (s *[Screen]Screen) SetSize(width, height int)`
  - `func (s *[Screen]Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd)`
  - `func (s *[Screen]Screen) View() string`
- Utilitários de Persistência (`save.go`):
  - `func SavePlayer(repo *storage.PlayerRepository, p *engine.Player) tea.Cmd`: Dispara um comando assíncrono do Bubble Tea para salvar o estado do jogador no banco.

## Invariantes
- Todas as telas implementam fielmente a interface `tea.Model` do Bubble Tea.
- Transições de tela são emitidas exclusivamente através de mensagens `ui.ChangeScreenMsg`.
- Modificações no estado do jogador disparam persistência assíncrona via `SavePlayer` ou `ui.PlayerUpdatedMsg`.

## Efeitos Colaterais
- Emite comandos do Bubble Tea que executam operações de I/O e escrita no banco de dados.

## Verificação
- Executar: `go test ./internal/tui/screens/...`
