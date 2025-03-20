# Tiger-Tail Microblog Helm Chart (ENV-based)

This Helm chart deploys the Tiger-Tail Microblog application on a Kubernetes cluster, using **environment variables** to configure the app (rather than JSON ConfigMaps). By default, this chart also deploys **PostgreSQL** and **Redis** via Bitnami subcharts for a complete out-of-the-box setup similar to `docker-compose.yaml`.


## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure (if using persistent storage)


## Getting Started

1. **Add the Bitnami repo** (for PostgreSQL/Redis subcharts):
   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm repo update
   ```

2. **Update subchart dependencies**:
   ```bash
   helm dependency update ./charts/tiger-tail-chart
   ```

3. **Install the chart** with the release name “my-tigertail”:
   ```bash
   helm install my-tigertail ./charts/tiger-tail-chart
   ```
   This deploys Tiger-Tail (with default environment variables) along with PostgreSQL and Redis enabled by default.

4. **Verify the deployment**:
   - Check pods: `kubectl get pods`
   - Port-forward or configure your Ingress to reach the service.  
   - Example (port-forward):
     ```bash
     kubectl port-forward svc/my-tigertail 8080:8080
     curl http://localhost:8080/api/posts
     ```


## Configuration

Below is a summary of the top-level parameters in `values.yaml` and their defaults.

| Parameter                                       | Description                                                                          | Default                         |
|-------------------------------------------------|--------------------------------------------------------------------------------------|---------------------------------|
| **Global/Chart Basics**                         |                                                                                      |                                 |
| `replicaCount`                                  | Number of replicas of the Tiger-Tail pod                                             | `1`                             |
| `image.repository`                              | Container image repository                                                           | `joobypm/tiger-tail`            |
| `image.tag`                                     | Container image tag                                                                  | `""` (uses `.Chart.AppVersion`) |
| `image.pullPolicy`                              | Container image pull policy                                                          | `IfNotPresent`                  |
| `imagePullSecrets`                              | Secret names for private registry credentials                                        | `[]`                            |
| `nameOverride`                                  | Override chart name in resources                                                     | `""`                            |
| `fullnameOverride`                              | Override full generated name                                                         | `""`                            |
| **Service Account**                             |                                                                                      |                                 |
| `serviceAccount.create`                         | Whether to create a new ServiceAccount                                               | `true`                          |
| `serviceAccount.annotations`                    | Annotations for the ServiceAccount                                                   | `{}`                            |
| `serviceAccount.name`                           | Name of the ServiceAccount                                                           | `""`                            |
| **Pod/Container Configuration**                 |                                                                                      |                                 |
| `podAnnotations`                                | Additional pod annotations                                                           | `{}`                            |
| `podSecurityContext.enabled`                    | Enable pod-level security context                                                    | `true`                          |
| `podSecurityContext.runAsUser`                  | User ID for the pod                                                                  | `65532`                         |
| `podSecurityContext.runAsGroup`                 | Group ID for the pod                                                                 | `65532`                         |
| `podSecurityContext.fsGroup`                    | File-system group ID for volumes                                                     | `65532`                         |
| `securityContext`                               | Container-level security context (capabilities, readOnlyRootFilesystem, etc.)        | See `values.yaml`               |
| **Service**                                     |                                                                                      |                                 |
| `service.type`                                  | Kubernetes Service type (ClusterIP, NodePort, LoadBalancer)                          | `ClusterIP`                     |
| `service.port`                                  | Tiger-Tail container listening port                                                  | `8080`                          |
| **Ingress**                                     |                                                                                      |                                 |
| `ingress.enabled`                               | Enable the Ingress resource                                                          | `false`                         |
| `ingress.className`                             | Ingress class name (if your cluster uses `ingressClassName` vs annotations)          | `""`                            |
| `ingress.annotations`                           | Additional Ingress annotations                                                       | `{}`                            |
| `ingress.hosts`                                 | Ingress host rules                                                                   | See `values.yaml`               |
| `ingress.tls`                                   | TLS configuration for the Ingress                                                    | `[]`                            |
| **Resources & Autoscaling**                     |                                                                                      |                                 |
| `resources`                                     | CPU/Memory requests & limits for the container                                       | `{}` (no defaults)              |
| `autoscaling.enabled`                           | Enable Horizontal Pod Autoscaler                                                     | `false`                         |
| `autoscaling.minReplicas`                       | HPA minimum replicas                                                                 | `1`                             |
| `autoscaling.maxReplicas`                       | HPA maximum replicas                                                                 | `100`                           |
| `autoscaling.targetCPUUtilizationPercentage`    | Target average CPU utilization (%)                                                   | `80`                            |
| `autoscaling.targetMemoryUtilizationPercentage` | Target average memory utilization (%)                                                | _not set by default_            |
| **Node Scheduling**                             |                                                                                      |                                 |
| `nodeSelector`                                  | Node labels for scheduling pods                                                      | `{}`                            |
| `tolerations`                                   | Tolerations for scheduling                                                           | `[]`                            |
| `affinity`                                      | Affinity rules                                                                       | `{}`                            |
| **ENV-based Config**                            | These environment variables control Tiger-Tail’s behavior (mirroring docker-compose) |                                 |
| `env.useRealDb`                                 | Whether to use real PostgreSQL (`"true"`) or stub (`"false"`)                        | `"true"`                        |
| `env.useRealRedis`                              | Whether to use real Redis (`"true"`) or stub (`"false"`)                             | `"true"`                        |
| `env.serverPort`                                | HTTP port inside the container                                                       | `"8080"`                        |
| `env.serverHost`                                | Host interface to listen on                                                          | `"0.0.0.0"`                     |
| `env.logLevel`                                  | Log level (`debug`, `info`, `warn`, etc.)                                            | `"info"`                        |
| **Secrets**                                     | Various credentials stored in a Kubernetes Secret                                    |                                 |
| `secrets.dbUser`                                | DB username                                                                          | `"postgres"`                    |
| `secrets.dbPassword`                            | DB password                                                                          | `"postgres"`                    |
| `secrets.dbName`                                | DB name                                                                              | `"tigertail"`                   |
| `secrets.authUsername`                          | Basic Auth username for admin post creation                                          | `"admin"`                       |
| `secrets.authPassword`                          | Basic Auth password for admin post creation                                          | `"password"`                    |
| `secrets.redisPassword`                         | Redis password (if needed)                                                           | `""`                            |
| **PostgreSQL Subchart**                         |                                                                                      |                                 |
| `postgresql.enabled`                            | Whether to deploy the Bitnami PostgreSQL subchart                                    | `true`                          |
| `postgresql.auth.username`                      | PostgreSQL user (mirrors `secrets.dbUser`)                                           | `"postgres"`                    |
| `postgresql.auth.password`                      | PostgreSQL password (mirrors `secrets.dbPassword`)                                   | `"postgres"`                    |
| `postgresql.auth.database`                      | PostgreSQL DB name (mirrors `secrets.dbName`)                                        | `"tigertail"`                   |
| `postgresql.primary.service.port`               | PostgreSQL service port                                                              | `5432`                          |
| **Redis Subchart**                              |                                                                                      |                                 |
| `redis.enabled`                                 | Whether to deploy the Bitnami Redis subchart                                         | `false`                         |
| `redis.auth.password`                           | Redis password (mirrors `secrets.redisPassword`)                                     | `""`                            |
| `redis.master.service.port`                     | Redis service port                                                                   | `6379`                          |


## Dependencies

- **[Bitnami PostgreSQL](https://github.com/bitnami/charts/tree/master/bitnami/postgresql)** (enabled by default)  
  Installs a PostgreSQL database. If disabled (`postgresql.enabled=false`), you must supply your own DB and override DB host/port via `.Values.secrets.db*` and `.Values.env.useRealDb`.
- **[Bitnami Redis](https://github.com/bitnami/charts/tree/master/bitnami/redis)** (disabled by default)  
  Installs a Redis cache. If you want caching, set `redis.enabled=true` and `env.useRealRedis="true"`.


## Security

This chart follows best practices:
- Runs as non-root (UID 65532)
- Read-only root filesystem
- Drops capabilities
- Proper security contexts


## Probes

Tiger-Tail includes two main probes by default:
- **Liveness**: `/livez` endpoint
- **Readiness**: `/readyz` endpoint

These are configured in the chart’s Deployment.


## Scaling

To enable Horizontal Pod Autoscaling (HPA), set `autoscaling.enabled=true`. You can then adjust the CPU/Memory targets under `autoscaling.targetCPUUtilizationPercentage` and `autoscaling.targetMemoryUtilizationPercentage`.


## Example Commands

1. **Basic installation**:
   ```bash
   helm dependency update ./charts/tiger-tail-chart
   helm install my-tigertail ./charts/tiger-tail-chart
   ```

2. **Disable PostgreSQL & Redis subcharts** (use external DB/cache) and override credentials:
   ```bash
   helm install my-tigertail ./charts/tiger-tail-chart \
     --set postgresql.enabled=false \
     --set redis.enabled=false \
     --set secrets.dbUser=myuser \
     --set secrets.dbPassword=secret \
     --set secrets.dbName=tigerdb \
     --set env.useRealDb="true" \
     --set env.useRealRedis="false"
   ```

3. **Override environment & secrets** in one-liner:
   ```bash
   helm install my-tigertail ./charts/tiger-tail-chart \
     --set secrets.authPassword="admin123" \
     --set env.logLevel="debug" \
     --set ingress.enabled=true
   ```

Install from OCI registry:
```bash
helm install tiger-tail-chart oci://registry-1.docker.io/joobypm/tiger-tail-chart -n tiger-tail \
  --create-namespace \
  --set postgresql.auth.username=postgres \
  --set postgresql.auth.password=postgres \
  --set postgresql.auth.database=tigertail \
  --set redis.auth.password="redispass" \
  --set secrets.redisPassword="redispass"
```

Update the chart:
```bash
helm repo update
helm upgrade tiger-tail-chart oci://registry-1.docker.io/joobypm/tiger-tail-chart \
  --namespace tiger-tail \
  --reuse-values \
  --set ingress.enabled=true \
  --set "ingress.hosts[0].host=lab.jooby.tv" \
  --set "ingress.hosts[0].paths[0].path=/" \
  --set "ingress.hosts[0].paths[0].pathType=ImplementationSpecific"
```

Uninstall the chart:
```bash
helm uninstall tiger-tail-chart -n tiger-tail
```
