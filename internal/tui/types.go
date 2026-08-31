package tui

import (
	"lotgd/internal/engine"
	"lotgd/internal/ui"
)

// Re-export identifiers from internal/ui for convenience
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

func RenderStatusBar(p *engine.Player, width int) string {
	return ui.RenderStatusBar(p, width)
}



