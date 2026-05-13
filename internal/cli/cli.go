package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AndroidPoet/shipkit/internal/config"
	"github.com/AndroidPoet/shipkit/internal/doctor"
	"github.com/AndroidPoet/shipkit/internal/install"
	"github.com/AndroidPoet/shipkit/internal/runner"
	"github.com/AndroidPoet/shipkit/internal/workflow"
)

func Run(args []string, stdout, stderr io.Writer) error {
	return runWith(context.Background(), runner.ExecRunner{}, args, stdout, stderr)
}

func runWith(ctx context.Context, r runner.Runner, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return nil
	}

	switch args[0] {
	case "install":
		return install.Run(ctx, r, stdout, stderr)
	case "init":
		appName := "MyApp"
		if len(args) > 1 {
			appName = strings.Join(args[1:], " ")
		}
		path, err := config.Write(".", config.Default(appName))
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Created %s\n", path)
		return nil
	case "doctor":
		return doctor.Print(ctx, r, stdout)
	case "ci":
		if len(args) < 2 || args[1] != "github" {
			return fmt.Errorf("usage: shipkit ci github")
		}
		path, err := workflow.WriteGitHub(".")
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Created %s\n", path)
		return nil
	case "release":
		return release(ctx, r, args[1:], stdout, stderr)
	case "launch-check":
		if err := doctor.Print(ctx, r, stdout); err != nil {
			return err
		}
		_, err := os.Stat(config.FileName)
		if err != nil {
			return fmt.Errorf("%s missing; run `shipkit init`", config.FileName)
		}
		fmt.Fprintln(stdout, "Launch config found.")
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func release(ctx context.Context, r runner.Runner, args []string, stdout, stderr io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: shipkit release android|ios|all")
	}

	switch args[0] {
	case "android":
		return r.Run(ctx, stdout, stderr, "gpc", "release", "--track", "internal")
	case "ios":
		return r.Run(ctx, stdout, stderr, "asc", "testflight", "upload")
	case "all":
		if err := release(ctx, r, []string{"android"}, stdout, stderr); err != nil {
			return err
		}
		return release(ctx, r, []string{"ios"}, stdout, stderr)
	default:
		return fmt.Errorf("usage: shipkit release android|ios|all")
	}
}

func printHelp(stdout io.Writer) {
	fmt.Fprint(stdout, `shipkit makes mobile release tooling feel like one command.

Usage:
  shipkit install            Install gpc, rc, and asc under the hood
  shipkit init [app name]    Create .shipkit.yaml
  shipkit doctor             Check required tools
  shipkit ci github          Generate a GitHub Actions release workflow
  shipkit release android    Run the Android release flow through gpc
  shipkit release ios        Run the iOS release flow through asc
  shipkit release all        Run Android then iOS release flows
  shipkit launch-check       Check local launch readiness

`)
}
