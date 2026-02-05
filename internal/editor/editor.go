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

func (e *Editor) NewFile() {
	e.lines = [][]rune{[]rune{}}
}
