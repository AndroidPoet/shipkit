package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

// Loaded holds the identifiers Shipkit reads back from a .shipkit.yaml.
type Loaded struct {
	Name           string
	IOSBundleID    string
	AndroidPackage string
}

// Load reads the shipkit-generated config. It is a deliberately small, dependency-free
// reader for the flat structure Render writes — not a general YAML parser. The three
// keys it extracts are unique in the file, so a line scan is sufficient and avoids
// pulling a YAML dependency into an audit-friendly tool.
func Load(dir string) (Loaded, error) {
	data, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		return Loaded{}, err
	}

	var loaded Loaded
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "name:"):
			loaded.Name = parseValue(trimmed[len("name:"):])
		case strings.HasPrefix(trimmed, "ios_bundle_id:"):
			loaded.IOSBundleID = parseValue(trimmed[len("ios_bundle_id:"):])
		case strings.HasPrefix(trimmed, "android_package:"):
			loaded.AndroidPackage = parseValue(trimmed[len("android_package:"):])
		}
	}
	return loaded, nil
}

func parseValue(raw string) string {
	value := strings.TrimSpace(raw)
	if unquoted, err := strconv.Unquote(value); err == nil {
		return unquoted
	}
	return strings.Trim(value, `"`)
}

func Write(dir string, cfg AppConfig) (string, error) {
	path := filepath.Join(dir, FileName)
	if _, err := os.Stat(path); err == nil {
		return path, fmt.Errorf("%s already exists", FileName)
	} else if !os.IsNotExist(err) {
		return path, err
	}
	return path, os.WriteFile(path, []byte(Render(cfg)), 0o644)
}
