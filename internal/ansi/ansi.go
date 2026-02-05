package ansi

const (
	ESC = "\x1b"

	// Screen buffer
	EnterAltScreen = ESC + "[?1049h"
	ExitAltScreen  = ESC + "[?1049l"

	// Screen control
	ClearScreen = ESC + "[2J"
	CursorHome  = ESC + "[H"

	// Line control
	ClearLine = ESC + "[K"

	// Cursor visibility
	ShowCursor = ESC + "[?25h"
	HideCursor = ESC + "[?25l"

	// Cursor shapes
	CursorBlinkBar = ESC + "[5 q"
	CursorDefault  = ESC + "[0 q"

	// Colors
	Red   = ESC + "[31m"
	Green = ESC + "[32m"
	Blue  = ESC + "[34m"
	Gray  = ESC + "[90m"
	Reset = ESC + "[0m"

	// Invert colors (status bar)
	Invert = ESC + "[7m"
)
