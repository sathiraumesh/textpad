package editor

import (
	"bufio"
	"os"
	"strings"

	"github.com/sathiraumesh/textpad/internal/ansi"
	"golang.org/x/term"
)

const DefaultScreenHeight = 80
const DefaultScreenWidth = 80

// Newline characters
const (
	LF = "\n"
)

// Key represents a keyboard input (regular character or special key)
type Key rune

// Special key codes (values above Unicode range to avoid conflicts)
const (
	KeyUp Key = iota + 1000
	KeyDown
	KeyLeft
	KeyRight
)

// Character key constants
const (
	KeyQuit Key = 'q'
)

type Editor struct {
	screenW, screenH int      // height and the width of the terminal
	filepath         string   // path to the file
	lines            [][]rune // list of
	isDirty          bool

	curX, curY int // cursor in buffer coordinates (rune index, line index)

	reader *bufio.Reader
}

func NewEditor() *Editor {
	return &Editor{
		curX:    0,
		curY:    0,
		isDirty: false,
		reader:  bufio.NewReader(os.Stdin),
	}
}

func (e *Editor) UpdateSize(fd int) {
	w, h, err := term.GetSize(fd)

	if err != nil || w <= 0 || h <= 0 {
		w, h = DefaultScreenHeight, DefaultScreenWidth
	}

	e.screenW, e.screenH = w, h
}

func (e *Editor) OpenFile(name string) error {
	b, err := os.ReadFile(name)

	if err != nil {
		if os.IsNotExist(err) {
			e.lines = [][]rune{[]rune{}}
			e.filepath = name
			return nil
		}
		return err
	}

	lines := strings.Split(string(b), LF)

	if len(lines) == 0 {
		e.lines = [][]rune{[]rune{}}
	} else {
		e.lines = make([][]rune, 0, len(lines))
		for _, line := range lines {
			e.lines = append(e.lines, []rune(line))
		}
	}
	return nil
}

// the actual start location of the cursor window is (1, 1) we use (0, 0) based initialization
func (e *Editor) CursorPosition() (row int, col int) {
	return e.curY + 1, e.curX + 1
}

func (e *Editor) MoveUp() {
	if e.curY > 0 {
		e.curY--
	}
}

func (e *Editor) MoveDown() {
	if e.curY < len(e.lines)-1 {
		e.curY++
	}
}

func (e *Editor) MoveRight() {
	if e.curX < len(e.lines[e.curY]) {
		e.curX++
	}
}

func (e *Editor) MoveLeft() {
	if e.curX > 0 {
		e.curX--
	}
}

func (e *Editor) NewFile() {
	e.lines = [][]rune{[]rune{}}
}

func (e *Editor) ReadKey() (Key, error) {
	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()

	if err != nil {
		return 0, err
	}

	// checking if the first byte is an ANSI ESC
	if r == ansi.ESC {
		seq1, _, err := reader.ReadRune()
		if err != nil {
			return Key(r), err
		}

		seq2, _, err := reader.ReadRune()
		if err != nil {
			return 0, err
		}

		// Check if it's "[" followed by A/B/C/D
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
		return Key(r), nil // Unknown escape sequence
	}
	return Key(r), nil
}

// HandleKey processes a key input and returns true if the editor should quit
func (e *Editor) HandleKey(key Key) bool {
	switch key {
	case KeyQuit:
		return true // signal to quit
	case KeyUp:
		e.MoveUp()
	case KeyDown:
		e.MoveDown()
	case KeyLeft:
		e.MoveLeft()
	case KeyRight:
		e.MoveRight()
	}
	return false
}
