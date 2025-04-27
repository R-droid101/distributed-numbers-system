#!/bin/bash

set -e

NAMESPACE=numbers-system

echo "Checking Kubernetes status for namespace: $NAMESPACE"

echo ""
echo "Pods:"
kubectl get pods -n $NAMESPACE

echo ""
echo "Services:"
kubectl get svc -n $NAMESPACE

echo ""
echo "Deployments:"
kubectl get deployments -n $NAMESPACE

echo ""
echo "Secrets:"
kubectl get secrets -n $NAMESPACE

echo ""
echo "✅ Status check complete."
