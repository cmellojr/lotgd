# Visual Overhaul ANSI v0.0.3

*   **Author:** Carlos Mello Jr
*   **Status:** Draft
*   **Date:** 2026-11-18

--------------------------------------------------------------------------------

## Overview

The `0.0.3` release of *The Legend of the Go Dragon* delivers a visual overhaul of the terminal user interface (TUI). It reproduces the visual aesthetic and color palette of 1990s BBS door games, specifically *The Legend of the Red Dragon* (LORD). The release establishes a fixed canonical virtual resolution of 80 columns by 24 lines (80x24) with a centered viewport, introduces a strict 16-color ANSI palette, renders block-art frames (`█`, `▓`, `▒`, `░`), and supports concurrent hybrid navigation (instant single-key shortcuts alongside arrow/vi-key cursor selection).

## The Problem

The initial TUI implementation relied on generic terminal styling: 24-bit TrueColor hex codes (`#D4AF37`, `#8B0000`), unbounded layouts that stretch awkwardly on modern wide displays (140+ columns), thin borders, and cursor-only list navigation. This breaks the authentic BBS door game visual identity, degrades compatibility on legacy terminal emulators or SSH clients, and causes line-wrapping defects when terminal windows resize.

The proposed solution introduces:
1. A virtual 80x24 bounding box with letterboxed centering via `lipgloss.Place` for large screens, plus an informative resize notice for screens narrower than 80 columns.
2. A strict 16-color ANSI palette (`ColorBlack` through `ColorBrightWhite`) matching DOS/BBS standards.
3. Modular block-art primitives and rhythm dividers (`-=-=-=-=-`) in `internal/ui`.
4. Concurrent hybrid navigation across all screens, allowing instant key execution without pressing Enter, while preserving arrow-key cursor selection.

## Goals

*   Establish an 80x24 canonical virtual bounding canvas rendered through `MainModel` with centered letterboxing.
*   Replace 24-bit hex color definitions in `internal/ui` with standard 16-color ANSI definitions.
*   Provide reusable visual primitives in `internal/ui`: block borders (`█▓▒░`), rhythm dividers (`-=-=-=-`), health gauges (`[████░░]`), and a two-column split layout.
*   Implement concurrent hybrid navigation (instant shortcut key or arrow cursor + Enter) across all interactive screens.
*   Maintain 100% test passing rate across the entire test suite (`go test ./...`).
*   Preserve clean architectural boundaries: zero modifications to domain rules (`internal/engine`) or localization catalogs (`internal/i18n`).

## Non-Goals

*   No modifications to game mechanics, combat math, turn limits, or SQLite schema.
*   No external binary `.ans` file parsing or external asset dependencies; all ANSI art and templates remain pure Go constants.
*   No dynamic skin or user-defined theme switching engines in this release.
*   No multi-player asynchronous features (PvP, clans, global rankings), which remain scheduled for `0.0.4`.

## Open Questions

*   None. All fundamental architectural trade-offs (80x24 bounding, 16-color ANSI palette, hybrid navigation) were resolved during the ideation and brainstorm phase.

## Critical User Journeys (CUJs)

### CUJ 1: Navigating the Town Square via Instant Shortcut
*   **Given**: A player logged into the game in a terminal larger than 80x24 (e.g., 120x40).
*   **When**: The player views the Town Square and presses `F` (or `f`).
*   **Then**: The screen instantly transitions to the Dark Forest without requiring `Enter`, and the view remains centered within the 80x24 canvas.

### CUJ 2: Navigating the Town Square via Arrow Cursor
*   **Given**: A player viewing the Town Square.
*   **When**: The player presses `↓` or `j` to highlight the Blacksmith option (`► [M]estre Torin (Ferraria)`) and presses `Enter`.
*   **Then**: The active selection indicator updates visually with high-contrast highlight, and pressing `Enter` transitions to the Smith screen.

### CUJ 3: Terminal Resize Below Minimum Resolution
*   **Given**: An active game session over SSH or local terminal.
*   **When**: The terminal window is resized below 80 columns or below 24 lines.
*   **Then**: The UI displays a clean retro message requesting a minimum 80x24 window size, preventing broken layout wraps.

--------------------------------------------------------------------------------

## Technical Plan (Architecture)

### Architectural Overview & Viewport Letterboxing

The TUI architecture preserves The Elm Architecture (TEA) pattern managed by `MainModel`. To enforce visual stability and prevent distortion, all active screens render inside a virtual fixed canvas of exactly **80 columns by 24 lines**.

```text
+-------------------------------------------------------------------------------+
| Terminal Window (e.g. 120 x 40)                                                |
|                                                                               |
|                 +-------------------------------------------+                 |
|                 | Virtual Canvas 80x24 (ANSI BBS Style)     |                 |
|                 | +---------------------------------------+ |                 |
|                 | | Compact Status Bar (HUD)              | |                 |
|                 | +---------------------------------------+ |                 |
|                 | | Left Column (24)  | Right Content(54) | |                 |
|                 | | [ANSI Block Art]  | [BBS Menu / Text] | |                 |
|                 | |                   |                   | |                 |
|                 | +---------------------------------------+ |                 |
|                 | | Rhythm Divider (-=-=-=-=-=-=-=-=-=-=-)| |                 |
|                 | | Command Prompt & Hotkeys              | |                 |
|                 | +---------------------------------------+ |                 |
|                 +-------------------------------------------+                 |
|                                                                               |
+-------------------------------------------------------------------------------+
```

When `MainModel.View()` is called:
1. If `m.width < 80 || m.height < 24`:
   Renders a standardized retro warning box centered in the available space:
   ```text
   +--------------------------------------------------+
   | AVISO: O terminal precisa de no mínimo 80x24.    |
   | Resolução atual: %dx%d. Redimensione a janela.   |
   +--------------------------------------------------+
   ```
2. If `m.width >= 80 && m.height >= 24`:
   Retrieves the active screen's 80x24 rendered string, and wraps it using:
   `lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, screenView)`

### Strict 16-Color ANSI Theme Primitives (`internal/ui`)

The package `internal/ui` replaces hex colors with `lipgloss.ANSIColor` (values `0` to `15`):

```go
const (
    // Standard 8 Colors (0-7)
    ColorBlack       = lipgloss.ANSIColor(0)
    ColorRedDark     = lipgloss.ANSIColor(1)
    ColorGreenDark   = lipgloss.ANSIColor(2)
    ColorBrownYellow = lipgloss.ANSIColor(3)
    ColorBlueDark    = lipgloss.ANSIColor(4)
    ColorMagentaDark = lipgloss.ANSIColor(5)
    ColorCyanDark    = lipgloss.ANSIColor(6)
    ColorLightGray   = lipgloss.ANSIColor(7)

    // High-Intensity 8 Colors (8-15)
    ColorDarkGray       = lipgloss.ANSIColor(8)
    ColorBrightRed      = lipgloss.ANSIColor(9)
    ColorBrightGreen    = lipgloss.ANSIColor(10)
    ColorBrightYellow   = lipgloss.ANSIColor(11)
    ColorBrightBlue     = lipgloss.ANSIColor(12)
    ColorBrightMagenta  = lipgloss.ANSIColor(13)
    ColorBrightCyan     = lipgloss.ANSIColor(14)
    ColorBrightWhite    = lipgloss.ANSIColor(15)
)
```

### Visual Components & Primitives (`internal/ui`)

1. **`ui.BlockBorder`**: Custom `lipgloss.Border` constructed using CP437/Unicode block characters:
   - Top/Bottom: `▀` and `▄`
   - Left/Right: `█` and `▌`
   - Shading: `▓`, `▒`, `░` for textured dialog accents.
2. **`ui.RhythmDivider(width int) string`**: Generates exact-width `-=-=-=-=-=-` separators in `ColorGreenDark` or `ColorDarkGray` preventing line overruns.
3. **`ui.TwoColumnLayout(leftArt string, rightContent string) string`**:
   Combines an ANSI art illustration (width: 24, height: 16) on the left with a narrative dialog and menu box (width: 54, height: 16) on the right using `lipgloss.JoinHorizontal`.
4. **`ui.HealthBar(current, max int, width int) string`**:
   Renders visual gauges in block characters:
   `[████████░░░░] 65/100` with color coding (Green > 50%, Yellow 25-50%, Red < 25%).
5. **`ui.HotkeyPrompt(prompt string, hotkeys []string) string`**:
   Formats bottom navigation bars showing active options, with hotkey characters rendered in `ColorBrightMagenta` and labels in `ColorBrightGreen`.

### Concurrent Hybrid Navigation Pattern

In all interactive screens (`TownScreen`, `SmithScreen`, `TavernScreen`, etc.), the input processing in `Update(msg tea.Msg)` supports both paradigms simultaneously without mode switching:

```go
switch msg := msg.(type) {
case tea.KeyMsg:
    switch msg.String() {
    // Arrow and vi-keys move the visual cursor
    case "up", "k":
        s.cursor = (s.cursor - 1 + len(s.options)) % len(s.options)
        return s, nil
    case "down", "j":
        s.cursor = (s.cursor + 1) % len(s.options)
        return s, nil
    case "enter", " ":
        return s, s.executeOption(s.options[s.cursor].action)

    // Instant hotkey triggers action immediately without Enter
    case "f", "F":
        return s, s.executeOption(ActionForest)
    case "t", "T":
        return s, s.executeOption(ActionTavern)
    case "m", "M":
        return s, s.executeOption(ActionSmith)
    // ...
    }
}
```

When rendering options, items format the shortcut letter in brackets with high-contrast color (e.g. `[F]loresta Sombria`), and prepend the retro cursor `► ` when highlighted by the active cursor index.

--------------------------------------------------------------------------------

## Alternatives Considered

### Alternative A: Responsive Fluid Layout Without Fixed Dimensions

*   **Pros**: Fills wide terminals completely without empty outer margins.
*   **Cons**: Stretches ASCII/ANSI art disproportionately, breaks block alignment, disrupts tabular dotted menus (`Item .............. Price`), and alters the authentic BBS door game visual identity.
*   **Reason for Rejection**: Rejected by author: Preserving the authentic 1990s BBS door game experience requires strict spatial composition. A fixed 80x24 virtual canvas centered via letterbox provides consistent visual proportions on any terminal size without line-wrap breakage.

### Alternative B: Command-Line Single-Key Shortcuts Only (Classic LORD Style)

*   **Pros**: Exact 1:1 emulation of original 1989 LORD mechanics; simplifies screen model code by removing cursor state indices.
*   **Cons**: Creates friction for modern terminal users who expect interactive cursor navigation (`↑`/`↓`, `j`/`k`, `Enter`).
*   **Reason for Rejection**: Rejected by author: Modern terminal ergonomic standards favor accessible keyboard navigation. Supporting both mechanisms simultaneously in the `Update` loop satisfies both nostalgia purists and contemporary TUI expectations without introducing input conflicts.

### Alternative C: External Binary `.ans` Files with SAUCE Metadata Parser

*   **Pros**: Allows loading standard BBS art files exported directly from ANSI drawing programs (e.g., PabloDraw, Moebius).
*   **Cons**: Introduces external file I/O dependencies, complex stateful terminal escape sequence parsing, security risks from unvalidated terminal escape injections, and packaging overhead.
*   **Reason for Rejection**: Rejected by author: YAGNI. Embedding ASCII/ANSI art as typed Go raw string constants inside `internal/ui` keeps binaries fully self-contained, eliminates runtime file missing errors, ensures compile-time safety, and avoids escape sequence parsing vulnerabilities.

--------------------------------------------------------------------------------

## Detailed Implementation Plan

### Target Files to Create/Modify

#### 1. `internal/ui/palette.go` (Create)
*   **Type / Scope**: Go source file (Constants & Color Theme).
*   **Proposed Changes**: Define the canonical 16-color ANSI constants (`ColorBlack` through `ColorBrightWhite`) using `lipgloss.ANSIColor(0)` to `lipgloss.ANSIColor(15)`. Expose semantic roles (`ColorText`, `ColorHighlight`, `ColorDanger`, `ColorNotice`, `ColorHeader`, `ColorBorder`).
*   **Rationale**: Decouples color definitions from `ui.go`, establishing strict 16-color palette boundaries without 24-bit TrueColor hex dependencies.

#### 2. `internal/ui/borders.go` (Create)
*   **Type / Scope**: Go source file (Visual Components).
*   **Proposed Changes**: Define `BlockBorder` using CP437/Unicode block glyphs (`▀`, `▄`, `█`, `▌`, `▓`, `▒`, `░`), and implement `RhythmDivider(width int) string` returning alternating `-=-=-=-=-` patterns.
*   **Rationale**: Centralizes reusable BBS-styled borders and rhythmic separators for dialogs and screen headers.

#### 3. `internal/ui/art.go` (Create)
*   **Type / Scope**: Go source file (Embedded ANSI Illustrations).
*   **Proposed Changes**: Store raw multi-line string constants for local scene illustrations formatted to fit 24 columns wide by 16 lines high:
    - `ArtTownSquare`: Village fountain and statue.
    - `ArtBlacksmith`: Anvil, hammer, and knight helmet.
    - `ArtTavern`: Mug of ale and tavern hearth.
    - `ArtChapel`: Chalice and cathedral stained glass.
    - `ArtGuild`: Crossed swords and warrior shield.
    - `ArtForest`: Spooky dead trees and foliage.
    - `ArtDragon`: Dragon head and smoking nostrils.
*   **Rationale**: Keeps all visual artwork compiled directly into the binary without external file I/O or SAUCE parsing overhead.

#### 4. `internal/ui/layout.go` (Create)
*   **Type / Scope**: Go source file (Layout Engine & Gauges).
*   **Proposed Changes**: Implement layout helper functions:
    - `TwoColumnLayout(leftArt string, rightContent string) string`: Joins a 24-col art block and 54-col content box within 80 columns.
    - `HealthBar(current, max int, width int) string`: Renders colored block health gauges (`[████░░░░] 50%`).
    - `HotkeyPrompt(prompt string, hotkeys []string) string`: Renders the bottom command line with bright hotkey letters.
    - `CenteredCanvas(content string, termWidth, termHeight int) string`: Wraps the 80x24 content inside `lipgloss.Place` or outputs the minimum dimension warning banner when below 80x24.
*   **Rationale**: Standardizes spatial composition so individual screens simply supply content without calculating terminal centering.

#### 5. `internal/ui/ui.go` (Modify)
*   **Type / Scope**: Go source file (Main UI Package).
*   **Proposed Changes**: Deprecate hex color definitions in favor of `palette.go`. Refactor `RenderStatusBar(p *engine.Player, width int)` to fit the top 80-column layout using ANSI 16 colors and block gauge indicators.
*   **Rationale**: Aligns existing status bar rendering with the new BBS aesthetic while maintaining backwards compatibility with caller signatures.

#### 6. `internal/ui/ui_test.go` (Modify)
*   **Type / Scope**: Unit Tests.
*   **Proposed Changes**: Add unit tests verifying:
    - `palette.go` color values remain within the 0–15 ANSI range.
    - `RhythmDivider` truncates exactly at the specified width without trailing line-breaks.
    - `HealthBar` handles 0, full, and negative/overflow edge cases correctly.
    - `CenteredCanvas` emits the warning message when terminal dimensions are `< 80x24` and centers properly when `>= 80x24`.
*   **Rationale**: Guarantees visual regression safety and AST integrity under Go Doctor standards.

#### 7. `internal/tui/tui.go` (Modify)
*   **Type / Scope**: TUI Main Model Coordinator.
*   **Proposed Changes**: Refactor `MainModel.View()` to pass the active screen string through `ui.CenteredCanvas(screenView, m.width, m.height)`. Remove raw unpadded rendering.
*   **Rationale**: Enforces the 80x24 virtual bounding box globally across both local terminal and SSH sessions.

#### 8. `internal/tui/screens/town.go` (Modify)
*   **Type / Scope**: Location Screen (Pilot).
*   **Proposed Changes**: Refactor layout to use `ui.TwoColumnLayout(ui.ArtTownSquare, menuBox)`. Implement hybrid navigation: instant single-key shortcuts (`F`, `T`, `C`, `M`, `G`, `D`, `S`, `Q`) and arrow cursor (`↑`/`↓`, `k`/`j`, `Enter`) with highlighted `► [X]...` indicator.
*   **Rationale**: Acts as the pilot screen validating the new BBS visual identity and hybrid navigation.

#### 9. `internal/tui/screens/smith.go`, `tavern.go`, `chapel.go`, `guild.go` (Modify)
*   **Type / Scope**: Location Screens (Town Facilities).
*   **Proposed Changes**: Apply `ui.TwoColumnLayout` with their respective ANSI art (`ArtBlacksmith`, `ArtTavern`, etc.). Implement hybrid hotkey/cursor navigation on all menus. Format shop listings with dotted dotted price alignments (`Item ................ Price`).
*   **Rationale**: Standardizes all town facility screens under the LORD aesthetic.

#### 10. `internal/tui/screens/forest.go`, `dragon.go`, `game_over.go` (Modify)
*   **Type / Scope**: Combat & Encounter Screens.
*   **Proposed Changes**: Incorporate `ui.HealthBar` for both player and adversary. Use high-contrast battle logs (Yellow for attacks, Red for damage, Green for potion heals). Bind combat actions (`[A]ttack`, `[P]otion`, `[F]lee`) directly to instant keys and cursor options.
*   **Rationale**: Brings high-contrast arcade tension and clear visual readability to turn-based battles.

#### 11. `internal/tui/screens/login.go` (Modify)
*   **Type / Scope**: Authentication Screen.
*   **Proposed Changes**: Wrap credentials dialog in `ui.BlockBorder` centered inside 80x24 with title banner.
*   **Rationale**: Ensures visual continuity starting from the very first player connection screen.

### Verification Plan

*   **Compilation**: `go build ./...` with zero compilation errors or warnings.
*   **Static Analysis & Linting**: `go vet ./...` clean without warnings.
*   **Unit Tests**: `go test ./internal/ui/... ./internal/tui/... -v` passing at 100%.
*   **Visual & Functional Verification**:
    1. Local execution (`go run ./cmd/lotgd`): Navigate town, enter forest, attack a monster, visit blacksmith, buy an item, verify instant hotkeys and arrow cursor.
    2. SSH session execution (`go run ./cmd/server` & `ssh localhost -p 2222`): Test resizing terminal window from 120x40 down to 70x20 to verify the warning screen, and expand back to verify centering.


