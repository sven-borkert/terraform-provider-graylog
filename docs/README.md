# Terraform Provider Graylog - Documentation

Welcome to the terraform-provider-graylog documentation. This directory contains comprehensive documentation for using, developing, and understanding the provider.

## 📚 Documentation Structure

```
docs/
├── index.md                          # Provider overview and quick start
├── changelog.md                      # Version history and changes
├── guides/                           # User guides
│   ├── local_usage.md               # Local testing and development
│   ├── migration_guide_v7.md        # Graylog 7.0 upgrade guide
│   ├── provider_naming.md           # Naming conventions
│   └── json-string-attribute.md     # Working with JSON attributes
├── reference/                        # Technical reference
│   ├── architecture.md              # Provider architecture with diagrams
│   ├── api_mapping.md               # API endpoint documentation
│   └── inventory.md                 # Component inventory
├── resources/                        # Resource documentation
│   ├── stream.md
│   ├── input.md
│   ├── pipeline.md
│   └── ...                          # 25+ resource types
├── data-sources/                     # Data source documentation
│   ├── stream.md
│   ├── index_set.md
│   └── ...
└── development/                      # Development documentation
    └── project_summary.md           # Project completion summary
```

## 🚀 Quick Navigation

### For Users

**Getting Started:**
- [Provider Overview](index.md) - Features, configuration, and quick start
- [Local Testing Guide](guides/local_usage.md) - Test the provider locally
- [Example Directory](../example-local-usage/) - Complete working examples

**Migration & Upgrades:**
- [Graylog 7.0 Migration Guide](guides/migration_guide_v7.md) - Upgrade from earlier versions
- [Changelog](changelog.md) - What's changed in each version

**Reference:**
- [All Resources](index.md#available-resources) - Complete resource list with links
- [Data Sources](index.md#available-data-sources) - Query existing resources
- [JSON Attributes Guide](guides/json-string-attribute.md) - Working with JSON configs

### For Developers

**Architecture & Design:**
- [Architecture Overview](reference/architecture.md) - Component structure with diagrams
- [API Mapping](reference/api_mapping.md) - Graylog API endpoints and changes
- [Component Inventory](reference/inventory.md) - Complete file listing

**Development:**
- [Local Testing](guides/local_usage.md) - Development workflow
- [Example Directory](../example-local-usage/README.md) - Comprehensive testing setup
- [Project Summary](development/project_summary.md) - Implementation details

## 📖 Documentation by Topic

### Provider Configuration

**Basic Setup:**
```hcl
provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = var.graylog_password
  api_version      = "v1"
}
```

**Documentation:**
- [Provider Configuration](index.md#provider-configuration) - All arguments and options
- [Authentication Methods](index.md#authentication-methods) - Password, token, session
- [Environment Variables](index.md#environment-variables) - Alternative configuration

### Resource Management

**Core Resources:**
- **Log Management:** [index_set](resources/index_set.md), [stream](resources/stream.md), [stream_rule](resources/stream_rule.md)
- **Data Inputs:** [input](resources/input.md), [extractor](resources/extractor.md)
- **Processing:** [pipeline](resources/pipeline.md), [pipeline_rule](resources/pipeline_rule.md), [grok_pattern](resources/grok_pattern.md)
- **Alerting:** [event_definition](resources/event_definition.md), [event_notification](resources/event_notification.md)
- **Access Control:** [user](resources/user.md), [role](resources/role.md)

### Graylog 7.0 Compatibility

**Key Changes:**
- Automatic entity wrapping for creation requests
- Computed field removal in update requests
- Zero configuration changes required

**Documentation:**
- [Migration Guide](guides/migration_guide_v7.md) - Step-by-step upgrade instructions
- [Compatibility Matrix](index.md#compatibility) - Version requirements
- [API Changes](reference/api_mapping.md) - Technical details

### Local Development

**Quick Start:**
```bash
# From example-local-usage directory
make setup    # Build provider and configure
make plan     # Test with terraform plan
make apply    # Apply changes
```

**Documentation:**
- [Local Usage Guide](guides/local_usage.md) - Multiple testing methods
- [Example Directory](../example-local-usage/README.md) - Complete setup guide
- [Architecture](reference/architecture.md) - Component structure

## 🎯 Common Tasks

### Creating a New Resource

1. Define your resource configuration
2. Reference data sources for dependencies
3. Apply with `terraform apply`

**Example:**
```hcl
resource "graylog_stream" "app_logs" {
  title                              = "Application Logs"
  description                        = "Stream for application messages"
  index_set_id                       = data.graylog_index_set.default.id
  remove_matches_from_default_stream = true
}
```

**See:** [Stream Resource Documentation](resources/stream.md)

### Importing Existing Resources

```bash
# Import a stream
terraform import graylog_stream.existing_stream <stream_id>

# Import with complex ID
terraform import graylog_stream_rule.existing_rule <stream_id>/<rule_id>
```

**See:** [Import Examples](../example-local-usage/imports.tf)

### Testing Locally

**Method 1: Example Directory (Recommended)**
```bash
cd example-local-usage/
make build && make plan
```

**Method 2: Developer Override**
```bash
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
./example-local-usage/use-dev-mode.sh
./example-local-usage/dev-plan.sh
```

**See:** [Local Testing Guide](guides/local_usage.md)

### Upgrading to Graylog 7.0

1. Upgrade Graylog server to 7.0+
2. Update provider version in Terraform:
   ```hcl
   terraform {
     required_providers {
       graylog = {
         version = "~> 3.0"
       }
     }
   }
   ```
3. Run `terraform init -upgrade`
4. Verify with `terraform plan`

**See:** [Migration Guide](guides/migration_guide_v7.md)

## 📋 Resource Categories

### By Function

**Log Management** (4 resources)
- Index Sets, Streams, Stream Rules, Stream Outputs

**Data Collection** (3 resources)
- Inputs, Input Static Fields, Extractors

**Processing** (4 resources)
- Pipelines, Pipeline Rules, Pipeline Connections, Grok Patterns

**Output & Forwarding** (1 resource)
- Outputs

**Alerting & Events** (4 resources)
- Event Definitions, Event Notifications, Alert Conditions*, Alarm Callbacks*

  *Deprecated - use Events System

**Visualization** (3 resources)
- Dashboards, Dashboard Widgets, Dashboard Widget Positions

**Access Control** (3 resources)
- Users, Roles, LDAP Settings

**Sidecar Management** (3 resources)
- Sidecar Configurations, Sidecar Collectors, Sidecars

**Total: 25+ Resource Types**

## 🔍 Finding Information

### By Use Case

| What do you want to do? | Where to look |
|-------------------------|---------------|
| Get started quickly | [Provider Overview](index.md) |
| Test locally | [Local Usage Guide](guides/local_usage.md) |
| Upgrade to Graylog 7.0 | [Migration Guide](guides/migration_guide_v7.md) |
| Create a stream | [Stream Resource](resources/stream.md) |
| Set up log input | [Input Resource](resources/input.md) |
| Process messages | [Pipeline Resource](resources/pipeline.md) |
| Configure alerts | [Event Definition](resources/event_definition.md) |
| Manage users | [User Resource](resources/user.md) |
| Import existing config | [Import Examples](../example-local-usage/imports.tf) |
| Understand architecture | [Architecture Guide](reference/architecture.md) |
| View API changes | [API Mapping](reference/api_mapping.md) |
| See what's new | [Changelog](changelog.md) |

### By Role

**Infrastructure Engineers:**
- [Provider Overview](index.md)
- [All Resources](index.md#available-resources)
- [Migration Guide](guides/migration_guide_v7.md)

**Developers:**
- [Local Testing](guides/local_usage.md)
- [Architecture](reference/architecture.md)
- [Example Directory](../example-local-usage/)

**DevOps/SRE:**
- [Quick Start](index.md#quick-start)
- [Resource Documentation](index.md#available-resources)
- [Troubleshooting](guides/local_usage.md#debugging)

## 📊 Diagrams & Visual Guides

- [Architecture Overview](reference/architecture.md#high-level-architecture) - Component structure
- [Data Flow](reference/architecture.md#data-flow) - Request/response flow
- [Resource Lifecycle](reference/architecture.md#resource-lifecycle) - State management
- [Graylog 7.0 Compatibility](reference/architecture.md#graylog-70-compatibility-layer) - API changes
- [Testing Strategy](reference/architecture.md#testing-strategy) - Development workflow

## 🔗 External Resources

### Official Documentation
- [Graylog Documentation](https://go2docs.graylog.org/)
- [Graylog 7.0 Upgrade Guide](https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm)
- [Graylog REST API](https://go2docs.graylog.org/current/setting_up_graylog/rest_api.html)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)

### Community
- [GitHub Repository](https://github.com/terraform-provider-graylog/terraform-provider-graylog)
- [Terraform Registry](https://registry.terraform.io/providers/terraform-provider-graylog/graylog)
- [Graylog Community](https://community.graylog.org/)
- [Issue Tracker](https://github.com/terraform-provider-graylog/terraform-provider-graylog/issues)

## 📝 Contributing to Documentation

Found an error or want to improve the documentation? Contributions are welcome!

1. Edit the relevant `.md` file
2. Follow the existing structure and style
3. Test any code examples
4. Submit a pull request

**Documentation Guidelines:**
- Use clear, concise language
- Include code examples where helpful
- Link to related documentation
- Keep examples up-to-date with current version

## 🆘 Getting Help

**Provider Issues:**
- Check the [Troubleshooting Guide](guides/local_usage.md#debugging)
- Search [GitHub Issues](https://github.com/terraform-provider-graylog/terraform-provider-graylog/issues)
- Create a new issue with details

**Graylog Questions:**
- [Graylog Community Forums](https://community.graylog.org/)
- [Graylog Documentation](https://go2docs.graylog.org/)

**Terraform Questions:**
- [Terraform Documentation](https://www.terraform.io/docs)
- [HashiCorp Community](https://discuss.hashicorp.com/)

## 📅 Documentation Updates

This documentation is maintained alongside the provider code. Major updates occur with each release:

- **v3.0.0** (2025-11-13) - Graylog 7.0 compatibility, architecture docs, examples
- **Earlier versions** - See git history

**Last Updated:** 2025-11-13