# ADR-0004: Build Terminal UI with Bubble Tea and Multi-User SSH Server with Wish

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon

## 1. Context
O *The Legend of the Go Dragon* revive a experiência nostálgica dos jogos de BBS door games das décadas de 1980 e 1990. Para proporcionar uma experiência moderna, o jogo requer interfaces ricas no terminal com suporte a cores ANSI, menus responsivos, atalhos de teclado rápidos e formatação visual elegante.

Adicionalmente, os jogadores devem ser capazes de se conectar remotamente ao servidor central via SSH nativo sem precisar de clientes customizados além de um emulador de terminal com cliente SSH (`ssh jogador@servidor -p 2222`), além da possibilidade de jogar offline/localmente via terminal interativo.

## 2. Decision
Adotamos o ecossistema da Charm (`charmbracelet`):
1. **Bubble Tea (`github.com/charmbracelet/bubbletea`)**: Framework TUI baseado em The Elm Architecture (TEA). Todo o fluxo de navegação e combate é gerenciado por uma máquina de estados finitos (`Model`, `Update`, `View`).
2. **Lip Gloss (`github.com/charmbracelet/lipgloss`)**: Utilizado para estilização visual ANSI declarativa, alinhamento e molduras.
3. **Wish (`github.com/charmbracelet/wish`)**: Servidor SSH modular em Go. Cada conexão SSH estabelecida autentica o jogador e inicializa uma sessão isolada de `tea.Program`.
4. **Desacoplamento Entrada Local vs Remota**:
   - `cmd/lotgd`: Inicia o `tea.NewProgram` conectando a `os.Stdin` e `os.Stdout`.
   - `cmd/server`: Inicia o listener SSH que delega o PTY da sessão remota para uma nova instância do mesmo modelo TUI.

Alternativas consideradas:
- **Telnet clássico**: Rejeitado por falta de criptografia, autenticação segura e gerenciamento nativo moderno de PTY.
- **Tview / Ccell**: Rejeitado por ter modelo de atualização imperativo e menor flexibilidade composicional quando comparado à Elm Architecture reativa do Bubble Tea.
- **Interface Web (HTML/Wasm)**: Rejeitado para o loop principal do jogo, visto que a identidade canônica do projeto é estritamente um jogo de terminal/BBS.

## 3. Consequences
### Positivas
- Código TUI 100% compartilhado entre cliente local e servidor SSH.
- Tratamento nativo e robusto de redimensionamento de janela (`tea.WindowSizeMsg`) e sequências ANSI.
- Autenticação e transporte criptografado seguro via SSH sem esforço adicional de infraestrutura.

### Negativas / Restrições
- Dependência do ciclo de mensagens assíncronas do Bubble Tea; operações lentas de banco de dados ou rede devem ser despachadas como comandos assíncronos (`tea.Cmd`) para não congelar o render da interface.
- Terminais antigos ou mal configurados sem suporte a UTF-8 ou ANSI completo podem apresentar problemas visuais.

## 4. Compliance and Verification
- Telas do TUI devem implementar a interface `tea.Model` e possuir testes de unidade para processamento de mensagens em `internal/tui/...`.
- Testes de concorrência com conexões simultâneas no servidor Wish.
- O linter e os testes automatizados devem rodar com `go test ./internal/tui/...`.
