#!/bin/bash

# Setup script to use Terraform with developer overrides for local provider

echo "==================================================="
echo "Setting up Terraform with Developer Overrides"
echo "==================================================="
echo ""

# Check if provider binary exists
PROVIDER_BINARY="../terraform-provider-graylog"
if [ ! -f "$PROVIDER_BINARY" ]; then
    echo "❌ Error: Provider binary not found at $PROVIDER_BINARY"
    echo "Please build the provider first:"
    echo "  cd .."
    echo "  go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog"
    exit 1
fi

echo "✅ Found provider binary"

# Clean up any existing Terraform directories
if [ -d ".terraform" ] || [ -f ".terraform.lock.hcl" ] || [ -f "terraform.tfstate" ]; then
    echo "Cleaning up existing Terraform files..."
    rm -rf .terraform .terraform.lock.hcl
    # Keep terraform.tfstate if it has real resources
    if [ -f "terraform.tfstate" ]; then
        if grep -q '"resources": \[\]' terraform.tfstate 2>/dev/null; then
            rm -f terraform.tfstate
        else
            echo "  Keeping existing terraform.tfstate (contains resources)"
        fi
    fi
fi

echo ""
echo "==================================================="
echo "✅ Setup Complete!"
echo "==================================================="
echo ""
echo "When using developer overrides mode:"
echo "  • NO 'terraform init' is needed"
echo "  • Terraform will use the local binary directly"
echo "  • You'll see a warning about using an unverified provider"
echo ""
echo "To use Terraform with the local provider:"
echo ""
echo "  export TF_CLI_CONFIG_FILE=./dev.tfrc"
echo "  terraform plan"
echo "  terraform apply"
echo ""
echo "Or use the convenience commands:"
echo "  ./dev-plan.sh     # Run terraform plan"
echo "  ./dev-apply.sh    # Run terraform apply"
echo "  ./dev-refresh.sh  # Refresh state from Graylog"
echo ""
echo "Note: The provider warning is expected and can be ignored."