# Terraform Provider: DirtCloud (Fake Local Cloud)

DirtCloud is a fake local cloud provider for learning and testing Terraform. It does not provision any real infrastructure. Instead, it simulates resources (projects, instances, metadata) and is paired with a local console that looks and behaves like a real cloud so you can practice Terraform workflows safely.

- Ideal for: learning Terraform, demos, workshops, CI validation, and module development without touching real infrastructure.
- Safety: zero cost, zero side‑effects. Everything is local and disposable.

See the generated provider docs in [`docs/index.md`](file:///Users/nicolas/terraform-provider-dirt/docs/index.md#L1-L40).

## Quick Start

1. Start a local DirtCloud API:
   - Option A: Run the mock API included here:
     ```bash
     ./mock-server.py
     # Serves on http://localhost:8080
     ```
   - Option B: Run the full server (if you have it): `~/dirtcloud-server` (listens on `http://localhost:8080/v1`).

2. Use the provider in Terraform:
   ```hcl
   terraform {
     required_providers {
       dirt = {
         source = "hypertf/dirt"
       }
     }
   }

   provider "dirt" {
     endpoint = "http://localhost:8080/v1" # defaults to this
     # token   = var.dirt_token            # optional
   }
   
   # Example simulated resources and data sources
   resource "dirt_project" "demo" {
     name = "demo-project"
   }

   data "dirt_metadata" "feature_flag" {
     path = "app/features/new_ui_enabled"
   }
   ```

3. Plan/apply like normal. The provider pretends to create resources and exposes them via the API/console, but nothing real is provisioned.

Explore more examples in [`examples/`](file:///Users/nicolas/terraform-provider-dirt/examples#L1).

## What This Is (and Isn’t)

- What it is:
  - A safe sandbox to practice Terraform workflows end‑to‑end (init/plan/apply/destroy, state, data sources, drift, etc.).
  - A teaching and demo tool with a console UI that looks like a real cloud.

- What it isn’t:
  - A real cloud provider. No real infra is created.
  - A drop‑in replacement for production providers.

## Use Cases

- Terraform tooling development
  - Build/test CLIs, wrappers, and orchestrators that drive `terraform init/plan/apply/destroy` without risking real infrastructure.
  - Exercise graph construction, parallelism, locks, refresh, and dependency handling against simulated CRUD resources.
- Module development and CI
  - Validate module variable validation, `for_each`/`count` expansions, and outputs with stable fake resources.
  - Run fast “acceptance-like” tests in CI with no cloud credentials or costs.
- Lifecycle and diff behavior
  - Practice `create_before_destroy`, `prevent_destroy`, `ignore_changes`, and `replace_triggered_by` using resources that behave like real ones.
  - Verify ForceNew semantics, unknowns during planning, computed + sensitive attributes, and custom diff edge cases.
- Drift, import, and state operations
  - Reproduce drift via the console/API, then test `terraform plan` detection and remediation.
  - Exercise `terraform import`, `state mv`, `state rm`, `taint`/`untaint`, and targeted plans.
- Error handling and retries (with mock server tweaks)
  - Introduce delays, 4xx/5xx responses, or pagination quirks in the mock API to test retry/backoff and user messaging.
  - Simulate eventual consistency and long-running operations to validate timeouts and `-parallel` behaviors.
- Data source graphing
  - Model data lookups that gate resource creation and verify ordering and dependency propagation.

Compared to `null_resource`
- `null_resource` doesn’t emulate provider CRUD, read-after-write, or server-side state, which limits testing of realistic behaviors.
- Dirt simulates real API interactions, IDs, and read paths, enabling more representative end-to-end workflows.

## Configuration

- `endpoint` (string): DirtCloud API endpoint. Defaults to `http://localhost:8080/v1`. Can also be set via `DIRT_ENDPOINT`.
- `token` (string, sensitive): Optional auth token. Can also be set via `DIRT_TOKEN`.

See the full schema in the provider docs: [`docs/index.md`](file:///Users/nicolas/terraform-provider-dirt/docs/index.md#L34-L40).

## Development

- Build:
  ```bash
  make build
  ```
- Test:
  ```bash
  make test
  ```
- Lint/Format:
  ```bash
  make lint
  make fmt
  ```
- Generate docs (required after schema/example changes):
  ```bash
  make generate
  ```

Local server notes:
- The provider expects a DirtCloud API at `http://localhost:8080/v1`.
- A minimal mock server is provided at [`mock-server.py`](file:///Users/nicolas/terraform-provider-dirt/mock-server.py#L1-L110).

## License

MPL-2.0

