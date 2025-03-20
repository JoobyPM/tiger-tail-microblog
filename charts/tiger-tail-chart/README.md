# Tiger-Tail Microblog Helm Chart

This Helm chart deploys the Tiger-Tail Microblog application on a Kubernetes cluster, using environment variables for configuration. By default, this chart also deploys PostgreSQL and Redis via Bitnami subcharts for a complete out-of-the-box setup.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support (if using persistent storage)

## Quick Start

```bash
# Add Bitnami repo for PostgreSQL/Redis subcharts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Update dependencies
helm dependency update ./charts/tiger-tail-chart

# Install the chart
helm install my-tigertail ./charts/tiger-tail-chart -n tiger-tail --create-namespace

# Verify and access
kubectl get pods -n tiger-tail
kubectl port-forward svc/my-tigertail 8080:8080 -n tiger-tail
curl http://localhost:8080/api/posts
```


## Key Configuration Parameters

For a complete list of configuration options, see the comments in `values.yaml`. Key parameters include:

| Category           | Key Parameters                                           |
|--------------------|----------------------------------------------------------|
| **Application**    | `replicaCount`, `image.*`, `env.*`                       |
| **Credentials**    | `secrets.*` (single source of truth for all credentials) |
| **Infrastructure** | `service.*`, `ingress.*`, `resources.*`, `autoscaling.*` |
| **Security**       | `podSecurityContext.*`, `securityContext.*`              |

## Dependencies

This chart uses the following Bitnami subcharts:

- **[PostgreSQL](https://github.com/bitnami/charts/tree/master/bitnami/postgresql)**: For database (enabled by default)
- **[Redis](https://github.com/bitnami/charts/tree/master/bitnami/redis)**: For caching (enabled by default)

See the official documentation for these charts for detailed configuration options.


## Common Usage Examples

```bash
# Basic installation with namespace creation
helm install my-tigertail ./charts/tiger-tail-chart -n tiger-tail --create-namespace

# Use external database
helm install my-tigertail ./charts/tiger-tail-chart \
  --set postgresql.enabled=false \
  --set secrets.dbUser=myuser \
  --set secrets.dbPassword=secret \
  --set secrets.dbName=tigerdb

# Install from OCI registry with custom settings
helm install tiger-tail-chart oci://registry-1.docker.io/joobypm/tiger-tail-chart \
  -n tiger-tail --create-namespace \
  --set secrets.authPassword="admin123" \
  --set env.logLevel="debug" \
  --set ingress.enabled=true \
  --set "ingress.hosts[0].host=vs1.jooby.pro" \
  --set "ingress.hosts[0].paths[0].path=/" \
  --set "ingress.hosts[0].paths[0].pathType=ImplementationSpecific"
```

For more examples and detailed configuration options, see the comments in `values.yaml`.
