# Avaliação Técnica e Proposta de Solução — Issue 42

**Issue:** `#42 - refactor(i18n): textos de interface fixos no código (cerca de 230)`
**Data:** 12 de Setembro de 2026
**Autor:** Jules (Agente de IA)
**Projeto:** *The Legend of the Go Dragon (LOTGD)*

---

## 1. Contexto e Objetivos

Esta avaliação técnica analisa o problema relatado na **Issue 42**, referente à presença de aproximadamente 230 strings de interface e mensagens ao jogador codificadas diretamente (*hardcoded*) no código-fonte Go, espalhadas entre a camada de apresentação (`internal/tui`), ajudantes visuais (`internal/ui`), regras de domínio (`internal/engine`) e geradores procedurais (`internal/bestiary`).

O objetivo deste documento é:
1. Realizar o diagnóstico quantitativo e qualitativo das ocorrências.
2. Confrontar o estado atual com as decisões registradas nas **ADRs** do projeto.
3. Avaliar os riscos sob a perspectiva das diretrizes do `AGENTS.md` e das *skills* técnicas de engenharia (`godoctor`, `engineering-flow` e `git-workflow-and-versioning`).
4. Propor a arquitetura de solução, o design das APIs do pacote `internal/i18n` e um plano de migração gradual com commits atômicos.

---

## 2. Conformidade com as ADRs e Diretrizes do Projeto

### 2.1. ADR-0005: Centralize Portuguese (PT-BR) Localization in Dedicated i18n Layer
- **Decisão Estabelecida:** Nenhuma string voltada ao jogador pode ser mantida de forma *hardcoded* nas camadas de domínio (`internal/engine`), dados (`internal/storage`) ou apresentação (`internal/tui`). A TUI e o motor devem consumir mensagens exclusivamente a partir do pacote `internal/i18n`.
- **Diagnóstico:** A presença de mais de 200 literais de texto em português espalhados por telas e regras de combate **viola diretamente a ADR-0005**. A falta de centralização dificulta revisões editoriais do tom de voz e impede a futura internacionalização (ex: suporte a `en_us.go`).

### 2.2. ADR-0002: Adopt Layered Clean Architecture and English Identifiers
- **Decisão Estabelecida:** Todos os identificadores em Go (structs, funções, variáveis, enums, constantes) devem utilizar **inglês idiomático**, enquanto o texto final exibido ao jogador fica retido nos dicionários do `internal/i18n`.
- **Diagnóstico:** Campos como `NamePTBR` em `bestiary.AffixModifier` misturam a língua portuguesa com identificadores de struct. Além disso, métodos de domínio em `internal/engine` geram strings formatadas em português (ex: `"VITÓRIA LENDÁRIA!..."`), acoplando a linguagem natural à lógica de combate.

### 2.3. AGENTS.md
- **Diretriz 2 (Interface em PT-BR):** "Textos e mensagens para o jogador devem ser gerenciados pela camada de i18n (`internal/i18n`), sem strings de exibição hardcoded na lógica de domínio ou TUI."

---

## 3. Inventário Detalhado das Strings Hardcoded

A varredura estática no código identificou aproximadamente **230 ocorrências** distribuídas da seguinte forma:

| Camada / Pacote | Arquivos Afetados | Qtd. Aprox. | Tipos de Strings Identificadas |
| :--- | :--- | :--- | :--- |
| **`internal/tui/screens`** | `town.go`, `forest.go`, `smith.go`, `tavern.go`, `chapel.go`, `guild.go`, `dragon.go`, `game_over.go`, `login.go` | ~160 | Rótulos de menu, descrições de opções, diálogos de NPCs, rodapés de ajuda (`[Enter] Confirmar`), avisos do sistema (`infoMsg`), botões e banners. |
| **`internal/ui`** | `ui.go` | ~20 | Rótulos do cabeçalho/barra de status (`"Herói:"`, `"Nível:"`, `"HP:"`, `"Ouro:"`, `"Banco:"`, `"Lutas Diárias:"`, `"Arma:"`, `"Armadura:"`, `"Poções:"`, `"ATK:"`, `"DEF:"`, `"XP:"`). |
| **`internal/engine`** | `items.go`, `combat.go`, `progression.go`, `economy.go` | ~35 | Descrições dos catálogos de armas/armaduras/poções, mensagens de log de turno de combate, motivos de impedimento em `CanLevelUp` e erros de validação financeira em `Deposit`/`Withdraw`. |
| **`internal/bestiary`** | `types.go`, `generator.go` | ~15 | Nomes dos afixos procedurais (`NamePTBR`: `"Feroz"`, `"Covarde"`, `"Gigantesco"`, etc.). |

---

## 4. Avaliação sob a Ótica das Skills Recomendadas

### 4.1. Skill `godoctor` (Segurança AST, Qualidade e Tipagem)
- **Type Safety via Enums/Aliases:** O refactoring deve expandir os aliases de tipo customizados em `internal/i18n` (`UIKey`, `ItemDescriptionID`, `AffixID`, `MessageKey`), prevenindo o uso acidental de strings genéricas e garantindo verificação em tempo de compilação.
- **Invariância de Compilação & Gates:** As alterações não devem introduzir regressões. Cada etapa do refactoring deve passar nos testes de cobertura e no gate do compilador (`go vet ./...`).

### 4.2. Skill `engineering-flow` (Higiene de Código & Política 0.x Zero-Debt)
- **Zero-Debt Policy:** Na versão `0.x`, refatorações não devem manter pontes de compatibilidade ou aliases legados para suportar o código antigo. As assinaturas e chamadores devem ser refatorados diretamente.
- **Broken Window Code Hygiene:** Remover todos os literais duplicados e textos inline. Erros e logs de combate devem utilizar chaves tipadas ou funções formatadoras da camada `internal/i18n`.

### 4.3. Skill `git-workflow-and-versioning` (Commits Atômicos & Versionamento)
- **Commits Atômicos & Pequenos:** Um refactoring deste porte (~230 strings) não deve ser feito em um único commit gigante. Ele deve ser fatiado em incrementos lógicos, isolados por pacote/camada.
- **Padrão Conventional Commits:**
  - `refactor(i18n): add ui and item description keys to i18n package`
  - `refactor(ui): localize status bar labels in internal/ui`
  - `refactor(engine): migrate item descriptions and combat messages to i18n`
  - `refactor(tui): localize town, forest and smith screens`
  - `refactor(tui): localize tavern, chapel, guild, dragon and login screens`

---

## 5. Arquitetura da Solução Proposta

### 5.1. Expansão do Pacote `internal/i18n`

Proposta de estruturação das chaves e dicionários em `internal/i18n`:

```go
// internal/i18n/types.go

type UIKey string
type MessageKey string
type AffixID string

// Rótulos de Interface e Barra de Status
const (
	UIHero       UIKey = "hero"
	UILevel      UIKey = "level"
	UIHealth     UIKey = "health"
	UIGold       UIKey = "gold"
	UIBank       UIKey = "bank"
	UIDailyFights UIKey = "daily_fights"
	UIWeapon     UIKey = "weapon"
	UIArmor      UIKey = "armor"
	UIPotions    UIKey = "potions"
	UIAttack     UIKey = "attack"
	UIDefense    UIKey = "defense"
	UIExperience UIKey = "experience"
	UIMaxLevel   UIKey = "max_level"
)

// Chaves de Ações e Rodapés Genéricos
const (
	UIFooterNav          UIKey = "footer_nav"
	UIFooterSelectConfirm UIKey = "footer_select_confirm"
	UIFooterCombat       UIKey = "footer_combat"
)
```

### 5.2. Gestão de Templates e Mensagens Dinâmicas

Para mensagens com parâmetros dinâmicos (ex: `"Você depositou %d moedas de ouro"`), o pacote `internal/i18n` disponibilizará funções utilitárias do tipo `GetMessage(key MessageKey, args ...any) string`:

```go
// internal/i18n/pt_br.go

var MessageTemplatesPTBR = map[MessageKey]string{
	MsgBankDepositSuccess: "Você depositou %d moedas de ouro no cofre com segurança!",
	MsgCombatVictory:      "Você derrotou %s e ganhou %d XP e %d moedas de ouro!",
	MsgCombatDefeat:       "%s desferiu um golpe mortal! Você sucumbiu na escuridão...",
	MsgHealSuccess:        "Frei Anselmo unge seus ferimentos com óleos sagrados. Vida restaurada por %d moedas!",
}

func GetMessage(key MessageKey, args ...any) string {
	tmpl, ok := MessageTemplatesPTBR[key]
	if !ok {
		return string(key)
	}
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}
```

### 5.3. Desacoplamento da Camada de Domínio (`internal/engine`)

1. **Catálogo de Itens (`items.go`):** As structs de `Item` não armazenarão a string `Description` hardcoded. Elas utilizarão `i18n.ItemID` para buscar a descrição traduzida via `i18n.GetItemDescription(item.ID)`.
2. **Afixos em `bestiary` (`types.go`):** Substituir `NamePTBR string` por `ID i18n.AffixID`, mapeando os nomes dos afixos em `i18n.AffixNamesPTBR`.
3. **Motor de Combate (`combat.go`):** O `CombatEngine` continuará retornando a struct `TurnResult`, mas as mensagens textuais serão formatadas chamando `i18n.GetMessage(...)` ou retornando chaves puras para a TUI formatar.

---

## 6. Plano de Execução e Cronograma de Migração

A migração será dividida em **4 fases atômicas**:

### Fase 1: Infraestrutura de Chaves e Dicionários em `internal/i18n`
- Criar chaves e mapeamentos em `internal/i18n/types.go` e `internal/i18n/pt_br.go` para:
  - Rótulos da barra de status e UI.
  - Descrições do catálogo de itens.
  - Nomes de afixos procedurais.
  - Templates de mensagens de combate, banco e navegação.
- Adicionar testes unitários em `internal/i18n/pt_br_test.go` garantindo 100% de cobertura de chaves.

### Fase 2: Localização de `internal/ui`, `internal/engine` e `internal/bestiary`
- Refatorar `internal/ui/ui.go` para consumir os rótulos do cabeçalho via `i18n.GetUIText(...)`.
- Refatorar `internal/engine/items.go`, `combat.go` e `economy.go`.
- Refatorar `internal/bestiary/types.go` e `generator.go`.
- Executar `go test ./internal/...` para validar que nenhum comportamento numérico ou de combate foi alterado.

### Fase 3: Localização das Telas TUI (Core Hubs & Combate)
- Refatorar `town.go`, `forest.go`, `smith.go` e `dragon.go`.
- Substituir todos os menus e avisos do sistema por chamadas a `i18n.GetUIText` / `i18n.GetMessage`.

### Fase 4: Localização das Telas TUI Restantes & Higiene Final
- Refatorar `tavern.go`, `chapel.go`, `guild.go`, `game_over.go` e `login.go`.
- Executar auditoria estática com `grep` confirmando eliminação de strings fixas.
- Rodar suíte completa de testes (`go test ./...`) e verificação do compilador (`go vet ./...`).

---

## 7. Conclusão

A Issue 42 aborda um débito técnico importante para a sustentabilidade arquitetural do *The Legend of the Go Dragon*. A implementação da solução proposta garantirá estrita conformidade com a **ADR-0005** e a **ADR-0002**, eliminando o acoplamento de idioma na lógica do jogo e preparando o projeto para suporte a múltiplos idiomas de forma limpa e performática.
