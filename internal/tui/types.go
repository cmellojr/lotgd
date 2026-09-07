package tui

import (
	"lotgd/internal/engine"
	"lotgd/internal/ui"
)

// Re-exportação de identificadores e aliases do pacote `internal/ui` para ergonomia e conveniência.
//
// Didática Go: Re-exportar tipos (`type ScreenID = ui.ScreenID`) e constantes permite que os consumidores
// do pacote `tui` acessem os identificadores das telas diretamente sem criar acoplamento circular
// entre subpacotes.
type ScreenID = ui.ScreenID

const (
	ScreenLogin    = ui.ScreenLogin
	ScreenTown     = ui.ScreenTown
	ScreenForest   = ui.ScreenForest
	ScreenTavern   = ui.ScreenTavern
	ScreenChapel   = ui.ScreenChapel
	ScreenSmith    = ui.ScreenSmith
	ScreenGuild    = ui.ScreenGuild
	ScreenDragon   = ui.ScreenDragon
	ScreenGameOver = ui.ScreenGameOver
)

type ChangeScreenMsg = ui.ChangeScreenMsg
type PlayerUpdatedMsg = ui.PlayerUpdatedMsg

// RenderStatusBar redireciona a chamada para a função de renderização no pacote `ui`.
func RenderStatusBar(p *engine.Player, width int) string {
	return ui.RenderStatusBar(p, width)
}
