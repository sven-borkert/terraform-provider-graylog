#!/bin/bash

# Build and install terraform-provider-graylog locally

set -e

echo "Building terraform-provider-graylog..."

# Build the provider in the current directory (repository root)
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

echo "Build completed successfully!"
echo "Provider binary created at: ./terraform-provider-graylog"
echo ""

# Ask user about installation options
echo "Installation options:"
echo "1) Install to ~/.terraform.d/plugins (for standard 'terraform init' workflow)"
echo "2) Keep in current directory (for developer override mode with example-local-usage/)"
echo "3) Both"
echo "4) Skip installation"
echo ""
read -p "Choose an option (1-4): " -n 1 -r
echo

case $REPLY in
    1|3)
        # Install to ~/.terraform.d/plugins
        OS=$(go env GOOS)
        ARCH=$(go env GOARCH)
        VERSION="3.0.0"

        PLUGIN_DIR="$HOME/.terraform.d/plugins/terraform-provider-graylog/graylog/$VERSION/${OS}_${ARCH}"
        echo "Creating plugin directory: $PLUGIN_DIR"
        mkdir -p "$PLUGIN_DIR"

        echo "Installing provider to ~/.terraform.d/plugins..."
        cp terraform-provider-graylog "$PLUGIN_DIR/"

        echo "✅ Provider installed to ~/.terraform.d/plugins"
        echo ""
        echo "For standard workflow, use in your Terraform configuration:"
        echo ""
        echo "terraform {"
        echo "  required_providers {"
        echo "    graylog = {"
        echo "      source  = \"terraform-provider-graylog/graylog\""
        echo "      version = \"$VERSION\""
        echo "    }"
        echo "  }"
        echo "}"
        echo ""
        if [[ $REPLY == "1" ]]; then
            echo "Then run: terraform init"
        fi
        ;;
esac

if [[ $REPLY == "2" ]] || [[ $REPLY == "3" ]]; then
    echo ""
    echo "✅ Provider binary kept at: ./terraform-provider-graylog"
    echo ""
    echo "For developer override mode:"
    echo "  cd example-local-usage/"
    echo "  ./use-dev-mode.sh"
    echo "  ./dev-plan.sh"
    echo ""
    echo "Or use the Makefile:"
    echo "  cd example-local-usage/"
    echo "  make plan"
fi

if [[ $REPLY == "4" ]]; then
    echo ""
    echo "Provider binary available at: ./terraform-provider-graylog"
    echo "You can manually install it later or use developer override mode."
fi