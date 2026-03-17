# CLAUDE.md

## What is this?

Helm chart for deploying Camunda 8 load test applications (starter + worker) to Kubernetes. Used for reliability and performance testing against self-managed or SaaS Camunda Platform 8 clusters.

The load test applications themselves live in `camunda/camunda` under `load-tests/load-tester/`. This repo is only the Helm chart.

## Repository Layout

- `charts/camunda-load-tests/` — the Helm chart (templates, values, tests)
- `charts/camunda-load-tests/templates/` — starter, workers, credentials, configmap templates
- `charts/camunda-load-tests/test/` — Go + Terratest tests (golden files + property tests)
- `Makefile` — all build, test, lint, and release targets
- `CONTRIBUTING.md` — full contribution guidelines, testing patterns, and code style

## Build and Test

```bash
make go.test              # Run unit tests (validates golden files)
make go.update-golden     # Regenerate golden files after template changes
make go.fmt               # Format Go files
make go.addlicense-run    # Add Apache license headers to new .go files (CI requires this)
make helm.template        # Render all templates for inspection
helm lint charts/camunda-load-tests  # Lint the chart
```

## Key Architecture Details

**Starter** (`templates/starter.yaml`) creates process instances at a configurable rate. **Workers** (`templates/workers.yaml`) is a range-based template — one Deployment per entry in `workers.*` values. In `workers.yaml`, use `$worker`/`$workerName` from the range loop, `$.Values` for globals.

**Dual configuration:** Templates set both old HOCON `-Dapp.*` flags in `JDK_JAVA_OPTIONS` and new Spring Boot env vars (`LOAD_TESTER_*`, `CAMUNDA_CLIENT_*`) from the same Helm values. The old app reads the `-D` flags; the new Spring Boot app reads the env vars. Keep both in sync when modifying templates.

**SaaS mode** (`saas.enabled=true`): sets `CAMUNDA_CLIENT_MODE=saas`, `CAMUNDA_CLIENT_AUTH_METHOD=oidc`, pulls credentials from a Secret, omits broker URL flags.

## Working with Tests

Two test types — see `CONTRIBUTING.md` for details:

1. **Golden file tests** — compare full rendered manifests against `test/golden/*.golden.yaml`. After any template change, run `make go.update-golden` then commit the updated golden files.
2. **Property tests** — render templates, unmarshal into K8s types, assert specific values. Use for non-default configurations.

All new `.go` test files need Apache 2.0 headers — run `make go.addlicense-run` before committing.

## Conventions

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/): `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `ci:`, `build:`, `style:`. See `CONTRIBUTING.md` for values documentation pattern and Helm best practices.
