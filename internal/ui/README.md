# ui

## Visão Geral
Utilitários compartilhados de interface e componentes visuais do terminal, incluindo o HUD/barra de status persistente que renderiza os atributos vitais do jogador, pontos de vida, nível, turnos restantes e ouro.

## Interface
- `type ScreenID string`: Identificadores únicos das telas do jogo.
- `type ChangeScreenMsg`: Mensagem do Bubble Tea que sinaliza requisições de transição entre telas.
- `type PlayerUpdatedMsg`: Mensagem do Bubble Tea para propagar atualizações de dados do jogador.
- `func RenderStatusBar(p *engine.Player, width int) string`: Renderiza a barra de status ANSI/lipgloss com os atributos vitais do aventureiro em todas as telas.

## Invariantes
- A função `RenderStatusBar` adapta-se a diferentes larguras de terminal sem quebrar a formatação das linhas.

## Efeitos Colaterais
- Nenhum. Produz exclusivamente strings formatadas com códigos de escape ANSI.

## Verificação
- Executar: `go test ./internal/ui/...`
