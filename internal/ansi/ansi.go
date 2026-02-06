package ansi

import "fmt"

const (
	ESC    rune   = '\x1b' // ESC as rune for comparisons
	escStr string = "\x1b" // ESC as string for concatenation

	// Screen buffer
	EnterAltScreen = escStr + "[?1049h"
	ExitAltScreen  = escStr + "[?1049l"

	// Screen control
	ClearScreen = escStr + "[2J"
	CursorHome  = escStr + "[H"

	// Line control
	ClearLine = escStr + "[K"

	// Cursor visibility
	ShowCursor = escStr + "[?25h"
	HideCursor = escStr + "[?25l"

	// Cursor shapes
	CursorBlinkBar = escStr + "[5 q"
	CursorDefault  = escStr + "[0 q"

	// Colors
	Red   = escStr + "[31m"
	Green = escStr + "[32m"
	Blue  = escStr + "[34m"
	Gray  = escStr + "[90m"
	Reset = escStr + "[0m"

	// Invert colors (status bar)
	Invert = escStr + "[7m"
)

// MoveCursor returns the escape sequence to position cursor at row, col (1-indexed)
func MoveCursor(row, col int) {
	fmt.Printf("%s[%d;%dH", escStr, row, col)
}

// InitScreen initializes the terminal for editor use
func InitScreen() {
	fmt.Print(EnterAltScreen)
	fmt.Print(ClearScreen)
	fmt.Print(CursorHome)
	// fmt.Print(CursorBlinkBar)
}

// RestoreScreen restores the terminal to its previous state
func RestoreScreen() {
	fmt.Print(CursorDefault)
	fmt.Print(ExitAltScreen)
}
