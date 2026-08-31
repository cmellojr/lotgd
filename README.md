# 🐉 The Legend of the Go Dragon (LOTGD)

<p align="center">
  <img src="docs/assets/lotgd_logo.jpg" alt="The Legend of the Go Dragon" width="700"/>
</p>

Um RPG clássico em modo texto retrô (BBS / Roguelike), construído 100% em **Go** com arquitetura em camadas, interface terminal com **Bubble Tea** & **Lip Gloss**, persistência ACID em **SQLite** (CGO-free) e suporte a servidor **SSH multi-usuário** via **Wish**.

---

## 🎯 Sumário

1. [Requisitos](#-requisitos)
2. [Estrutura do Projeto](#-estrutura-do-projeto)
3. [Como Jogar Localmente (CLI Standalone)](#-como-jogar-localmente-cli-standalone)
4. [Como Executar e Conectar ao Servidor SSH BBS](#-como-executar-e-conectar-ao-servidor-ssh-bbs)
5. [Controles e Navegação](#-controles-e-navega%C3%A7%C3%A3o)
6. [Executando Testes](#-executando-testes)
7. [Arquitetura & Engenharia](#-arquitetura--engenharia)

---

## 📦 Requisitos

- **Go 1.22+** (desenvolvido e testado em Go 1.25)
- Terminal com suporte a cores ANSI (Windows Terminal, iTerm2, Alacritty, Kitty, Bash, etc.)
- Cliente SSH padrão (`ssh`) para conexões remotas ao servidor BBS.

---

## 🗺️ Estrutura do Projeto

```text
lotgd/
├── cmd/
│   ├── lotgd/                 # Executável CLI standalone (single-player offline)
│   │   └── main.go
│   └── server/                # Servidor SSH BBS multi-usuário (Charmbracelet Wish)
│       └── main.go
├── docs/                      # Documentação de design, arquitetura e roadmap
│   ├── GDD.md
│   ├── architecture.md
│   ├── roadmap.md
│   └── universo-e-prompt.md
├── internal/
│   ├── bestiary/              # Bestiário procedural (Tiers 1..4, Afixos e Dragão Diário)
│   ├── engine/                # Regras de negócio, combate por turnos e progressão de RPG
│   ├── i18n/                  # Dicionário de localização PT-BR
│   ├── storage/               # Camada de banco de dados SQLite puro (modernc.org/sqlite)
│   ├── tui/                   # Modelos e máquina de estados Bubble Tea (The Elm Architecture)
│   │   └── screens/           # Telas do vilarejo, floresta, taverna, capela, ferreiro, etc.
│   └── ui/                    # Design system ANSI e componentes de renderização (Lip Gloss)
├── AGENTS.md                  # Diretrizes pedagógicas e guia de estilo de código
├── go.mod
├── go.sum
└── README.md
```

---

## 🎮 Como Jogar Localmente (CLI Standalone)

Para iniciar o jogo diretamente no seu terminal:

```bash
# A partir da pasta lotgd/
go run ./cmd/lotgd
```

### Flags Disponíveis

| Flag | Padrão | Descrição |
|---|---|---|
| `-db` | `lotgd.db` | Caminho personalizado para o arquivo de banco de dados SQLite. |

Exemplo com banco customizado:
```bash
go run ./cmd/lotgd -db=meu_save.db
```

---

## 🌐 Como Executar e Conectar ao Servidor SSH BBS

O servidor permite que múltiplos jogadores se conectem simultaneamente pela rede usando qualquer cliente SSH padrão, compartilhando o mesmo mundo persistente do vilarejo.

### 1. Iniciar o Servidor

```bash
# A partir da pasta lotgd/
go run ./cmd/server
```

### Flags do Servidor

| Flag | Padrão | Descrição |
|---|---|---|
| `-host` | `0.0.0.0` | Endereço IP de escuta da interface de rede. |
| `-port` | `2222` | Porta TCP do serviço SSH BBS. |
| `-host-key` | `.wish_host_key` | Caminho para armazenar/ler a chave privada do host SSH. |
| `-db` | `lotgd.db` | Caminho para o banco de dados compartilhado. |

Exemplo de inicialização customizada:
```bash
go run ./cmd/server -host=127.0.0.1 -port=2222 -db=lotgd.db
```

### 2. Conectar via SSH

Em outro terminal (ou de outra máquina na rede):

```bash
ssh localhost -p 2222
```

> **Nota:** Não é necessária chave pública SSH de cliente para autenticação inicial; o login e a criação do aventureiro são realizados diretamente na interface TUI do jogo.

---

## 🕹️ Controles e Navegação

A interface utiliza um sistema híbrido intuitivo:

- **Navegação por Listas:** Use as setas `[↑]` e `[↓]` (ou `k`/`j`) para selecionar opções e pressione `[Enter]`.
- **Atalhos Rápidos de Menu (Estilo BBS Clássico):** Pressione a tecla destacada entre colchetes para ação instantânea:
  - `[F]` — Floresta Sombria
  - `[T]` — Taverna da Dona Rosalinda
  - `[C]` — Capela do Frei Anselmo
  - `[M]` — Ferraria do Mestre Torin
  - `[G]` — Guilda do Mestre Tobias
  - `[D]` — Covil do Dragão
  - `[S]` — Salvar e Sair
- **Combate na Floresta:**
  - `[A]` — Atacar o monstro
  - `[P]` — Tomar poção de cura
  - `[F]` — Tentar fugir para o vilarejo
- **Sair do Jogo:** `[Ctrl+C]` salva o progresso e finaliza a sessão.

---

## 🧪 Executando Testes

Para rodar toda a suíte de testes unitários do motor de RPG, banco de dados e geradores:

```bash
go test ./... -v
```

Para verificar análise estática e padronização:

```bash
go vet ./...
```

---

## 🏛️ Arquitetura & Engenharia

- **The Elm Architecture (TEA):** Cada tela é um componente desacoplado que implementa `Init()`, `Update(msg)` e `View()`.
- **Persistência Concorrente:** SQLite configurado com `PRAGMA journal_mode=WAL;` e `PRAGMA busy_timeout=5000;`, permitindo acesso multi-thread e multi-sessão seguro.
- **Isolamento de Efeitos Colaterais:** O motor de combate (`CombatEngine`) opera como uma função de transição de estado pura e testável.
- **Didática de Código:** Conforme estabelecido no [AGENTS.md](file:///c:/GitHub/go/lotgd/AGENTS.md), todo o código-fonte contém comentários conceituais profundos para estudantes e desenvolvedores Go.
