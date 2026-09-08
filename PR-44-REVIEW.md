# 🔍 Relatório de Code Review — Pull Request #44

**Repositório:** `cmellojr/lotgd`
**PR Original:** [#44 — fix(tui): bloquear saída lateral da tela de Game Over em estado de derrota](https://github.com/cmellojr/lotgd/pull/44)
**Autor:** `@fernandomozone`
**Branch de Origem:** `fernandomozone:fix/tui-game-over-bypass`
**Branch de Destino:** `cmellojr:develop`
**Issue Relacionada:** Fecha [#37](https://github.com/cmellojr/lotgd/issues/37)
**Veredito:** **APROVADO (APPROVED) ✅**

---

## 📌 1. Resumo Executivo

O Pull Request #44 corrige um *exploit* no fluxo TUI de derrota onde o jogador conseguia burlar a tela de **Game Over** e evitar a penalidade moratória de perda de ouro na bolsa e pontos de experiência (XP).

A solução é elegante, direta e segura. Foram adicionados testes de regressão automatizados bem estruturados e as alterações respeitam integralmente as diretrizes arquiteturais do projeto estabelecidas em `AGENTS.md`.

---

## 🐞 2. Contexto e Causa Raiz do Bug

### O Problema
No *The Legend of the Go Dragon*, quando a vida do jogador chega a `0` em combate na Floresta Sombria ou no Covil do Dragão, a tela transiciona localmente para o estado de derrota (`forestStateDefeat` ou `dragonStateDefeat`).

No entanto, antes do bloco `switch s.state` ser avaliado nos métodos `Update()` de `ForestScreen` e `DragonScreen`, havia uma checagem global de atalhos de navegação para retornar à cidade (`[V]`, `[Esc]` ou `[C]`).

### O Efeito de Bypass
Como essa checagem acontecia no topo do método `Update()` e só bloqueava quando o estado fosse explicitamente de combate ativo (`forestStateCombat` / `dragonStateCombat`), pressionar `[V]`, `[Esc]` ou `[C]` enquanto a mensagem de derrota estivesse visível na tela acionava imediatamente a transição para `ui.ScreenTown`:

```go
// Comportamento anterior (com bug)
if k == "V" || k == "ESC" || (s.state == forestStateExploring && k == "C") {
    if s.state == forestStateCombat {
        // Bloqueava apenas DURANTE o combate
        return s, nil
    }
    SavePlayer(s.db, s.player)
    return s, func() tea.Msg {
        return ui.ChangeScreenMsg{Screen: ui.ScreenTown}
    }
}
```

Ao retornar direto para a cidade (`ScreenTown`), a transição para `ScreenGameOver` nunca ocorria no roteador principal (`MainModel`). Como a aplicação da penalidade de morte (`econ.ProcessDeathPenalty(m.player)`) é executada no `MainModel.Update` exclusivamente quando a mensagem `ChangeScreenMsg` indica a tela `ScreenGameOver`, o jogador retornava para a Praça Central com 0 HP, mas **sem perder ouro nem XP**.

---

## 🛠️ 3. Análise Detalhada das Alterações

### 3.1. `internal/tui/screens/forest.go`
```go
<<<<<<< SEARCH
		// Retornar à cidade
		if k == "V" || k == "ESC" || (s.state == forestStateExploring && k == "C") {
=======
		// Retornar à cidade.
		// Em estado de derrota nenhuma tecla escapa: o fluxo tem de passar pela
		// tela de Game Over, onde a penalidade de morte é aplicada.
		if s.state != forestStateDefeat &&
			(k == "V" || k == "ESC" || (s.state == forestStateExploring && k == "C")) {
>>>>>>> REPLACE
```
* **Avaliação:** Adicionar a guarda `s.state != forestStateDefeat` impede que atalhos de saída sejam processados em estado de derrota, fazendo com que qualquer tecla pressionada caia no `case forestStateDefeat:` do `switch s.state`, que redireciona corretamente para `ui.ScreenGameOver`.

### 3.2. `internal/tui/screens/dragon.go`
```go
<<<<<<< SEARCH
		if k == "V" || k == "ESC" {
=======
		// Em estado de derrota nenhuma tecla escapa: o fluxo tem de passar pela
		// tela de Game Over, onde a penalidade de morte é aplicada.
		if (k == "V" || k == "ESC") && s.state != dragonStateDefeat {
>>>>>>> REPLACE
```
* **Avaliação:** Mesma lógica aplicada à tela do Covil do Dragão. O bloqueio garante paridade de comportamento entre as duas telas de combate do jogo.

### 3.3. `internal/tui/screens/defeat_bypass_test.go`
Foi criado um arquivo dedicado para testes de regressão com três casos principais:
1. `TestForestScreen_DefeatHasNoSideExit`: Testa se no estado de derrota na floresta, as teclas `"v"`, `"V"`, `"esc"` e `"c"` resultam estritamente em transição para `ui.ScreenGameOver`.
2. `TestDragonScreen_DefeatHasNoSideExit`: Testa se no estado de derrota no dragão, as teclas `"v"`, `"V"` e `"esc"` resultam em `ui.ScreenGameOver`.
3. `TestForestScreen_ExploringStillReturnsToTown`: Garante que, fora da derrota (durante exploração normal), a tecla `"v"` continua retornando para `ui.ScreenTown`.

Os testes utilizam a função auxiliar `screenAfterKey`, que simula a emissão de `tea.KeyMsg` e captura a mensagem de comando resultante (`ui.ChangeScreenMsg`). A criação de banco de dados temporário em isolamento (`t.TempDir()` com `t.Cleanup()`) garante execução limpa e sem efeitos colaterais.

---

## 📐 4. Conformidade com as Diretrizes do Projeto (`AGENTS.md`)

| Critério / Diretriz | Status | Observações |
| :--- | :---: | :--- |
| **Código em Inglês** | ✅ SIM | Nomes de funções, variáveis e testes em inglês idiomático (`screenAfterKey`, `newDefeatedPlayer`, `TestForestScreen_DefeatHasNoSideExit`). |
| **Documentação/Comentários em PT-BR** | ✅ SIM | Comentários em Go no padrão do projeto, explicativos e didáticos em português. |
| **Estilo Go & Idiomaticidade** | ✅ SIM | Código limpo, tratamento adequado de tipos e alinhado com as convenções do Google Go Style Guide. |
| **Arquitetura Bubble Tea (TEA)** | ✅ SIM | Manipulação pura e correta de `tea.Model` e `tea.Cmd`. |
| **Ausência de Erros Silenciados** | ✅ SIM | Nenhuma violação ou supressão indevida de erros (`_ = err`). |
| **Verificações de Qualidade (`go vet`, `go test`)** | ✅ SIM | `go vet ./...` e `go test ./...` executam e passam sem falhas. |

---

## 🎯 5. Pontos Fortes

1. **Correção Pontual e Precisa:** A alteração mexe exatamente onde o problema ocorria, sem reescritas desnecessárias ou impacto colateral na navegação normal.
2. **Excelente Cobertura de Testes:** O teste de regressão cobre não apenas as teclas minúsculas/maiúsculas de saída e ESC, mas também valida a não-regressão do fluxo normal de exploração.
3. **Didática e Clareza:** Os comentários adicionados explicam o motivo do bloqueio para futuros mantenedores da base de código.

---

## 💡 6. Sugestões e Recomendações (Nitpicks Opcionais)

Não há nenhum fator impeditivo para o merge. Abaixo estão apenas duas pequenas sugestões de melhoria contínua para consideração futura:

1. **Normalização de Teclas no Helper de Teste:**
   Em `defeat_bypass_test.go`, a função `screenAfterKey` faz um switch manual para `"esc"`. Se no futuro outros atalhos especiais forem testados (como `enter` ou `space`), pode ser interessante expandir o helper ou utilizar mapeamento padronizado de `tea.KeyMsg`.
2. **Formatação no `db.go`:**
   Como apontado pelo próprio autor do PR na descrição, o `gofmt` acusa apenas a pendência legada no `internal/storage/db.go` (registrada na Issue #10), estando todas as telas e os novos testes 100% formatados.

---

## 🏆 7. Veredito Final

O PR **#44** cumpre perfeitamente seu propósito, resolve a falha de bypass de Game Over e adiciona uma suíte de testes robusta.

**Recomendação:** Aprovado para merge no branch `develop`.
