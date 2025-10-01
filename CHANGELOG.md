## 0.4.0 (October 02, 2025)

FEATURES:
* **dirt_bucket resource**: Manage buckets with create/read/update/delete and import
* **dirt_object resource**: Manage bucket-scoped objects, base64 content handling, and import `{bucket_id}/{object_id}`

ENHANCEMENTS:
* **client**: Added Bucket and Object APIs (CRUD), nested object endpoints under `/v1/bucket/{bucket_id}/objects`
* **docs/examples**: Generated docs and added usage examples for new resources

## 0.3.2 (September 29, 2025)

ENHANCEMENTS:
* **docs**: Clarified provider purpose as a fake local cloud for Terraform learning and tooling development
* **README**: Added detailed use cases, quick start, and comparison with `null_resource`

## 0.3.1 (September 18, 2025)

BUG FIXES:
* **documentation**: Fixed provider name in documentation generation from "scaffolding" to "dirt"
* **linting**: Added missing periods to Go comments for godot linter compliance

## 0.3.0 (September 18, 2025)

FEATURES:
* **instance resource**: Added immutable field support for `image` attribute - changing image now triggers destroy/recreate operation instead of in-place update

ENHANCEMENTS:
* **provider**: Enhanced error handling and validation for immutable field updates
* **server**: Improved error messages for field validation with detailed current/requested values
* **client**: Better parsing and display of server error responses

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
