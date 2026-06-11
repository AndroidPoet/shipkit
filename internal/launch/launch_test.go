package launch

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/AndroidPoet/shipkit/internal/config"
)

type fakeRunner struct {
	paths map[string]string
}

func (f fakeRunner) LookPath(name string) (string, error) {
	if path, ok := f.paths[name]; ok {
		return path, nil
	}
	return "", errors.New("missing")
}

func (fakeRunner) Run(context.Context, io.Writer, io.Writer, string, ...string) error {
	return nil
}

func allToolsInstalled() fakeRunner {
	return fakeRunner{paths: map[string]string{
		"gpc": "/bin/gpc",
		"rc":  "/bin/rc",
		"asc": "/bin/asc",
	}}
}

func writeConfig(t *testing.T, cfg config.AppConfig) {
	t.Helper()
	t.Chdir(t.TempDir())
	if _, err := config.Write(".", cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func Test_Evaluate_readyWhenToolsAndRealIdentifiersPresent(t *testing.T) {
	writeConfig(t, config.AppConfig{
		Name:           "Launch Pad",
		IOSBundleID:    "com.acme.launchpad",
		AndroidPackage: "com.acme.launchpad",
	})

	report := Evaluate(context.Background(), allToolsInstalled())

	if !report.Ready {
		t.Fatalf("expected ready, got: %#v", report.Checks)
	}
}

func Test_Evaluate_notReadyOnPlaceholderIdentifiers(t *testing.T) {
	// config.Default leaves the com.company.* placeholder, which is not shippable.
	writeConfig(t, config.Default("Launch Pad"))

	report := Evaluate(context.Background(), allToolsInstalled())

	if report.Ready {
		t.Fatal("expected not ready while identifiers are placeholders")
	}
	if !hasFailingCheck(report, "ios bundle id") || !hasFailingCheck(report, "android package") {
		t.Fatalf("expected placeholder identifier checks to fail: %#v", report.Checks)
	}
}

func Test_Evaluate_notReadyWhenToolsMissing(t *testing.T) {
	writeConfig(t, config.AppConfig{
		Name:           "Launch Pad",
		IOSBundleID:    "com.acme.launchpad",
		AndroidPackage: "com.acme.launchpad",
	})

	report := Evaluate(context.Background(), fakeRunner{})

	if report.Ready {
		t.Fatal("expected not ready when provider tools are missing")
	}
}

func hasFailingCheck(report Report, name string) bool {
	for _, check := range report.Checks {
		if check.Name == name && !check.OK {
			return true
		}
	}
	return false
}
