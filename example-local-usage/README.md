# Example: Using the Local Graylog Provider

This directory contains a simple example of using the locally-built Graylog provider.

## Quick Start

1. **Build and install the provider** (from repository root):
   ```bash
   cd ..
   ./build-local.sh
   ```

2. **Configure your Graylog credentials**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your Graylog server details
   ```

3. **Initialize Terraform**:
   ```bash
   terraform init
   ```

4. **Review the plan**:
   ```bash
   terraform plan
   ```

5. **Apply the configuration**:
   ```bash
   terraform apply
   ```

## Using Environment Variables

Instead of using `terraform.tfvars`, you can set environment variables:

```bash
export TF_VAR_graylog_endpoint="https://your-graylog-server.com/api"
export TF_VAR_graylog_username="admin"
export TF_VAR_graylog_password="your-password"

terraform plan
terraform apply
```

## What This Example Does

This example:
1. Connects to your Graylog server
2. Retrieves information about the default index set
3. Creates a new stream called "Test Stream from Terraform"
4. Outputs the ID of the created stream

## Cleanup

To remove the created resources:

```bash
terraform destroy
```

## Troubleshooting

If you encounter issues:

1. **Provider not found**: Make sure you ran `./build-local.sh` and selected "y" to install
2. **Authentication failed**: Verify your credentials and that the API endpoint ends with `/api`
3. **Permission denied**: Ensure your Graylog user has admin privileges
4. **Connection refused**: Check that your Graylog server is accessible and HTTPS certificate is valid

For Graylog 7.0+, ensure your user has the necessary API permissions.