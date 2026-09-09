# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/),
e este projeto segue [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Não publicado]

### Adicionado
- Framework formal de registros de decisão de arquitetura em `design/adr/` com índice `README.md` e cinco registros iniciais seguindo o padrão `adr-template`:
  - **ADR-0001**: Adoção de Architecture Decision Records (ADRs) imutáveis.
  - **ADR-0002**: Arquitetura em camadas canônicas e identificadores 100% em inglês idiomático.
  - **ADR-0003**: Persistência com SQLite puro em Go (CGO-free), modo WAL e transações ACID.
  - **ADR-0004**: Terminal UI com Bubble Tea e servidor SSH multi-usuário com Wish.
  - **ADR-0005**: Camada centralizada de localização (i18n) em PT-BR.
- Suíte de testes unitários para a barra de status e estilos ANSI em `internal/ui`.
- Mecanismo de auto-save do estado do jogador no encerramento do programa e ao desconectar da sessão SSH (`b14be93`).
- Suporte a navegação Vim (`j`/`k`) e atalhos por abas na Ferraria (`ScreenSmith`) (`c6a6334`).
- Documento `docs/plano-fixes-auditoria.md` rastreando a resolução dos 27 achados do relatório de auditoria de código.

### Alterado
- Consolidação do design system ANSI em `internal/ui`: remoção de estilos duplicados em `internal/tui` e centralização da paleta e componentes.
- Atualização do `AGENTS.md` com referências canônicas a `design/adr/` e catálogo de skills recomendadas (`adr-template`, `godoctor`, `engineering-flow`, `git-workflow-and-versioning`, `game-design`, `latest-version`).
- Enriquecimento do `README.md` com atmosfera temática de RPG e nostalgia BBS dos anos 1980-90.
- Lógica de "Novo Dia" consolidada no `TurnManager` (`internal/engine`) como fonte única da verdade para renovação de turnos.
- Geração determinística do Dragão do Dia unificada via `bestiary.GenerateDragonOfDay` injetada na camada TUI.

### Corrigido
- **GameOver**: Corrigida a duplicação de penalidade de morte em `GameOverScreen` (`9b41902`).
- **Persistência de Poções**: Adicionada persistência da quantidade de poções (`PotionsCount`) no banco SQLite com migração de schema automática (`efe58ba`).
- **Segurança de Armazenamento**: Configurados `busy_timeout` e `foreign_keys` no DSN de conexão SQLite, transação em `RecordDragonSlayed` e tratamento adequado de erros de gravação (`e7765e4`).
- **Autenticação**: Diferenciação clara na tela de login entre senha inválida e conta inexistente (`e3c5e50`, `fe97252`).
- **Taverna**: Implementado limite diário para flerte com Cassandra e atualização das mensagens do Cavaleiro Vermelho (`d43db26`).

### Planejado
- Módulo de interface TUI completo para todas as telas do vilarejo (Fase 3 do roadmap).
- PvP assíncrono e mecânica do Cavaleiro Vermelho na taverna.
- Sistema de clãs e ranking persistente no banco compartilhado.
- Internacionalização: suporte a múltiplos idiomas adicionais além de PT-BR.

## [0.0.1] - 2026-08-31

### Adicionado
- Estrutura inicial do projeto `lotgd/`, módulo Go com `go.mod` e `go.sum` fixando as versões estáveis das dependências (`bubbletea v1.3.10`, `wish v1.4.7`, `lipgloss v1.1.0`, `modernc.org/sqlite v1.57.0`).
- Módulo `internal/i18n` com o dicionário canônico de localização em PT-BR, cobrindo menus, diálogos de NPCs, nomes de monstros, armas, armaduras e itens consumíveis.
- Módulo `internal/storage` implementando o driver SQLite puro (CGO-free via `modernc.org/sqlite`) com modo WAL, `busy_timeout` configurado, migrations automáticas e CRUD de contas de jogadores (`CreatePlayer`, `Authenticate`, `SaveProgress`).
- Módulo `internal/bestiary` com o gerador procedural de monstros organizado em 4 Tiers, sistema de afixos de prefixo (*Feroz*, *Covarde*, *Sortudo*, entre outros) e gerador dinâmico do **Dragão do Dia** baseado na data do servidor.
- Módulo `internal/engine` com modelos puros em Go (`Player`, `Monster`, `Weapon`, `Armor`, `Item`), motor de combate por turnos determinístico (`CombatEngine`) com cálculo de dano, chance de fuga e poções, sistema de economia (ouro na bolsa versus ouro no banco), progressão por experiência e gerenciamento dos 15 turnos diários.
- Módulo `internal/tui` com a máquina de estados principal em Bubble Tea (The Elm Architecture) e telas parciais do vilarejo, floresta, taverna, capela, ferreiro, guilda e covil do dragão.
- Módulo `internal/ui` com o design system ANSI para o tema retrô BBS via Lip Gloss, incluindo a paleta de cores e os componentes reutilizáveis.
- Dois binários de entrada: `cmd/lotgd` (jogo local single-player offline) e `cmd/server` (servidor SSH BBS multi-usuário via Charmbracelet Wish, com graceful shutdown).
- Documentação base: `README.md`, `AGENTS.md` (diretrizes para agentes de IA) e a pasta `docs/` contendo o Game Design Document, a descrição de universo narrativo, o documento de arquitetura e o roadmap de desenvolvimento.
- Suíte de testes unitários para os módulos `bestiary`, `engine`, `i18n`, `storage` e `tui` (todos verdes via `go test ./...`).
- Arquivo `CHANGELOG.md` (este arquivo), `CONTRIBUTING.md` com guia de contribuição e a estratégia oficial de versionamento e branches documentada no `README.md`.

### Alterado
- Formatação do código fonte padronizada via `gofmt` e `goimports` em 18 arquivos (apenas realinhamento mecânico de colunas em `struct` e agrupamento de imports, sem mudanças semânticas).
- `lotgd.db` removido do controle de versão: o banco de save local é regenerado dinamicamente e passa a ser listado no `.gitignore` para evitar commits acidentais.
- `.gitignore` expandido para cobrir binários Go, perfis de cobertura, metadados de IDE, metadados de sistema operacional, chave privada SSH do servidor Wish e artefatos auxiliares do SQLite WAL.

### Corrigido
- Atribuição inefetiva (`ineffectual assignment`) em `MonsterGenerator.GenerateForPlayer` (detectada por `golangci-lint`/`ineffassign`): a inicialização `tier := 1` era sobrescrita em todos os ramos do `switch` subsequente. Substituída por `var tier int`, preservando o comportamento original.

### Segurança
- A chave privada SSH do host (`.wish_host_key`) foi adicionada ao `.gitignore` para impedir versionamento acidental de material criptográfico.
