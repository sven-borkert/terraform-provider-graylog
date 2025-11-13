#!/bin/bash
# Run terraform plan with developer overrides

export TF_CLI_CONFIG_FILE=./dev.tfrc
terraform plan "$@"