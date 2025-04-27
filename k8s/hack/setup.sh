#!/bin/bash

set -e

NAMESPACE="numbers-system"

echo "Creating namespace: $NAMESPACE"
kubectl create namespace $NAMESPACE || echo "✅ Namespace already exists"

echo ""
echo "Creating ghcr image pull secret"
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=$GITHUB_USERNAME \
  --docker-password=$GITHUB_TOKEN \
  --docker-email=$GITHUB_EMAIL \
  -n $NAMESPACE || echo "✅ ghcr-secret already exists"

echo ""
echo "Creating DB secret"
kubectl create secret generic db-secret \
  --from-literal=username=$DB_USERNAME \
  --from-literal=password=$DB_PASSWORD \
  -n $NAMESPACE || echo "✅ db-secret already exists"

echo ""
echo "Creating Auth Token secret"
kubectl create secret generic auth-token \
  --from-literal=token=$AUTH_TOKEN \
  -n $NAMESPACE || echo "✅ auth-token already exists"

echo ""
echo "Installing Redis via Helm"
helm upgrade --install redis ../charts/redis --namespace $NAMESPACE

echo ""
echo "Installing Postgres via Helm"
helm upgrade --install postgres ../charts/postgres --namespace $NAMESPACE

echo "Installing DB migration job..."
helm install migrate ../charts/migrate --namespace numbers-system || echo "Migration job may already exist."

kubectl wait --for=condition=complete --timeout=60s job/migrate-db -n numbers-system

echo ""
echo "Installing Consumer via Helm"
helm upgrade --install consumer ../charts/consumer --namespace $NAMESPACE

echo ""
echo "Installing Publishers via Helm"
helm upgrade --install publishers ../charts/publisher --namespace $NAMESPACE

echo ""
echo "Setup complete! "
