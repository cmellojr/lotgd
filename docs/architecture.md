# Architecture & Engineering Document — The Legend of the Go Dragon

> **Documento Técnico de Arquitetura de Software**  
> **Referência Narrativa:** [universo-e-lore.md](universo-e-lore.md)
> **Documento de Game Design:** [GDD.md](GDD.md)
> **Diretrizes para Agentes:** [AGENTS.md](../AGENTS.md)

---

## 1. Visão Geral da Arquitetura

O **The Legend of the Go Dragon** é construído utilizando uma arquitetura modular em camadas, seguindo o padrão canônico da comunidade Go (`cmd/`, `internal/`):

```text
┌────────────────────────────────────────────────────────┐
│                   Entradas (Ponto de Acesso)           │
│         cmd/lotgd (CLI Local)  │  cmd/server (SSH BBS) │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│          Apresentação / TUI (internal/tui)             │
│        Bubble Tea (State Machine) + Lip Gloss (ANSI)   │
└───────────────────────────┬────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────┐
│            Domínio & Regras (internal/engine)          │
│      Player, CombatEngine, Items, Inventory, Turns     │
└──────────────┬─────────────────────────┬───────────────┘
               │                         │
┌──────────────▼──────────────┐   ┌──────▼───────────────┐
│ Bestiário (internal/bestiary)│   │ Localização (internal│
│  Tiers 1..4, Afixos, Dragão  │   │  /i18n): pt_br.go    │
└──────────────┬──────────────┘   └──────────────────────┘
               │
┌──────────────▼─────────────────────────────────────────┐
│        Persistência / Dados (internal/storage)         │
│          modernc.org/sqlite (CGO-free puro Go)         │
└────────────────────────────────────────────────────────┘
```

---

## 2. Padrões por Camada

### 2.1 Camada de Apresentação (`internal/tui`)
- **The Elm Architecture (TEA)**:
  - `Model`: Estado imutável da cena atual.
  - `Update`: Processamento determinístico de mensagens (`tea.KeyMsg`, `tea.WindowSizeMsg`, etc.).
  - `View`: Renderização ANSI pura com **Lip Gloss**.
- **Gerenciador de Cenas (Screen Router)**:
  - Enum `ScreenID` para rotear entre Praça, Floresta, Taverna, Capela, Ferraria, Guilda e Covil do Dragão.

### 2.2 Camada de Domínio (`internal/engine`)
- **Separação Rigorosa de Idioma**: Identificadores de código (structs, methods, enums, fields) estritamente em **inglês**.
- **Isolamento de Efeitos Colaterais**: A `CombatEngine` recebe instâncias de atacante e defensor e executa turnos de forma determinística/testável, sem dependência de I/O ou banco.

### 2.3 Camada de Bestiário (`internal/bestiary`)
- Gerador dinâmico de monstros por tier de dificuldade.
- Sistema de afixos estocásticos para variedade de atributos e nomes compostos.
- Geração do **Dragão do Dia** com seed baseada na data do servidor (`YYYY-MM-DD`).

### 2.4 Camada de Localização (`internal/i18n`)
- Mapeia identificadores internos para textos legíveis e formatados em português do Brasil (`PT-BR`).
- Garante desacoplamento para futura expansão de idiomas se necessário.

### 2.5 Camada de Persistência (`internal/storage`)
- **Driver SQLite CGO-free** (`modernc.org/sqlite`): Portabilidade sem dependência de GCC/Clang no host.
- **Transações ACID & Concorrência**:
  - **Write-Ahead Logging (WAL mode)**: Ativado via `PRAGMA journal_mode=WAL`, permitindo leituras simultâneas em paralelo sem bloquear leitores nem ser bloqueado por leituras durante escritas.
  - **Busy Timeout**: Configurado via DSN (`busy_timeout(5000)`), garantindo que transações de escrita aguardem resiliamente até 5 segundos para adquirir o lock de escrita se outra escrita estiver ativa.
  - **Pool de Conexões Calibrado**: Configuramos `SetMaxOpenConns(10)` e `SetMaxIdleConns(10)`. No Go SQL pool em modo WAL, conexões abertas > 1 removem o gargalo de enfileiramento em leituras (`SELECT`), permitindo que múltiplas sessões leiam concorrentemente em paralelo.
- **Schema**:
  - `players`: Contas, senha hasheada, stats, ouro, nível, turnos restantes.
  - `village_state`: Data do dia atual, status do Dragão, ranking diário.
  - `news`: Fofocas e anúncios de vitórias no vilarejo.

---

## 3. Concorrência & Servidor SSH BBS (`cmd/server`)

- Utiliza o framework **Wish** (`github.com/charmbracelet/wish`) sobre o protocolo SSH padrão.
- Cada conexão SSH autenticada instancia uma sessão isolada de `tea.Program`, compartilhando a mesma instância de banco de dados SQLite (`*storage.DB`).
- **Garantias de Concorrência Multi-Sessão**:
  - **Isolamento de Estado de Usuário**: Cada sessão TUI mantém seu estado local em memória e sincroniza com o banco via repositórios otimizados.
  - **Prevenção de Condições de Corrida**: Operações críticas de escrita (como registro do abate do Dragão do Dia) utilizam controle de concorrência otimista (`UPDATE ... WHERE day_date = ? AND dragon_alive = 1`) dentro de transações explícitas (`BeginTx`), garantindo que apenas um herói receba os prêmios do dia em caso de abates simultâneos.
  - **Leituras Sem Gargalo**: Graças ao modo WAL e ao pool calibrado (`SetMaxOpenConns(10)`), operações frequentes de leitura (como atualizações de ranking, mural de notícias e consultas de perfil) são atendidas em paralelo sem encavalamento.
- Suporte a desconexão limpa e auto-save de progresso do aventureiro via interceptação de `tea.QuitMsg`.
