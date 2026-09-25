# ADR-0008: Adotar Sistema de Design ANSI Fiel à Época

- Status: Proposed
- Date: 2026-09-22
- Author(s): Carlos Mello Jr.
- Deciders: Carlos Mello Jr.

## 1. Context

### Estado Atual do Projeto
O pacote `internal/ui` utiliza bordas de caixa (`DoubleBorder()`, `RoundedBorder()`) em componentes como `TitleStyle`, `ContentBoxStyle` e `CombatBoxStyle`. Essa abordagem resulta em uma estética de formulário de banco de dados corporativo (estilo Clipper ou FoxPro) em vez de um jogo de BBS da era dos terminais.

A paleta de cores atual faz uso de `lipgloss.Color` com valores hexadecimais em TrueColor. Essa especificação é incompatível com o suporte histórico de clientes e emuladores de terminal retrô de BBS, que operam estritamente sobre a paleta ANSI padrão de 16 cores (3-bit/4-bit).

### Padrões Confirmados no Jogo Original
A análise de 26 capturas de tela reais do jogo original (*Legend of the Red Dragon*, 1989/1991-1996, versão 4.00a) estabeleceu as seguintes características de interface:

1. **Uso Seletivo de Bordas**: As bordas não estão totalmente ausentes, mas aparecem exclusivamente em duas estruturas específicas:
   - Painéis de interação com NPCs, combinando arte ASCII/retrato com a caixa de diálogo da loja (exemplo: *Abdul's Armour*).
   - Tabelas numeradas de catálogo de itens e preços (exemplo: *Weapons List*).
   Todas as demais telas (navegação, combate, narrativa, estatísticas, ranking, diálogos sociais e criação de personagem) utilizam texto plano sem bordas delimitadoras.
2. **Divisor Tracejado (`-=-=-=-=-=-`)**: O divisor tracejado posicionado sob um título centralizado é o principal elemento de identidade visual do jogo. O elemento ocorre em duas variantes dimensionais:
   - Divisor longo (extensão correspondente à largura da tela) para separação de seções inteiras sob seus títulos.
   - Divisor curto (~7 caracteres) para separação de itens individuais dentro de uma lista ou feed (exemplo: entradas das Notícias Diárias).
3. **Cor Temática por Tela e Herança de Cor**: A cor do divisor tracejado e dos elementos de ênfase varia conforme a localização. Cada tela possui uma cor temática própria (verde na maioria das telas de navegação, ciano na tela do bardo Seth Able).
4. **Cores Temáticas de Lojas**: A Loja de Armas utiliza tons ciano/azul, enquanto a Loja de Armaduras utiliza amarelo/dourado. Não existe uma paleta única aplicada uniformemente a todos os estabelecimentos.
5. **Convenção Semântica Global de Cores**:
   - Verde: texto narrativo e descrições padrão.
   - Magenta: opções e atalhos de comando entre parênteses, ex.: `(F)orest`.
   - Vermelho: dano recebido, perigos, ações inimigas e estado de morte (`Dead`).
   - Branco em negrito: destaques e resultados significativos (subida de nível, confirmação de eliminação de inimigos, valores numéricos de estatísticas).
6. **Ênfase Inline Leve (`**texto**`)**: Para frases curtas de impacto no corpo do texto (exemplo: `** Welcome to the realm, new warrior! **` ou `** ULTRA POWERFUL MOVE **`), o jogo utiliza asteriscos colados ao texto em vez de divisores de linha inteira.
7. **Cor de Identidade Visual por Jogador**: No feed de notícias, cada jogador mantém uma cor ANSI fixa atribuída ao seu nome de usuário em todas as menções (exemplo: `bene` em verde, `MischiefTMaker` em ciano, `Sleepy(Sheep)` em azul, `ECUpirate05` em magenta).
8. **Menu de Navegação em Duas Colunas**: Os hubs de navegação (como a Praça Central) organizam seus comandos em duas colunas paralelas, seguidas por uma linha sintética de atalhos e a linha de prompt no formato `Your command, {nome}? [tempo restante] :`.

### Fora de Escopo
A aplicação prática da cor de identidade por jogador depende de interfaces ainda inexistentes no projeto (Notícias Diárias, Correio de Mensagens, programadas para a release `0.0.4`). Esta ADR formaliza a especificação técnica e o contrato da função determinística de cor, porém a implementação nas telas correspondentes ocorrerá de forma progressiva com a construção dessas interfaces.

---

## 2. Decision

Ficam estabelecidas as seguintes diretrizes para o sistema de design ANSI do `internal/ui`:

1. **Paleta Restrita a 16 Cores ANSI**: Substituição de `lipgloss.Color(hex)` por `lipgloss.ANSIColor(0-15)` em todo o módulo `internal/ui`.
2. **Política de Bordas Restritiva**: Remoção de bordas genéricas das telas TUI. O uso de bordas é restrito a:
   - Componente de painel de NPC (retrato ASCII + diálogo).
   - Componente de tabela numerada de itens e preços.
   Telas de navegação, combate, narrativa, estatísticas e ranking devem utilizar texto plano.
3. **Componente `SectionDivider`**: Implementação de componente para exibição do divisor tracejado longo sob título centralizado, herdando a cor temática configurada para a tela.
4. **Componente `ItemDivider`**: Implementação de componente para exibição do divisor tracejado curto (~7 caracteres) para separação de itens dentro de listas e feeds.
5. **Sistema de Cor Temática por Tela**: Introdução do conceito de tema visual por local no `internal/ui`. A cor padrão é verde, podendo ser sobrescrita por telas específicas (exemplo: ciano para Loja de Armas, amarelo/dourado para Loja de Armaduras).
6. **Helper de Ênfase Inline (`** texto **`)**: Criação de função de formatação para frases de impacto curtas no corpo do texto com asteriscos delimitadores e cor destacada.
7. **Função Determinística de Cor por Jogador (`GetPlayerColor`)**: Implementação de função pura que mapeia um nome de usuário para um dos 16 índices ANSI via hash determinístico. O contrato da função é disponibilizado em `internal/ui`, com utilização nas telas diferida para a release `0.0.4`.
8. **Componente de Menu em Duas Colunas (`TwoColumnMenu`)**: Criação de componente de layout para hubs de navegação (como a Praça Central), organizando opções em duas colunas com rodapé de atalhos e linha de prompt padronizada.

---

## 3. Consequences

### Ganhos e Aprendizados
- Aumento da fidelidade visual em relação ao *Legend of the Red Dragon* (v4.00a).
- Registro do aprendizado metodológico de que análises fundamentadas em evidências diretas de capturas de tela do software original possuem maior precisão técnica do que premissas baseadas em memória ou abstrações genéricas.

### Trade-offs e Custos de Arquitetura
- Aumento da quantidade de componentes fundamentais a serem construídos em `internal/ui` (`SectionDivider`, `ItemDivider`, `NPCPanel`, `ItemTable`, `TwoColumnMenu`) antes da migração visual das telas.
- Necessidade de adicionar o conceito de tema visual por tela ao módulo `internal/ui`.
- Disponibilização da função `GetPlayerColor` sem chamadores ativos na release `0.0.3`, exigindo documentação para evitar sinalização como código morto em análises estáticas.

---

## 4. Compliance and Verification

1. **Testes de Regressão da Paleta ANSI**: Suíte de testes unitários em `internal/ui` garantindo que todos os estilos e componentes utilizem exclusivamente tipos `lipgloss.ANSIColor` com valores entre 0 e 15.
2. **Teste Unitário Determinístico para `GetPlayerColor`**: Validação automatizada confirmando que a função retorna consistentemente o mesmo índice de cor ANSI para um determinado nome de jogador.
3. **Checklist de Inspeção de Bordas por Tela**: Cada Pull Request referente à conversão de telas de UI deve ser auditado para confirmar:
   - Ausência total de bordas em telas que não sejam painéis de NPC ou tabelas numeradas.
   - Emprego dos componentes padrão em telas que exijam painéis de NPC ou tabelas numeradas.
4. **Verificação de Divisores**: Testes de renderização para assegurar o uso correto e não intercambiável de `SectionDivider` (seções/cabeçalhos) e `ItemDivider` (itens em listas).
