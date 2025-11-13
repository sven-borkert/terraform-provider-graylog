#!/bin/bash
# Refresh Terraform state with developer overrides

export TF_CLI_CONFIG_FILE=./dev.tfrc
terraform refresh "$@"