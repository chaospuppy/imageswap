#!/bin/bash
deploy_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Start KinD and deploy local registry
$deploy_dir/kind-with-registry.sh

# Deploy webhook resources
kustomize build $deploy_dir | kubectl apply -f -

# Create TLS keypair and CA bundle for webhook to use
$deploy_dir/webhook-create-signed-cert.sh

# Create mutating webhook with CA
cat $deploy_dir/mutatingwebhook.yaml | $deploy_dir/webhook-patch-ca-bundle.sh | kubectl apply -f -
