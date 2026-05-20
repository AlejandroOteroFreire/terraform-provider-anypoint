# Installation Guide

How to install and use `terraform-provider-anypoint` from this fork.

This provider is distributed via [GitHub Releases](https://github.com/AlejandroOteroFreire/terraform-provider-anypoint/releases) — **not** the public Terraform Registry. There are three installation paths:

1. [End users — install from a GitHub Release](#1-end-users--install-from-a-github-release)
2. [Developers — local build with `dev_overrides`](#2-developers--local-build-with-dev_overrides)
3. [CI/CD — install in a GitLab/GitHub Actions pipeline](#3-cicd--install-in-a-gitlabgithub-actions-pipeline)

Then:

- [Migrating from `mulesoft-anypoint/anypoint` 1.8.x](#migrating-from-mulesoft-anypointanypoint-18x)
- [Configuring the provider](#configuring-the-provider)
- [Troubleshooting](#troubleshooting)

---

## 1. End users — install from a GitHub Release

### Step 1: Pick the right zip for your OS/arch

Go to https://github.com/AlejandroOteroFreire/terraform-provider-anypoint/releases/latest and download the zip matching your platform:

| OS / Arch | Zip name (example for v2.0.2) |
|-----------|-------------------------------|
| macOS Apple Silicon | `terraform-provider-anypoint_2.0.2_darwin_arm64.zip` |
| macOS Intel         | `terraform-provider-anypoint_2.0.2_darwin_amd64.zip` |
| Linux x86_64        | `terraform-provider-anypoint_2.0.2_linux_amd64.zip`  |
| Linux ARM64         | `terraform-provider-anypoint_2.0.2_linux_arm64.zip`  |
| Windows x86_64      | `terraform-provider-anypoint_2.0.2_windows_amd64.zip`|

### Step 2: Install into Terraform's plugin directory

Terraform expects providers under a specific directory layout:

```
<plugin-dir>/registry.terraform.io/alejandrooterofreire/anypoint/<VERSION>/<OS_ARCH>/terraform-provider-anypoint_v<VERSION>[.exe]
```

- **macOS / Linux**: `<plugin-dir>` = `~/.terraform.d/plugins`
- **Windows**: `<plugin-dir>` = `%APPDATA%\terraform.d\plugins`

Pick the snippet that matches your OS:

#### 🍎 macOS

```bash
VERSION=2.0.2
OS_ARCH=darwin_arm64   # use darwin_amd64 on Intel Macs
NAMESPACE=alejandrooterofreire
TARGET="$HOME/.terraform.d/plugins/registry.terraform.io/${NAMESPACE}/anypoint/${VERSION}/${OS_ARCH}"

mkdir -p "$TARGET"
curl -L -o /tmp/anypoint.zip \
  "https://github.com/${NAMESPACE}/terraform-provider-anypoint/releases/download/v${VERSION}/terraform-provider-anypoint_${VERSION}_${OS_ARCH}.zip"

# (optional) verify checksum
curl -sL "https://github.com/${NAMESPACE}/terraform-provider-anypoint/releases/download/v${VERSION}/terraform-provider-anypoint_${VERSION}_SHA256SUMS" \
  | grep "${OS_ARCH}" | shasum -a 256 -c -

unzip -o /tmp/anypoint.zip -d "$TARGET"
chmod +x "$TARGET/terraform-provider-anypoint_v${VERSION}"
echo "✓ Installed to $TARGET"
```

#### 🐧 Linux

```bash
VERSION=2.0.2
OS_ARCH=linux_amd64    # use linux_arm64 on ARM hosts (e.g. AWS Graviton, RPi)
NAMESPACE=alejandrooterofreire
TARGET="$HOME/.terraform.d/plugins/registry.terraform.io/${NAMESPACE}/anypoint/${VERSION}/${OS_ARCH}"

mkdir -p "$TARGET"
curl -L -o /tmp/anypoint.zip \
  "https://github.com/${NAMESPACE}/terraform-provider-anypoint/releases/download/v${VERSION}/terraform-provider-anypoint_${VERSION}_${OS_ARCH}.zip"

# (optional) verify checksum
curl -sL "https://github.com/${NAMESPACE}/terraform-provider-anypoint/releases/download/v${VERSION}/terraform-provider-anypoint_${VERSION}_SHA256SUMS" \
  | grep "${OS_ARCH}" | sha256sum -c -

unzip -o /tmp/anypoint.zip -d "$TARGET"
chmod +x "$TARGET/terraform-provider-anypoint_v${VERSION}"
echo "✓ Installed to $TARGET"
```

#### 🪟 Windows — PowerShell

```powershell
$VERSION   = "2.0.2"
$OS_ARCH   = "windows_amd64"
$NAMESPACE = "alejandrooterofreire"
$TARGET    = "$env:APPDATA\terraform.d\plugins\registry.terraform.io\$NAMESPACE\anypoint\$VERSION\$OS_ARCH"

New-Item -ItemType Directory -Force -Path $TARGET | Out-Null

$Url = "https://github.com/$NAMESPACE/terraform-provider-anypoint/releases/download/v$VERSION/terraform-provider-anypoint_${VERSION}_$OS_ARCH.zip"
$Zip = "$env:TEMP\anypoint.zip"
Invoke-WebRequest -Uri $Url -OutFile $Zip

# (optional) verify checksum
$SumsUrl  = "https://github.com/$NAMESPACE/terraform-provider-anypoint/releases/download/v$VERSION/terraform-provider-anypoint_${VERSION}_SHA256SUMS"
$Expected = (Invoke-WebRequest -Uri $SumsUrl -UseBasicParsing).Content -split "`n" |
            Where-Object { $_ -match "$OS_ARCH" } | ForEach-Object { ($_ -split '\s+')[0] }
$Actual   = (Get-FileHash -Algorithm SHA256 $Zip).Hash.ToLower()
if ($Expected -and $Actual -ne $Expected) {
    Write-Error "Checksum mismatch (expected $Expected, got $Actual)"; exit 1
}

Expand-Archive -Path $Zip -DestinationPath $TARGET -Force
Remove-Item $Zip
Write-Host "✓ Installed to $TARGET"
```

> On Windows the provider binary is `terraform-provider-anypoint_v2.0.2.exe` (with `.exe` extension). The zip already contains it; you don't need to rename anything.

#### 🪟 Windows — Git Bash / WSL

Use the **Linux** snippet above. On WSL, `~` is the WSL home (e.g. `/home/<user>`). On Git Bash, replace `$HOME/.terraform.d/plugins` with `"$APPDATA/terraform.d/plugins"` so the binary lands where the Windows-native Terraform expects it.

### Step 3: Declare the provider in your Terraform config

```hcl
terraform {
  required_providers {
    anypoint = {
      source  = "alejandrooterofreire/anypoint"
      version = "2.0.2"
    }
  }
}

provider "anypoint" {
  client_id     = var.anypoint_client_id
  client_secret = var.anypoint_client_secret
  cplane        = "us"
}
```

### Step 4: Configure the Terraform CLI config file to use the local mirror

Even with the binary in the plugin directory, **`terraform init` hits the public Terraform Registry by default** — and since this fork isn't published there, you'll get:

```
Could not retrieve the list of available versions for provider
registry.terraform.io/alejandrooterofreire/anypoint: provider registry
registry.terraform.io does not have a provider named ...
```

To fix this, add a `filesystem_mirror` block so Terraform looks at the local plugin directory first for `alejandrooterofreire/*`. The CLI config file location depends on your OS:

| OS | Config file path |
|----|------------------|
| macOS  | `~/.terraformrc` |
| Linux  | `~/.terraformrc` |
| Windows | `%APPDATA%\terraform.rc` |

#### 🍎 macOS — `~/.terraformrc`

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/Users/YOUR_USERNAME/.terraform.d/plugins"   # absolute path required
    include = ["alejandrooterofreire/*"]
  }
  direct {
    exclude = ["alejandrooterofreire/*"]
  }
}
```

#### 🐧 Linux — `~/.terraformrc`

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/home/YOUR_USERNAME/.terraform.d/plugins"   # absolute path required
    include = ["alejandrooterofreire/*"]
  }
  direct {
    exclude = ["alejandrooterofreire/*"]
  }
}
```

#### 🪟 Windows — `%APPDATA%\terraform.rc`

> Use **forward slashes** in the path even on Windows — Terraform's HCL parser doesn't accept backslashes in string values without escaping. Otherwise you get cryptic `An argument or block definition is required here` errors.

```hcl
provider_installation {
  filesystem_mirror {
    path    = "C:/Users/YOUR_USERNAME/AppData/Roaming/terraform.d/plugins"
    include = ["alejandrooterofreire/*"]
  }
  direct {
    exclude = ["alejandrooterofreire/*"]
  }
}
```

Or generate it from PowerShell:

```powershell
$RcPath = "$env:APPDATA\terraform.rc"
@"
provider_installation {
  filesystem_mirror {
    path    = "$($env:APPDATA -replace '\\','/')/terraform.d/plugins"
    include = ["alejandrooterofreire/*"]
  }
  direct {
    exclude = ["alejandrooterofreire/*"]
  }
}
"@ | Set-Content -Encoding ASCII $RcPath
Write-Host "✓ Wrote $RcPath"
```

What this does (any OS):
- `filesystem_mirror { include = ["alejandrooterofreire/*"] }` — `alejandrooterofreire/*` providers come **only** from the local plugin directory
- `direct { exclude = ["alejandrooterofreire/*"] }` — every other provider (`hashicorp/random`, `hashicorp/time`, etc.) goes to the public registry as usual

### Step 5: Run it

```bash
# If you already had a .terraform/ from a previous version, clear it first
rm -rf .terraform .terraform.lock.hcl

terraform init -upgrade
terraform plan
```

Expected output during init:

```
- Installing alejandrooterofreire/anypoint v2.0.0...
- Installed alejandrooterofreire/anypoint v2.0.0 (unauthenticated)
```

> The `(unauthenticated)` warning is expected — we don't ship GPG-signed artifacts. The SHA256 checksums in the GitHub Release are your integrity check.

---

## 2. Developers — local build with `dev_overrides`

For contributors and people building the provider from source:

### Step 1: Clone and build

```bash
git clone https://github.com/AlejandroOteroFreire/terraform-provider-anypoint
cd terraform-provider-anypoint
go build -o terraform-provider-anypoint
```

### Step 2: Configure the Terraform CLI config file

The file is `~/.terraformrc` on macOS/Linux, `%APPDATA%\terraform.rc` on Windows.

#### 🍎🐧 macOS / Linux

```hcl
provider_installation {
  dev_overrides {
    "alejandrooterofreire/anypoint" = "/absolute/path/to/terraform-provider-anypoint"
  }
  direct {}
}
```

#### 🪟 Windows

Forward slashes even on Windows:

```hcl
provider_installation {
  dev_overrides {
    "alejandrooterofreire/anypoint" = "C:/Users/YOUR_USERNAME/code/terraform-provider-anypoint"
  }
  direct {}
}
```

With this in place, Terraform picks up your local binary directly. **You can skip `terraform init`** for any project using `alejandrooterofreire/anypoint`.

### Step 3: Use the same `required_providers` block

```hcl
terraform {
  required_providers {
    anypoint = {
      source  = "alejandrooterofreire/anypoint"
    }
  }
}
```

> Note: with `dev_overrides`, Terraform shows a warning every run — that's expected.

---

## 3. CI/CD — install in a GitLab/GitHub Actions pipeline

Both setups need to:

1. Download the provider zip and extract to the local plugin directory
2. Write a `~/.terraformrc` that uses a `filesystem_mirror` for `alejandrooterofreire/*` (otherwise `terraform init` hits the public registry and fails)

### GitLab CI (`.gitlab-ci.yml`)

```yaml
variables:
  ANYPOINT_PROVIDER_VERSION: "2.0.0"
  ANYPOINT_PROVIDER_OS_ARCH: "linux_amd64"
  ANYPOINT_PROVIDER_NAMESPACE: "alejandrooterofreire"

before_script:
  - |
    # 1. Download and install the provider binary
    TARGET="${HOME}/.terraform.d/plugins/registry.terraform.io/${ANYPOINT_PROVIDER_NAMESPACE}/anypoint/${ANYPOINT_PROVIDER_VERSION}/${ANYPOINT_PROVIDER_OS_ARCH}"
    mkdir -p "$TARGET"
    curl -L \
      "https://github.com/${ANYPOINT_PROVIDER_NAMESPACE}/terraform-provider-anypoint/releases/download/v${ANYPOINT_PROVIDER_VERSION}/terraform-provider-anypoint_${ANYPOINT_PROVIDER_VERSION}_${ANYPOINT_PROVIDER_OS_ARCH}.zip" \
      -o /tmp/anypoint.zip
    unzip -o /tmp/anypoint.zip -d "$TARGET"
    chmod +x "$TARGET/terraform-provider-anypoint_v${ANYPOINT_PROVIDER_VERSION}"

    # 2. Write ~/.terraformrc so Terraform finds the binary locally
    cat > "$HOME/.terraformrc" <<EOF
    provider_installation {
      filesystem_mirror {
        path    = "${HOME}/.terraform.d/plugins"
        include = ["${ANYPOINT_PROVIDER_NAMESPACE}/*"]
      }
      direct {
        exclude = ["${ANYPOINT_PROVIDER_NAMESPACE}/*"]
      }
    }
    EOF

stages:
  - plan

terraform-plan:
  stage: plan
  image: hashicorp/terraform:1.7
  script:
    - terraform init
    - terraform plan
```

### GitHub Actions

```yaml
- name: Install anypoint provider
  env:
    VERSION: 2.0.0
    OS_ARCH: linux_amd64
    NAMESPACE: alejandrooterofreire
  run: |
    # 1. Download + install
    TARGET="${HOME}/.terraform.d/plugins/registry.terraform.io/${NAMESPACE}/anypoint/${VERSION}/${OS_ARCH}"
    mkdir -p "$TARGET"
    curl -L "https://github.com/${NAMESPACE}/terraform-provider-anypoint/releases/download/v${VERSION}/terraform-provider-anypoint_${VERSION}_${OS_ARCH}.zip" -o /tmp/a.zip
    unzip -o /tmp/a.zip -d "$TARGET"
    chmod +x "$TARGET/terraform-provider-anypoint_v${VERSION}"

    # 2. Configure filesystem_mirror
    cat > "$HOME/.terraformrc" <<EOF
    provider_installation {
      filesystem_mirror {
        path    = "${HOME}/.terraform.d/plugins"
        include = ["${NAMESPACE}/*"]
      }
      direct {
        exclude = ["${NAMESPACE}/*"]
      }
    }
    EOF
```

---

## Migrating from `mulesoft-anypoint/anypoint` 1.8.x

If you were using the upstream `mulesoft-anypoint/anypoint` provider and want to move to this fork:

### Step 1: Update `required_providers`

```hcl
# before
required_providers {
  anypoint = {
    source  = "mulesoft-anypoint/anypoint"
    version = "1.8.2"
  }
}

# after
required_providers {
  anypoint = {
    source  = "alejandrooterofreire/anypoint"
    version = "2.0.0"
  }
}
```

### Step 2: Rewrite state to use the new provider

This is the **critical step** — without it, Terraform will try to recreate everything.

```bash
terraform state replace-provider \
  registry.terraform.io/mulesoft-anypoint/anypoint \
  registry.terraform.io/alejandrooterofreire/anypoint
```

Confirm with `yes` when prompted. The command rewrites internal state metadata only — no real resources are touched.

### Step 3: Migrate the rebrand (Flex Gateway → Self-Managed Omni Gateway)

If you have any of these resources in your state, you also need a `state mv`:

| Old name | New name |
|---|---|
| `anypoint_apim_flexgateway` | `anypoint_self_managed_omni_gateway` |
| `anypoint_secretgroup_tlscontext_flexgateway` | `anypoint_secretgroup_tlscontext_self_managed_omni_gateway` |

```bash
# Find affected resources first
terraform state list | grep -E "flexgateway"

# For each one:
terraform state mv \
  'anypoint_apim_flexgateway.example' \
  'anypoint_self_managed_omni_gateway.example'

terraform state mv \
  'anypoint_secretgroup_tlscontext_flexgateway.example' \
  'anypoint_secretgroup_tlscontext_self_managed_omni_gateway.example'
```

For `for_each` / `count` instances, repeat per key:

```bash
terraform state mv \
  'anypoint_apim_flexgateway.example["sandbox"]' \
  'anypoint_self_managed_omni_gateway.example["sandbox"]'
```

Then update your `.tf` files to use the new names too.

### Step 4: Refresh & plan

```bash
terraform init -upgrade
terraform plan
```

Expected: **No changes** (or very small drift). If you see destruction of resources, **don't apply** — recheck the state replace steps above.

---

## Configuring the provider

### Three auth modes

```hcl
# 1. Connected App (default, recommended for automation)
provider "anypoint" {
  client_id     = var.client_id
  client_secret = var.client_secret
  cplane        = "us"
}

# 2. User-on-behalf (required for Access Management resources)
provider "anypoint" {
  alias         = "admin"
  auth_type     = "user"
  client_id     = var.admin_client_id    # Connected App "acts on behalf of a user"
  client_secret = var.admin_client_secret
  username      = var.admin_username     # service account, MFA disabled
  password      = var.admin_password
  cplane        = "us"
}

# 3. Pre-signed token
provider "anypoint" {
  access_token = var.access_token
  cplane       = "us"
}
```

### Which resources need `auth_type = "user"`?

Resources that hit Access Management endpoints:

- `anypoint_team_roles`
- `anypoint_team_member`
- `anypoint_connected_app_scopes`
- `anypoint_private_space_association`

For those, set `provider = anypoint.admin` on the resource block.

See [`docs/index.md`](docs/index.md) and [`README.md`](README.md) for the full reference.

---

## Troubleshooting

### `Error: Failed to query available provider packages` / `provider registry.terraform.io does not have a provider named ...`

The provider is not on the public Registry. The error means your `~/.terraformrc` doesn't declare a `filesystem_mirror` for `alejandrooterofreire/*` so Terraform tries the registry and 404s.

**Fix**: configure the mirror (Step 4 above). The key bit:

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/Users/<YOUR_USERNAME>/.terraform.d/plugins"
    include = ["alejandrooterofreire/*"]
  }
  direct {
    exclude = ["alejandrooterofreire/*"]
  }
}
```

After saving, clear any stale init state and re-run:

```bash
rm -rf .terraform .terraform.lock.hcl
terraform init -upgrade
```

### `Provider produced inconsistent result after apply`

Usually a state drift after the rebrand migration. Run:

```bash
terraform refresh
terraform plan
```

If drift persists, double-check that `state replace-provider` succeeded and that the resource names match the new schema.

### `Error: Unable to authenticate on behalf of user (password grant)` — 401

- The Connected App is not configured for `Resource Owner Password Credentials` grant. It must be type **"App acts on behalf of a user"** with that grant enabled.
- The user has **MFA enabled** — password grant doesn't work with MFA. Use a service account without MFA.
- Wrong `cplane` (e.g., your org is in EU but you set `"us"`).

### `terraform plan` wants to recreate resources after upgrading

You likely skipped `terraform state replace-provider` — that's the only way Terraform learns the resources belong to the new provider. Roll back the `.tf` changes, run the replace-provider command, then re-update the `.tf`.

### My local build doesn't take effect

Check that `dev_overrides` in `~/.terraformrc` points to the **absolute** path of your build directory (where the `terraform-provider-anypoint` binary lives). Relative paths don't work.
