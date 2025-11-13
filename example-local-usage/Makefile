# Makefile for Graylog Terraform Project

.PHONY: help build init plan apply destroy refresh import-help clean validate fmt

# Default target
help:
	@echo "Graylog Terraform Project - Available targets:"
	@echo ""
	@echo "  make build        - Build the provider binary"
	@echo "  make init         - Initialize Terraform providers"
	@echo "  make validate     - Validate Terraform configuration"
	@echo "  make fmt          - Format Terraform files"
	@echo "  make plan         - Show planned changes"
	@echo "  make apply        - Apply Terraform changes"
	@echo "  make destroy      - Destroy all managed resources (CAUTION!)"
	@echo "  make refresh      - Refresh state from Graylog"
	@echo "  make discover     - Discover existing Graylog resources"
	@echo "  make import-help  - Generate import commands"
	@echo "  make clean        - Remove Terraform files and state"
	@echo ""
	@echo "  make setup        - Initial setup (build provider, copy tfvars, init)"
	@echo ""

# Build the provider
build:
	@echo "Building the Graylog provider..."
	@cd .. && go build -o terraform-provider-graylog ./cmd/terraform-provider-graylog
	@echo "✅ Provider built successfully at ../terraform-provider-graylog"

# Initial setup
setup: build
	@if [ ! -f terraform.tfvars ]; then \
		echo "Creating terraform.tfvars from example..."; \
		cp terraform.tfvars.example terraform.tfvars; \
		echo "Please edit terraform.tfvars with your Graylog credentials"; \
	else \
		echo "terraform.tfvars already exists"; \
	fi
	@$(MAKE) init

# Initialize Terraform
init:
	@echo "Setting up local provider with developer overrides..."
	@./use-dev-mode.sh

# Validate configuration
validate:
	@echo "Validating Terraform configuration..."
	@TF_CLI_CONFIG_FILE=./dev.tfrc terraform validate

# Format Terraform files
fmt:
	@echo "Formatting Terraform files..."
	@terraform fmt -recursive

# Plan changes
plan:
	@echo "Planning Terraform changes..."
	@TF_CLI_CONFIG_FILE=./dev.tfrc terraform plan

# Apply changes
apply:
	@echo "Applying Terraform changes..."
	@TF_CLI_CONFIG_FILE=./dev.tfrc terraform apply

# Destroy resources
destroy:
	@echo "WARNING: This will destroy all managed resources!"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		TF_CLI_CONFIG_FILE=./dev.tfrc terraform destroy; \
	else \
		echo "Cancelled"; \
	fi

# Refresh state
refresh:
	@echo "Refreshing Terraform state from Graylog..."
	@TF_CLI_CONFIG_FILE=./dev.tfrc terraform refresh

# Discover resources
discover: refresh
	@echo ""
	@echo "=== Provider Status ==="
	@TF_CLI_CONFIG_FILE=./dev.tfrc terraform output provider_status 2>/dev/null || echo "Provider not configured"
	@echo ""
	@echo "Note: The provider doesn't have 'list all' data sources."
	@echo "You need to know specific resource IDs to import them."
	@echo ""
	@echo "To discover resources, you'll need to:"
	@echo "1. Get resource IDs from the Graylog web UI"
	@echo "2. Add data sources to main.tf with those IDs"
	@echo "3. Run 'make plan' to query them"

# Generate import commands
import-help:
	@echo "Generating import commands..."
	@./generate-imports.sh

# Clean Terraform files
clean:
	@echo "Cleaning Terraform files..."
	@echo "This will remove:"
	@echo "  - .terraform directory"
	@echo "  - terraform.tfstate files"
	@echo "  - generated import scripts"
	@read -p "Continue? (y/n): " confirm; \
	if [ "$$confirm" = "y" ]; then \
		rm -rf .terraform .terraform.lock.hcl; \
		rm -f terraform.tfstate terraform.tfstate.backup; \
		rm -f import-commands.sh imported-resources.txt; \
		echo "Cleaned!"; \
	else \
		echo "Cancelled"; \
	fi

# Check if credentials are configured
check-creds:
	@if [ ! -f terraform.tfvars ]; then \
		echo "ERROR: terraform.tfvars not found!"; \
		echo "Run 'make setup' first"; \
		exit 1; \
	fi

# Targets that need credentials
plan: check-creds
apply: check-creds
refresh: check-creds
destroy: check-creds