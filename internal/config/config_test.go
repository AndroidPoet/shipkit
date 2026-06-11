package config

import (
	"strings"
	"testing"
)

func Test_Load_roundTripsWrittenConfig(t *testing.T) {
	dir := t.TempDir()

	if _, err := Write(dir, AppConfig{
		Name:           "Launch Pad",
		IOSBundleID:    "com.acme.launchpad",
		AndroidPackage: "com.acme.launchpad.android",
	}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != "Launch Pad" {
		t.Errorf("Name = %q", loaded.Name)
	}
	if loaded.IOSBundleID != "com.acme.launchpad" {
		t.Errorf("IOSBundleID = %q", loaded.IOSBundleID)
	}
	if loaded.AndroidPackage != "com.acme.launchpad.android" {
		t.Errorf("AndroidPackage = %q", loaded.AndroidPackage)
	}
}

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
