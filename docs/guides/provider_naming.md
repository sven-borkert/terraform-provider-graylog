# Terraform Provider Naming - Where It's Defined

The Terraform provider name "graylog" is defined through multiple conventions and locations rather than a single configuration value. Here's a comprehensive breakdown:

## 1. Binary Name Convention

The provider name comes from the binary filename pattern:
```
terraform-provider-{NAME}
```

In our case:
- Binary: `terraform-provider-graylog`
- Provider name extracted: **`graylog`**

This is defined when building:
```bash
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
```

## 2. Local Plugin Directory Structure

When installed locally, the provider name appears in the directory path:
```
~/.terraform.d/plugins/{NAMESPACE}/{NAME}/{VERSION}/{OS}_{ARCH}/
~/.terraform.d/plugins/terraform-provider-graylog/graylog/3.0.0/linux_amd64/
                        └── namespace ──┘  └─name─┘
```

- **Namespace**: `terraform-provider-graylog` (organization/publisher)
- **Name**: `graylog` (the actual provider name)

## 3. Terraform Configuration

In your Terraform configuration files:

```hcl
terraform {
  required_providers {
    graylog = {  # <-- Provider local name (can be aliased)
      source  = "terraform-provider-graylog/graylog"
      #          └── namespace ──┘  └─name─┘
      version = "3.0.0"
    }
  }
}

provider "graylog" {  # <-- Must match the local name above
  # Configuration...
}
```

## 4. Go Module Name (Not Provider Name)

The Go module name in `go.mod` is different - it's the import path:
```go
module github.com/sven-borkert/terraform-provider-graylog
```

This is NOT the provider name - it's just for Go imports.

## 5. Provider Schema Code

The provider itself doesn't declare its name internally. The schema just defines configuration:

```go
// graylog/provider.go
func Provider() *schema.Provider {
    return &schema.Provider{
        Schema:         provider.SchemaMap(),  // Configuration schema
        ResourcesMap:   resourceMap,           // Available resources
        DataSourcesMap: dataSourcesMap,        // Available data sources
        ConfigureFunc:  provider.Configure,    // Configuration function
    }
}
```

## 6. Terraform Registry (If Published)

If published to the Terraform Registry, the full address would be:
```
registry.terraform.io/{NAMESPACE}/{NAME}
registry.terraform.io/terraform-provider-graylog/graylog
```

## How to Change the Provider Name

If you wanted to change the provider name from "graylog" to something else, you would need to:

1. **Change the binary name**:
   ```bash
   go build -o terraform-provider-mynewname ./cmd/terraform-provider-graylog
   ```

2. **Update the installation path**:
   ```
   ~/.terraform.d/plugins/mynamespace/mynewname/3.0.0/{OS}_{ARCH}/
   ```

3. **Update Terraform configurations**:
   ```hcl
   terraform {
     required_providers {
       mynewname = {
         source  = "mynamespace/mynewname"
         version = "3.0.0"
       }
     }
   }

   provider "mynewname" {
     # Configuration...
   }
   ```

## Current Setup Summary

- **Provider Name**: `graylog`
- **Namespace**: `terraform-provider-graylog`
- **Full Source**: `terraform-provider-graylog/graylog`
- **Binary Name**: `terraform-provider-graylog`
- **Go Module**: `github.com/sven-borkert/terraform-provider-graylog` (unrelated to provider name)

## Important Notes

1. The provider name is derived from conventions, not explicitly configured in code
2. The binary must follow the `terraform-provider-{name}` pattern
3. The directory structure must match what Terraform expects
4. The Go module name is independent of the Terraform provider name
5. When using locally, the "namespace" can be arbitrary but should be consistent

## Environment Variables

The provider uses `GRAYLOG_` prefixed environment variables, which conventionally match the provider name:
- `GRAYLOG_WEB_ENDPOINT_URI`
- `GRAYLOG_AUTH_NAME`
- `GRAYLOG_AUTH_PASSWORD`
- `GRAYLOG_X_REQUESTED_BY`

These are defined in `graylog/provider/provider.go` in the schema.