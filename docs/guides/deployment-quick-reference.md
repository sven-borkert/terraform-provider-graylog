# Provider Deployment Quick Reference

Quick decision guide for choosing a deployment method for terraform-provider-graylog.

## Decision Tree

```
Is your provider open source?
├─ Yes → Official Terraform Registry (best discoverability)
└─ No
   ├─ Do you have existing infrastructure?
   │  ├─ Artifactory → Use Terraform repository
   │  ├─ GitLab → Use GitLab Terraform registry
   │  └─ Terraform Cloud/Enterprise → Use private registry
   └─ No existing infrastructure
      ├─ Need simple solution → GitHub Releases + manual install
      ├─ Have cloud account → S3/GCS/Azure Storage
      └─ Need air-gapped → Network mirror
```

## Quick Setup Commands

### Option 1: Official Registry (Public)

```bash
# 1. Create GitHub org and transfer repo
# 2. Add .goreleaser.yml
# 3. Create and push tag
git tag v3.0.0
git push origin v3.0.0

# 4. Run goreleaser
goreleaser release --clean

# 5. Register at https://registry.terraform.io/publish
```

**Users use:**
```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "~> 3.0"
    }
  }
}
```

### Option 2: GitHub Releases (Private)

```bash
# 1. Build
VERSION=3.0.0
GOOS=linux GOARCH=amd64 go build -o terraform-provider-graylog_${VERSION}_linux_amd64 ./cmd/terraform-provider-graylog

# 2. Create zip
zip terraform-provider-graylog_${VERSION}_linux_amd64.zip terraform-provider-graylog_${VERSION}_linux_amd64

# 3. Generate checksums
sha256sum *.zip > terraform-provider-graylog_${VERSION}_SHA256SUMS

# 4. Create GitHub release and upload files
gh release create v${VERSION} *.zip *.SHA256SUMS
```

**Users install:**
```bash
# Download and extract to local plugins
wget https://github.com/YOUR_ORG/terraform-provider-graylog/releases/download/v3.0.0/terraform-provider-graylog_3.0.0_linux_amd64.zip
unzip -d ~/.terraform.d/plugins/terraform-provider-graylog/graylog/3.0.0/linux_amd64/ terraform-provider-graylog_3.0.0_linux_amd64.zip
```

### Option 3: Artifactory

```bash
# 1. Package
zip terraform-provider-graylog_3.0.0_linux_amd64.zip terraform-provider-graylog

# 2. Upload
curl -u user:pass -T terraform-provider-graylog_3.0.0_linux_amd64.zip \
  "https://artifactory.company.com/artifactory/terraform/terraform-provider-graylog/3.0.0/"
```

**Users configure ~/.terraformrc:**
```hcl
provider_installation {
  network_mirror {
    url = "https://artifactory.company.com/artifactory/terraform/"
  }
}
```

### Option 4: S3 Bucket

```bash
# 1. Upload to S3
aws s3 cp terraform-provider-graylog_3.0.0_linux_amd64.zip \
  s3://terraform-providers/graylog/3.0.0/

# 2. Set bucket policy for access
aws s3api put-bucket-policy --bucket terraform-providers --policy file://policy.json
```

**Users download:**
```bash
aws s3 cp s3://terraform-providers/graylog/3.0.0/terraform-provider-graylog_3.0.0_linux_amd64.zip .
# Then manually install
```

### Option 5: Local Testing Only

```bash
# Build and test locally with dev overrides
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

# Create dev.tfrc
cat > dev.tfrc <<EOF
provider_installation {
  dev_overrides {
    "terraform-provider-graylog/graylog" = "$(pwd)"
  }
  direct {}
}
EOF

# Use with
export TF_CLI_CONFIG_FILE=./dev.tfrc
terraform plan
```

## Comparison Matrix

| Method | Setup Time | Maintenance | User Experience | Cost | Security | Best For |
|--------|------------|-------------|-----------------|------|----------|----------|
| **Official Registry** | 2-4 hours | Low | Excellent (automatic) | Free | Public | Open source projects |
| **GitHub Releases** | 30 min | Low | Manual install | Free | Private repos | Small teams |
| **Artifactory** | 1-2 hours | Medium | Good (with .terraformrc) | Licensed | Enterprise | Large organizations |
| **Terraform Cloud** | 1 hour | Low | Excellent | Subscription | Enterprise | TFC/TFE users |
| **S3/Cloud Storage** | 30 min | Low | Manual install | Storage costs | IAM controlled | Cloud-native teams |
| **Network Mirror** | 2-4 hours | High | Good (with .terraformrc) | Infrastructure | Full control | Air-gapped environments |
| **Dev Overrides** | 5 min | None | Development only | Free | Local only | Development/testing |

## Required Files by Method

### All Methods Need:
- ✅ Compiled provider binaries
- ✅ SHA256 checksums

### Official Registry Also Needs:
- ✅ GPG signed checksums
- ✅ .goreleaser.yml
- ✅ GitHub organization (not personal)
- ✅ Public repository
- ✅ Documentation in docs/

### Private Registry Needs:
- ✅ Registry-specific metadata
- ✅ Authentication setup
- ✅ Network access configuration

## Multi-Platform Build Script

Save as `build-all.sh`:

```bash
#!/bin/bash
VERSION=${1:-3.0.0}

# Platforms to build for
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  OUTPUT="terraform-provider-graylog_${VERSION}_${GOOS}_${GOARCH}"

  if [ "$GOOS" = "windows" ]; then
    OUTPUT="${OUTPUT}.exe"
  fi

  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -o "release/${OUTPUT}" ./cmd/terraform-provider-graylog

  # Create zip
  (cd release && zip "${OUTPUT}.zip" "${OUTPUT}")
done

# Generate checksums
(cd release && sha256sum *.zip > "terraform-provider-graylog_${VERSION}_SHA256SUMS")

echo "Build complete! Files in release/"
```

## Security Considerations

### For Production Deployments

1. **Always provide checksums**
   ```bash
   sha256sum terraform-provider-graylog_* > SHA256SUMS
   ```

2. **Sign releases (for public distribution)**
   ```bash
   gpg --detach-sign SHA256SUMS
   ```

3. **Use HTTPS only** for distribution

4. **Implement access controls**
   - Private repos for proprietary code
   - IAM/RBAC for cloud storage
   - Authentication for private registries

5. **Audit provider usage**
   - Track downloads
   - Monitor for unauthorized access
   - Log provider versions in use

## Common Pitfalls to Avoid

❌ **Don't** commit binaries to git repos
❌ **Don't** use HTTP for distribution
❌ **Don't** skip checksums
❌ **Don't** use personal GitHub accounts for official registry
❌ **Don't** forget to test on all target platforms

✅ **Do** use semantic versioning
✅ **Do** provide comprehensive documentation
✅ **Do** test upgrade paths
✅ **Do** automate releases with CI/CD
✅ **Do** keep signing keys secure

## Getting Help

- **Terraform Registry Issues**: https://github.com/hashicorp/terraform/issues
- **GoReleaser**: https://goreleaser.com/documentation/
- **Provider Development**: https://developer.hashicorp.com/terraform/plugin
- **This Provider**: File issues in this repository