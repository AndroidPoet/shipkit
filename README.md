<div align="center">

# Shipkit

### The release cockpit for mobile apps

One AI-agent-friendly command surface for Google Play, App Store Connect, RevenueCat, and CI release automation.

[![CI](https://github.com/AndroidPoet/shipkit/actions/workflows/ci.yml/badge.svg)](https://github.com/AndroidPoet/shipkit/actions/workflows/ci.yml)
[![Release](https://github.com/AndroidPoet/shipkit/actions/workflows/release.yml/badge.svg)](https://github.com/AndroidPoet/shipkit/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/AndroidPoet/shipkit.svg)](https://pkg.go.dev/github.com/AndroidPoet/shipkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

```bash
shipkit guide
shipkit install
shipkit doctor
shipkit agent --json
```

</div>

---

## Why Shipkit Exists

Mobile release work is scattered across too many places.

You need one tool for Android releases, another for iOS releases, another for subscriptions, and then CI glue to make the whole thing repeatable. Every new app recreates the same setup.

Shipkit does not replace the specialist CLIs. It makes them feel like one product.

| You need | Shipkit uses | Binary |
| --- | --- | --- |
| Android release automation | [`playconsole-cli`](https://github.com/AndroidPoet/playconsole-cli) | `gpc` |
| RevenueCat products, offerings, and paywalls | [`revenuecat-cli`](https://github.com/AndroidPoet/revenuecat-cli) | `rc` |
| App Store Connect and TestFlight | [`App-Store-Connect-CLI`](https://github.com/rorkai/App-Store-Connect-CLI) | `asc` |

Shipkit owns the workflow. Provider CLIs own their APIs.

---

## The Old Way

```bash
brew tap AndroidPoet/tap
brew install playconsole-cli
brew install revenuecat-cli
brew install asc

gpc setup --auto
rc login
asc auth login

# Remember which commands go with which provider.
# Rebuild CI by hand.
# Explain all of this again to every teammate and agent.
```

## The Shipkit Way

```bash
shipkit guide
shipkit install
shipkit doctor
shipkit ci github
shipkit agent --json
```

Readable for humans. Structured for agents. Small enough to trust.

---

## Install

### Homebrew

After the first tagged release:

```bash
brew tap AndroidPoet/tap
brew install shipkit
```

### Go

```bash
go install github.com/AndroidPoet/shipkit/cmd/shipkit@latest
```

### Install Script

```bash
curl -fsSL https://raw.githubusercontent.com/AndroidPoet/shipkit/main/install.sh | sh
```

Custom directory:

```bash
INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://raw.githubusercontent.com/AndroidPoet/shipkit/main/install.sh)"
```

---

## Quick Start

Start with the guided flow:

```bash
shipkit guide
```

Example:

```text
Shipkit Guide
Answer a few questions and Shipkit will give you the shortest setup path.

App name [My App]: LaunchKit
Platforms (both/android/ios) [both]: both
Use RevenueCat (yes/no) [yes]: yes
CI provider (github/local) [github]: github

Recommended setup
  shipkit init "LaunchKit"
  shipkit install
  shipkit doctor
  shipkit ci github

Provider auth
  gpc setup --auto
  rc login
  asc auth login

Release commands
  shipkit release android
  shipkit release ios
  shipkit release all
```

Or run the commands directly:

```bash
shipkit init "LaunchKit"
shipkit install
shipkit doctor
shipkit ci github
```

---

## Commands

| Command | Purpose |
| --- | --- |
| `shipkit guide` | Interactive setup guide |
| `shipkit agent --json` | AI-agent-friendly project context |
| `shipkit install` | Install `gpc`, `rc`, and `asc` under the hood |
| `shipkit init "My App"` | Create `.shipkit.yaml` |
| `shipkit doctor` | Check local tool readiness |
| `shipkit ci github` | Generate a GitHub Actions workflow |
| `shipkit release android` | Run Android release flow through `gpc` |
| `shipkit release ios` | Run iOS release flow through `asc` |
| `shipkit release all` | Run Android then iOS release flows |
| `shipkit launch-check` | Check launch readiness |
| `shipkit version` | Print build metadata |

---

## AI-Agent-Friendly Output

Shipkit does not call an AI API. It does not need an API key. It does not send your project data anywhere.

Instead, it exposes deterministic context that AI agents, scripts, and CI can parse:

```bash
shipkit agent --json
```

Example:

```json
{
  "schema_version": "1",
  "goal": "Make mobile release setup deterministic for humans and AI agents.",
  "tools": [
    {
      "name": "Google Play Console CLI",
      "command": "gpc",
      "installed": true,
      "path": "/opt/homebrew/bin/gpc",
      "install_url": "https://github.com/AndroidPoet/playconsole-cli"
    }
  ],
  "config": {
    "file": ".shipkit.yaml",
    "present": false
  },
  "next_actions": [
    "shipkit init \"My App\"",
    "shipkit doctor",
    "gpc setup --auto",
    "rc login",
    "asc auth login",
    "shipkit ci github"
  ]
}
```

This is the useful AI integration: predictable CLI responses that external agents can understand.

---

## What `shipkit install` Does

Shipkit checks for the executable names first:

```text
gpc
rc
asc
```

If any are missing, it installs the underlying CLIs with Homebrew:

```bash
brew tap AndroidPoet/tap
brew install playconsole-cli
brew install revenuecat-cli
brew install asc
```

Then authenticate each provider CLI directly:

```bash
gpc setup --auto
rc login
asc auth login
```

Shipkit keeps the provider auth flows in the provider tools where they belong.

---

## Project Config

```bash
shipkit init "LaunchKit"
```

Creates:

```yaml
app:
  name: "LaunchKit"
  ios_bundle_id: "com.company.launchkit"
  android_package: "com.company.launchkit"

tools:
  google_play: gpc
  revenuecat: rc
  app_store_connect: asc

release:
  android_track: internal
  ios_testflight: true
  revenuecat_enabled: true
```

The config starts small on purpose. It will become the source of truth for release tracks, TestFlight behavior, RevenueCat checks, metadata paths, and CI secret validation.

---

## GitHub Actions

Generate a starter workflow:

```bash
shipkit ci github
```

Creates:

```text
.github/workflows/mobile-release.yml
```

The generated workflow:

- installs Shipkit
- checks local release tooling
- runs `shipkit release` for `android`, `ios`, or `all`

Manual dispatch input:

```yaml
platform:
  type: choice
  options:
    - all
    - android
    - ios
```

---

## Release Commands

```bash
shipkit release android
shipkit release ios
shipkit release all
```

Current mapping:

```bash
shipkit release android  # gpc release --track internal
shipkit release ios      # asc testflight upload
```

These are deliberately thin wrappers. Advanced users can always drop down to `gpc`, `rc`, or `asc` directly.

---

## Launch Readiness

```bash
shipkit launch-check
```

The product direction is to answer one question:

```text
Can this app ship today?
```

Planned checks:

- Android package name matches Play Console setup
- iOS bundle ID matches App Store Connect setup
- RevenueCat product IDs exist for both stores
- CI secrets are present
- release notes exist
- store metadata exists
- screenshots are present
- internal track or TestFlight target is configured
- output is available as text and JSON

---

## Environment

Shipkit itself does not require secrets for local status checks.

Provider tools may require their own credentials:

| Tool | Typical auth command |
| --- | --- |
| `gpc` | `gpc setup --auto` |
| `rc` | `rc login` |
| `asc` | `asc auth login` |

CI workflows should store provider credentials in GitHub Actions secrets. Shipkit will add explicit secret validation in a future release.

---

## Security Model

- No AI API calls
- No telemetry
- No provider credentials stored by Shipkit
- No hidden network calls in `shipkit agent --json`
- Provider auth remains inside the provider CLIs
- Generated CI workflows are visible files you can review

---

## Repository Release Setup

This repo ships with release automation:

```text
.github/workflows/ci.yml
.github/workflows/release.yml
.goreleaser.yml
Makefile
install.sh
LICENSE
```

Local checks:

```bash
make test
make build
```

Optional GoReleaser checks:

```bash
make release-check
make snapshot
```

Production release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow publishes GitHub release artifacts and updates the Homebrew tap.

Required GitHub secret:

```text
HOMEBREW_TAP_GITHUB_TOKEN
```

That token must be able to push to:

```text
AndroidPoet/homebrew-tap
```

---

## Architecture

```text
cmd/shipkit
  main.go              binary entrypoint and version injection

internal/cli
  cli.go               command routing and user-facing behavior

internal/agent
  agent.go             deterministic context for AI agents and scripts

internal/guide
  guide.go             interactive setup wizard

internal/install
  install.go           Homebrew-backed install orchestration

internal/doctor
  doctor.go            local tool readiness checks

internal/config
  config.go            .shipkit.yaml rendering

internal/workflow
  github.go            generated GitHub Actions workflow

internal/runner
  runner.go            command execution boundary
```

Small codebase. Clear boundaries. No duplicate provider API clients.

---

## Roadmap

- `shipkit doctor --json`
- `shipkit launch-check --json`
- provider auth validation, not only executable checks
- GitHub secret checklist generation
- release commands driven by `.shipkit.yaml`
- RevenueCat product consistency checks across iOS and Android
- store metadata and screenshot readiness checks
- CI summary comments for release readiness

---

## Philosophy

Shipkit should feel like the missing control panel for mobile release work.

It should be:

- easy for humans
- predictable for scripts
- parseable for agents
- small enough to audit
- flexible enough to step aside when the provider CLI is the better tool

---

## License

MIT
