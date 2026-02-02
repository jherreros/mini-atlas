#!/bin/bash

set -o errexit

cd "$(dirname "$0")"

kind create cluster --config kind-config.yaml --name shoulders