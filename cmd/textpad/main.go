package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sathiraumesh/textpad/internal/ansi"
	"github.com/sathiraumesh/textpad/internal/editor"
	"golang.org/x/term"
)

func main() {
	inFd := int(os.Stdin.Fd())

	if !term.IsTerminal(inFd) {
		slog.Error("Run this in a terminal (TTy).")
		return
	}

	// make the terminal to raw mode and gets the old state of the terminal to restore back
	inOldFd, err := term.MakeRaw(inFd)
	if err != nil {
		slog.Error("Terminal raw mode error", "err", err)
	}
	defer term.Restore(inFd, inOldFd)

	// get the stdout fd to get setup the size of the window
	oFd := int(os.Stdout.Fd())
	ed := editor.NewEditor()
	ed.UpdateSize(oFd)

	args := os.Args

	if len(args) >= 2 {
		err = ed.OpenFile(args[1])
		if err != nil {
			os.Exit(1)
		}
	} else {
		ed.NewFile()
	}

	// save the current sceeen
	fmt.Print(ansi.EnterAltScreen)

	fmt.Print(ansi.ClearScreen)
	fmt.Print(ansi.CursorHome)
	fmt.Print(ansi.CursorBlinkBar)

	defer fmt.Print(ansi.ExitAltScreen)
	defer fmt.Print(ansi.CursorDefault)

	for {
		b, err := readKey()
		if err != nil {
			break
		}

		fmt.Print(ansi.ShowCursor)
		if b == 'q' {
			return
		}
	}
}
