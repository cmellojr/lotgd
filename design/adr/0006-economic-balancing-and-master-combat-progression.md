# ADR-0006: Modelo de Balanceamento Econômico e Progressão por Combate com Mestres

- Status: Approved
- Date: 2026-09-12
- Author(s): Carlos Mello Jr
- Deciders: Carlos Mello Jr
- Related Issues: #35

## 1. Contexto e Fundamentação Histórica

A issue #35 solicita a definição do balanceamento econômico entre recompensas (ouro e experiência) e custos (equipamentos, curas e avanço de nível) no *The Legend of the Go Dragon*.

A implementação inicial em `internal/engine/progression.go` exigia simultaneamente **XP E Ouro** para avançar de nível, utilizando uma tabela estática de custos financeiros (`CostGold`) sem qualquer componente de combate, risco de derrota ou limite diário de tentativas.

### 1.1 Contexto Histórico e Referência Canônica
O projeto possui como premissa ser um **clone fiel do jogo de BBS original de 1989**, *Legend of the Red Dragon* (LORD), criado por Seth Able Robinson — e **não** do *Legend of the Green Dragon* (LoGD). Embora o LoGD seja uma referência técnica útil no ecossistema Go/Web, ele expandiu significativamente o escopo original (múltiplas cidades, vida após a morte, clãs e sistema complexo de raças e classes).

Conforme documentado historicamente (ver *Wikipedia: Legend of the Red Dragon*, seção sobre *Turgon's Warrior Training*), no jogo original de 1989 **subir de nível NÃO é uma transação financeira**. O fluxo canônico do LORD funciona da seguinte forma:
1. O jogador acumula a quantidade de Experiência (XP) exigida para o próximo nível combatendo na floresta.
2. Ao atingir o limiar de XP, o jogador viaja até o centro de treinamento de guerreiros (*Turgon's Warrior Training* / Guilda) e **desafia o seu Mestre de Nível**.
3. A promoção ocorre através de um combate individual por turnos.
4. O combate contra o Mestre possui um limite estrito de **1 tentativa por dia** (por Mestre).
5. Se o jogador for derrotado pelo Mestre, **não há morte nem perda de HP**; contudo, o jogador perde a oportunidade do dia e precisa aguardar o próximo amanhecer ("Novo Dia") para tentar novamente.
6. O **Ouro** no jogo original é reservado estritamente para a compra de equipamentos (armas e armaduras na Ferraria/Loja), cura na Capela, serviços na Taverna/Pousada e transações bancárias, e **não** para pagar taxas de avanço de nível.

### 1.2 Princípio de Fidelidade vs. RPGs de Exploração (ex: Skyrim)
Um "clone fiel do LORD" diferencia-se fundamentalmente de RPGs de exploração para um jogador (como *Skyrim*). Enquanto RPGs de exploração simulam um mundo extenso para um único jogador habitar sozinho sem restrições de tempo real, o LORD simula uma **comunidade assíncrona competindo por recursos diários escassos** (*Forest Fights* limitados a 15 por dia).

Toda mecânica de jogo em *The Legend of the Go Dragon* deve ser projetada para reforçar a escassez de turnos e a competição entre jogadores reais no servidor, em vez de incentivar um "grind" ilimitado de moedas ou cliques indolores de compra. O combate contra o Mestre reintroduz o risco estratégico e a valorização do tempo real do jogador.

---

## 2. Decisão

Adotamos a **Opção B (Modelo do LORD 1989)** para a economia e a progressão do jogo:

### Opção B: Modelo Canônico do LORD (1989) — XP + Desafio de Combate contra Mestre
- **Descrição**: A progressão de nível é destravada ao atingir a XP exigida e vencer o combate por turnos contra o Mestre correspondente em *Turgon's Warrior Training* (Guilda dos Aventureiros).
- **Estrutura de 12 Níveis**: O jogo possui exatamente **12 níveis de maestria**. O Nível 12 é o topo do jogo, onde se encontra o Mestre Turgon. O jogador **só pode buscar e desafiar o Dragão Ancestral após vencer o Mestre Turgon no Nível 12**.
- **Regras do Combate com Mestre**:
  - Limite de **1 tentativa por dia por Mestre**.
  - Em caso de derrota contra o Mestre: o herói **não morre nem perde pontos de vida (HP)**, mas consome a tentativa do dia, exigindo que o jogador retorne no dia seguinte.
- **Papel do Ouro**: O ouro é 100% liberado de taxas de promoção de nível e passa a ser o dreno primário da economia na compra de armas e armaduras escalonadas na Ferraria do Mestre Torin, serviços curativos na Capela e depósitos no Banco.

---

## 3. Detalhamento da Tabela de Progressão (12 Níveis & Mestres)

A economia e a curva de experiência foram calibradas considerando **15 lutas diárias na floresta** (`engine.DailyForestFights = 15`) e a necessidade de equipamentos da Ferraria para superar os atributos dos Mestres.

| Nível Alvo | XP Exigido | Mestre de Treinamento (NPC) | Stats do Mestre (HP / ATK / DEF) | Recompensas de Nível (+HP / +ATK / +DEF) | Tier de Equipamento Esperado (Ferraria) | Custo Estimado do Equipamento | Dias Acumulados Estimados (Ritmo Esperado) |
|:---:|:---:|:---|:---:|:---:|:---:|:---:|:---:|
| **1** | 0 | *(Início do Jogo)* | - | - | Tier 1 (Bastão / Roupas) | 0 Ouro | Dia 1 |
| **2** | 100 | Mestre Halder | 25 HP / 7 ATK / 3 DEF | +15 HP / +2 ATK / +2 DEF | Tier 2 (Adaga / Couro Leve) | ~90 Ouro | Dia 1 - 2 |
| **3** | 300 | Mestre Tobias | 45 HP / 12 ATK / 6 DEF | +20 HP / +3 ATK / +2 DEF | Tier 3 (Espada Madeira / Couro Batido) | ~190 Ouro | Dias 3 - 4 |
| **4** | 700 | Mestre Kaelen | 75 HP / 18 ATK / 10 DEF | +25 HP / +4 ATK / +3 DEF | Tier 4 (Espada Curta / Cota Malha) | ~380 Ouro | Dias 5 - 7 |
| **5** | 1.500 | Mestra Vanya | 110 HP / 25 ATK / 15 DEF | +30 HP / +5 ATK / +4 DEF | Tier 5 (Machado Guerra / Cota Escamas) | ~750 Ouro | Dias 8 - 10 |
| **6** | 3.000 | Mestre Roderick | 150 HP / 34 ATK / 21 DEF | +35 HP / +6 ATK / +5 DEF | Tier 6 (Espada Longa / Placa) | ~1.300 Ouro | **Dias 10 - 12** (Metade da Jornada) |
| **7** | 5.500 | Mestra Elora | 200 HP / 44 ATK / 28 DEF | +40 HP / +7 ATK / +6 DEF | Tier 7 (Martelo Guerra / Placa Completa) | ~2.200 Ouro | Dias 13 - 15 |
| **8** | 9.000 | Mestre Thorgrim | 260 HP / 55 ATK / 36 DEF | +45 HP / +8 ATK / +7 DEF | Tier 8 (Lâmina Rúnica / Placa Mithril) | ~3.600 Ouro | Dias 16 - 18 |
| **9** | 14.000 | Mestre Valerius | 330 HP / 68 ATK / 45 DEF | +50 HP / +10 ATK / +8 DEF | Tier 9 (Machado Mithril / Adamantina) | ~6.300 Ouro | Dias 19 - 22 |
| **10** | 22.000 | Mestre Arthorian | 420 HP / 83 ATK / 55 DEF | +60 HP / +12 ATK / +10 DEF | Tier 10 (Lâmina Mística / Placa Rúnica) | ~9.000 Ouro | Dias 23 - 25 |
| **11** | 35.000 | Mestra Ignis | 530 HP / 100 ATK / 68 DEF | +75 HP / +15 ATK / +12 DEF | Tier 11 (Espada Titânica / Placa Obscura) | ~13.000 Ouro | Dias 26 - 28 |
| **12** | 55.000 | Mestre Turgon | 680 HP / 120 ATK / 85 DEF | +90 HP / +18 ATK / +15 DEF | Tier 12 (Lendária / Placa de Dragão) | ~18.000 Ouro | **Dias 25 - 30** (Qualificação para o Dragão) |

### 3.1 Metas Explícitas de Ritmo de Progressão (*Pacing*)
- **Nível 1 ao 6 (Early/Mid Game)**: O jogador atinge o Nível 6 em aproximadamente **10 a 12 dias reais** de atividade diária constante (acumulando ~150 a 180 lutas na floresta e atualizando seus equipamentos na Ferraria).
- **Nível 6 ao 12 (End Game & Maestria)**: A curva desacelera progressivamente, exigindo entre **25 e 30 dias acumulados** para atingir o Nível 12 e vencer o Mestre Turgon.
- **Portão do Covil do Dragão**: O portão de qualificação para o Dragão Ancestral fica **fixado no Nível 12**. A verificação em `internal/tui/screens/dragon.go` será ajustada para `Level < 12`, garantindo que o confronto final ocorra apenas após o herói provar sua maestria contra Turgon.

---

## 4. Proposta de Teste / Simulação Determinística contra Regressões

Para garantir que a curva de dificuldade dos Mestres e a importância dos equipamentos permaneçam equilibradas em refatorações futuras, propõe-se a criação da suíte de teste determinística `TestMasterCombatWinRateProgression` em `internal/engine`:

```go
// TestMasterCombatWinRateProgression simula 1.000 combates entre um jogador adequadamente
// equipado para o seu nível e o Mestre de Treinamento correspondente.
//
// Critério de Aceite:
// - Jogador com equipamento recomendado do nível: Taxa de vitória entre 65% e 85%.
// - Jogador sem equipamentos (atributos base apenas): Taxa de vitória < 35% (exigindo a compra de armas/armaduras na Ferraria).
func TestMasterCombatWinRateProgression(t *testing.T) {
    // 1. Instanciar CombatEngine determinístico com seed fixa.
    // 2. Para cada nível de 1 a 11:
    //    a. Criar Player com estatísticas acumuladas do nível e equipamentos do Tier correspondente.
    //    b. Carregar o Mestre do próximo nível da tabela MasterCatalog.
    //    c. Executar 1.000 combates e contabilizar vitórias.
    //    d. Validar se winRate >= 0.65 e winRate <= 0.85.
}
```

Esta simulação funcionará nos mesmos moldes do teste `TestLevelOneWinRateVsFerozMonsters` em `internal/bestiary/bestiary_test.go`.

---

## 5. Consequências & Extensibilidade Futura

### 5.1 O que fica fiel ao LORD agora (v0.0.x / v0.1)
- 12 níveis de progressão com 12 Mestres canônicos (com Mestre Turgon no Nível 12).
- Requisito do Nível 12 para desafiar o Dragão Ancestral.
- Limite de 1 tentativa diária por Mestre sem penalidade de morte/perda de HP em caso de derrota.
- Remoção total do custo em ouro para subir de nível (`CostGold = 0`).

### 5.2 Extensibilidade e Arquitetura Futura
- **Dados Configuráveis**: A tabela de Mestres e os níveis (`MasterCatalog` / `LevelTable`) devem ser implementados como estruturas de dados puras (slices/maps de structs em Go), sem constantes dispersas no código.
- **Portas Abertas**: Caso o projeto no futuro decida ir além do escopo do LORD (ex: adicionar um 13º nível, múltiplos mestres por nível ou vilas alternativas como no LoGD), isso será estritamente uma **alteração de dados/configuração**, sem necessidade de refatorar o motor de combate ou as regras de negócio.
- *Nota de Escopo*: Recursos do LoGD (como clãs, montarias, companheiros, vida após a morte e raças/classes) não fazem parte do escopo atual do LORD 1989 e permanecem fora do roadmap v0.0.3.

### 5.3 Sinalização Explícita de Mudança de Escopo (Scope Change)
Transformar a Guilda de uma tela de "compra de nível" para uma interface interativa de **combate por turnos contra Mestres com controle de tentativas diárias** representa uma alteração funcional de escopo.

**Por este motivo, NENHUM CÓDIGO de produção (Go) é alterado nesta tarefa.** Esta decisão de arquitetura estabelece a especificação necessária para que uma **issue de implementação separada** seja aberta e executada a seguir.

---

## 6. Mapeamento de Atualizações na Documentação de Suporte

Nas próximas etapas da documentação, as seguintes seções serão atualizadas:

1. **`docs/GDD.md`**:
   - **Seção 1.2 (Core Gameplay Loop)**: Atualizar o diagrama Mermaid com "Alcançou Nível 12?" e rota "Desafiar Mestre na Guilda (Turgon's Warrior Training)".
   - **Seção 2.1 (Economia de Recursos)**: Atualizar definicão de Ouro (exclusivo para Ferraria, Capela, Banco) e XP (libera desafio contra Mestre de Nível).
   - **Seção 2.3 / 3 (Guilda / Treinamento de Guerreiros)**: Detalhar a tabela de 12 Mestres, limite de 1 tentativa/dia e ausência de pena de morte na derrota do Mestre.

2. **`docs/roadmap.md`**:
   - **Fase 2 (Game Engine)**: Atualizar subitem de progressão para "Sistema de experiência (12 níveis), catálogo de 12 tiers de equipamentos e combate de promoção contra Mestres de Treinamento".
   - **Fase 3 (TUI / Guilda & Dragão)**: Atualizar especificações da tela da Guilda (combate contra Mestre e tentativas diárias) e Covil do Dragão (trava de Nível 12).
