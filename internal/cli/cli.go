package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/AndroidPoet/shipkit/internal/agent"
	"github.com/AndroidPoet/shipkit/internal/config"
	"github.com/AndroidPoet/shipkit/internal/doctor"
	"github.com/AndroidPoet/shipkit/internal/guide"
	"github.com/AndroidPoet/shipkit/internal/install"
	"github.com/AndroidPoet/shipkit/internal/launch"
	"github.com/AndroidPoet/shipkit/internal/runner"
	"github.com/AndroidPoet/shipkit/internal/workflow"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, build BuildInfo) error {
	return runWith(context.Background(), runner.ExecRunner{}, args, stdin, stdout, stderr, build)
}

func runWith(ctx context.Context, r runner.Runner, args []string, stdin io.Reader, stdout, stderr io.Writer, build BuildInfo) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return nil
	}

	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "shipkit %s (commit %s, built %s)\n", build.Version, build.Commit, build.Date)
		return nil
	case "guide":
		_, err := guide.Run(stdin, stdout)
		return err
	case "agent":
		return agent.Print(ctx, r, stdout, hasFlag(args[1:], "--json"))
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
		if hasFlag(args[1:], "--json") {
			return doctor.PrintJSON(ctx, r, stdout)
		}
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
		return launch.Print(ctx, r, stdout, hasFlag(args[1:], "--json"))
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

var errReleaseUsage = errors.New("usage: shipkit release android|ios|all [--dry-run]")

// releaseCommands maps a release target to the ordered provider commands it runs.
// Keeping it as data (rather than inline calls) lets `--dry-run` preview the exact
// commands and lets tests assert the mapping without executing anything.
func releaseCommands(target string) ([][]string, error) {
	switch target {
	case "android":
		return [][]string{{"gpc", "release", "--track", "internal"}}, nil
	case "ios":
		return [][]string{{"asc", "testflight", "upload"}}, nil
	case "all":
		android, _ := releaseCommands("android")
		ios, _ := releaseCommands("ios")
		return append(android, ios...), nil
	default:
		return nil, errReleaseUsage
	}
}

func release(ctx context.Context, r runner.Runner, args []string, stdout, stderr io.Writer) error {
	dryRun := hasFlag(args, "--dry-run")

	targets := make([]string, 0, len(args))
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			targets = append(targets, arg)
		}
	}
	if len(targets) != 1 {
		return errReleaseUsage
	}

	commands, err := releaseCommands(targets[0])
	if err != nil {
		return err
	}

	for _, command := range commands {
		if dryRun {
			fmt.Fprintf(stdout, "[dry-run] %s\n", strings.Join(command, " "))
			continue
		}
		if err := r.Run(ctx, stdout, stderr, command[0], command[1:]...); err != nil {
			return err
		}
	}
	return nil
}

func printHelp(stdout io.Writer) {
	fmt.Fprint(stdout, `Shipkit
  The release cockpit for mobile apps.

  One command surface for Google Play, App Store Connect, RevenueCat,
  and the CI glue that makes releases repeatable.

Usage:
  shipkit version            Print build information
  shipkit guide              Interactive setup guide
  shipkit agent [--json]     AI-agent-friendly local context
  shipkit install            Install gpc, rc, and asc under the hood
  shipkit init [app name]    Create .shipkit.yaml
  shipkit doctor [--json]    Check required tools
  shipkit ci github          Generate a GitHub Actions release workflow
  shipkit release android    Run the Android release flow through gpc
  shipkit release ios        Run the iOS release flow through asc
  shipkit release all        Run Android then iOS release flows
  shipkit release ... --dry-run  Print the provider commands without running them
  shipkit launch-check [--json]  Check local launch readiness

Start:
  shipkit init "My App"
  shipkit install
  shipkit doctor
  shipkit ci github

`)
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
