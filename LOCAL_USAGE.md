# Using the Terraform Provider Graylog from Local Repository

This guide explains how to build and use the terraform-provider-graylog directly from this git repository.

## Method 1: Build and Install Locally (Recommended)

### Step 1: Build the Provider

```bash
# From the repository root directory
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
```

### Step 2: Install to Local Terraform Plugin Directory

For Terraform 0.13+, install the provider to your local plugin directory:

```bash
# Determine your OS and architecture
OS=$(go env GOOS)
ARCH=$(go env GOARCH)

# Create the plugin directory structure
mkdir -p ~/.terraform.d/plugins/terraform-provider-graylog/graylog/3.0.0/${OS}_${ARCH}

# Copy the built provider
cp terraform-provider-graylog ~/.terraform.d/plugins/terraform-provider-graylog/graylog/3.0.0/${OS}_${ARCH}/
```

### Step 3: Configure Terraform to Use the Local Provider

Create a `versions.tf` file in your Terraform configuration:

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "3.0.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "https://your-graylog-server.com/api"
  auth_name        = "admin"
  auth_password    = "your-password"
  api_version      = "v1"
}
```

## Method 2: Development Override (For Active Development)

If you're actively developing the provider, use Terraform's development overrides:

### Step 1: Create a Terraform CLI Configuration

Create or edit `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "terraform-provider-graylog/graylog" = "/home/sborkert/Repos/terraform-provider-graylog"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### Step 2: Build the Provider

```bash
# Build the provider in the repository root
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
```

### Step 3: Use in Your Terraform Configuration

With dev_overrides, you don't need to run `terraform init`. Just use the provider:

```hcl
provider "graylog" {
  web_endpoint_uri = "https://your-graylog-server.com/api"
  auth_name        = "admin"
  auth_password    = "your-password"
  api_version      = "v1"
}

# Example resource
resource "graylog_stream" "example" {
  title                              = "Example Stream"
  description                        = "Stream for example logs"
  index_set_id                       = "default"
  remove_matches_from_default_stream = true
}
```

## Method 3: Using Go Install

You can also install directly using go:

```bash
# Install the provider globally
go install github.com/terraform-provider-graylog/terraform-provider-graylog/cmd/terraform-provider-graylog@latest

# The binary will be available in $GOPATH/bin or $HOME/go/bin
```

## Example Terraform Configuration

Here's a complete example configuration file (`main.tf`):

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "3.0.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = var.graylog_web_endpoint_uri
  auth_name        = var.graylog_auth_name
  auth_password    = var.graylog_auth_password
  api_version      = "v1"
}

variable "graylog_web_endpoint_uri" {
  description = "Graylog API endpoint URI"
  type        = string
}

variable "graylog_auth_name" {
  description = "Graylog authentication username"
  type        = string
}

variable "graylog_auth_password" {
  description = "Graylog authentication password"
  type        = string
  sensitive   = true
}

# Example: Create an index set
resource "graylog_index_set" "example" {
  title                = "Example Index Set"
  description          = "An example index set for testing"
  index_prefix         = "example"
  shards               = 1
  replicas             = 0
  rotation_strategy    = "count"
  rotation_strategy_details = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.MessageCountRotationStrategy"
    max_docs_per_index = 20000000
  })
  retention_strategy   = "delete"
  retention_strategy_details = jsonencode({
    type                    = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
    max_number_of_indices   = 10
  })
  index_analyzer       = "standard"
  index_optimization_max_num_segments = 1
  field_type_refresh_interval = 5000
}

# Example: Create a stream
resource "graylog_stream" "example" {
  title                              = "Example Stream"
  description                        = "Stream for example application logs"
  index_set_id                       = graylog_index_set.example.id
  remove_matches_from_default_stream = true
}

# Example: Create a stream rule
resource "graylog_stream_rule" "example" {
  stream_id    = graylog_stream.example.id
  field        = "application"
  value        = "example-app"
  type         = 1  # EXACT match
  inverted     = false
  description  = "Match logs from example-app"
}
```

## Environment Variables

You can also use environment variables for authentication:

```bash
export GRAYLOG_WEB_ENDPOINT_URI="https://your-graylog-server.com/api"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="your-password"
```

Then in your Terraform configuration, the provider block can be simplified:

```hcl
provider "graylog" {
  # Authentication will be read from environment variables
  api_version = "v1"
}
```

## Testing Your Configuration

1. Build the provider (if using local build):
   ```bash
   go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
   ```

2. Initialize Terraform (skip if using dev_overrides):
   ```bash
   terraform init
   ```

3. Plan your changes:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```

## Troubleshooting

### Provider Not Found

If Terraform cannot find the provider, check:
- The provider binary is built and in the correct location
- The plugin directory structure matches: `~/.terraform.d/plugins/<namespace>/<name>/<version>/<os>_<arch>/`
- Your `versions.tf` specifies the correct source and version

### Authentication Issues

Ensure:
- The API endpoint includes `/api` (e.g., `https://graylog.example.com/api`)
- Your user has the necessary permissions in Graylog
- For Graylog 7.0+, ensure your user has the required API permissions

### Development Override Not Working

If using `.terraformrc`:
- Ensure the path in `dev_overrides` is absolute
- The provider binary must be named `terraform-provider-graylog` and be in the specified directory
- You cannot use `terraform init` with dev_overrides active

## Graylog 7.0 Compatibility

This provider now fully supports Graylog 7.0. The provider automatically handles:
- CreateEntityRequest wrapper for entity creation
- Removal of computed fields in update requests
- All API changes required for Graylog 7.0

Your existing Terraform configurations will work without modification when upgrading to Graylog 7.0.