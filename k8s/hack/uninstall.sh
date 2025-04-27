#!/bin/bash

set -e

echo "🗑️ Uninstalling Distributed Numbers System..."

# 1. Uninstall Helm Releases
helm uninstall publisher-1 --namespace numbers-system || echo "Publisher-1 may not be installed."
helm uninstall publisher-2 --namespace numbers-system || echo "Publisher-2 may not be installed."
helm uninstall publisher-3 --namespace numbers-system || echo "Publisher-3 may not be installed."
helm uninstall publisher-4 --namespace numbers-system || echo "Publisher-4 may not be installed."
helm uninstall publisher-5 --namespace numbers-system || echo "Publisher-5 may not be installed."
helm uninstall consumer --namespace numbers-system || echo "Consumer may not be installed."
helm uninstall postgres --namespace numbers-system || echo "Postgres may not be installed."
helm uninstall redis --namespace numbers-system || echo "Redis may not be installed."
helm uninstall migrate --namespace numbers-system || echo "Migrate job may not be installed."

# 2. Delete Secrets
kubectl delete secret db-secret -n numbers-system || echo "db-secret may not exist."
kubectl delete secret auth-token -n numbers-system || echo "auth-token may not exist."
kubectl delete secret ghcr-secret -n numbers-system || echo "ghcr-secret may not exist."

# 3. Delete Namespace
kubectl delete namespace numbers-system || echo "Namespace may already be deleted."

echo "Uninstall complete!"
