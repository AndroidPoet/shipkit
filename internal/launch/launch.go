package launch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AndroidPoet/shipkit/internal/config"
	"github.com/AndroidPoet/shipkit/internal/doctor"
	"github.com/AndroidPoet/shipkit/internal/runner"
)

// Check is a single launch-readiness signal.
type Check struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

// Report answers one question: can this app ship today?
type Report struct {
	Ready  bool    `json:"ready"`
	Checks []Check `json:"checks"`
}

// Evaluate runs the checks that can be verified locally today: required tools are
// installed, the config exists, and the store identifiers are filled in (not left as
// the generated `com.company.*` placeholder). Store/network checks remain on the
// roadmap and are intentionally not faked here.
func Evaluate(ctx context.Context, r runner.Runner) Report {
	report := Report{Ready: true}
	add := func(name string, ok bool, detail string) {
		report.Checks = append(report.Checks, Check{Name: name, OK: ok, Detail: detail})
		if !ok {
			report.Ready = false
		}
	}

	for _, result := range doctor.Check(ctx, r) {
		detail := result.Message
		if result.Ready {
			detail = "installed at " + result.Message
		} else {
			detail = "missing; run `shipkit install`"
		}
		add("tool: "+result.Tool.Name, result.Ready, detail)
	}

	if _, err := os.Stat(config.FileName); err != nil {
		add("config", false, config.FileName+" missing; run `shipkit init`")
		return report
	}
	add("config", true, config.FileName+" present")

	loaded, err := config.Load(".")
	if err != nil {
		add("config readable", false, err.Error())
		return report
	}

	add("app name", loaded.Name != "", identifierDetail("app name", loaded.Name))
	add("ios bundle id", validIdentifier(loaded.IOSBundleID), identifierDetail("ios_bundle_id", loaded.IOSBundleID))
	add("android package", validIdentifier(loaded.AndroidPackage), identifierDetail("android_package", loaded.AndroidPackage))

	return report
}

func validIdentifier(value string) bool {
	return value != "" && !strings.HasPrefix(value, "com.company.")
}

func identifierDetail(field, value string) string {
	switch {
	case value == "":
		return field + " is empty"
	case strings.HasPrefix(value, "com.company."):
		return field + " still uses the placeholder " + value
	default:
		return value
	}
}

func Print(ctx context.Context, r runner.Runner, stdout io.Writer, jsonOutput bool) error {
	report := Evaluate(ctx, r)

	if jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			return err
		}
	} else {
		for _, check := range report.Checks {
			mark := "✓"
			if !check.OK {
				mark = "✗"
			}
			fmt.Fprintf(stdout, "%s %s: %s\n", mark, check.Name, check.Detail)
		}
		fmt.Fprintln(stdout)
		if report.Ready {
			fmt.Fprintln(stdout, "Ready: this app can ship today.")
		} else {
			fmt.Fprintln(stdout, "Not ready: resolve the items marked ✗ above.")
		}
	}

	if !report.Ready {
		return fmt.Errorf("launch readiness checks failed")
	}
	return nil
}
