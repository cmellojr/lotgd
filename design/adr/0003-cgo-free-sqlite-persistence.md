# ADR-0003: Use CGO-Free SQLite and WAL Mode for Persistence

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon

## 1. Context
Como jogo multiusuário persistente estilo BBS door game, múltiplos aventureiros acessam simultaneamente os mesmos dados: ranking diário, estado do vilarejo, dados do herói e o status do Dragão do Dia.

O projeto precisa rodar em diversos ambientes (servidores Linux, desktops Windows e macOS), de preferência gerando um único executável estático sem exigir ferramentas de compilação C (`gcc`/`clang`), preservando a ergonomia do Go cross-compilation (`GOOS=linux go build`).

Ao mesmo tempo, acessos concorrentes em SQLite tradicional podem sofrer contenção de travas de arquivo (`database is locked`) se a estratégia de concorrência não for desenhada corretamente.

## 2. Decision
Adotamos o driver SQLite puro em Go `modernc.org/sqlite` (sem CGO) com o seguinte modelo operacional:

1. **Driver 100% Go**: Elimina qualquer dependência de CGO, permitindo compilação cruzada trivial para qualquer arquitetura ou SO.
2. **Modo WAL (Write-Ahead Logging)**: Habilitado para permitir leituras concorrentes simultâneas com escrita sem bloqueio de leitura.
3. **Parâmetros de Conexão no DSN**: Configuração de `busy_timeout` de 5000ms e integridade referencial via `PRAGMA foreign_keys = ON` embutidos no DSN para assegurar que todas as conexões do pool recebam os pragmas.
4. **Transações ACID Explícitas**: Operações críticas como abate do Dragão (`RecordDragonSlayed`) e compras na ferraria devem usar transações atômicas com `tx.Commit()` / `tx.Rollback()`.

Alternativas consideradas:
- **`mattn/go-sqlite3`**: Rejeitado porque requer CGO ativo, exigindo toolchains C instaladas no host e impedindo cross-compile simples no GitHub Actions e máquinas Windows/Linux heterogêneas.
- **Bancos Cliente-Servidor (PostgreSQL/MySQL)**: Rejeitado por violar o princípio de simplicidade auto-contida dos clássicos BBS games, introduzindo dependência operacional pesada desnecessária para a escala do jogo.
- **Arquivos Flat JSON/YAML**: Rejeitado por falta de garantias ACID, alto risco de corrupção em crashes e incapacidade de lidar com concorrência segura entre conexões SSH simultâneas.

## 3. Consequences
### Positivas
- Binários 100% estáticos gerados com `CGO_ENABLED=0`.
- Confiabilidade em escritas concorrentes por meio de transações e WAL mode.
- Facilidade de backup: banco reside em um único arquivo `.db`.

### Negativas / Restrições
- `modernc.org/sqlite` possui desempenho computacional em consultas massivas ligeiramente inferior à compilação C nativa (diferença desprezível para a carga de RPG textual).
- Requer gerenciamento cuidadoso do pool de conexões do `database/sql` para não esgotar descritores de arquivos em execuções prolongadas.

## 4. Compliance and Verification
- Suíte de testes unitários e de integração (`internal/storage/...`) executada sob concorrência com `-race`.
- Validação no CI de compilação sem CGO (`CGO_ENABLED=0 go build ./...`).
- Verificação do schema versioning via `PRAGMA user_version` nas rotinas de inicialização do banco.
