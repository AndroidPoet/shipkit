package config

import (
	"strings"
	"testing"
)

func Test_Render_includesToolCommands(t *testing.T) {
	cfg := Default("Launch Pad")

	rendered := Render(cfg)

	for _, want := range []string{"google_play: gpc", "revenuecat: rc", "app_store_connect: asc"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("Render() missing %q\n%s", want, rendered)
		}
	}
}

func Test_Default_buildsBundleIdentifiersFromAppName(t *testing.T) {
	cfg := Default("Launch Pad")

	if cfg.IOSBundleID != "com.company.launchpad" {
		t.Fatalf("IOSBundleID = %q", cfg.IOSBundleID)
	}
	if cfg.AndroidPackage != "com.company.launchpad" {
		t.Fatalf("AndroidPackage = %q", cfg.AndroidPackage)
	}
}
