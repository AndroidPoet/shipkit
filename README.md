# Shipkit

Shipkit is the easy button for mobile release tooling.

It does not replace the provider CLIs. It installs and orchestrates them:

- `gpc` from [`playconsole-cli`](https://github.com/AndroidPoet/playconsole-cli)
- `rc` from [`revenuecat-cli`](https://github.com/AndroidPoet/revenuecat-cli)
- `asc` from [`App-Store-Connect-CLI`](https://github.com/rorkai/App-Store-Connect-CLI)

## Install

```bash
go install github.com/AndroidPoet/shipkit/cmd/shipkit@latest
```

Then install the underlying tools:

```bash
shipkit install
```

## Quick start

```bash
shipkit init "My App"
shipkit install
shipkit doctor
shipkit ci github
```

Shipkit creates a small `.shipkit.yaml` file and a GitHub Actions release workflow.

## Commands

```bash
shipkit install
```

Installs the required provider CLIs with Homebrew when they are missing.

```bash
shipkit init "My App"
```

Creates `.shipkit.yaml` with iOS, Android, RevenueCat, and release defaults.

```bash
shipkit doctor
```

Checks whether `gpc`, `rc`, and `asc` are available on your `PATH`.

```bash
shipkit ci github
```

Creates `.github/workflows/mobile-release.yml`.

```bash
shipkit release android
shipkit release ios
shipkit release all
```

Runs the release flow through the underlying provider CLIs.

## Design

Shipkit owns the workflow. Provider CLIs own their APIs.

That keeps this project small, useful, and maintainable:

- no duplicate Google Play API client
- no duplicate App Store Connect API client
- no duplicate RevenueCat API client
- one clean onboarding path for mobile developers

## Current status

This is an early CLI scaffold. The first goal is to make setup and release checks painless. The next useful features are:

- provider auth validation
- secret checks for CI
- product ID consistency checks across stores and RevenueCat
- richer launch readiness reports
