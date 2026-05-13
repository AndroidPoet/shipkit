package install

import (
	"context"
	"fmt"
	"io"

	"github.com/AndroidPoet/shipkit/internal/runner"
)

type Tool struct {
	Name        string
	Executable  string
	BrewTap     string
	BrewFormula string
	InstallURL  string
}

var Tools = []Tool{
	{Name: "Google Play Console CLI", Executable: "gpc", BrewTap: "AndroidPoet/tap", BrewFormula: "playconsole-cli", InstallURL: "https://github.com/AndroidPoet/playconsole-cli"},
	{Name: "RevenueCat CLI", Executable: "rc", BrewTap: "AndroidPoet/tap", BrewFormula: "revenuecat-cli", InstallURL: "https://github.com/AndroidPoet/revenuecat-cli"},
	{Name: "App Store Connect CLI", Executable: "asc", BrewFormula: "asc", InstallURL: "https://github.com/rorkai/App-Store-Connect-CLI"},
}

func Run(ctx context.Context, r runner.Runner, stdout, stderr io.Writer) error {
	if _, err := r.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is required for automatic install. Install the tools manually from their GitHub repos, then run shipkit doctor")
	}

	for _, tool := range Tools {
		if _, err := r.LookPath(tool.Executable); err == nil {
			fmt.Fprintf(stdout, "✓ %s already installed (%s)\n", tool.Name, tool.Executable)
			continue
		}

		fmt.Fprintf(stdout, "Installing %s...\n", tool.Name)
		if tool.BrewTap != "" {
			if err := r.Run(ctx, stdout, stderr, "brew", "tap", tool.BrewTap); err != nil {
				return err
			}
		}
		if err := r.Run(ctx, stdout, stderr, "brew", "install", tool.BrewFormula); err != nil {
			fmt.Fprintf(stdout, "Could not install %s with Homebrew. Manual install: %s\n", tool.Name, tool.InstallURL)
			return err
		}
	}

	fmt.Fprintln(stdout, "All mobile shipping tools are installed.")
	return nil
}
