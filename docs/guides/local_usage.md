# Local Provider Testing and Development

This guide explains how to test and develop the Graylog Terraform provider locally.

## Quick Start

For comprehensive examples and tools for local provider testing, see the **[example-local-usage](../../example-local-usage/)** directory in the repository root.

The example-local-usage directory provides:
- Complete working examples
- Multiple testing methods (standard, developer override, Makefile)
- Build scripts and helper tools
- Import examples and data source queries
- Comprehensive documentation

## Testing Methods

### Method 1: Using the Example Directory (Recommended)

The easiest way to test the provider locally is to use the example directory:

```bash
# Navigate to the example directory
cd example-local-usage/

# Build the provider and set up for testing
make setup

# Run terraform plan
make plan
```

See the [example-local-usage README](../../example-local-usage/README.md) for detailed instructions.

### Method 2: Manual Build and Install

If you prefer to manually build and install:

```bash
# Build the provider
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

# Install to local plugin directory
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
mkdir -p ~/.terraform.d/plugins/sven-borkert/graylog/3.0.0/${OS}_${ARCH}
cp terraform-provider-graylog ~/.terraform.d/plugins/sven-borkert/graylog/3.0.0/${OS}_${ARCH}/

# Use in your Terraform configuration
terraform init
```

### Method 3: Developer Override Mode

For active development with fast iteration:

```bash
# Build the provider
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

# Create or edit ~/.terraformrc
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "sven-borkert/graylog" = "$(pwd)"
  }
  direct {}
}
EOF

# No terraform init needed - just run terraform plan
terraform plan
```

## Development Workflow

### Standard Development Cycle

1. Make changes to the provider code
2. Build the provider: `go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog`
3. Test your changes: `terraform plan` or `terraform apply`
4. Iterate as needed

### Using the Example Directory

The example-local-usage directory provides the best development experience:

```bash
cd example-local-usage/

# Rebuild provider after code changes
make build

# Test with developer overrides (no init needed)
make plan
make apply

# Clean up when done
make clean
```

## Testing Against Real Graylog

### Prerequisites

- Running Graylog 7.0+ instance
- Admin credentials
- Network access to Graylog API

### Configuration

Create a `terraform.tfvars` file:

```hcl
graylog_endpoint    = "https://your-graylog-server.com/api"
graylog_username    = "admin"
graylog_password    = "your-password"
graylog_api_version = "v1"
```

Or use environment variables:

```bash
export TF_VAR_graylog_endpoint="https://your-graylog-server.com/api"
export TF_VAR_graylog_username="admin"
export TF_VAR_graylog_password="your-password"
export TF_VAR_graylog_api_version="v1"
```

### Example Test Configuration

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "sven-borkert/graylog"
      version = "3.0.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = var.graylog_endpoint
  auth_name        = var.graylog_username
  auth_password    = var.graylog_password
  api_version      = var.graylog_api_version
}

# Test with a simple stream
resource "graylog_stream" "test" {
  title                              = "Test Stream"
  description                        = "Created for provider testing"
  index_set_id                       = data.graylog_index_set.default.id
  remove_matches_from_default_stream = false
}

data "graylog_index_set" "default" {
  index_set_id = "default"
}
```

## Debugging

### Enable Debug Logging

```bash
export TF_LOG=DEBUG
terraform plan
```

### Provider Logs

The provider logs detailed information about API calls when debug logging is enabled.

### Common Issues

**Provider not found:**
- Ensure the binary is built and in the correct location
- Check that the plugin directory structure is correct
- Verify dev_overrides path is absolute (if using)

**Authentication failures:**
- Verify the API endpoint includes `/api`
- Check user has necessary permissions in Graylog 7.0+
- Test API access with curl:
  ```bash
  curl -u admin:password https://your-graylog.com/api/system
  ```

**API errors:**
- Check Graylog version compatibility (requires 7.0+)
- Review Graylog logs for detailed error messages
- Ensure resources exist before querying them

## Additional Resources

- **[Example Directory](../../example-local-usage/)** - Complete working examples
- **[Migration Guide](migration_guide_v7.md)** - Upgrading to Graylog 7.0
- **[API Mapping](../reference/api_mapping.md)** - API changes documentation
- **[Provider Naming](provider_naming.md)** - Resource naming conventions

## Quick Reference

| Task | Command |
|------|---------|
| Build provider | `go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog` |
| Build from example dir | `cd example-local-usage && make build` |
| Run quick test | `cd example-local-usage && make plan` |
| Install locally | `./build-local.sh` (interactive) |
| Setup example | `cd example-local-usage && make setup` |
| Clean test environment | `cd example-local-usage && make clean` |

## See Also

For the most comprehensive and up-to-date testing documentation, always refer to:
- **[example-local-usage/README.md](../../example-local-usage/README.md)**