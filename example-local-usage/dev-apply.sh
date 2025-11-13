#!/bin/bash
# Run terraform apply with developer overrides

export TF_CLI_CONFIG_FILE=./dev.tfrc
terraform apply "$@"