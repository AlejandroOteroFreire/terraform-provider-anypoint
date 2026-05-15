# Contributing

Thanks for your interest in contributing to `terraform-provider-anypoint`!

## Development workflow

### Prerequisites

- Go 1.20+
- Terraform 1.0+
- (optional) GitHub CLI (`gh`) for releases

### Local setup

```bash
git clone https://github.com/AlejandroOteroFreire/terraform-provider-anypoint
cd terraform-provider-anypoint
go build -o terraform-provider-anypoint
```

Configure `~/.terraformrc` to use the local build:

```hcl
provider_installation {
  dev_overrides {
    "AlejandroOteroFreire/anypoint" = "/absolute/path/to/terraform-provider-anypoint"
  }
  direct {}
}
```

With `dev_overrides` you can skip `terraform init` — Terraform picks up your local binary directly.

## Adding a new resource

1. Add the Go code:
   - `anypoint/resource_<name>.go` (CRUD + schema)
   - `anypoint/data_source_<name>.go` (optional, for read-only data source)
2. Register it in the corresponding map:
   - `anypoint/provider_resources.go`
   - `anypoint/provider_datasources.go`
3. If the underlying API requires a new Go client module, add it under `../anypoint-client-go-fork/<module>/` and:
   - Add the module to `go.work`
   - Add the `require` and `replace` directives in `go.mod`
   - Add the client to `provider_clients.go`
4. Write docs by hand (see the [Documentation](#documentation) section below)
5. Write an example in `examples/resources/<full_name>/`

## Documentation

Per-resource docs live in `docs/resources/<name>.md` and `docs/data-sources/<name>.md`. They are **handwritten** — follow the format of existing docs (e.g. [`docs/resources/private_space_association.md`](docs/resources/private_space_association.md)).

Each doc should contain:

- YAML front-matter with `page_title`, `subcategory`, `description`
- `## Example Usage` section (at least one working example)
- `## Schema` section listing Required / Optional / Read-Only fields
- `## Import` section if the resource supports import

For each resource also add a working `examples/resources/<name>/resource.tf` and `examples/resources/<name>/import.sh` (if import is supported).

The provider-level `docs/index.md` is regenerated from [`templates/index.md.tmpl`](templates/index.md.tmpl) via:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
tfplugindocs generate
```

> Long-term goal: convert handwritten resource docs into `templates/{resources,data-sources}/<name>.md.tmpl` files so `tfplugindocs` can keep the `## Schema` section in sync with the Go schema automatically.

## Releasing

This fork distributes binaries via GitHub Releases (not the public Terraform Registry).

### Cutting a release

1. Bump the version in [`CHANGELOG.md`](CHANGELOG.md) under the `[Unreleased]` section. Move the section header to the actual version.
2. Commit and push to `master`.
3. Tag and push:
   ```bash
   git tag v2.0.0
   git push origin v2.0.0
   ```
4. The [`.github/workflows/release.yml`](.github/workflows/release.yml) workflow runs `goreleaser` and creates a GitHub Release with binaries for:
   - linux/amd64, linux/386, linux/arm, linux/arm64
   - darwin/amd64, darwin/arm64
   - windows/amd64, windows/386
   - freebsd/amd64, freebsd/386, freebsd/arm, freebsd/arm64
5. Users install by downloading the zip matching their platform — see the [Installation section in the README](README.md#installation).

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** — incompatible changes (resource/data-source renames or removals, breaking schema changes)
- **MINOR** — new resources, data sources, or backward-compatible features
- **PATCH** — bug fixes only

When in doubt, err on the side of MINOR rather than PATCH.

## Code style

- Standard Go formatting (`go fmt ./...`)
- Run `go build ./...` before pushing — no broken code on `master`
- Match the existing style for schema definitions, error handling, and naming

## Issues and PRs

- Open an issue first if you're going to change something significant
- Keep PRs focused — one feature/fix per PR
- Always update `CHANGELOG.md` under `[Unreleased]` with your change
