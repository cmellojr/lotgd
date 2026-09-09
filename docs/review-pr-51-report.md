# Relatório de Code Review — Pull Request #51

**Repositório:** `cmellojr/lotgd`
**PR:** [#51 - docs(GDD): align keyboard shortcuts table with code implementation](https://github.com/cmellojr/lotgd/pull/51)
**Branch Origem:** `fix/43-gdd-keyboard-mappings` -> **Branch Destino:** `develop`
**Commit Analisado:** `c65e25cf168d255b3b29eadcd69932cc30d9cef2`
**Avaliador:** AI Software Engineer Agent (Jules)
**Data:** 8 de Março de 2026

---

## 📋 1. Visão Geral & Resumo Executivo

A Pull Request #51 tem como objetivo resolver a discrepância entre a documentação de atalhos de teclado e comandos TUI no [Game Design Document (GDD.md)](docs/GDD.md) e a implementação real no código em Go (localizado em `internal/tui/screens/`).

### Diagnóstico Geral
* **Status:** **Aprovado com Recomendação Menor de Atualização (Minor Improvement)**.
* **Qualidade do Commit:** Segue o padrão *Conventional Commits* (`docs(GDD): ...`), com descrição clara no body.
* **Impacto no Projeto:** Alto impacto de qualidade de documentação (*Broken Window Code Hygiene*). Elimina divergências que confundiam novos jogadores e desenvolvedores.

---

## 🔎 2. Análise Detalhada por Tela (GDD vs. Código Go)

Abaixo consta o mapeamento auditado de cada tela descrita na tabela do GDD contra a implementação nos modelos Bubble Tea do pacote `internal/tui/screens/`:

### 2.1. Login / Criação (`ScreenLogin`)
* **GDD PR #51:** `Digitar nome/senha, [Enter] confirma, [Tab] alterna`
* **Código (`login.go`):**
  * `Tab` / `Down`: Alterna o foco para o próximo campo/botão (`nextFocus()`).
  * `Shift+Tab` / `Up`: Alterna o foco para o campo anterior (`prevFocus()`).
  * `Enter`: Confirma autenticação ou avança foco.
* **Avaliação:** **100% Correto**.

---

### 2.2. Praça do Vilarejo (`ScreenTown`)
* **GDD Antigo:** `[F] Floresta, [T] Taverna, [C] Capela, [M] Ferraria, [G] Guilda, [D] Dragão, [S] Status, [Q] Sair`
* **GDD PR #51:** `[F] Floresta, [T] Taverna, [C] Capela, [M] Ferraria, [G] Guilda, [D] Dragão, [S] Salvar e Sair (Logout)`
* **Código (`town.go`):**
  * Teclas diretas: `F`, `T`, `C`, `M`, `G`, `D`, `B` (Banco do Vilarejo) e `S` (Salvar e Sair / Logout).
  * Removidos do código anteriormente: `[S]` como Status (agora `[S]` faz Logout) e `[Q]` como Sair (descontinuado).
* **Avaliação:** **Excelente correção**, pois removeu `[Q]` e corrigiu a ação da tecla `[S]`.
* **💡 Oportunidade de Melhoria:** A tecla `[B]` (Banco do Vilarejo) está presente no menu do hub e implementada interativamente no código (`bankMode`), porém não foi incluída na lista de atalhos da Praça no GDD.
  * *Sugestão:* Incluir `[B] Banco` na descrição da Praça no GDD.

---

### 2.3. Floresta Sombria (`ScreenForest`)
* **GDD PR #51:** `[P] Procurar monstro, [A] Atacar, [F] Fugir, [V] Voltar à vila`
* **Código (`forest.go`):**
  * Explorando: `P` / `E` / `Enter` para procurar monstro.
  * Em Combate: `A` / `Enter` (Atacar), `F` (Fugir), `P` / `U` (Usar Poção).
  * Retorno: `V` / `Esc` (Voltar à vila).
* **Avaliação:** **100% Correto** para os comandos principais do fluxo de navegação.

---

### 2.4. Ferraria (`ScreenSmith`)
* **GDD Antigo:** `[1..5] Comprar Armas, [6..0] Comprar Armaduras`
* **GDD PR #51:** `[1] Armas, [2] Armaduras, [3] Poções, [Tab] Trocar aba`
* **Código (`smith.go`):**
  * `1`, `2`, `3`: Alterna diretamente para as abas de Armas, Armaduras e Poções.
  * `Tab` / `Right` / `Left`: Alterna navegação entre abas.
  * `Up` / `Down` / `j` / `k`: Move seleção no catálogo.
  * `Enter` / `C`: Realiza compra.
  * `V` / `Esc`: Volta para a praça.
* **Avaliação:** **100% Correto**. Corrigiu o mapeamento incorreto antigo `1..5` / `6..0`.

---

### 2.5. Capela (`ScreenChapel`)
* **GDD Antigo:** `[C] Curar ferimentos, [B] Pedir bênção`
* **GDD PR #51:** `[C] Curar ferimentos, [D] Doação (10ouro), [M] Meditar`
* **Código (`chapel.go`):**
  * `C`: Ativa diretamente a opção de cura.
  * `D`: Ativa doação de 10 moedas de ouro.
  * `M`: Ativa meditação.
  * `V` / `Esc`: Volta para a praça.
* **Avaliação:** **100% Correto**. A opção `[B]` ("Pedir bênção") não existia no código.

---

### 2.6. Taverna (`ScreenTavern`)
* **GDD Antigo:** `[O] Ouvir fofocas, [F] Flertar com Cassandra`
* **GDD PR #51:** `[R]/[F] Ouvir fofocas, [C] Flertar com Cassandra`
* **Código (`tavern.go`):**
  * `R` / `F`: Fofocas com Dona Rosalinda.
  * `C`: Flertar/interagir com Cassandra.
  * `D`: Duelo com Cavaleiro Vermelho.
  * `N` / `M`: Ler mural de notícias.
  * `V` / `Esc`: Volta para a praça.
* **Avaliação:** **100% Correto**. A tecla antiga `[O]` foi corrigida para `[R]/[F]`.

---

### 2.7. Covil do Dragão (`ScreenDragon`)
* **GDD PR #51:** `[D] Desafiar o Dragão do Dia`
* **Código (`dragon.go`):**
  * `D` / `L` / `Enter`: Desafiar o Dragão.
  * `A` / `Enter`: Atacar.
  * `P` / `U`: Poção.
  * `F`: Fugir.
  * `V` / `Esc`: Voltar.
* **Avaliação:** **100% Correto**.

---

### 2.8. Guilda dos Aventureiros (`ScreenGuild`)
* **GDD Antigo:** *(Linha ausente na tabela do GDD)*
* **GDD PR #51:** `[A]/[T]` Avançar/consultar, `[R]/[C]` Consultar tabela, `[P]/[L]` Ler pergaminhos
* **Código (`guild.go`):**
  * `A` / `T`: Avançar de Nível (Treinamento com Mestre Tobias).
  * `R` / `C`: Consultar requisitos de níveis.
  * `P` / `L`: Ler pergaminhos antigos.
  * `V` / `Esc`: Voltar.
* **Avaliação:** **Excelente Adição!** A inclusão da Guilda corrige uma omissão na tabela.

---

## 📐 3. Conformidade com Diretrizes do Projeto (`AGENTS.md` & Skills)

1. **Idioma & Convenção PT-BR (`docs/universo-e-lore.md` & `AGENTS.md`):**
   * A documentação mantida em Português do Brasil (PT-BR) respeita integralmente a divisão arquitetural: textos e docs para usuário/jogador em PT-BR, identificadores e código em inglês.
2. **Higiene de Código & Débito Técnico Zero (`engineering-flow`):**
   * A PR elimina incongruências entre documentação e código real, prevenindo confusões para futuros colaboradores e jogadores.
3. **Commit Único Atômico (`git-workflow-and-versioning`):**
   * Commit bem estruturado: `docs(GDD): align keyboard shortcuts table with code implementation`.

---

## 💡 4. Sugestão de Ajuste Menor (Opcional)

Para que a tabela do GDD reflita **100%** de todas as funcionalidades principais presentes na praça central (`ScreenTown`), sugere-se adicionar o atalho `[B] Banco`:

```markdown
| **Praça do Vilarejo** | `ScreenTown` | `[F]` Floresta, `[T]` Taverna, `[C]` Capela, `[M]` Ferraria, `[G]` Guilda, `[D]` Dragão, `[B]` Banco, `[S]` Salvar e Sair (Logout) |
```

---

## 🏁 5. Conclusão e Veredito

* **Veredito:** **APROVADO (APPROVED)** ✅
* **Resumo:** A PR #51 é um excelente refinamento de documentação. Todos os testes automatizados da aplicação continuam passando sem regressões (`go test ./...` ok), e as correções no GDD refletem fielmente o comportamento atual do sistema TUI.
