# Local Provider Testing and Development

This guide explains how to test and develop the Graylog Terraform provider locally.

## Quick Start

```bash
# Build the provider
make build

# Test with local provider (uses bin/terraform-dev wrapper)
cd examples/graylog7-e2e
../../bin/terraform-dev plan
../../bin/terraform-dev apply
```

## Testing Methods

### Method 1: Using terraform-dev Wrapper (Recommended)

The `bin/terraform-dev` script automatically configures Terraform to use the locally built provider:

```bash
# Build the provider
make build

# Use terraform-dev instead of terraform
cd examples/graylog7-e2e
../../bin/terraform-dev plan
../../bin/terraform-dev apply
../../bin/terraform-dev destroy
```

The wrapper:
- Creates a `.terraformrc-dev` config file in the project root
- Sets `TF_CLI_CONFIG_FILE` to use that config
- Passes all arguments to terraform

No `terraform init` is needed when using dev overrides.

### Method 2: Using the Registry Version

To test against the published registry version:

```bash
cd examples/graylog7-e2e
terraform init
terraform plan
terraform apply
```

### Method 3: Manual Dev Override

If you prefer manual configuration:

```bash
# Build the provider
make build

# Create ~/.terraformrc with dev override
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "sven-borkert/graylog" = "/path/to/terraform-provider-graylog/bin"
  }
  direct {}
}
EOF

# Run terraform (no init needed)
terraform plan
```

## Development Workflow

### Standard Development Cycle

1. Make changes to the provider code
2. Build: `make build`
3. Test: `bin/terraform-dev plan` (from examples directory)
4. Iterate as needed

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build provider binary to `bin/` |
| `make test` | Run unit tests with race detection |
| `make acc-test` | Run acceptance tests (requires Graylog) |
| `make lint` | Run golangci-lint |
| `make fmt` | Format Go code |
| `make clean` | Remove build artifacts |

## Testing Against Real Graylog

### Prerequisites

- Running Graylog 7.0+ instance
- Admin credentials
- Network access to Graylog API

### Configuration

Create `examples/graylog7-e2e/graylog.auto.tfvars`:

```hcl
graylog_web_endpoint_uri = "https://your-graylog-server.com/api"
graylog_auth_name        = "admin"
graylog_auth_password    = "your-password"
```

This file is gitignored and will not be committed.

## Debugging

### Enable Debug Logging

```bash
export TF_LOG=DEBUG
../../bin/terraform-dev plan
```

### Common Issues

**Provider not found:**
- Ensure `make build` completed successfully
- Check that `bin/terraform-provider-graylog` exists

**Authentication failures:**
- Verify the API endpoint includes `/api`
- Check user has necessary permissions
- Test API access with curl:
  ```bash
  curl -u admin:password https://your-graylog.com/api/system
  ```

**API errors:**
- Check Graylog version compatibility (requires 7.0+)
- Review Graylog server logs for detailed error messages

## Quick Reference

| Task | Command |
|------|---------|
| Build provider | `make build` |
| Test with local provider | `cd examples/graylog7-e2e && ../../bin/terraform-dev plan` |
| Test with registry provider | `cd examples/graylog7-e2e && terraform init && terraform plan` |
| Run unit tests | `make test` |
| Run linter | `make lint` |
| Format code | `make fmt` |

## Additional Resources

- **[API Mapping](../reference/api_mapping.md)** - API changes documentation
- **[Provider Naming](provider_naming.md)** - Resource naming conventions
