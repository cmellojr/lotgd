# ADR-0007: Add Equipment Sale and Trade-In Credit at the Smith

- Status: Proposed
- Date: 2026-09-14
- Author(s): Carlos Mello Jr.
- Deciders: Carlos Mello Jr.
- Related Issues: #40

## 1. Context

A issue #40 reporta que a Ferraria do Mestre Torin (`internal/tui/screens/smith.go`) apresenta duas lacunas funcionais relativas à gestão de equipamentos:
1. **Ausência de mecanismo de venda avulsa**: O jogador não possui opção de vender a arma ou armadura atualmente equipada para obter moedas de ouro sem realizar uma nova compra.
2. **Ausência de crédito de troca (*trade-in credit*)**: A aquisição de uma nova arma ou armadura descarta o equipamento anterior sem conceder retorno financeiro ou abater o valor do novo item.

### Scope Control (Escopo Restrito desta ADR)
Esta ADR cobre exclusivamente os itens 1 e 2 da issue #40 (venda de equipamento inexistente e troca sem crédito de abatimento). O item 3 da issue #40 (exibição de aviso/confirmação em substituições que resultem em *downgrade* de atributos) aponta para um mecanismo de confirmação de interface estruturalmente diferente e independente do cálculo financeiro de troca/venda. Por essa razão, o item 3 foi definido como fora de escopo para este documento e será tratado em uma issue de acompanhamento e ADR dedicada posteriormente.

### Contexto Histórico e Referência de Implementação
- **LORD Original (1989)**: No jogo de BBS original *Legend of the Red Dragon* por Seth Able Robinson, a venda do equipamento antigo não apenas era permitida como documentada como etapa necessária para adquirir um novo item. O comerciante da loja oferecia um valor de recompra gerado aleatoriamente a cada tentativa. O jogador podia recusar a oferta, sair da loja, retornar e receber uma nova oferta com valor diferente pela mesma arma ou armadura.
- **LoGD (Fork de Referência)**: O *Legend of the Green Dragon* (implementado nos arquivos `weapons.php` e `armor.php`) simplificou a mecânica original aplicando uma taxa fixa de crédito de troca correspondente a 75% do valor de compra do item equipado, sem sistema de barganha ou ofertas aleatórias.
- **Estado Atual do Projeto**: O código atual em `internal/tui/screens/smith.go` não contempla venda nem crédito de troca; a compra de um item novo substitui o item equipado e zera o equipamento anterior sem qualquer compensação financeira.

A adição desses mecanismos trata-se de um ajuste de fidelidade às regras do jogo original, corrigindo uma lacuna frente ao comportamento do LORD (1989) e do LoGD.

---

## 2. Decision

Esta decisão avalia quatro alternativas para o modelo de venda e crédito de troca na Ferraria. O status desta ADR permanece como **Proposed**, pois a decisão final entre as opções cabe ao mantenedor.

### Alternativas Consideradas

#### Opção A: Apenas aviso de confirmação na compra de itens inferiores (Item 3 da Issue #40)
- **Descrição**: Adicionar uma janela modal ou confirmação de confirmação ao tentar comprar um item com menor atributo de ATK/DEF do que o equipado.
- **Prós**: Previne compras acidentais de itens inferiores.
- **Contras**: Não resolve a venda de itens nem concede crédito de troca. Rejeitada como solução isolada para esta ADR por não atender ao escopo dos itens 1 e 2.

#### Opção B: Crédito fixo de troca (*Trade-in Credit* percentual no ato da compra)
- **Descrição**: Ao comprar um novo equipamento, o valor do item equipado é calculado a um percentual fixo (ex.: 50% a 75%, como no LoGD) e abatido diretamente do custo do novo item.
- **Prós**: Implementação simples (cálculo determinístico via constante), fluxo de UI direto e sem atrito para o jogador.
- **Contras**: Não resolve a venda avulsa de equipamentos quando o jogador quer apenas converter seu item em ouro sem comprar outro.

#### Opção C: Menu explícito de Venda `[V]ender` com percentual fixo de recompra
- **Descrição**: Adicionar a tecla `[V]ender` no menu da Ferraria para vender o item equipado a uma fração fixa (ex.: 50% do valor de catálogo). Combinada com a Opção B, cobre compra com crédito e venda avulsa.
- **Prós**: Cobre tanto a venda avulsa quanto o abatimento no ato da compra; regras de cálculo simples e determinísticas.
- **Contras**: Desvia do modelo histórico do LORD 1989, que utilizava ofertas aleatórias e barganha em vez de tabela fixa.

#### Opção D (Recomendada): Recompra com valor aleatório a cada oferta e opção de recusa (Modelo LORD 1989)
- **Descrição**: O ferreiro oferece um valor de recompra flutuante (gerado aleatoriamente dentro de uma faixa percentual configurável, por exemplo, entre 40% e 80% do valor original do item). O jogador pode aceitar a oferta (vendendo o item ou usando o valor como crédito de troca) ou recusar. Caso recuse, pode tentar novamente em uma nova interação para obter uma cotação diferente.
- **Prós**: Fidelidade máxima ao comportamento do LORD original (1989). Adiciona um elemento estratégico de negociação/barganha na Ferraria.
- **Contras**: Maior complexidade de implementação. Requer gerenciamento de estado da oferta atual em `SmithScreen`, gerador de números aleatórios com faixa configurável e suíte de testes determinística com seed fixa.

### Recomendação
Recomenda-se a **Opção D** por estar alinhada ao princípio de fidelidade ao LORD (1989) adotado no projeto (ADR-0006). Fica ressaltado que a Opção D apresenta o maior esforço de implementação técnico e de estado de UI entre as alternativas apresentadas.

---

## 3. Consequences

### Ganhos e Benefícios
- **Fidelidade de Mecânica**: Alinhamento com o comportamento do LORD (1989) em caso de escolha da Opção D, ou com o LoGD em caso de escolha das Opções B/C.
- **Recuperação de Ouro**: Permite ao jogador reaproveitar parte do investimento financeiro em equipamentos antigos ao progredir de tier ou ao necessitar de ouro líquido para cura/banco.

### Trade-offs e Custos de Desenvolvimento
- **Complexidade de UI na Ferraria**: A tela `SmithScreen` (`internal/tui/screens/smith.go`) deixará de ter um fluxo direto de seleção e compra instantânea (`Enter`), passando a gerenciar o estado de cotação/oferta e confirmação de venda.
- **Gestão de Estado**: Necessidade de armazenar e validar o estado temporário da proposta do ferreiro no modelo TEA durante a navegação.

### Mudanças Estruturais Neutras
- O item 3 da issue #40 (aviso de downgrade) permanece não implementado e não bloqueado por esta ADR, devendo ser detalhado e resolvido em uma issue e ADR de acompanhamento.

---

## 4. Compliance and Verification

Para garantir o cumprimento desta decisão após a aprovação e implementação, os seguintes critérios devem ser validados:

### 1. Suíte de Testes Automatizados
- **Testes Determinísticos de Cálculo e Faixa**: Implementação de testes unitários em `internal/engine` com seed estática de gerador aleatório, validando:
  - Aplicação do crédito de troca e venda avulsa.
  - Respeito aos limites inferior e superior da faixa percentual de oferta (em caso de Opção D).
  - Atualização do saldo de ouro e remoção/substituição do item no inventário do `Player`.

### 2. Atualização da Documentação do Projeto
- **`docs/GDD.md`**: Atualizar a tabela de telas da TUI (Seção 3), corrigindo os atalhos da Ferraria (`ScreenSmith`), uma vez que o documento atual faz referência a teclas `[1..12]` inexistentes na interface de abas do Bubble Tea.
- **`docs/roadmap.md`**: Atualizar a descrição da Fase 3 para explicitar os mecanismos de venda e crédito de troca na Ferraria do Mestre Torin.
- **`docs/universo-e-lore.md`**: Atualizar o papel de Mestre Torin na Seção 3 para registrar as funções de compra, venda e negociação de equipamentos.
