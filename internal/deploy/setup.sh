#!/bin/bash

# Start KinD and deploy local registry
./kind-with-registry.sh

# Deploy webhook resources
kustomize build . | kubectl apply -f -

# Create TLS keypair and CA bundle for webhook to use
./webhook-create-signed-cert.sh

# Create mutating webhook with CA
cat mutatingwebhook.yaml | ./webhook-patch-ca-bundle.sh | kubectl apply -f -

# Create test pod to see if the MWH works
kubectl run alpine --image=alpine --restart=Never -n default --command -- sleep infinity
