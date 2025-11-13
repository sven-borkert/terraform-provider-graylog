# Provider Deployment Guide

This guide explains how to deploy the terraform-provider-graylog to various platforms and distribution methods, from local development to production registries.

## Table of Contents

1. [Local Development](#local-development)
2. [Official Terraform Registry](#official-terraform-registry)
3. [Private Registries](#private-registries)
4. [GitHub Releases](#github-releases)
5. [Network Mirror](#network-mirror)
6. [Cloud Storage](#cloud-storage)
7. [Checksum and Signing](#checksum-and-signing)

## Prerequisites

Before deploying the provider, ensure you have:

- Go 1.15+ installed
- The provider source code
- Built binaries for target platforms
- (Optional) GPG key for signing releases

### Building the Provider

Build for multiple platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o terraform-provider-graylog_linux_amd64 ./cmd/terraform-provider-graylog

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o terraform-provider-graylog_linux_arm64 ./cmd/terraform-provider-graylog

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o terraform-provider-graylog_windows_amd64.exe ./cmd/terraform-provider-graylog

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o terraform-provider-graylog_darwin_amd64 ./cmd/terraform-provider-graylog

# macOS ARM64 (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o terraform-provider-graylog_darwin_arm64 ./cmd/terraform-provider-graylog
```

## Local Development

For local development and testing without any registry.

### Method 1: Developer Overrides

See [graylog-terraform-import/README.md](../../graylog-terraform-import/README.md) for detailed instructions.

### Method 2: Manual Installation

```bash
# Build the provider
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

# Install to local plugin directory
VERSION=3.0.0
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
PLUGIN_DIR=~/.terraform.d/plugins/terraform-provider-graylog/graylog/$VERSION/${OS}_${ARCH}

mkdir -p $PLUGIN_DIR
cp terraform-provider-graylog $PLUGIN_DIR/
```

## Official Terraform Registry

Deploy to the public Terraform Registry at registry.terraform.io.

### Prerequisites

1. **GitHub Organization** (not personal account)
2. **Public repository** named `terraform-provider-graylog`
3. **GPG signing key** for releases
4. **Proper documentation** structure

### Step 1: Prepare Repository Structure

```
terraform-provider-graylog/
├── .goreleaser.yml
├── docs/
│   ├── index.md
│   ├── data-sources/
│   │   └── *.md
│   └── resources/
│       └── *.md
├── examples/
├── main.go
└── go.mod
```

### Step 2: Create .goreleaser.yml

```yaml
version: 2

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    mod_timestamp: '{{ .CommitTimestamp }}'
    flags:
      - -trimpath
    ldflags:
      - '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}}'
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: darwin
        goarch: '386'
    binary: '{{ .ProjectName }}_v{{ .Version }}'

archives:
  - format: zip
    name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}'
    files:
      - LICENSE
      - README.md
      - CHANGELOG.md

checksum:
  name_template: '{{ .ProjectName }}_{{ .Version }}_SHA256SUMS'
  algorithm: sha256

signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"

release:
  github:
    owner: terraform-provider-graylog  # Your GitHub org
    name: terraform-provider-graylog
  draft: false
  prerelease: auto
  name_template: "v{{ .Version }}"
  disable: false

publishers:
  - name: upload
    ids:
      - terraform-registry
    cmd: |
      echo "Release {{ .Version }} published"
    env:
      - API_TOKEN={{ .Env.GITHUB_TOKEN }}

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
```

### Step 3: Create GitHub Release

```bash
# Set up GPG
export GPG_FINGERPRINT="YOUR_GPG_FINGERPRINT"

# Tag the release
git tag v3.0.0
git push origin v3.0.0

# Create release with goreleaser
goreleaser release --clean
```

### Step 4: Register on Terraform Registry

1. Go to https://registry.terraform.io/publish/provider
2. Sign in with GitHub
3. Select your organization and repository
4. Follow the setup wizard
5. The registry will automatically detect new releases

### Step 5: User Configuration

Users can now use:

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "~> 3.0"
    }
  }
}

provider "graylog" {
  # Configuration
}
```

## Private Registries

### Terraform Cloud/Enterprise

Terraform Cloud and Enterprise include private registry functionality.

#### Publishing to Terraform Cloud

```bash
# Configure API token
export TFE_TOKEN="your-tfc-api-token"

# Upload provider
curl \
  --header "Authorization: Bearer $TFE_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @payload.json \
  https://app.terraform.io/api/v2/organizations/YOUR_ORG/registry-providers
```

#### User Configuration

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "app.terraform.io/YOUR_ORG/graylog"
      version = "~> 3.0"
    }
  }
}
```

### Artifactory

JFrog Artifactory supports Terraform providers.

#### Setup Repository

1. Create a Terraform repository in Artifactory
2. Configure authentication

#### Upload Provider

```bash
# Package the provider
VERSION=3.0.0
zip terraform-provider-graylog_${VERSION}_linux_amd64.zip terraform-provider-graylog

# Upload to Artifactory
curl -u username:password \
  -T terraform-provider-graylog_${VERSION}_linux_amd64.zip \
  "https://artifactory.company.com/artifactory/terraform-local/terraform-provider-graylog/${VERSION}/terraform-provider-graylog_${VERSION}_linux_amd64.zip"
```

#### User Configuration

```hcl
# ~/.terraformrc
provider_installation {
  network_mirror {
    url = "https://artifactory.company.com/artifactory/terraform-local/"
  }
}
```

### GitLab Terraform Registry

GitLab (13.0+) includes a Terraform module registry.

#### Publish to GitLab

```bash
# Using GitLab CI/CD
deploy:
  stage: deploy
  script:
    - |
      curl --header "PRIVATE-TOKEN: $CI_JOB_TOKEN" \
           --upload-file terraform-provider-graylog_${VERSION}_linux_amd64.zip \
           "https://gitlab.company.com/api/v4/projects/${CI_PROJECT_ID}/packages/terraform/modules/graylog/provider/${VERSION}/file"
```

#### User Configuration

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "gitlab.company.com/YOUR_GROUP/graylog"
      version = "~> 3.0"
    }
  }
}
```

## GitHub Releases

Distribute via GitHub releases without a formal registry.

### Step 1: Create Release Structure

```bash
# Create release directory
VERSION=3.0.0
mkdir -p release

# Build for multiple platforms
for OS in linux darwin windows; do
  for ARCH in amd64 arm64; do
    GOOS=$OS GOARCH=$ARCH go build \
      -o release/terraform-provider-graylog_${VERSION}_${OS}_${ARCH} \
      ./cmd/terraform-provider-graylog
  done
done

# Create archives
cd release
for file in terraform-provider-graylog_*; do
  zip ${file}.zip $file
done

# Generate checksums
shasum -a 256 *.zip > terraform-provider-graylog_${VERSION}_SHA256SUMS
```

### Step 2: Create GitHub Release

```bash
# Using GitHub CLI
gh release create v${VERSION} \
  --title "v${VERSION}" \
  --notes "Release notes here" \
  release/*.zip \
  release/terraform-provider-graylog_${VERSION}_SHA256SUMS
```

### Step 3: User Installation

Users must manually download and install:

```bash
# Download release
VERSION=3.0.0
OS=linux  # or darwin, windows
ARCH=amd64  # or arm64

wget https://github.com/YOUR_ORG/terraform-provider-graylog/releases/download/v${VERSION}/terraform-provider-graylog_${VERSION}_${OS}_${ARCH}.zip

# Extract and install
unzip terraform-provider-graylog_${VERSION}_${OS}_${ARCH}.zip
mkdir -p ~/.terraform.d/plugins/terraform-provider-graylog/graylog/${VERSION}/${OS}_${ARCH}
mv terraform-provider-graylog_${VERSION}_${OS}_${ARCH} ~/.terraform.d/plugins/terraform-provider-graylog/graylog/${VERSION}/${OS}_${ARCH}/terraform-provider-graylog

# Make executable (Unix/macOS)
chmod +x ~/.terraform.d/plugins/terraform-provider-graylog/graylog/${VERSION}/${OS}_${ARCH}/terraform-provider-graylog
```

## Network Mirror

Set up a network mirror for organization-wide distribution.

### Step 1: Create Mirror Structure

```
/terraform-mirror/
├── registry.terraform.io/
│   └── terraform-provider-graylog/
│       └── graylog/
│           ├── index.json
│           ├── 3.0.0.json
│           └── 3.0.0/
│               ├── terraform-provider-graylog_3.0.0_SHA256SUMS
│               ├── terraform-provider-graylog_3.0.0_SHA256SUMS.sig
│               ├── terraform-provider-graylog_3.0.0_linux_amd64.zip
│               ├── terraform-provider-graylog_3.0.0_linux_arm64.zip
│               ├── terraform-provider-graylog_3.0.0_darwin_amd64.zip
│               └── terraform-provider-graylog_3.0.0_windows_amd64.zip
```

### Step 2: Create Metadata Files

**index.json:**
```json
{
  "versions": {
    "3.0.0": {}
  }
}
```

**3.0.0.json:**
```json
{
  "archives": {
    "darwin_amd64": {
      "url": "terraform-provider-graylog_3.0.0_darwin_amd64.zip",
      "hashes": [
        "h1:Wrf6gbP8uz7jGy4Yz6fbeZqN3jhJnx4OHbsdjSnN6xbs="
      ]
    },
    "darwin_arm64": {
      "url": "terraform-provider-graylog_3.0.0_darwin_arm64.zip",
      "hashes": [
        "h1:xyz..."
      ]
    },
    "linux_amd64": {
      "url": "terraform-provider-graylog_3.0.0_linux_amd64.zip",
      "hashes": [
        "h1:abc..."
      ]
    },
    "linux_arm64": {
      "url": "terraform-provider-graylog_3.0.0_linux_arm64.zip",
      "hashes": [
        "h1:def..."
      ]
    },
    "windows_amd64": {
      "url": "terraform-provider-graylog_3.0.0_windows_amd64.zip",
      "hashes": [
        "h1:ghi..."
      ]
    }
  }
}
```

### Step 3: Serve via HTTP

```nginx
# Nginx configuration
server {
    listen 443 ssl;
    server_name terraform-mirror.company.com;

    location / {
        root /var/www/terraform-mirror;
        autoindex on;
    }
}
```

### Step 4: User Configuration

```hcl
# ~/.terraformrc
provider_installation {
  network_mirror {
    url = "https://terraform-mirror.company.com/"
  }
  direct {
    exclude = ["terraform-provider-graylog/*"]
  }
}
```

## Cloud Storage

Use cloud storage services for simple distribution.

### Amazon S3

#### Upload to S3

```bash
VERSION=3.0.0

# Create bucket structure
aws s3api create-bucket --bucket terraform-providers-company

# Upload providers
for file in release/*.zip; do
  aws s3 cp $file s3://terraform-providers-company/terraform-provider-graylog/${VERSION}/
done

# Upload checksums
aws s3 cp terraform-provider-graylog_${VERSION}_SHA256SUMS \
  s3://terraform-providers-company/terraform-provider-graylog/${VERSION}/

# Make public (optional)
aws s3api put-bucket-policy --bucket terraform-providers-company --policy '{
  "Statement": [{
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::terraform-providers-company/*"
  }]
}'
```

#### User Installation

```bash
# Download from S3
aws s3 cp s3://terraform-providers-company/terraform-provider-graylog/3.0.0/terraform-provider-graylog_3.0.0_linux_amd64.zip .

# Or via HTTPS (if public)
wget https://terraform-providers-company.s3.amazonaws.com/terraform-provider-graylog/3.0.0/terraform-provider-graylog_3.0.0_linux_amd64.zip
```

### Google Cloud Storage

```bash
VERSION=3.0.0

# Create bucket
gsutil mb gs://terraform-providers-company

# Upload providers
gsutil cp release/*.zip gs://terraform-providers-company/terraform-provider-graylog/${VERSION}/

# Make public (optional)
gsutil iam ch allUsers:objectViewer gs://terraform-providers-company
```

### Azure Blob Storage

```bash
VERSION=3.0.0

# Create container
az storage container create --name terraform-providers

# Upload providers
az storage blob upload-batch \
  --destination terraform-providers \
  --destination-path terraform-provider-graylog/${VERSION} \
  --source release/

# Generate SAS URL for access
az storage container generate-sas \
  --name terraform-providers \
  --permissions r \
  --expiry 2024-12-31
```

## Checksum and Signing

### Generate Checksums

```bash
# SHA256 checksums
sha256sum terraform-provider-graylog_* > terraform-provider-graylog_${VERSION}_SHA256SUMS

# SHA512 (alternative)
sha512sum terraform-provider-graylog_* > terraform-provider-graylog_${VERSION}_SHA512SUMS
```

### GPG Signing

```bash
# Import GPG key
gpg --import private-key.asc

# Sign the checksum file
gpg --detach-sign terraform-provider-graylog_${VERSION}_SHA256SUMS

# Verify signature
gpg --verify terraform-provider-graylog_${VERSION}_SHA256SUMS.sig
```

### Terraform Registry Requirements

For the official Terraform Registry, you must:

1. Generate SHA256 checksums
2. Sign the checksum file with GPG
3. Include both in the GitHub release
4. Configure your GPG public key in the registry

## Automation

### GitHub Actions

Create `.github/workflows/release.yml`:

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v4
        with:
          go-version: 1.19

      - name: Import GPG key
        uses: crazy-max/ghaction-import-gpg@v5
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.GPG_PASSPHRASE }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v4
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ secrets.GPG_FINGERPRINT }}
```

### GitLab CI

Create `.gitlab-ci.yml`:

```yaml
release:
  stage: release
  only:
    - tags
  script:
    - apt-get update && apt-get install -y gpg
    - echo "$GPG_PRIVATE_KEY" | gpg --import
    - goreleaser release --clean
  variables:
    GITLAB_TOKEN: $CI_JOB_TOKEN
    GPG_FINGERPRINT: $GPG_FINGERPRINT
```

## Best Practices

### Version Naming

- Use semantic versioning: `MAJOR.MINOR.PATCH`
- Prefix with 'v' for tags: `v3.0.0`
- Include pre-release for testing: `v3.0.0-rc1`

### Documentation

Always include:
- README.md with installation instructions
- CHANGELOG.md with version history
- LICENSE file
- Examples directory
- Comprehensive provider documentation

### Testing

Before release:
1. Test on multiple platforms
2. Verify upgrades from previous versions
3. Test with different Terraform versions
4. Validate documentation examples

### Security

- Sign all releases with GPG
- Use checksums for integrity verification
- Rotate signing keys periodically
- Use secure channels for distribution
- Implement access controls for private registries

## Troubleshooting

### Common Issues

**"Provider not found"**
- Check the source path in required_providers
- Verify the provider is installed in the correct directory
- Ensure the version constraint matches

**"Checksum mismatch"**
- Re-generate checksums
- Ensure binary hasn't been modified after checksum generation
- Check for corruption during download

**"Permission denied"**
- Make the provider executable: `chmod +x terraform-provider-graylog`
- Check file ownership and permissions

**"Incompatible provider version"**
- Check Terraform version compatibility
- Verify the provider version constraint
- Update the provider version if needed

## Summary

Choose your deployment method based on your needs:

| Method | Best For | Complexity | Cost |
|--------|----------|------------|------|
| Official Registry | Open source projects | Medium | Free |
| GitHub Releases | Small teams, private repos | Low | Free |
| Artifactory | Enterprises with existing infrastructure | Medium | Licensed |
| Network Mirror | Air-gapped environments | High | Self-hosted |
| Cloud Storage | Simple distribution | Low | Storage costs |
| Terraform Cloud | Teams using TFC/TFE | Low | Subscription |

For most organizations, start with GitHub Releases for simplicity, then move to a more sophisticated solution as needs grow.