## 0.2.2 (September 17, 2025)

BUG FIXES:
* **ci**: Fixed release workflow to not require passphrase for GPG key import

## 0.2.1 (September 17, 2025)

BUG FIXES:
* **linting**: Fixed golangci-lint configuration compatibility issues
* **client**: Properly handle `resp.Body.Close()` errors to satisfy errcheck linter

ENHANCEMENTS:
* **development**: Added golangci-lint v2 installation instructions to AGENTS.md
* **ci**: Updated golangci-lint configuration to use version 2 schema

## 0.2.0 (September 17, 2025)

FEATURES:
* **metadata data source**: Changed to use `path` instead of `id` as the required parameter

ENHANCEMENTS:
* **client**: Added `GetMetadataByPath` method for path-based metadata retrieval

## 0.1.0 (Unreleased)

FEATURES:
