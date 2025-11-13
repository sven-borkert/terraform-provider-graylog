#!/bin/bash

# Build and install terraform-provider-graylog locally

set -e

echo "Building terraform-provider-graylog..."

# Build the provider
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

echo "Build completed successfully!"

# Ask if user wants to install locally
read -p "Do you want to install the provider to ~/.terraform.d/plugins? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Detect OS and architecture
    OS=$(go env GOOS)
    ARCH=$(go env GOARCH)
    VERSION="3.0.0"

    # Create plugin directory
    PLUGIN_DIR="$HOME/.terraform.d/plugins/terraform-provider-graylog/graylog/$VERSION/${OS}_${ARCH}"
    echo "Creating plugin directory: $PLUGIN_DIR"
    mkdir -p "$PLUGIN_DIR"

    # Copy provider
    echo "Installing provider..."
    cp terraform-provider-graylog "$PLUGIN_DIR/"

    echo "Provider installed successfully!"
    echo ""
    echo "You can now use the provider in your Terraform configuration with:"
    echo ""
    echo "terraform {"
    echo "  required_providers {"
    echo "    graylog = {"
    echo "      source  = \"terraform-provider-graylog/graylog\""
    echo "      version = \"$VERSION\""
    echo "    }"
    echo "  }"
    echo "}"
else
    echo ""
    echo "Provider binary available at: ./terraform-provider-graylog"
    echo "You can manually copy it to your Terraform plugins directory or use development overrides."
fi