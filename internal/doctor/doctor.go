package doctor

import (
	"context"
	"fmt"
	"io"

	"github.com/AndroidPoet/shipkit/internal/install"
	"github.com/AndroidPoet/shipkit/internal/runner"
)

type Result struct {
	Tool    install.Tool
	Ready   bool
	Message string
}

func Check(ctx context.Context, r runner.Runner) []Result {
	results := make([]Result, 0, len(install.Tools))
	for _, tool := range install.Tools {
		if path, err := r.LookPath(tool.Executable); err == nil {
			results = append(results, Result{Tool: tool, Ready: true, Message: path})
		} else {
			results = append(results, Result{Tool: tool, Ready: false, Message: "missing"})
		}
	}
	return results
}

func Print(ctx context.Context, r runner.Runner, stdout io.Writer) error {
	results := Check(ctx, r)
	missing := 0
	for _, result := range results {
		if result.Ready {
			fmt.Fprintf(stdout, "✓ %s: %s\n", result.Tool.Name, result.Message)
			continue
		}
		missing++
		fmt.Fprintf(stdout, "✗ %s: missing (%s)\n", result.Tool.Name, result.Tool.Executable)
	}

	if missing > 0 {
		fmt.Fprintf(stdout, "\nRun `shipkit install` or install missing tools manually.\n")
		return fmt.Errorf("%d required tools missing", missing)
	}

	fmt.Fprintln(stdout, "\nCore tools are installed. Run each provider auth flow next:")
	fmt.Fprintln(stdout, "  gpc setup --auto")
	fmt.Fprintln(stdout, "  rc login")
	fmt.Fprintln(stdout, "  asc auth login")
	return nil
}
