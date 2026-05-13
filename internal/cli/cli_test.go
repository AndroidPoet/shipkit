package cli

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

func Test_Run_doctorReportsMissingTools(t *testing.T) {
	var stdout bytes.Buffer
	r := &fakeRunner{paths: map[string]string{"gpc": "/bin/gpc"}}

	err := runWith(context.Background(), r, []string{"doctor"}, &stdout, io.Discard, BuildInfo{})

	if err == nil {
		t.Fatal("expected doctor to fail when tools are missing")
	}
	if !strings.Contains(stdout.String(), "RevenueCat CLI: missing") {
		t.Fatalf("doctor output missing RevenueCat status:\n%s", stdout.String())
	}
}

func Test_Run_releaseAllRunsAndroidThenIOS(t *testing.T) {
	var stdout bytes.Buffer
	r := &fakeRunner{}

	err := runWith(context.Background(), r, []string{"release", "all"}, &stdout, io.Discard, BuildInfo{})

	if err != nil {
		t.Fatal(err)
	}
	want := []string{"gpc release --track internal", "asc testflight upload"}
	if strings.Join(r.runs, "\n") != strings.Join(want, "\n") {
		t.Fatalf("runs = %#v", r.runs)
	}
}
