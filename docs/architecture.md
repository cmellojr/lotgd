# Architecture & Engineering Document — The Legend of the Go Dragon

> **Documento Técnico de Arquitetura de Software**  
> **Referência Narrativa:** [universo-e-prompt.md](file:///c:/GitHub/go/lotgd/docs/universo-e-prompt.md)  
> **Documento de Game Design:** [GDD.md](file:///c:/GitHub/go/lotgd/docs/GDD.md)  
> **Padrões de Código e Ensino:** [AGENTS.md](file:///c:/GitHub/go/lotgd/AGENTS.md)

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
- **Transações ACID**: Proteção contra corrupção em concorrência multi-sessão SSH.
- **Schema**:
  - `players`: Contas, senha hasheada, stats, ouro, nível, turnos restantes.
  - `village_state`: Data do dia atual, status do Dragão, ranking diário.
  - `graveyard`: Histórico de heróis caídos.
  - `news`: Fofocas e anúncios de vitórias no vilarejo.

---

## 3. Concorrência & Servidor SSH BBS (`cmd/server`)

- Utiliza o framework **Wish** (`github.com/charmbracelet/wish`) sobre o protocolo SSH padrão.
- Cada conexão SSH autenticada instancia uma sessão isolada de `tea.Program`, compartilhando o mesmo pool de banco de dados SQLite com controle de concorrência (`WAL mode`).
- Suporte a desconexão limpa e auto-save de progresso do aventureiro.
