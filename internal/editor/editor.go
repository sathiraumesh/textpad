package editor

import (
	"os"
	"strings"

	"golang.org/x/term"
)

const DefaultScreenHeight = 80
const DefaultScreenWidth = 80

// Newline characters
const (
	LF = "\n"
)

type Editor struct {
	screenW, screenH int      // height and the width of the terminal
	filepath         string   // path to the file
	lines            [][]rune // list of
	isDirty          bool

	curX, curY int // cursor in buffer coordinates (rune index, line index)
}

func NewEditor() *Editor {
	return &Editor{
		curX:    0,
		curY:    0,
		isDirty: false,
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
	// TODO just for testing move down need to change based on number of lines and ennter
	e.curY++
}

func (e *Editor) MoveRight() {
	e.curX++
	// TODO just for testing move right need to change based length of the line and space
}

func (e *Editor) MoveLeft() {
	if e.curX > 0 {
		e.curX--
	}
}

func (e *Editor) NewFile() {
	e.lines = [][]rune{[]rune{}}
}
