package install

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type fakeRunner struct {
	paths map[string]string
	runs  []string
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if path, ok := f.paths[name]; ok {
		return path, nil
	}
	return "", errors.New("missing")
}

func (f *fakeRunner) Run(_ context.Context, _ io.Writer, _ io.Writer, name string, args ...string) error {
	f.runs = append(f.runs, strings.Join(append([]string{name}, args...), " "))
	return nil
}

func Test_Run_installsMissingToolsWithHomebrew(t *testing.T) {
	var stdout bytes.Buffer
	r := &fakeRunner{paths: map[string]string{"brew": "/opt/homebrew/bin/brew"}}

	err := Run(context.Background(), r, &stdout, io.Discard)

	if err != nil {
		t.Fatal(err)
	}
	want := strings.Join([]string{
		"brew tap AndroidPoet/tap",
		"brew install playconsole-cli",
		"brew tap AndroidPoet/tap",
		"brew install revenuecat-cli",
		"brew install asc",
	}, "\n")
	if strings.Join(r.runs, "\n") != want {
		t.Fatalf("runs = %#v", r.runs)
	}
}
