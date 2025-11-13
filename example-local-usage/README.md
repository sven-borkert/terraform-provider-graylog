# Local Graylog Provider Testing Examples

This directory contains comprehensive examples and tools for testing the Graylog Terraform provider locally without pushing to a registry.

## Table of Contents

- [Quick Start](#quick-start)
- [Testing Methods](#testing-methods)
- [File Structure](#file-structure)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Prerequisites

- Terraform >= 0.13
- Go >= 1.15 (for building the provider)
- Access to a Graylog server (version 7.0+ recommended)
- Admin credentials for Graylog

### Setup Steps

1. **Build the provider**:

   **Option A: From this directory**
   ```bash
   ./build-provider.sh
   ```

   **Option B: From repository root**
   ```bash
   cd ..
   go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
   # Or use the interactive build script:
   ./build-local.sh
   ```

2. **Configure credentials**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your Graylog server details
   ```

3. **Choose your testing method**:

   **Option A: Standard Method (with terraform init)**
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

   **Option B: Developer Override Method (no init needed)**
   ```bash
   ./use-dev-mode.sh  # Setup developer overrides
   ./dev-plan.sh      # Run terraform plan
   ./dev-apply.sh     # Run terraform apply
   ```

   **Option C: Using Makefile**
   ```bash
   make setup         # Initial setup
   make plan          # Show planned changes
   make apply         # Apply changes
   ```

## Testing Methods

### Method 1: Standard Terraform Workflow

This uses the locally installed provider from `~/.terraform.d/plugins/`:

```bash
terraform init
terraform plan
terraform apply
```

### Method 2: Developer Override Mode (Recommended for Development)

This method bypasses the need for `terraform init` and uses the provider binary directly:

```bash
# Setup (only needed once)
./use-dev-mode.sh

# Use with environment variable
export TF_CLI_CONFIG_FILE=./dev.tfrc
terraform plan
terraform apply

# Or use convenience scripts
./dev-plan.sh
./dev-apply.sh
./dev-refresh.sh
```

**Benefits of Developer Override Mode:**
- No `terraform init` needed
- No checksum verification
- Direct use of local binary
- Faster iteration during development

### Method 3: Makefile Commands

The Makefile provides convenient commands for all common operations:

```bash
make help          # Show all available commands
make setup         # Initial setup
make validate      # Validate configuration
make fmt           # Format Terraform files
make plan          # Show planned changes
make apply         # Apply changes
make destroy       # Destroy resources (with confirmation)
make refresh       # Refresh state from Graylog
make import-help   # Generate import commands
make clean         # Clean up Terraform files
```

## File Structure

### Core Configuration Files

- **`main.tf`** - Main configuration with provider setup and basic test resources
- **`terraform.tfvars.example`** - Example variables file (copy to terraform.tfvars)
- **`data-sources.tf`** - Examples of querying existing Graylog resources
- **`imports.tf`** - Templates for importing existing resources
- **`examples.tf`** - Comprehensive resource examples (streams, inputs, pipelines, etc.)

### Developer Tools

- **`dev.tfrc`** - Terraform CLI config for developer override mode
- **`use-dev-mode.sh`** - Setup script for developer override mode
- **`dev-plan.sh`** - Run terraform plan with developer overrides
- **`dev-apply.sh`** - Run terraform apply with developer overrides
- **`dev-refresh.sh`** - Refresh state with developer overrides
- **`generate-imports.sh`** - Generate import commands for existing resources

### Supporting Files

- **`Makefile`** - Convenient commands for all operations
- **`.gitignore`** - Excludes sensitive and temporary files

## Common Tasks

### Creating Resources

1. Edit `main.tf` or create new `.tf` files
2. Add your resource definitions
3. Run `terraform plan` to preview changes
4. Run `terraform apply` to create resources

### Importing Existing Resources

1. Identify the resource ID from Graylog UI
2. Add an empty resource block to `imports.tf`:
   ```hcl
   resource "graylog_stream" "existing" {
     # Configuration will be filled after import
   }
   ```
3. Run the import command:
   ```bash
   terraform import graylog_stream.existing <STREAM_ID>
   ```
4. Run `terraform plan` to see the imported configuration
5. Update the resource block to match the imported state

### Discovering Resources

Use the import helper script to discover resources:

```bash
./generate-imports.sh
# Or
make import-help
```

This will attempt to discover available resources and generate import commands.

### Using Data Sources

Query existing resources without managing them:

```hcl
data "graylog_index_set" "default" {
  index_set_id = "default"
}

output "index_details" {
  value = data.graylog_index_set.default
}
```

## Environment Variables

Instead of using `terraform.tfvars`, you can set environment variables:

```bash
# Provider configuration
export TF_VAR_graylog_endpoint="https://your-graylog.com/api"
export TF_VAR_graylog_username="admin"
export TF_VAR_graylog_password="secure-password"
export TF_VAR_graylog_api_version="v1"

# Alternative (direct provider env vars)
export GRAYLOG_WEB_ENDPOINT_URI="https://your-graylog.com/api"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="secure-password"
```

## Resource Examples

The `examples.tf` file contains comprehensive examples of:

- **Index Sets** - Managing log indices
- **Streams** - Creating and configuring log streams
- **Stream Rules** - Filtering and routing messages
- **Inputs** - Configuring data sources (Syslog, GELF, Beats, etc.)
- **Pipelines** - Message processing pipelines
- **Pipeline Rules** - Processing rules for pipelines
- **Users & Roles** - User management and RBAC
- **Dashboards** - Creating monitoring dashboards
- **Alerts** - Setting up alert conditions

## Troubleshooting

### Provider Not Found

```bash
# Ensure the provider is built
cd ..
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

# For standard method, install it
./build-local.sh
# Choose 'y' to install locally

# For developer mode, ensure binary exists at ../terraform-provider-graylog
ls -la ../terraform-provider-graylog
```

### Authentication Failed

- Verify credentials in `terraform.tfvars`
- Ensure the API endpoint ends with `/api`
- Check user has admin privileges
- For Graylog 7.0+, verify API permissions

### Connection Refused

- Verify Graylog server is accessible
- Check HTTPS certificate validity
- Test with curl:
  ```bash
  curl -u admin:password https://your-graylog.com/api/system/sessions
  ```

### Developer Override Warning

When using developer override mode, you'll see:
```
Warning: Provider development overrides are in effect
```

This is expected and can be ignored. It indicates you're using the local binary.

### State Issues

If you encounter state conflicts:
```bash
# Refresh state from Graylog
terraform refresh
# Or with developer mode
./dev-refresh.sh

# Force unlock if locked
terraform force-unlock <LOCK_ID>

# Start fresh (caution: loses state)
rm -f terraform.tfstate*
```

## Best Practices

1. **Always use version control** - Commit your `.tf` files but not `.tfvars` or `.tfstate`
2. **Test in non-production first** - Use a development Graylog instance
3. **Keep credentials secure** - Never commit passwords or API keys
4. **Use data sources** - Reference existing resources instead of hardcoding IDs
5. **Plan before apply** - Always review changes with `terraform plan`
6. **Document your resources** - Use descriptions in your resource definitions
7. **Use consistent naming** - Follow a naming convention for resources

## Additional Resources

- [Graylog API Documentation](https://docs.graylog.org/en/latest/pages/configuration/rest_api.html)
- [Terraform Provider Development](https://www.terraform.io/docs/extend/testing/index.html)
- [Provider Documentation](../docs/)