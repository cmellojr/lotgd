# Reformulação Visual ANSI (v0.0.3) — The Legend of the Go Dragon

* **Author:** Carlos Mello Jr
* **Date:** 2026-11-17

## Overview

A versão `0.0.3` de *The Legend of the Go Dragon* implementa uma reformulação visual completa da interface em modo texto (TUI), resgatando a fidelidade estética e cromática dos clássicos *door games* de BBS dos anos 1990 (notadamente *The Legend of the Red Dragon / LORD*). O objetivo é transformar o visual da TUI em uma experiência imersiva nostálgica de alto contraste, mantendo 100% da integridade das regras de jogo, sem alterar structs de domínio em `internal/engine` ou quebrar o isolamento de localização em `internal/i18n`.

A interface opera em uma resolução virtual canônica de **80 colunas por 24 linhas (80x24)** com viewport centralizado, compatível com sessões SSH remotas (`charmbracelet/wish`) e terminais locais (`cmd/lotgd`), suportando navegação mista transparente (atalhos diretos clássicos sem confirmação combinados com cursor interativo por setas e vi-keys).

## Analysis

### 1. Resolução Virtual 80x24 e Centralização no Viewport
- **Comportamento em Terminais Largos (>= 80x24)**: O canvas da tela é composto rigorosamente dentro de uma caixa delimitadora de 80 colunas por 24 linhas. Quando o terminal do jogador for maior (ex: 120x40 ou 200x60), o `MainModel` utiliza `lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, canvas80x24)` para manter o jogo perfeitamente centralizado com fundo escuro limpo, sem distorcer as proporções da arte ANSI.
- **Comportamento em Terminais Estreitos (< 80 colunas ou < 24 linhas)**: Renderiza um banner retrô informando a necessidade de redimensionamento (`[ATENÇÃO: Terminal requer no mínimo 80 colunas por 24 linhas]`), prevenindo quebras de linha espúrias (*line wrapping*).

### 2. Paleta Canônica ANSI 16 Cores de Alto Contraste
Substituição dos códigos hexadecimais de 24-bit atuais em `internal/ui` pela paleta canônica de 16 cores do DOS/BBS via índices ANSI padrão (`0` a `15`):
- **Cores Padrão (0–7)**: Preto, Vermelho, Verde, Amarelo/Marrom, Azul, Magenta, Ciano, Cinza Claro.
- **Cores de Alta Intensidade (8–15)**: Cinza Escuro, Vermelho Brilhante, Verde Brilhante, Amarelo Brilhante, Azul Brilhante, Magenta Brilhante, Ciano Brilhante, Branco Puro.
- **Hierarquia Visual**:
  - Destaques de atalho e comandos: Magenta Brilhante / Amarelo Brilhante.
  - Títulos e nomes de locais: Amarelo Brilhante / Dourado.
  - Texto narrativo e diálogos: Verde Brilhante ou Cinza Claro.
  - Avisos de combate, sangue e perigo: Vermelho Brilhante.
  - Tabelas de preços e estatísticas: Ciano / Branco.

### 3. Componentes Reutilizáveis e Blocos Gráficos ANSI (`internal/ui`)
- **Bordas em Blocos e Degradês**: Componentes de moldura utilizando caracteres de bloco CP437/Unicode (`█`, `▀`, `▄`, `▌`, `▐`) e texturas sombreadas (`▓`, `▒`, `░`).
- **Divisores Rítmicos**: Utilitário `ui.Divider(width int)` gerando réplicas do padrão LORD `-=*=-=*=-=*=-` com truncamento exato.
- **Layout de Duas Colunas (`TwoColumnLayout`)**:
  - **Coluna Esquerda (24 colunas)**: Arte ANSI temática do local (elmo de armadura na ferraria, árvore na floresta, cálice na capela, estátua na praça, dragão no covil).
  - **Coluna Direita (54 colunas)**: Diálogo, tabela pontilhada de opções (`1. Stick ................ 200`) e histórico recente.
- **Medidores de Barra em Blocos**: Indicadores de vida (HP) e combate no formato `[████████░░] 80%`.

### 4. Navegação Mista Concorrente (Sem Conflitos)
- **Atalhos Diretos Clássicos**: Teclas alfanuméricas (`F`, `T`, `C`, `M`, `G`, `D`, `1`–`9`, `Q`, `0`) disparam a ação ou transição imediatamente, sem a necessidade de teclar `Enter`.
- **Navegação Interativa com Cursor**: Teclas direcionais (`↑`/`↓` e `k`/`j`) movem o seletor visual destacado por `► [F]loresta Sombria` com fundo destacado, executando a ação ao teclar `Enter` ou `Espaço`.

## Goals vs. Non-Goals

* **In-Scope Goals:**
  - Redesenhar todas as telas em `internal/tui/screens` no padrão visual split-screen 80x24 com paleta de 16 cores ANSI.
  - Implementar no pacote `internal/ui` as primitivas visuais: paleta ANSI 16, bordas de blocos `█▓▒░`, divisores `-=*=-`, barras de HP e centralizador de canvas.
  - Implementar navegação mista transparente (atalhos diretos instantâneos + cursor com setas/vi-keys).
  - Integrar artes em blocos ANSI para os locais principais (Vilarejo, Ferraria, Capela, Taverna, Guilda, Floresta e Covil do Dragão).
  - Garantir estabilidade de redimensionamento e centralização nas sessões SSH do `cmd/server` e terminal local `cmd/lotgd`.
  - Cobrir o pacote `internal/ui` com testes unitários focados em renderização de layout e limites de largura.

* **Out-of-Scope Non-Goals:**
  - Não alterar regras de negócio, tabelas de XP ou fórmulas de combate em `internal/engine`.
  - Não adicionar dependências externas de renderização gráfica além de Bubble Tea e Lip Gloss.
  - Não implementar parser de arquivos binários `.ans` externos nesta fase (usar strings Go embutidas).
  - Não alterar o esquema do banco de dados SQLite (`internal/storage`).
  - Recursos de PvP assíncrono, clãs e ranking persistente (reservados para a versão `0.0.4`).

## Alternatives Weighed

- **Abordagem Fluida / Responsiva sem Limite Fixo**: Rejeitada. Esticar artes ANSI e layouts de BBS em terminais de 140+ colunas distorce os desenhos em blocos, quebra os alinhamentos pontilhados e descaracteriza a identidade de *door game*. A abordagem canônica de viewport virtual 80x24 centralizado preserva o aspecto original em qualquer resolução.
- **Apenas Atalhos Clássicos por Linha de Comando**: Rejeitada. Forçar exclusivamente digitação direta alienaria jogadores modernos de TUI. A navegação mista simultânea atende aos dois perfis de jogadores sem causar conflitos de entrada.
- **Parser de Arquivos `.ans` com Sauce Metadata**: Rejeitado por YAGNI na fase `0.0.3`. Carregar arquivos externos aumentaria a complexidade de empacotamento do binário e abriria brechas de escape sequence. Constantes Go estruturadas embutidas atendem com simplicidade e tipagem estrita.

## Key Decisions & Staging for Next Steps

1. **Decisão 1**: Adotar 80x24 como resolução canônica da aplicação com letterbox centralizado via `lipgloss.Place`.
2. **Decisão 2**: Migrar a paleta em `internal/ui` para os 16 códigos ANSI clássicos do DOS/BBS.
3. **Decisão 3**: Estruturar o plano de execução em 4 etapas sequenciais:
   - **Etapa 1**: Primitivas gráficas, paleta e bordas de bloco em `internal/ui`.
   - **Etapa 2**: Refatoração da Praça do Vilarejo (`TownScreen`) como tela-piloto da navegação mista e split 80x24.
   - **Etapa 3**: Refatoração das telas de comércio e combate (`SmithScreen`, `TavernScreen`, `ChapelScreen`, `ForestScreen`, `DragonScreen`, `GuildScreen`).
   - **Etapa 4**: Centralização no `MainModel`, testes de redimensionamento via SSH e validação visual de ponta a ponta.

