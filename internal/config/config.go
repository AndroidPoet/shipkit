package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FileName = ".shipkit.yaml"

type AppConfig struct {
	Name           string
	IOSBundleID    string
	AndroidPackage string
	RevenueCat     bool
}

func Default(appName string) AppConfig {
	slug := strings.ToLower(strings.NewReplacer(" ", "", "-", "", "_", "").Replace(appName))
	if slug == "" {
		slug = "myapp"
	}
	return AppConfig{
		Name:           appName,
		IOSBundleID:    "com.company." + slug,
		AndroidPackage: "com.company." + slug,
		RevenueCat:     true,
	}
}

func Render(cfg AppConfig) string {
	revenueCat := "true"
	if !cfg.RevenueCat {
		revenueCat = "false"
	}

	return fmt.Sprintf(`app:
  name: %q
  ios_bundle_id: %q
  android_package: %q

tools:
  google_play: gpc
  revenuecat: rc
  app_store_connect: asc

release:
  android_track: internal
  ios_testflight: true
  revenuecat_enabled: %s
`, cfg.Name, cfg.IOSBundleID, cfg.AndroidPackage, revenueCat)
}

func Write(dir string, cfg AppConfig) (string, error) {
	path := filepath.Join(dir, FileName)
	if _, err := os.Stat(path); err == nil {
		return path, fmt.Errorf("%s already exists", FileName)
	} else if !os.IsNotExist(err) {
		return path, err
	}
	return path, os.WriteFile(path, []byte(Render(cfg)), 0644)
}
