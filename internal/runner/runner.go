package runner

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Runner interface {
	LookPath(name string) (string, error)
	Run(ctx context.Context, stdout, stderr io.Writer, name string, args ...string) error
}

type ExecRunner struct{}

func (ExecRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (ExecRunner) Run(ctx context.Context, stdout, stderr io.Writer, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}
