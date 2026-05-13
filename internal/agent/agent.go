package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AndroidPoet/shipkit/internal/doctor"
	"github.com/AndroidPoet/shipkit/internal/runner"
)

type Context struct {
	SchemaVersion string       `json:"schema_version"`
	Goal          string       `json:"goal"`
	Tools         []ToolStatus `json:"tools"`
	Config        ConfigStatus `json:"config"`
	NextActions   []string     `json:"next_actions"`
}

type ToolStatus struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	Installed  bool   `json:"installed"`
	Path       string `json:"path,omitempty"`
	InstallURL string `json:"install_url"`
}

type ConfigStatus struct {
	File    string `json:"file"`
	Present bool   `json:"present"`
}

func BuildContext(ctx context.Context, r runner.Runner) Context {
	results := doctor.Check(ctx, r)
	tools := make([]ToolStatus, 0, len(results))
	missing := false
	for _, result := range results {
		status := ToolStatus{
			Name:       result.Tool.Name,
			Command:    result.Tool.Executable,
			Installed:  result.Ready,
			InstallURL: result.Tool.InstallURL,
		}
		if result.Ready {
			status.Path = result.Message
		} else {
			missing = true
		}
		tools = append(tools, status)
	}

	configPresent := false
	if _, err := os.Stat(".shipkit.yaml"); err == nil {
		configPresent = true
	}

	nextActions := []string{}
	if !configPresent {
		nextActions = append(nextActions, `shipkit init "My App"`)
	}
	if missing {
		nextActions = append(nextActions, "shipkit install")
	}
	nextActions = append(nextActions, "shipkit doctor")
	if !missing {
		nextActions = append(nextActions, "gpc setup --auto", "rc login", "asc auth login", "shipkit ci github")
	}

	return Context{
		SchemaVersion: "1",
		Goal:          "Make mobile release setup deterministic for humans and AI agents.",
		Tools:         tools,
		Config: ConfigStatus{
			File:    ".shipkit.yaml",
			Present: configPresent,
		},
		NextActions: nextActions,
	}
}

func Print(ctx context.Context, r runner.Runner, stdout io.Writer, jsonOutput bool) error {
	context := BuildContext(ctx, r)
	if jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(context)
	}

	fmt.Fprintln(stdout, "Shipkit Agent Context")
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Tools:")
	for _, tool := range context.Tools {
		if tool.Installed {
			fmt.Fprintf(stdout, "- %s (%s): installed at %s\n", tool.Name, tool.Command, tool.Path)
		} else {
			fmt.Fprintf(stdout, "- %s (%s): missing\n", tool.Name, tool.Command)
		}
	}
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "Config: %s present=%t\n", context.Config.File, context.Config.Present)
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Next actions:")
	for _, action := range context.NextActions {
		fmt.Fprintf(stdout, "- %s\n", action)
	}
	return nil
}
