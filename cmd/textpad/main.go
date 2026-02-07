package main

import (
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

	ansi.InitScreen()
	defer ansi.RestoreScreen()

	for {
		row, col := ed.CursorPosition()
		ansi.MoveCursor(row, col)
		key, err := ed.ReadKey()
		if err != nil {
			break
		}

		if ed.HandleKey(key) {
			break // quit
		}

	}
}
