# Guia de Contribuição

Obrigado por considerar contribuir com **The Legend of the Go Dragon (LOTGD)**! Este documento descreve como preparar o ambiente, o estilo de código esperado, o fluxo de contribuição via *pull request* e as convenções que mantemos para garantir a saúde do projeto a longo prazo.

O projeto é distribuído sob a licença **GNU General Public License v3.0 (GPL-3.0)** — toda contribuição herda automaticamente essa mesma licença. Para mais detalhes, consulte o arquivo [`LICENSE`](LICENSE).

---

## 📜 Índice

1. [Código de Conduta](#-código-de-conduta)
2. [Preparando o Ambiente](#-preparando-o-ambiente)
3. [Estratégia de Branches](#-estratégia-de-branches)
4. [Convencional Commits e DCO](#-convencional-commits-e-dco)
5. [Estilo de Código](#-estilo-de-código)
6. [Testes e Verificações Locais](#-testes-e-verificações-locais)
7. [Abrindo um Pull Request](#-abrindo-um-pull-request)
8. [Reportando Bugs e Sugerindo Funcionalidades](#-reportando-bugs-e-sugerindo-funcionalidades)

---

## 🤝 Código de Conduta

Esperamos que toda interação seja respeitosa, técnica e construtiva. Discutimos ideias, não pessoas. Críticas devem ser fundamentadas em evidências (código, documentação, métricas) e acompanhadas de propostas de melhoria.

---

## 🛠️ Preparando o Ambiente

### Pré-requisitos

- **Go 1.22+** (recomendamos a versão 1.25 listada no `go.mod`).
- Um terminal com suporte a cores ANSI (Windows Terminal, iTerm2, Alacritty, Kitty, GNOME Terminal).
- Cliente SSH padrão (`ssh`) caso queira testar o servidor BBS.
- Ferramentas opcionais para o gate de qualidade completo:
  - `gofmt` e `goimports` (já incluídos na distribuição do Go e em `golang.org/x/tools`).
  - `golangci-lint` — agregador oficial de linters para Go.

### Clonando o Repositório

```bash
git clone https://github.com/cmellojr/lotgd.git
cd lotgd
go mod download
```

### Rodando Localmente

```bash
# Jogo single-player offline
go run ./cmd/lotgd

# Servidor SSH BBS multi-usuário (em outro terminal: ssh localhost -p 2222)
go run ./cmd/server
```

---

## 🌿 Estratégia de Branches

O projeto adota o **GitFlow simplificado** com duas branches de longa duração e branches efêmeras por mudança:

| Branch | Função | Política de merge |
|---|---|---|
| `main` | Recebe apenas releases finalizadas; sempre estável. | Merge de `release/*` ou `hotfix/*` via PR. Cada merge em `main` gera uma tag SemVer. |
| `develop` | Integração contínua de funcionalidades. | PRs de `feature/*`, `fix/*`, `chore/*` entram aqui. |
| `feature/<escopo>-<descrição>` | Desenvolvimento de uma funcionalidade nova. | Curta (1 a 3 dias). Merge via PR contra `develop`. |
| `fix/<escopo>-<descrição>` | Correção de bug. | Curta (1 a 3 dias). Merge via PR contra `develop` (ou contra `main` se for hotfix). |
| `chore/<escopo>-<descrição>` | Tarefas de manutenção (dependências, CI, docs). | Curta (1 a 3 dias). Merge via PR contra `develop`. |

### Convenção de Nomes

- Use kebab-case em ASCII: `feature/cavaleiro-vermelho`, `fix/poison-potion-overflow`.
- Mantenha os nomes curtos e descritivos; o *escopo* (ex.: `combat`, `tui`, `storage`) é opcional, mas ajuda no `git log --oneline`.

### Fluxo Típico

```bash
# 1. Garanta que está em develop e atualize
git checkout develop
git pull --ff-only origin develop

# 2. Crie sua branch a partir de develop
git checkout -b feature/minha-funcionalidade

# 3. Trabalhe em commits pequenos e atômicos
# (...)

# 4. Antes de abrir o PR, sincronize com develop
git fetch origin
git rebase origin/develop

# 5. Empururre e abra o PR contra develop
git push -u origin feature/minha-funcionalidade
```

> **Regra de ouro:** PRs vão contra `develop`, nunca contra `main` (a menos que seja um *hotfix* urgente).

---

## 📝 Convencional Commits e DCO

### Mensagens de Commit

Seguimos [Conventional Commits](https://www.conventionalcommits.org/pt-br/) em português ou inglês (preferimos português para coerência com a documentação):

```text
<tipo>(<escopo opcional>): <descrição curta em até 72 caracteres>

<corpo opcional: parágrafos adicionais explicando o "porquê", não o "o quê">

<rodapé opcional: referências a issues, breaking changes>
```

**Tipos permitidos:**

| Tipo | Quando usar |
|---|---|
| `feat` | Nova funcionalidade visível ao usuário. |
| `fix` | Correção de bug. |
| `refactor` | Mudança de código que não corrige bug nem adiciona feature. |
| `perf` | Otimização de desempenho. |
| `test` | Adição ou ajuste de testes. |
| `docs` | Apenas documentação (incluindo este arquivo). |
| `chore` | Tarefas de manutenção (dependências, CI, build, `.gitignore`). |
| `style` | Ajustes puramente estéticos (formatação, espaçamento). |
| `build` | Mudanças no sistema de build ou em dependências externas. |
| `ci` | Mudanças em pipelines de CI. |

### Assinatura DCO (Developer Certificate of Origin)

Para garantir compatibilidade legal com a GPL-3.0, **todo commit deve incluir o sign-off DCO** (`Signed-off-by:` no rodapé). Isso certifica que você tem o direito de submeter a contribuição sob a licença do projeto.

Adicione o sign-off automaticamente com:

```bash
git commit -s -m "feat(combat): adicionar mecânica de contra-ataque"
```

Para configurar seu nome e e-mail globalmente (uma única vez):

```bash
git config --global user.name "Seu Nome"
git config --global user.email "seu.email@exemplo.com"
```

---

## 🎨 Estilo de Código

### Princípios Gerais

- Todo o projeto segue o [**Google Go Style Guide**](https://google.github.io/styleguide/go/) e as diretrizes do `AGENTS.md` na raiz.
- **Identificadores em inglês** (variáveis, funções, tipos, constantes, pacotes, métodos) — o inglês é a língua franca do ecossistema Go.
- **Comentários e documentação em português do Brasil** (PT-BR), para manter coerência com a interface do jogo e com a documentação geral.
- **Sem comentários supérfluos**: código idiomático se explica; comente apenas a *intenção* quando ela não é óbvia pelo próprio código.

### Formatação e Imports

```bash
# Formatação canônica (deve passar sem listar arquivos)
gofmt -l .

# Organização de imports (stdlib, third-party, internal)
goimports -l .
```

Ambos os comandos devem retornar saída vazia antes de qualquer commit.

### Padrões Recomendados

- Retorne erros como último valor e envolva-os com contexto:

  ```go
  if err := db.Ping(); err != nil {
      return fmt.Errorf("storage: falha ao validar conexão: %w", err)
  }
  ```

- Prefira `var x T` quando o valor zero for semanticamente significativo; caso contrário, use `:=` com valor explícito.
- Não use `panic` em bibliotecas ou lógica de negócio — retorne erros explícitos.
- Evite nomes stutter (`storage.Storage` → `storage.Repository` ou similar).
- Use `context.Context` em funções que façam I/O ou que possam ser canceladas.
- Prefira composição a herança; mantenha interfaces pequenas e definidas pelo consumidor (*accept interfaces, return structs*).

---

## 🧪 Testes e Verificações Locais

Antes de enviar o PR, rode a cadeia completa de verificações:

```bash
# Análise estática
go vet ./...

# Formatação
gofmt -l .
goimports -l .

# Suíte de testes
go test ./... -v

# Linter agregado (recomendado)
golangci-lint run ./... --timeout=10m
```

> Os comandos `gofmt` e `goimports` **devem** retornar saída vazia. O `go vet` e o `go test` **devem** terminar com `exit 0`. O `golangci-lint` pode emitir sugestões, mas nenhum warning pode permanecer sem justificativa.

Escreva testes para qualquer mudança de comportamento. Cobertura é um sinal de saúde, não um objetivo em si — foque em testar o que pode quebrar.

---

## 🔀 Abrindo um Pull Request

1. **Garanta que sua branch está atualizada** com `develop`:

   ```bash
   git fetch origin
   git rebase origin/develop
   ```

2. **Confirme que os testes passam** localmente (ver seção anterior).

3. **Abra o PR no GitHub** descrevendo:
   - O que foi alterado (1 a 3 frases).
   - Por que a mudança é necessária.
   - Como testar manualmente (passos reproduzíveis).
   - Capturas de tela, se a mudança afetar a interface TUI.

4. **Referencie issues** usando palavras-chave como `Closes #123` ou `Refs #456` no corpo do PR.

5. **Aguarde a revisão**: pedimos ao menos 1 aprovação para merges em `main` e CI verde para merges em `develop`.

6. **Squash ou merge?** Preferimos *merge commit* para preservar o histórico de branches de feature, exceto para commits triviais de um único autor.

---

## 🐛 Reportando Bugs e Sugerindo Funcionalidades

Use as **Issues** do GitHub com os templates apropriados:

- **Bug report:** passos para reproduzir, comportamento esperado, comportamento observado, ambiente (SO, versão do Go, terminal).
- **Feature request:** problema que a feature resolve, proposta de solução, alternativas consideradas.

Antes de abrir uma issue, procure por issues similares já abertas para evitar duplicação.

---

Obrigado por ajudar a manter o vilarejo de LOTGD vivo! 🐉
