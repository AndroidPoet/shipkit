package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteGitHub(t *testing.T) {
	dir := t.TempDir()

	path, err := WriteGitHub(dir)
	if err != nil {
		t.Fatalf("WriteGitHub: %v", err)
	}

	want := filepath.Join(dir, ".github", "workflows", "mobile-release.yml")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	// The readiness gate must install the provider CLIs before doctor runs,
	// otherwise doctor fails and the gated release job is always skipped.
	if strings.Count(content, "shipkit install") != 2 {
		t.Errorf("expected `shipkit install` in both jobs, got:\n%s", content)
	}
	for _, want := range []string{"shipkit doctor", "shipkit release"} {
		if !strings.Contains(content, want) {
			t.Errorf("workflow missing %q", want)
		}
	}
}
