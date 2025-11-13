#!/bin/bash

# Build the Graylog provider for local testing
# This script builds the provider in the parent directory for use with developer override mode

set -e

echo "==================================================="
echo "Building Graylog Provider for Local Testing"
echo "==================================================="
echo ""

# Move to parent directory to build
cd ..

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Are you in the right directory?"
    echo "This script should be run from the example-local-usage directory."
    exit 1
fi

echo "Building provider binary..."
go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "Provider binary created at: ../terraform-provider-graylog"
    echo ""
    echo "You can now use the provider with developer override mode:"
    echo "  ./use-dev-mode.sh"
    echo "  ./dev-plan.sh"
    echo ""
    echo "Or with the Makefile:"
    echo "  make plan"
else
    echo "❌ Build failed. Please check the error messages above."
    exit 1
fi

# Return to example directory
cd example-local-usage/