# Graylog 7 E2E Test Config

This folder will hold the Terraform smoke/E2E setup for the Graylog provider.

## Credentials
- Copy `graylog.auto.tfvars.example` to `graylog.auto.tfvars` and fill in real admin creds for the test Graylog instance.
- The `.gitignore` protects `graylog.auto.tfvars` from commits.

## Usage
1) Copy creds: `cp graylog.auto.tfvars.example graylog.auto.tfvars` and edit with real admin values.
2) Build provider locally (from repo root): `go build -o ./bin/terraform-provider-graylog_v0.0.0 ./cmd/terraform-provider-graylog`
3) Tell Terraform to use the local build (one-time `$HOME/.terraformrc` snippet):
   ```hcl
   provider_installation {
     dev_overrides {
       "terraform-provider-graylog/graylog" = "<REPO>/bin"
     }
     direct {}
   }
   ```
4) In this folder: `terraform init && terraform plan` (optionally `terraform apply` to create test resources).

## Notes
- MCP access is read-only; all create/update/delete validation must go through Terraform.
- Resources here are intentionally minimal (index set) to validate CRUD against Graylog 7; extend with inputs/streams/pipelines as refactor progresses.
- Rotate or replace admin credentials after runs; `graylog.auto.tfvars` stays git-ignored.
