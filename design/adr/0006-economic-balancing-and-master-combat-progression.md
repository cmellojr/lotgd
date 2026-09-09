# ADR-0006: Modelo de Balanceamento Econômico e Progressão por Combate com Mestres

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon
- Related Issues: #35

## 1. Contexto

A issue #35 solicita a definição do balanceamento econômico entre recompensas (ouro e experiência) e custos (equipamentos, curas e avanço de nível) no *The Legend of the Go Dragon*.

A implementação atual em `internal/engine/progression.go` exige simultaneamente **XP E Ouro** para avançar de nível, utilizando uma tabela estática de custos financeiros (`CostGold`) sem qualquer componente de combate, risco de derrota ou limite diário de tentativas.

### Contexto Histórico e Requisito de Fidelidade

O projeto possui como premissa ser um **clone fiel do jogo de BBS original de 1989**, *Legend of the Red Dragon* (LORD), criado por Seth Able Robinson — e não do *Legend of the Green Dragon* (LoGD). Embora o LoGD seja uma referência técnica útil, ele expandiu significativamente o escopo original (múltiplas cidades, sistema de vida após a morte, clãs, e sistema complexo de raças e classes).

Conforme documentado historicamente (ver *Wikipedia: Legend of the Red Dragon*, seção sobre *Turgon's Warrior Training*), no jogo original de 1989 **subir de nível NÃO é uma transação financeira**. O fluxo original do LORD funciona da seguinte forma:
1. O jogador acumula a quantidade de Experiência (XP) exigida para o próximo nível.
2. Ao atingir a XP necessária, o jogador viaja até o centro de treinamento de guerreiros (*Turgon's Warrior Training* / Guilda) e **desafia o seu Mestre**.
3. A promoção ocorre através de um combate individual com risco real de derrota (que envia o jogador ferido para a Capela/Necrotério) e com limite de tentativas diárias.
4. O **Ouro** no jogo original é reservado estritamente para a compra de equipamentos (armas e armaduras na Ferraria/Loja), cura na Capela, hospedagem na Pousada e transações bancárias, e **não** para desbloquear níveis.

O arquivo `train.php` do próprio LoGD preserva esse mecanismo herdado do LORD, confirmando que a progressão por combate contra o Mestre é a mecânica canônica original.

---

## 2. Decisão

Avaliamos duas opções principais de design para o modelo econômico e de progressão:

### Opção A: Manter o modelo econômico atual (XP + Ouro) e recalibrar valores
- **Descrição**: O jogador acumula XP e paga uma taxa fixa em ouro para subir de nível instantaneamente na Guilda.
- **Prós**:
  - Implementação simples e já existente no código.
  - Zero risco de derrota ao subir de nível.
- **Contras**:
  - Diverge diretamente da fidelidade ao *Legend of the Red Dragon* (1989).
  - Cria um "dreno de ouro" artificial que compete com a compra de armas/armaduras, forçando grind excessivo de moedas.
  - Remove a tensão dramática e a sensação de conquista das lutas de graduação.

### Opção B (RECOMENDADA E APROVADA): Adotar o modelo original do LORD (1989) — XP + Desafio de Combate contra Mestre
- **Descrição**: A progressão de nível é destravada ao atingir a XP exigida e vencer o combate contra o Mestre de Treinamento correspondente em *Turgon's Warrior Training* (Guilda dos Aventureiros). O combate contra o Mestre consome uma tentativa diária de Mestre (limite diário de 1 a 3 tentativas) e possui risco real de derrota. O Ouro é totalmente liberado para compra de equipamentos, serviços do vilarejo e economia do banco.
- **Prós**:
  - Paridade e fidelidade total ao *Legend of the Red Dragon* (1989).
  - Separação clara na economia do jogo: XP destrava níveis/mestres, enquanto Ouro compra poder (equipamentos) e utilidades.
  - Transforma a mudança de nível em um evento de gameplay marcante e estratégico.
- **Contras**:
  - Requer expansão na lógica da Guilda para suportar encontros de combate e controle de tentativas diárias.

**Decisão do Mantenedor**: Aprovada a **Opção B**.

*Nota de Escopo*: Recursos adicionais introduzidos posteriormente pelo LoGD (como sistema elaborado de raças/classes, clãs, múltiplas vilas e reencarnação/vida após a morte) não fazem parte desta decisão e permanecem fora do escopo do v0.0.3, conforme planejado no roadmap.

---

## 3. Detalhamento do Modelo Recomendado (Opção B)

### 3.1 Tabela de Requisitos, Mestres, Recompensas e Ritmo Esperado

A economia e a curva de experiência foram calibradas considerando **15 lutas diárias na floresta** (`engine.DailyForestFights = 15`) e a disponibilidade de equipamentos por Tier na Ferraria do Mestre Torin.

| Nível Alvo | XP Exigido | Mestre de Treinamento (NPC) | Stats do Mestre (HP / ATK / DEF) | Recompensas de Nível (+HP / +ATK / +DEF) | Tier de Equipamento Esperado | Dias de Jogo Estimados (Pace Acumulado) |
|:---:|:---:|:---|:---:|:---:|:---:|:---:|
| **1** | 0 | *(Início do jogo)* | - | - | Tier 1 (Bastão / Roupas) | Dia 1 |
| **2** | 100 | Mestre Halder | 25 HP / 7 ATK / 3 DEF | +15 HP / +2 ATK / +2 DEF | Tier 1 (Adaga / Couro Leve) | Dia 1 (15 lutas) |
| **3** | 300 | Mestre Tobias | 45 HP / 12 ATK / 6 DEF | +20 HP / +3 ATK / +2 DEF | Tier 1 (Espada de Madeira / Couro Batido) | Dias 2 - 3 |
| **4** | 700 | Mestre Kaelen | 75 HP / 18 ATK / 10 DEF | +25 HP / +4 ATK / +3 DEF | Tier 2 (Espada Curta / Cota de Malha) | Dias 4 - 5 |
| **5** | 1.500 | Mestra Vanya | 110 HP / 25 ATK / 15 DEF | +30 HP / +5 ATK / +4 DEF | Tier 2 (Machado de Guerra / Cota de Escamas) | **Dias 6 - 7** (~1 semana de jogo) |
| **6** | 3.000 | Mestre Roderick | 150 HP / 34 ATK / 21 DEF | +35 HP / +6 ATK / +5 DEF | Tier 3 (Espada Longa / Armadura de Placas) | Dias 9 - 10 |
| **7** | 5.500 | Mestra Elora | 200 HP / 44 ATK / 28 DEF | +40 HP / +7 ATK / +6 DEF | Tier 3 (Martelo de Guerra / Placa Completa) | Dias 12 - 13 |
| **8** | 9.000 | Mestre Thorgrim | 260 HP / 55 ATK / 36 DEF | +45 HP / +8 ATK / +7 DEF | Tier 4 (Lâmina Rúnica / Placa Mithril) | Dias 15 - 16 |
| **9** | 14.000 | Mestre Valerius | 330 HP / 68 ATK / 45 DEF | +50 HP / +10 ATK / +8 DEF | Tier 4 (Machado Mithril / Armadura Adamantina) | Dias 18 - 20 |
| **10** | 22.000 | Mestre Arthorian | 420 HP / 83 ATK / 55 DEF | +60 HP / +12 ATK / +10 DEF | Tier 4 (Espada Lendária / Placa de Dragão) | **Dias 22 - 25** (~3 a 4 semanas) |

#### Meta Explícita de Ritmo de Progressão:
- **Nível 1 ao 5 (Iniciação / Early Game)**: O jogador atinge o Nível 5 em aproximadamente **6 a 7 dias reais** de jogo diário (acumulando ~100 lutas na floresta e atualizando seus equipamentos na Ferraria).
- **Nível 5 ao 10 (Maturidade / End Game)**: A curva de XP desacelera progressivamente, exigindo entre **22 e 25 dias acumulados** de atividade constante para atingir o Nível 10 e se qualificar para enfrentar o **Dragão do Dia**.

---

### 3.2 Proposta de Teste / Simulação Determinística contra Regressão

Para garantir que o balanceamento dos Mestres e a curva de progressão não sofram regressões em refatorações futuras, propõe-se a criação do teste determinístico `TestMasterCombatWinRateProgression` no pacote `internal/engine` (ou `internal/bestiary`):

```go
// TestMasterCombatWinRateProgression simula 1.000 combates entre um jogador adequadamente
// equipado para o seu nível e o Mestre de Treinamento correspondente.
//
// Critério de Aceite:
// - Jogador com equipamento recomendado do nível: Taxa de vitória entre 65% e 85%.
// - Jogador sem equipamentos (atributos base apenas): Taxa de vitória < 35% (exigindo compra de armas/armaduras).
func TestMasterCombatWinRateProgression(t *testing.T) {
    // 1. Instanciar CombatEngine determinístico com seed fixa.
    // 2. Para cada nível de 1 a 9:
    //    a. Criar Player com estatísticas acumuladas do nível e equipamentos do Tier correspondente.
    //    b. Carregar o Mestre do próximo nível da tabela MasterCatalog.
    //    c. Executar 1.000 combates e contabilizar vitórias.
    //    d. Validar se winRate >= 0.65 e winRate <= 0.85.
}
```

Essa simulação funcionará nos mesmos moldes do teste `TestLevelOneWinRateVsFerozMonsters` em `internal/bestiary/bestiary_test.go`, servindo como trava de segurança na suíte de testes contínuos (`go test ./...`).

---

### 3.3 Documentos a Serem Atualizados Pós-Aprovação

Quando a implementação técnica desta decisão for iniciada, as seguintes seções dos documentos de design e planejamento deverão ser atualizadas:

1. **`docs/GDD.md`**:
   - **Seção 1.2 (Core Gameplay Loop)**: Atualizar o diagrama Mermaid substituindo *"Tem Nível para o Dragão?"* e a rota da Guilda por *"Desafiar Mestre de Treinamento (Turgon's Warrior Training / Guilda)"*.
   - **Seção 2.1 (Economia de Recursos)**:
     - Atualizar a definição de **Ouro**: desvinculá-lo do avanço de nível e explicitar seu uso em Equipamentos (Ferraria), Cura (Capela) e Banco.
     - Atualizar a definição de **Experiência (XP)**: pontuar que ao atingir o limiar de XP, o jogador habilita a opção de desafiar o Mestre de Nível na Guilda.
   - **Seção 2.3 / 3 (Guilda / Treinamento de Guerreiros)**: Incluir a tabela de Mestres, a limitação de tentativas diárias e as regras de combate de promoção.

2. **`docs/roadmap.md`**:
   - **Fase 2 (Game Engine / `internal/engine`)**: Atualizar a descrição do subitem de progressão de *"Sistema de economia (ouro na bolsa vs ouro no banco), experiência e avanço de nível"* para *"Sistema de experiência, curva de nível e combates de promoção contra Mestres de Treinamento"*.
   - **Fase 3 (TUI / `guild.go`)**: Atualizar o item da Guilda para especificar a interface de desafio ao Mestre e a exibição das tentativas diárias de Mestre restantes.

---

## 4. Consequências

### Positivas
- **Fidelidade Histórica**: O jogo alinha-se 100% com o funcionamento do clássico *Legend of the Red Dragon* (1989).
- **Clareza Econômica**: O ouro deixa de ser uma barreira arbitrária de nível e passa a ter utilidade real na progressão de equipamentos (armas/armaduras na Ferraria) e sobrevivência (poções e cura na Capela).
- **Engajamento e Desafio**: Subir de nível torna-se um marco comemorativo de combate individual contra NPCs lendários (Mestres).
- **Proteção contra Regressões**: A inclusão de testes de simulação de taxa de vitória contra os Mestres garante estabilidade nos atributos ao longo de refatorações.

### Negativas / Restrições
- **Esforço de Desenvolvimento**: Exige a criação do catálogo de Mestres em Go (`internal/engine`), atualização das structs de jogador/persistência no SQLite para registrar tentativas diárias de Mestre, e adaptação da tela da Guilda (`internal/tui/screens/guild.go`).

---

## 5. Conformidade e Verificação

- **Leitura Obrigatória**: Desenvolvedores e agentes de IA devem consultar este ADR antes de modificar `internal/engine/progression.go` ou `internal/tui/screens/guild.go`.
- **Validação Automatizada**: A implementação futura deverá incluir o teste `TestMasterCombatWinRateProgression` e manter 100% dos testes passando em `go test ./...`.
