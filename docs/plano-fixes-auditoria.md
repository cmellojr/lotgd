# Plano de Implementação — Fixes do Relatório de Auditoria v0.0.1

**Data**: 2026-08-31
**Estado**: Ativo
**Repositório**: `develop` branch

Este documento define o plano de correção dos 27 achados restantes do relatório de auditoria de August 2026, organizado por fases atômicas com critérios de aceite verificáveis.

---

## Premissas

1. Cada fase resulta num commit atômico no `develop`
2. Fases anteriores não quebram as posteriores (cada uma é independente e verificável)
3. `go test ./...`, `go vet`, `gofmt -l`, `goimports -l` passam no final de cada fase
4. Código em inglês; tudo exibido ao jogador em PT-BR
5. Sem dependências externas novas

---

## Fase 1 — C1: Penalidade de morte duplicada (Crítico)

**Problema**: `syncPlayerState()` chama `GameOverScreen.SetPlayer()` em cada mudança de tela; `SetPlayer` executa `ProcessDeathPenalty()` + `SavePlayer`. O construtor também o faz. Navegar entre telas re-aplica a penalidade.

**Arquivos**:
- `internal/tui/screens/game_over.go`
- `internal/tui/tui.go`
- `internal/tui/tui_test.go`

**Plano**:

1. `GameOverScreen.SetPlayer()` → transformar num setter puro: `s.player = p; s.lostGold = 0; s.lostXP = 0`. Sem `ProcessDeathPenalty`, sem `Save`.
2. `NewGameOverScreen()` → aceitar `lostGold`, `lostXP` como parâmetros em vez de calcular internamente.
3. `tui.go` → quando a tela de destino é `ScreenGameOver`, calcular `ProcessDeathPenalty` no `MainModel.Update()` antes de `syncPlayerState()`, gravar uma única vez, e passar os valores ao construtor.
4. `syncPlayerState()` → não chamar `SetPlayer()` no `GameOverScreen` (ele já tem os dados do construtor).
5. Teste: simular transição para `ScreenGameOver` → verificar penalidade exata uma única vez; navegar após morte não altera Gold/XP.

**Critério de aceite**: `ProcessDeathPenalty` é chamado exatamente 1 vez por morte; navegar entre telas após morte não altera Gold/XP.

---

## Fase 2 — G5: Derrota contornável (Grave)

**Estado**: já corrigido no `develop`. `dragonStateDefeat` e `forestStateDefeat` redirecionam imediatamente para `ScreenGameOver`.

**Ação**: verificação de regressão. Nenhuma alteração de código necessária.

---

## Fase 3 — G1+M7: Poções + migração de schema (Grave + Moderado)

**Problema**: `NewPlayerFromStorage` fixa `PotionsCount: 0`; `Save` não inclui a coluna; schema não tem a coluna; sem mecanismo de migração para DBs existentes.

**Arquivos**:
- `internal/storage/db.go`
- `internal/storage/models.go`
- `internal/engine/models.go`

**Plano**:

1. **Mecanismo de migração** (`db.go`): antes do `CREATE TABLE IF NOT EXISTS`, ler `PRAGMA user_version`. Se `< 2`, executar:
   ```sql
   ALTER TABLE players ADD COLUMN potions_count INTEGER NOT NULL DEFAULT 0;
   PRAGMA user_version = 2;
   ```
2. **Schema** (`db.go`): adicionar `potions_count INTEGER NOT NULL DEFAULT 0` ao `CREATE TABLE players`.
3. **`storage.Player`** (`models.go`): adicionar `PotionsCount int`.
4. **`scanPlayer`** e **`Save`** (`player_repo.go`): incluir `potions_count` no SELECT e no UPDATE.
5. **`NewPlayerFromStorage`** (`engine/models.go`): `PotionsCount: sp.PotionsCount` em vez de `0`.
6. **`ToStorage`** (`engine/models.go`): incluir `PotionsCount: p.PotionsCount`.

**Critério de aceite**: comprar poções → gravar → re-login → poções mantêm-se; DB existente (v0.0.1) recebe a coluna automaticamente.

---

## Fase 4 — C2+G2+G3: Persistência robusta (Crítico + 2 Graves)

**Problema**: PRAGMAs aplicados a uma única conexão do pool; erros de `SavePlayer` descartados em ~36 locais; zero transacções; `RecordDragonSlayed` sem `AND dragon_alive = 1`; `GetOrCreateTodayState` read-then-insert sem proteção.

**Arquivos**:
- `internal/storage/db.go`
- `internal/storage/village_repo.go`
- `internal/tui/tui.go`
- Todas as telas com `_ = s.db.SavePlayer(...)` (~15 arquivos)

**Plano**:

1. **PRAGMAs no DSN** (`db.go`): alterar o DSN para incluir parâmetros de pragma:
   ```
   file:lotgd.db?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)
   ```
   O `journal_mode=WAL` mantém-se via `db.Exec` (persiste no arquivo). Isto garante que todas as ligações do pool recebem os PRAGMAs.

2. **Transação em `RecordDragonSlayed`** (`village_repo.go`):
   ```go
   func (r *VillageRepository) RecordDragonSlayed(ctx context.Context, slayerName string) error {
       tx, err := r.db.BeginTx(ctx, nil)
       if err != nil { return err }
       defer tx.Rollback()

       today := time.Now().Format("2006-01-02")
       res, err := tx.ExecContext(ctx,
           "UPDATE village_state SET dragon_alive=0, dragon_hp=0, slayer_name=? WHERE day_date=? AND dragon_alive=1",
           slayerName, today)
       if err != nil { return err }

       rows, _ := res.RowsAffected()
       if rows == 0 {
           return fmt.Errorf("dragon already slain or not found for %s", today)
       }

       news := fmt.Sprintf("GLÓRIA AO HERÓI! %s derrotou o Dragão e salvou o Vilarejo hoje!", slayerName)
       _, err = tx.ExecContext(ctx, "INSERT INTO news (message, created_at) VALUES (?, CURRENT_TIMESTAMP)", news)
       if err != nil { return err }

       return tx.Commit()
   }
   ```

3. **`GetOrCreateTodayState`** (`village_repo.go`): usar transação com `INSERT OR IGNORE` + SELECT:
   ```go
   tx, _ := r.db.BeginTx(ctx, nil)
   defer tx.Rollback()
   tx.ExecContext(ctx, "INSERT OR IGNORE INTO village_state ...", ...)
   row := tx.QueryRowContext(ctx, "SELECT ... FROM village_state WHERE day_date=?", today)
   // scan e retornar
   return &state, tx.Commit()
   ```

4. **Erros de SavePlayer**: criar função helper no pacote `tui`:
   ```go
   func savePlayer(db *storage.DB, p *engine.Player) {
       if err := db.SavePlayer(p.ToStorage()); err != nil {
           log.Printf("WARN: failed to save player %s: %v", p.Username, err)
       }
   }
   ```
   Substituir `_ = s.db.SavePlayer(...)` em todas as telas por `savePlayer(s.db, s.player)`.

5. **`SetMaxOpenConns(1)`** (`db.go`): limitar o pool a 1 conexão dado o perfil de carga (BBS de terminal).

**Critério de aceite**: 2 sessões SSH simultâneas; uma matando o dragão → a outra vê erro controlado; `busy_timeout` funciona em todas as ligações; erros de Save são logados.

---

## Fase 5 — G6: Erros de login (Grave)

**Problema**: `login.go:136-145` trata toda a falha de `Authenticate` como conta inexistente e tenta `CreatePlayer`; `ErrUserExists`/`ErrInvalidPass` são código morto no consumo.

**Arquivos**:
- `internal/tui/screens/login.go`
- `internal/storage/player_repo.go`

**Plano**:

1. **`Register`** (`player_repo.go`): mapear a violação UNIQUE para `ErrUserExists`:
   ```go
   if strings.Contains(err.Error(), "UNIQUE constraint failed") {
       return nil, ErrUserExists
   }
   ```

2. **`login.go`**: distinguir os três cenários:
   ```go
   sp, err := s.db.AuthenticatePlayer(user, pass)
   if err != nil {
       switch {
       case errors.Is(err, storage.ErrInvalidPass):
           s.errMsg = "Senha incorreta para o aventureiro."
           return s, nil
       case errors.Is(err, storage.ErrPlayerNotFound):
           // Criar conta nova
           newSP, createErr := s.db.CreatePlayer(user, pass)
           if createErr != nil {
               if errors.Is(createErr, storage.ErrUserExists) {
                   s.errMsg = "Este nome de aventureiro já está registrado."
               } else {
                   s.errMsg = fmt.Sprintf("Erro ao criar conta: %v", createErr)
               }
               return s, nil
           }
           sp = newSP
       default:
           s.errMsg = fmt.Sprintf("Erro ao autenticar: %v", err)
           return s, nil
       }
   }
   ```

**Critério de aceite**: senha errada → "Senha incorreta"; nome novo → cria conta; nome existente → "já está registrado".

---

## Fase 6 — G7: Auto-save na desconexão SSH (Grave)

**Problema**: `cmd/server/main.go` não registra hook de fim de sessão; o Wish faz `program.Quit()`/`Kill()` que nunca gera a tecla `ctrl+c` da única gravação de saída.

**Arquivos**:
- `cmd/server/main.go`
- `internal/tui/tui.go`

**Plano**:

1. **`MainModel.Save() error`** público em `tui.go`:
   ```go
   func (m *MainModel) Save() error {
       if m.player == nil || m.db == nil {
           return nil
       }
       return m.db.SavePlayer(m.player.ToStorage())
   }
   ```

2. **`cmd/server/main.go`**: usar `tea.WithFilter` para gravar antes de `tea.Quit`:
   ```go
   model := tui.NewMainModel(db, dragonGen)
   filter := func(msg tea.Msg) (tea.Msg, tea.Cmd) {
       if _, ok := msg.(tea.QuitMsg); ok {
           model.Save()
       }
       return msg, nil
   }
   return model, []tea.ProgramOption{
       tea.WithAltScreen(),
       tea.WithFilter(filter),
   }
   ```

**Critério de aceite**: sessão SSH morta com SIGKILL → a DB refletia o estado da última ação do jogador.

---

## Fase 7 — G9+G11: Ferraria e navegação (2 Graves)

**Problema**: `[1][2][3]` não funcionam na ferraria; `k`/`j` não funciona em telas além da praça; `k` na taverna dispara duelo.

**Arquivos**:
- `internal/tui/screens/smith.go`
- `internal/tui/screens/tavern.go`
- `internal/tui/screens/guild.go`
- `internal/tui/screens/chapel.go`
- `internal/tui/screens/dragon.go`

**Plano**:

1. **Ferraria** (`smith.go`): tratar `"1"`, `"2"`, `"3"` como atalhos de tab:
   ```go
   case "1":
       s.tab = smithTabWeapons; s.cursor = 0; return s, nil
   case "2":
       s.tab = smithTabArmors; s.cursor = 0; return s, nil
   case "3":
       s.tab = smithTabPotions; s.cursor = 0; return s, nil
   ```

2. **Todas as telas** (`tavern.go`, `guild.go`, `smith.go`, `chapel.go`, `dragon.go`): adicionar `case "j":` e `case "k":` antes da normalização `strings.ToUpper`. Adotar o padrão de `town.go`:
   ```go
   // Antes de k := strings.ToUpper(msg.String())
   switch msg.String() {
   case "up", "k":
       // navegar para cima
   case "down", "j":
       // navegar para baixo
   }
   ```

3. **Taverna** (`tavern.go`): remover `"K"` do `case "D", "K":` para que `k` só funcione como navegação.

**Critério de aceite**: `k`/`j` navegam em todas as telas; `[1][2][3]` trocam tabs na ferraria; `k` na taverna navega, não dispara duelo.

---

## Fase 8 — G8+G10+M3: PvP, economia e turnos (Graves/Moderados)

**Problema**: PvP é fachada; economia desequilibrada; flerte vende turnos ilimitados.

**Arquivos**:
- `internal/tui/screens/tavern.go`

**Plano**:

1. **Flerte** (`tavern.go:119-120`): limitar a 1 turno extra por dia:
   ```go
   if s.player.ForestFights > engine.DailyForestFights {
       s.infoMsg = "Cassandra já te inspirou hoje. Volte amanhã para nova dose de coragem!"
   } else if s.player.Gold >= 10 {
       s.player.Gold -= 10
       s.player.ForestFights++
       // ...
   }
   ```

2. **Cavaleiro Vermelho** (`tavern.go:127`): melhorar mensagem:
   ```go
   s.infoMsg = "O Cavaleiro Vermelho ergue a viseira: 'Você ainda não possui a tempera necessária para cruzar lâminas comigo. Volte quando estiver mais experiente!'"
   ```

3. **Economia**: decisão de design pendente — subir recompensas de ouro (~20-30%) ou baixar custos de treino (~20%). Requer aprovação do autor.

**Critério de aceite**: flerte dá no máximo 1 turno extra por dia; mensagem do Cavaleiro é coerente.

---

## Fase 9 — Documentação e limpeza (Moderados + Menores)

### Documentação

| ID | Ação | Arquivo |
|---|---|---|
| M1 | Documentar crítico 10% no GDD §2.2 | `docs/GDD.md` |
| M2 | Alinhar GDD §2.2 com fuga real (50% fixo, 80% Covarde, 20% Dragão) | `docs/GDD.md` |
| M5 | Desmarcar PvP/Hall da Fama no roadmap como "não implementado" | `docs/roadmap.md` |
| M8 | Atualizar rodapés de todas as telas com atalhos reais | `internal/tui/screens/*.go` |
| M9 | Alinhar praça com código: "Requer nível 5" em vez de "máximo" | `internal/tui/screens/town.go` |

### Limpeza de código

| ID | Ação | Arquivo |
|---|---|---|
| m1 | Remover tabela `graveyard` do schema (ou implementar API) | `internal/storage/db.go` |
| m2 | Atualizar `UpdatedAt` no struct após `Save` | `internal/engine/models.go` |
| m3 | Alinhar fusos horários: tudo em UTC ou tudo em local | `internal/storage/*.go` |
| m4 | Remover `internal/tui/styles.go` (cópia morta de `ui/ui.go`) | `internal/tui/styles.go` |
| m5 | Remover "Ponteiro Nulo" do catálogo de armas | `internal/engine/items.go` |
| m6 | Remover `NPCYolanda` e `MonsterPrefixesPTBR` (não usados) | `internal/i18n/pt_br.go`, `internal/i18n/types.go` |
| m7 | Corrigir `go mod tidy` e `gofmt -l .` | projeto inteiro |
| m8 | Atualizar links em AGENTS.md (remover caminhos absolutos Windows) | `AGENTS.md` |

---

## Ordem de execução e dependências

```
Fase 1 (C1)  ──┐
                ├──→ Fase 4 (C2+G2+G3) ──→ Fase 6 (G7)
Fase 3 (G1)  ──┘
Fase 5 (G6)     ──→ independente
Fase 7 (G9+G11) ──→ independente
Fase 8 (G10)    ──→ precisa decisão de design
Fase 9          ──→ última, documentação e limpeza
```

Fases 1, 3, 5 e 7 são independentes e podem ser feitas em paralelo.
Fase 4 depende de Fase 1 e Fase 3 porque altera os mesmos arquivos de persistência.
Fase 6 depende de Fase 4 (já que grava na DB).

---

## Estimativas de esforço

| Fase | Arquivos | Linhas estimadas |
|---|---|---|
| Fase 1 (C1) | 3 | ~40 |
| Fase 3 (G1+M7) | 3 | ~30 |
| Fase 4 (C2+G2+G3) | ~15 | ~110 |
| Fase 5 (G6) | 2 | ~20 |
| Fase 6 (G7) | 2 | ~15 |
| Fase 7 (G9+G11) | 5 | ~25 |
| Fase 8 (G10) | 2 | ~10 |
| Fase 9 | ~10 | ~50 |
| **Total** | **~40** | **~300** |

---

## Achados já corrigidos no develop (referência)

| ID | Achado | Commit |
|---|---|---|
| G4 | Dragão do Dia unificado via `bestiary.GenerateDragonOfDay` | `e17e54a`, `8222adf`, `ec35654`, `7eb2c7b` |
| M4 | Novo Dia consolidado no `TurnManager` | `25dd09d` |
| G5 | Derrota contornável fechada | `dragon.go:137-141`, `forest.go:113-117` |

---

## Achados que permanecem abertos (27)

| Severidade | IDs |
|---|---|
| Crítico | C1, C2 |
| Grave | G1, G2, G3, G6, G7, G8, G9, G10, G11 |
| Moderado | M1, M2, M3, M5, M6, M7, M8, M9 |
| Menor | m1, m2, m3, m4, m5, m6, m7, m8 |
