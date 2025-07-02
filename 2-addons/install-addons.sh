#!/bin/bash

set -o errexit

# This script installs FluxCD.

if ! command -v flux &> /dev/null
then
    echo "Flux CLI not found. Installing..."
    curl -s https://fluxcd.io/install.sh | sudo bash
fi

if ! flux check --pre &> /dev/null
then
    echo "Flux pre-check failed. Please check your environment."
    exit 1
fi

if ! flux get kustomization flux-system &> /dev/null
then
    echo "Installing FluxCD..."
    flux install
    kubectl apply -f flux/
else
    echo "FluxCD already bootstrapped. Reconciling..."
    flux reconcile kustomization flux-system --with-source
fi
