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
