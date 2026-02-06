package main

import (
	"bufio"
	"os"

	"github.com/sathiraumesh/textpad/internal/ansi"
)

const (
	KeyUp    = 1000
	KeyDown  = 1001
	KeyLeft  = 1002
	KeyRight = 1003
)

func readKey() (rune, error) {
	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}

	// checking if the frist byte is a ansi ESC
	if r == ansi.ESC {

		seq1, _, err := reader.ReadRune()
		if err != nil {
			return r, err
		}

		seq2, _, err := reader.ReadRune()
		if err != nil {
			return 0, err
		}

		// Step 4: Check if it's "[" followed by A/B/C/D

		if seq1 == '[' {
			switch seq2 {
			case 'A':
				return KeyUp, nil
			case 'B':
				return KeyDown, nil
			case 'C':
				return KeyRight, nil
			case 'D':
				return KeyLeft, nil
			}
		}
		return r, nil // Unknown escape sequence
	}
	return r, nil
}
