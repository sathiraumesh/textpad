package editor

import (
	"os"
	"testing"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "editor-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmp.Name()) })

	_, err = tmp.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	return tmp.Name()
}

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantLines []string
	}{
		{
			name:      "single line no newline",
			content:   "Hello",
			wantLines: []string{"Hello"},
		},
		{
			name:      "multiple lines",
			content:   "Hello\nWorld",
			wantLines: []string{"Hello", "World"},
		},
		{
			name:      "trailing newline",
			content:   "Hello\n",
			wantLines: []string{"Hello", ""},
		},
		{
			name:      "unicode emoji",
			content:   "Go🎉\n编程\n",
			wantLines: []string{"Go🎉", "编程", ""},
		},
		{
			name:      "empty file",
			content:   "",
			wantLines: []string{""},
		},
		{
			name:      "only newlines",
			content:   "\n\n\n",
			wantLines: []string{"", "", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := createTempFile(t, tt.content)

			editor := NewEditor()
			err := editor.OpenFile(filePath)
			if err != nil {
				t.Fatalf("OpenFile failed: %v", err)
			}

			if len(editor.lines) != len(tt.wantLines) {
				t.Errorf("got %d lines, want %d", len(editor.lines), len(tt.wantLines))
				return
			}

			for i, want := range tt.wantLines {
				got := string(editor.lines[i])
				if got != want {
					t.Errorf("line %d = %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestOpenFileDoesNotExist(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
	}{
		{
			name:     "simple filename",
			filepath: "doesnt-exist.txt",
		},
		{
			name:     "path with directories",
			filepath: "/tmp/nonexistent/file.txt",
		},
		{
			name:     "filename with spaces",
			filepath: "my new file.txt",
		},
		{
			name:     "filename with unicode",
			filepath: "文件.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()

			err := editor.OpenFile(tt.filepath)
			if err != nil {
				t.Fatalf("OpenFile(%q) failed: %v", tt.filepath, err)
			}

			// Should have exactly one empty line
			if len(editor.lines) != 1 {
				t.Errorf("expected 1 line, got %d", len(editor.lines))
			}

			if len(editor.lines[0]) != 0 {
				t.Errorf("expected empty line, got %q", string(editor.lines[0]))
			}

			// Should store the filepath for later save
			if editor.filepath != tt.filepath {
				t.Errorf("filepath = %q, want %q", editor.filepath, tt.filepath)
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	editor := NewEditor()
	editor.NewFile()

	if len(editor.lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(editor.lines))
	}

	if len(editor.lines[0]) != 0 {
		t.Errorf("expected empty line, got %q", string(editor.lines[0]))
	}
}

func TestCursorPosition(t *testing.T) {
	editor := NewEditor()

	// Initial position should be (1, 1) in terminal coordinates
	row, col := editor.CursorPosition()
	if row != 1 || col != 1 {
		t.Errorf("initial position = (%d, %d), want (1, 1)", row, col)
	}
}

func TestMoveUp(t *testing.T) {
	tests := []struct {
		name   string
		startY int
		wantY  int
	}{
		{
			name:   "move up from row 0 stays at 0",
			startY: 0,
			wantY:  0,
		},
		{
			name:   "move up from row 1 goes to 0",
			startY: 1,
			wantY:  0,
		},
		{
			name:   "move up from row 5 goes to 4",
			startY: 5,
			wantY:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()
			editor.curY = tt.startY

			editor.MoveUp()

			if editor.curY != tt.wantY {
				t.Errorf("curY = %d, want %d", editor.curY, tt.wantY)
			}
		})
	}
}

func TestMoveDown(t *testing.T) {
	tests := []struct {
		name      string
		lines     [][]rune
		startY    int
		wantY     int
	}{
		{
			name:   "move down from row 0 to row 1",
			lines:  [][]rune{[]rune("line1"), []rune("line2"), []rune("line3")},
			startY: 0,
			wantY:  1,
		},
		{
			name:   "move down from middle row",
			lines:  [][]rune{[]rune("line1"), []rune("line2"), []rune("line3")},
			startY: 1,
			wantY:  2,
		},
		{
			name:   "move down from last row stays at last row",
			lines:  [][]rune{[]rune("line1"), []rune("line2"), []rune("line3")},
			startY: 2,
			wantY:  2,
		},
		{
			name:   "move down with single line stays at 0",
			lines:  [][]rune{[]rune("only line")},
			startY: 0,
			wantY:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()
			editor.lines = tt.lines
			editor.curY = tt.startY

			editor.MoveDown()

			if editor.curY != tt.wantY {
				t.Errorf("curY = %d, want %d", editor.curY, tt.wantY)
			}
		})
	}
}

func TestMoveLeft(t *testing.T) {
	tests := []struct {
		name   string
		startX int
		wantX  int
	}{
		{
			name:   "move left from col 0 stays at 0",
			startX: 0,
			wantX:  0,
		},
		{
			name:   "move left from col 1 goes to 0",
			startX: 1,
			wantX:  0,
		},
		{
			name:   "move left from col 10 goes to 9",
			startX: 10,
			wantX:  9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()
			editor.curX = tt.startX

			editor.MoveLeft()

			if editor.curX != tt.wantX {
				t.Errorf("curX = %d, want %d", editor.curX, tt.wantX)
			}
		})
	}
}

func TestMoveRight(t *testing.T) {
	tests := []struct {
		name   string
		lines  [][]rune
		startX int
		startY int
		wantX  int
	}{
		{
			name:   "move right from col 0",
			lines:  [][]rune{[]rune("Hello")},
			startX: 0,
			startY: 0,
			wantX:  1,
		},
		{
			name:   "move right from middle of line",
			lines:  [][]rune{[]rune("Hello")},
			startX: 2,
			startY: 0,
			wantX:  3,
		},
		{
			name:   "move right at end of line stays at end",
			lines:  [][]rune{[]rune("Hello")},
			startX: 5,
			startY: 0,
			wantX:  5,
		},
		{
			name:   "move right on empty line stays at 0",
			lines:  [][]rune{[]rune("")},
			startX: 0,
			startY: 0,
			wantX:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()
			editor.lines = tt.lines
			editor.curX = tt.startX
			editor.curY = tt.startY

			editor.MoveRight()

			if editor.curX != tt.wantX {
				t.Errorf("curX = %d, want %d", editor.curX, tt.wantX)
			}
		})
	}
}

func TestHandleKey(t *testing.T) {
	// Create lines for testing (10 lines, each 10 chars long)
	testLines := make([][]rune, 10)
	for i := range testLines {
		testLines[i] = []rune("0123456789")
	}

	tests := []struct {
		name     string
		key      Key
		startX   int
		startY   int
		wantX    int
		wantY    int
		wantQuit bool
	}{
		{
			name:     "KeyUp moves cursor up",
			key:      KeyUp,
			startX:   5,
			startY:   5,
			wantX:    5,
			wantY:    4,
			wantQuit: false,
		},
		{
			name:     "KeyDown moves cursor down",
			key:      KeyDown,
			startX:   5,
			startY:   5,
			wantX:    5,
			wantY:    6,
			wantQuit: false,
		},
		{
			name:     "KeyLeft moves cursor left",
			key:      KeyLeft,
			startX:   5,
			startY:   5,
			wantX:    4,
			wantY:    5,
			wantQuit: false,
		},
		{
			name:     "KeyRight moves cursor right",
			key:      KeyRight,
			startX:   5,
			startY:   5,
			wantX:    6,
			wantY:    5,
			wantQuit: false,
		},
		{
			name:     "KeyQuit signals quit",
			key:      KeyQuit,
			startX:   5,
			startY:   5,
			wantX:    5,
			wantY:    5,
			wantQuit: true,
		},
		{
			name:     "unknown key does nothing",
			key:      Key('x'),
			startX:   5,
			startY:   5,
			wantX:    5,
			wantY:    5,
			wantQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := NewEditor()
			editor.lines = testLines
			editor.curX = tt.startX
			editor.curY = tt.startY

			quit := editor.HandleKey(tt.key)

			if quit != tt.wantQuit {
				t.Errorf("quit = %v, want %v", quit, tt.wantQuit)
			}
			if editor.curX != tt.wantX {
				t.Errorf("curX = %d, want %d", editor.curX, tt.wantX)
			}
			if editor.curY != tt.wantY {
				t.Errorf("curY = %d, want %d", editor.curY, tt.wantY)
			}
		})
	}
}
