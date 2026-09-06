# Guia de Agentes de IA — The Legend of the Go Dragon

Diretrizes de desenvolvimento e contexto de atuação para agentes de IA no projeto.

---

## 🧭 Referências Canônicas

Antes de gerar código ou conteúdo, consulte os documentos especializados:

- **Universo e Lore**: [`docs/universo-e-lore.md`](docs/universo-e-lore.md) — Tom de voz, NPCs, bestiário e nomenclatura.
- **Game Design Document**: [`docs/GDD.md`](docs/GDD.md) — Core loop, combate, economia de turnos e fluxo TUI.
- **Arquitetura & Engenharia**: [`docs/architecture.md`](docs/architecture.md) — Camadas, SQLite (CGO-free), Bubble Tea e Wish (SSH).
- **Decisões de Arquitetura (ADRs)**: [`design/adr/`](design/adr/) — Registro imutável de decisões arquiteturais do projeto.

---

## 🛠️ Diretrizes Principais

1. **Código em Inglês**: Structs, métodos, funções, enums, variáveis e pacotes devem ser 100% em **inglês idiomático**.
2. **Interface em Português (PT-BR)**: Textos e mensagens para o jogador devem ser gerenciados pela camada de i18n (`internal/i18n`), sem strings de exibição hardcoded na lógica de domínio ou TUI.
3. **Estilo Go Idiomático**: Siga as recomendações do [Google Go Style Guide](https://google.github.io/styleguide/go/). Mantenha a documentação Go padrão (`godoc`) para tipos e funções exportadas.

---

## 🤖 Regras de Automação com `gh` CLI

Ao criar ou editar Pull Requests e Issues via `gh` CLI no PowerShell, **NUNCA use `--body` com texto inline**. O PowerShell faz double-escaping de backticks (`` ` ``), corrompendo a formatação Markdown no GitHub.

**Sempre escreva o corpo em um arquivo temporário e use `--body-file`:**

```powershell
# ✅ Correto — escrever corpo em arquivo, depois usar --body-file
$body = @"
## Descrição
Texto com `backticks` e **formatação** funcionando corretamente.
"@
$body | Out-File -FilePath "$env:TEMP\pr_body.md" -Encoding utf8
gh pr create --title "título" --body-file "$env:TEMP\pr_body.md"

# ❌ Errado — --body inline causa double-escaping de backticks
gh pr create --title "título" --body "Texto com `backticks`"
```

Isso se aplica a: `gh pr create`, `gh pr edit`, `gh issue create`, `gh issue edit`.

---

## 🧰 Skills & Tooling Recomendados

Para manter a consistência, a ergonomia e a qualidade de engenharia, os agentes devem consultar e seguir as diretrizes das seguintes skills quando disponíveis no ambiente ou via URL canônica:

- **[`adr-template`](https://skills.danicat.dev/standards/adr-template/SKILL.md)**: Padrão imutável para documentação de decisões técnicas e trade-offs arquiteturais.
- **[`godoctor`](https://skills.danicat.dev/coding/godoctor/SKILL.md)**: Boas práticas de Go idiomático, AST integrity, testes automatizados e compilação limpa via `go vet` e linters.
- **[`engineering-flow`](https://skills.danicat.dev/coding/engineering-flow/SKILL.md)**: Decisões técnicas fundamentadas, política Zero-Debt em `0.x`, proibição de silenciamento de erros (`_ = err`) e higiene de código (*Broken Window Code Hygiene*).
- **[`git-workflow-and-versioning`](https://github.com/addyosmani/agent-skills/blob/main/skills/git-workflow-and-versioning/SKILL.md)**: Commits atômicos, mensagens no padrão Conventional Commits, PRs curtas e branches de vida curta.
- **[`game-design`](https://skills.danicat.dev/game-dev/game-design/SKILL.md)**: Diretrizes de mecânica de jogo, balanceamento de atributos, progressão e loops de gameplay.
- **[`latest-version`](https://skills.danicat.dev/coding/latest-version/SKILL.md)**: Verificação de versões estáveis de dependências em manifestos (`go.mod`) sem adivinhação.

