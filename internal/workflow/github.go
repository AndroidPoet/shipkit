package workflow

import (
	"os"
	"path/filepath"
)

const releaseWorkflow = `name: Mobile Release

on:
  workflow_dispatch:
    inputs:
      platform:
        description: Platform to release
        required: true
        default: all
        type: choice
        options:
          - all
          - android
          - ios

jobs:
  release:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install Shipkit
        run: |
          go install github.com/AndroidPoet/shipkit/cmd/shipkit@latest
      - name: Check tools
        run: shipkit doctor
      - name: Release
        run: shipkit release "${{ inputs.platform }}"
`

func WriteGitHub(dir string) (string, error) {
	workflowDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(workflowDir, "mobile-release.yml")
	return path, os.WriteFile(path, []byte(releaseWorkflow), 0644)
}
