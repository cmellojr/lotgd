# Roadmap de Desenvolvimento — The Legend of the Go Dragon

> **Documento Estratégico de Fases & Entregas**  
> **Status:** Aprovado via Entrevista Interativa (`/grill-me`)  
> **Alinhamento:** [GDD.md](file:///c:/GitHub/go/lotgd/docs/GDD.md) · [universo-e-prompt.md](file:///c:/GitHub/go/lotgd/docs/universo-e-prompt.md) · [architecture.md](file:///c:/GitHub/go/lotgd/docs/architecture.md)

---

## 🎯 Decisões de Design Consolidadas

1. **Autenticação & Contas**: Usuário e senha armazenados no SQLite (`modernc.org/sqlite`), garantindo paridade total entre execução standalone local (`cmd/lotgd`) e servidor SSH (`cmd/server`).
2. **Mecânica de "Novo Dia"**: Verificação automática sob demanda (ao fazer login, se a data do servidor mudou, o sistema restaura os 15 turnos diários e regenera o Dragão do Dia com novas estatísticas).
3. **PvP & Rival Cavaleiro Vermelho**: Modelo clássico assíncrono de BBS — você pode desafiar clones/fantasmas de outros aventureiros e do Cavaleiro Vermelho na taverna.
4. **Navegação na Interface (TUI)**: Modo híbrido no **Bubble Tea** (navegação fluida por setas `[↑/↓/Enter]` combinada com teclas de atalho rápido estilo BBS `[F]loresta`, `[T]averna`, `[C]apela`, `[M]estre Torin`, etc.).

---

## 🗺️ Fases do Roadmap

### 📦 Fase 1 — Fundação, Localização & Persistência (Core Foundation)
- [x] Inicialização do `go.mod` em `lotgd/` com versões estáveis (`bubbletea v1.3.10`, `wish v1.4.7`, `lipgloss v1.1.0`, `sqlite v1.57.0`).
- [x] **Módulo `internal/i18n`**:
  - Dicionário `pt_br.go` com todos os textos canônicos de menus, diálogos de NPCs e nomes de monstros/itens.
- [x] **Módulo `internal/storage`**:
  - Driver SQLite puro (CGO-free) com modo WAL e migrations automáticas.
  - CRUD de contas de jogadores (`CreatePlayer`, `Authenticate`, `SaveProgress`).
  - Gerenciamento do estado do vilarejo e ciclo do Dragão diário.
- [x] Testes unitários para o repositório e migrations.

---

### ⚔️ Fase 2 — Domínio, Bestiário & Engine de RPG (Game Engine)
- [x] **Módulo `internal/bestiary`**:
  - Gerador procedural de monstros (Tiers 1 a 4) com afixos de prefixo (*"Feroz"*, *"Covarde"*, *"Sortudo"*).
  - Gerador dinâmico do **Dragão do Dia** baseado na data do servidor.
- [x] **Módulo `internal/engine`**:
  - Modelos puros em Go (`Player`, `Monster`, `Weapon`, `Armor`, `Item`).
  - Motor de combate por turnos determinístico (`CombatEngine`), com cálculo de dano, chance de fuga e poções.
  - Sistema de economia (ouro na bolsa vs ouro no banco), experiência e avanço de nível.
  - Gerenciamento de turnos diários (*15 Forest Fights*).
- [x] Testes unitários com cobertura para todas as regras de negócio de combate e progressão.

---

### 🎨 Fase 3 — Interface de Usuário Terminal (TUI & The Elm Architecture)
- [x] **Estilização com Lip Gloss (`internal/ui/ui.go`)**:
  - Paleta retrô ANSI (bordas, títulos temáticos, caixas de status do jogador e barras de HP).
- [x] **State Machine & Telas (`internal/tui/screens/`)**:
  - `login.go`: Tela de login / criação de aventureiro.
  - `town.go`: Praça principal do vilarejo com atalhos e seleção por setas.
  - `forest.go`: Exploração da Floresta Sombria e interface de combate por turnos.
  - `tavern.go`: Taverna da Dona Rosalinda (fofocas, flerte com Cassandra, duelos PvP assíncronos).
  - `chapel.go`: Capela do Frei Anselmo (cura de HP e bênçãos).
  - `smith.go`: Ferraria do Mestre Torin (compra e venda de armas e armaduras).
  - `guild.go`: Guilda & Biblioteca do Mestre Tobias (promoção de nível e lore).
  - `dragon.go`: Covil do Dragão (o grande confronto final do dia).
  - `game_over.go`: Tela de morte com penalidades e ressurreição na capela.

---

### 🚀 Fase 4 — Executáveis & Servidor SSH Multi-usuário
- [x] **Executável Local (`cmd/lotgd/main.go`)**:
  - Inicialização do jogo direto no terminal para jogar offline.
- [x] **Servidor SSH BBS (`cmd/server/main.go`)**:
  - Servidor Wish com suporte a conexões SSH remotas simultâneas em porta configurável (padrão: `2222`).
- [x] Documentação de execução e comandos no `README.md` de `lotgd/`.

---

## 📈 Critérios de Aceite e Validação Final

| Critério | Método de Verificação |
|---|---|
| **Zero Erros de Compilação & Lint** | `go build ./...` e `go vet ./...` sem advertências |
| **Sucesso nos Testes Unitários** | `go test ./internal/... -v` passando com 100% de sucesso |
| **Jogabilidade Completa (Single-player)** | Fluxo completo do novo personagem até a luta com o Dragão executável via `go run ./cmd/lotgd` |
| **Acesso Multi-usuário SSH** | Conexão de 2 ou mais terminais simultâneos via `ssh localhost -p 2222` sem conflitos no SQLite |
