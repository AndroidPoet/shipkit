# Shipkit

The release cockpit for mobile apps.

Shipkit gives indie teams and mobile engineers one command surface for the release work that usually lives across Google Play Console, App Store Connect, RevenueCat, and CI.

```bash
shipkit init "My App"
shipkit install
shipkit doctor
shipkit ci github
```

That is the whole idea: one front door, with the specialist tools still doing the specialist work underneath.

Want the guided version instead?

```bash
shipkit guide
```

Want Shipkit to read your local status and explain what to do next?

```bash
shipkit ai
```

## What It Does

Shipkit installs and orchestrates the CLIs that already know how to talk to each provider:

| Area | Tool Shipkit Uses | Binary |
| --- | --- | --- |
| Android releases | [`playconsole-cli`](https://github.com/AndroidPoet/playconsole-cli) | `gpc` |
| Subscriptions and paywalls | [`revenuecat-cli`](https://github.com/AndroidPoet/revenuecat-cli) | `rc` |
| iOS releases and TestFlight | [`App-Store-Connect-CLI`](https://github.com/rorkai/App-Store-Connect-CLI) | `asc` |

Shipkit does not hide those tools. It makes them easier to set up, easier to verify, and easier to wire into a repeatable mobile release workflow.

## The Problem

Every serious mobile app eventually needs the same setup:

- a Google Play CLI for Android bundles, tracks, release notes, and store operations
- an App Store Connect CLI for TestFlight, builds, certificates, profiles, and app metadata
- RevenueCat tooling for products, offerings, paywalls, metrics, and subscription checks
- CI workflows that know which commands to run and which secrets are required
- a repeatable way to answer: "Can this app ship today?"

The tools exist. The annoying part is connecting them every time.

Shipkit is that connection layer.

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

Custom install directory:

```bash
INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://raw.githubusercontent.com/AndroidPoet/shipkit/main/install.sh)"
```

## Quick Start

Create project config:

```bash
shipkit init "My App"
```

Or use the interactive guide:

```bash
shipkit guide
```

Install the underlying mobile release CLIs:

```bash
shipkit install
```

Check local readiness:

```bash
shipkit doctor
```

Generate GitHub Actions:

```bash
shipkit ci github
```

Run a release:

```bash
shipkit release android
shipkit release ios
shipkit release all
```

## Under The Hood

`shipkit install` currently maps to:

```bash
brew tap AndroidPoet/tap
brew install playconsole-cli
brew install revenuecat-cli
brew install asc
```

Shipkit checks whether `gpc`, `rc`, and `asc` already exist on your `PATH` before installing anything.

After install, authenticate the provider CLIs directly:

```bash
gpc setup --auto
rc login
asc auth login
```

Then use Shipkit as the daily workflow layer.

## Config

`shipkit init "My App"` creates `.shipkit.yaml`:

```yaml
app:
  name: "My App"
  ios_bundle_id: "com.company.myapp"
  android_package: "com.company.myapp"

tools:
  google_play: gpc
  revenuecat: rc
  app_store_connect: asc

release:
  android_track: internal
  ios_testflight: true
  revenuecat_enabled: true
```

The first version writes sensible defaults. Future releases will read this file to drive release tracks, TestFlight behavior, RevenueCat checks, metadata paths, and CI secret validation.

## Commands

### `shipkit version`

```bash
shipkit version
```

Prints build metadata. Homebrew uses this for formula testing.

### `shipkit install`

```bash
shipkit install
```

Installs missing provider CLIs through Homebrew.

### `shipkit guide`

```bash
shipkit guide
```

Starts an interactive setup guide:

```text
Shipkit Guide
Answer a few questions and Shipkit will give you the shortest setup path.

App name [My App]:
Platforms (both/android/ios) [both]:
Use RevenueCat (yes/no) [yes]:
CI provider (github/local) [github]:
```

Then it prints the exact commands for your setup, including provider auth and release commands.

### `shipkit ai`

```bash
shipkit ai
```

Reads local Shipkit status and gives next steps.

Without an API key, it stays local and prints deterministic guidance.

With an API key, it asks OpenAI for a concise setup/release plan:

```bash
export OPENAI_API_KEY="sk-..."
shipkit ai
```

This keeps the base CLI useful without AI, while making the guided mode smarter for teams that want an interactive release copilot.

### `shipkit init`

```bash
shipkit init "My App"
```

Creates `.shipkit.yaml`.

### `shipkit doctor`

```bash
shipkit doctor
```

Checks whether the required provider CLIs are installed and shows the next auth commands.

### `shipkit ci github`

```bash
shipkit ci github
```

Creates:

```text
.github/workflows/mobile-release.yml
```

The generated workflow installs Shipkit, checks the release tooling, and runs `shipkit release`.

### `shipkit release`

```bash
shipkit release android
shipkit release ios
shipkit release all
```

Current default command mapping:

```bash
shipkit release android  # gpc release --track internal
shipkit release ios      # asc testflight upload
```

These are intentionally thin wrappers. Advanced users can still run `gpc`, `rc`, or `asc` directly whenever they need full provider control.

### `shipkit launch-check`

```bash
shipkit launch-check
```

Checks whether the local project has Shipkit config and required tooling.

The direction for this command is the most important part of the product:

```text
Can this mobile app ship today?
```

Future checks should cover:

- Android package name matches Play setup
- iOS bundle ID matches App Store Connect setup
- RevenueCat product IDs exist for both stores
- CI secrets are present
- release notes exist
- store metadata exists
- screenshots are present
- internal track or TestFlight release target is configured

## GitHub Release Setup

The repo ships with:

```text
.github/workflows/ci.yml
.github/workflows/release.yml
.goreleaser.yml
Makefile
install.sh
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

## Homebrew Formula

GoReleaser generates the formula automatically:

```yaml
brews:
  - name: shipkit
    repository:
      owner: AndroidPoet
      name: homebrew-tap
```

Expected install after release:

```bash
brew tap AndroidPoet/tap
brew install shipkit
```

Formula smoke test:

```bash
shipkit version
```

## Architecture

```text
cmd/shipkit
  main.go              binary entrypoint and version injection

internal/cli
  cli.go               command routing and user-facing behavior

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

The code intentionally stays small:

- no duplicate Google Play API client
- no duplicate RevenueCat API client
- no duplicate App Store Connect API client
- no framework-heavy CLI abstraction
- no hidden magic when direct provider commands are better

## Product Direction

Shipkit should become the release readiness layer for mobile apps.

Near-term:

- read `.shipkit.yaml` during release commands
- validate provider auth status, not just executable presence
- generate GitHub secret checklists
- add `shipkit doctor --json`
- add `shipkit launch-check --json`

High-value:

- compare RevenueCat product IDs against App Store and Play Store products
- validate Android and iOS bundle identifiers
- verify TestFlight group and Play internal track readiness
- check release notes and metadata paths
- generate a complete launch report for CI comments

## Philosophy

Shipkit owns the workflow. Provider CLIs own their APIs.

That is the architecture. That is the product.

## License

MIT
