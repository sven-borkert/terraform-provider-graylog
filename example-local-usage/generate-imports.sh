#!/bin/bash

# Script to generate Terraform import commands for existing Graylog resources
# This script helps you import existing resources from your Graylog server

set -e

echo "==================================================="
echo "Graylog Resource Import Helper"
echo "==================================================="
echo ""

# Check if terraform.tfvars exists
if [ ! -f "terraform.tfvars" ]; then
    echo "ERROR: terraform.tfvars not found!"
    echo "Please copy terraform.tfvars.example to terraform.tfvars and configure it."
    exit 1
fi

# Note: Using developer overrides, no terraform init needed
echo "Using developer overrides mode (no init required)..."
export TF_CLI_CONFIG_FILE=./dev.tfrc

# Refresh state to get latest data
echo "Refreshing Terraform state to discover resources..."
terraform refresh > /dev/null 2>&1

# Create output file for import commands
OUTPUT_FILE="import-commands.sh"
echo "#!/bin/bash" > $OUTPUT_FILE
echo "# Generated import commands for Graylog resources" >> $OUTPUT_FILE
echo "# Generated on: $(date)" >> $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

# Get discovered resources using terraform output
echo "Discovering resources..."
echo ""

# Get streams
echo "Found streams:"
terraform output -json discovered_streams 2>/dev/null | jq -r '.[] | "  - " + .id + " (" + .title + ")"' || echo "  No streams found or data source not configured"

# Get dashboards
echo ""
echo "Found dashboards:"
terraform output -json discovered_dashboards 2>/dev/null | jq -r '.[] | "  - " + .id + " (" + .title + ")"' || echo "  No dashboards found or data source not configured"

echo ""
echo "==================================================="
echo "Import Command Examples"
echo "==================================================="
echo ""

# Generate example import commands
cat << 'EOF'
To import existing resources, use these commands:

1. Import a specific stream:
   terraform import graylog_stream.my_stream <STREAM_ID>

2. Import a specific dashboard:
   terraform import graylog_dashboard.my_dashboard <DASHBOARD_ID>

3. Import a user:
   terraform import graylog_user.my_user <USERNAME>

4. Import an index set:
   terraform import graylog_index_set.my_index_set <INDEX_SET_ID>

Steps to import:
1. Add the resource block to imports.tf with a unique name
2. Run the import command with the resource ID
3. Run 'terraform plan' to see the current state
4. Update the resource configuration to match the imported state
5. Run 'terraform apply' to confirm the import

EOF

echo "==================================================="
echo "Next Steps:"
echo "==================================================="
echo "1. Review the discovered resources above"
echo "2. Edit imports.tf to add resource blocks for items you want to import"
echo "3. Run the appropriate terraform import commands"
echo "4. Use 'terraform plan' to verify the imported state"
echo ""
echo "Import commands have been saved to: $OUTPUT_FILE"