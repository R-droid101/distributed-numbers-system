# Kubernetes Deployment Guide for Distributed Numbers System

This document explains how to deploy the **Distributed Numbers System** into any Kubernetes cluster.


# Prerequistes 

- You have a Kubernetes cluster (local k3d/k3s or cloud like EKS, GKE, Civo)
- `kubectl` is configured to point to the correct cluster
- `helm` is installed (Helm 3.x)


# Setup Steps

You can set it up manually or you can use the created scripts to setup the containers. I've listed down instructions for both:
## Using Setup and Uninstall Scripts

All helper scripts are inside `/k8s/hack/`:

Ensure to set the following env variables before running the scripts:
```
# 1. Set env vars
export GITHUB_USERNAME=yourname
export GITHUB_TOKEN=ghp_abc123...
export GITHUB_EMAIL=you@example.com
export DB_USERNAME=user
export DB_PASSWORD=pass
export AUTH_TOKEN=changeme
```

- `setup.sh` — install everything
- `uninstall.sh` — tear down everything

Usage:

```bash
./setup.sh
./uninstall.sh
```

## Setting up manually

### 1. Create Namespace

```bash
kubectl apply -f base/namespace.yaml
```

Namespace `numbers-system` will be created.


### 2. Create Secrets

Manually create the required Kubernetes Secrets:

#### Database Secret (`db-secret`)

```bash
kubectl create secret generic db-secret \
  --from-literal=username=your-db-username \
  --from-literal=password=your-db-password \
  -n numbers-system
```

#### Authorization Token Secret (`auth-token`)

```bash
kubectl create secret generic auth-token \
  --from-literal=token=your-secure-auth-token \
  -n numbers-system
```

**Note:** See [`secrets/README.md`](./secrets/README.md) for detailed explanation.


### 3. Install Charts (Helm Deployments)

#### Install Redis

```bash
helm install redis ./charts/redis --namespace numbers-system
```

#### Install Postgres

```bash
helm install postgres ./charts/postgres --namespace numbers-system
```

#### Setup migrations and wait till it is done
```bash
helm install migrate ../charts/migrate --namespace numbers-system || echo "Migration job may already exist."
kubectl wait --for=condition=complete --timeout=60s job/migrate-db -n numbers-system
```

#### Install Consumer

```bash
helm install consumer ./charts/consumer --namespace numbers-system
```

#### Install Publisher

```bash
helm install publisher ./charts/publisher --namespace numbers-system
```

---

# Useful Helm Commands

### Upgrade an Existing Release

```bash
helm upgrade <release-name> ../charts/<chart-folder> --namespace numbers-system
```
Example:
```bash
helm upgrade consumer ../charts/consumer --namespace numbers-system
```

### Uninstall a Release

```bash
helm uninstall <release-name> --namespace numbers-system
```
Example:
```bash
helm uninstall consumer --namespace numbers-system
```

### See All Installed Releases

```bash
helm list -n numbers-system
```

### See All Resources

```bash
kubectl get all -n numbers-system
```

---

# Reset Entire Deployment

If you want to delete everything and reset:

```bash
kubectl delete namespace numbers-system
```

This will clean up **all resources, secrets, pods, services** in one shot.