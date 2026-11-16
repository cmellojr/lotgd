# tui

## Visão Geral
Orquestrador da interface de usuário em modo texto construída com Charm Bubble Tea (`charmbracelet/bubbletea`), responsável pelo roteamento de telas, captura de atalhos globais de teclado, estado da sessão e persistência de dados.

## Interface
- Modelo Principal (`tui.go`):
  - `type MainModel`: Modelo raiz do Bubble Tea que gerencia as telas ativas e a pilha de navegação.
  - `func NewMainModel(db *storage.DB, width, height int) *MainModel`: Inicializa o modelo principal com todas as telas configuradas.
  - `func (m *MainModel) Init() tea.Cmd`: Inicializa a tela ativa.
  - `func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)`: Processa eventos de navegação, redimensionamento de janela e entradas do teclado.
  - `func (m *MainModel) View() string`: Renderiza a tela ativa combinada com a barra de status persistente.
  - `func (m *MainModel) Save() tea.Cmd`: Executa a persistência do estado atual do jogador.
- Subdiretórios:
  - [screens](./screens/README.md) — Telas individuais das localidades do jogo (Vilarejo, Floresta, Taverna, Ferraria, Guilda, Capela, Dragão, Login).

## Invariantes
- Apenas uma tela permanece ativa e visível por vez.
- Eventos de redimensionamento (`tea.WindowSizeMsg`) são propagados para todas as telas registradas.
- Atalhos globais de saída gravam com segurança o estado do jogador antes de encerrar o programa.

## Efeitos Colaterais
- Emite comandos do Bubble Tea para persistência em disco (`Save`).
- Controla o buffer de exibição do terminal.

## Verificação
- Executar: `go test ./internal/tui/...`
