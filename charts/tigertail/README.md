# Tiger-Tail Microblog Helm Chart

This Helm chart deploys the Tiger-Tail Microblog application on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure (if using persistent storage)

## Getting Started

### Add the Helm repository

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### Install the chart

```bash
# Update dependencies first
helm dependency update ./charts/tigertail

# Install the chart with the release name "my-tigertail"
helm install my-tigertail ./charts/tigertail
```

## Configuration

The following table lists the configurable parameters of the Tiger-Tail chart and their default values.

| Parameter                                    | Description                                           | Default                        |
|----------------------------------------------|-------------------------------------------------------|--------------------------------|
| `replicaCount`                               | Number of replicas                                    | `1`                            |
| `image.repository`                           | Image repository                                      | `yourorg/tiger-tail-microblog` |
| `image.tag`                                  | Image tag                                             | `""`                           |
| `image.pullPolicy`                           | Image pull policy                                     | `IfNotPresent`                 |
| `imagePullSecrets`                           | Image pull secrets                                    | `[]`                           |
| `nameOverride`                               | String to partially override the fullname template    | `""`                           |
| `fullnameOverride`                           | String to fully override the fullname template        | `""`                           |
| `serviceAccount.create`                      | Specifies whether a service account should be created | `true`                         |
| `serviceAccount.annotations`                 | Annotations to add to the service account             | `{}`                           |
| `serviceAccount.name`                        | The name of the service account to use                | `""`                           |
| `podAnnotations`                             | Pod annotations                                       | `{}`                           |
| `podSecurityContext.enabled`                 | Enable pod security context                           | `true`                         |
| `podSecurityContext.runAsUser`               | User ID for the pod                                   | `65532`                        |
| `podSecurityContext.runAsGroup`              | Group ID for the pod                                  | `65532`                        |
| `podSecurityContext.fsGroup`                 | Group ID for the volume                               | `65532`                        |
| `securityContext`                            | Container security context                            | See values.yaml                |
| `service.type`                               | Kubernetes Service type                               | `ClusterIP`                    |
| `service.port`                               | Service HTTP port                                     | `8080`                         |
| `ingress.enabled`                            | Enable ingress controller resource                    | `false`                        |
| `ingress.className`                          | IngressClass that will be used                        | `""`                           |
| `ingress.annotations`                        | Ingress annotations                                   | `{}`                           |
| `ingress.hosts`                              | Ingress accepted hostnames                            | See values.yaml                |
| `ingress.tls`                                | Ingress TLS configuration                             | `[]`                           |
| `resources`                                  | CPU/Memory resource requests/limits                   | `{}`                           |
| `autoscaling.enabled`                        | Enable autoscaling                                    | `false`                        |
| `autoscaling.minReplicas`                    | Minimum number of replicas                            | `1`                            |
| `autoscaling.maxReplicas`                    | Maximum number of replicas                            | `100`                          |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization percentage                     | `80`                           |
| `nodeSelector`                               | Node labels for pod assignment                        | `{}`                           |
| `tolerations`                                | Tolerations for pod assignment                        | `[]`                           |
| `affinity`                                   | Affinity for pod assignment                           | `{}`                           |
| `config.server.port`                         | Application server port                               | `8080`                         |
| `config.server.host`                         | Application server host                               | `"0.0.0.0"`                    |
| `config.server.baseURL`                      | Application base URL                                  | `"http://localhost:8080"`      |
| `config.database.host`                       | Database host                                         | `"postgres"`                   |
| `config.database.port`                       | Database port                                         | `5432`                         |
| `config.database.user`                       | Database user                                         | `"postgres"`                   |
| `config.database.password`                   | Database password                                     | `"postgres"`                   |
| `config.database.name`                       | Database name                                         | `"tigertail"`                  |
| `config.database.sslMode`                    | Database SSL mode                                     | `"disable"`                    |
| `config.cache.enabled`                       | Enable cache                                          | `false`                        |
| `config.cache.host`                          | Cache host                                            | `"redis"`                      |
| `config.cache.port`                          | Cache port                                            | `6379`                         |
| `config.cache.password`                      | Cache password                                        | `""`                           |
| `config.cache.db`                            | Cache database index                                  | `0`                            |
| `postgresql.enabled`                         | Deploy PostgreSQL                                     | `true`                         |
| `postgresql.auth.username`                   | PostgreSQL username                                   | `"postgres"`                   |
| `postgresql.auth.password`                   | PostgreSQL password                                   | `"postgres"`                   |
| `postgresql.auth.database`                   | PostgreSQL database                                   | `"tigertail"`                  |
| `postgresql.service.port`                    | PostgreSQL service port                               | `5432`                         |
| `redis.enabled`                              | Deploy Redis                                          | `false`                        |
| `redis.auth.password`                        | Redis password                                        | `""`                           |
| `redis.master.service.port`                  | Redis service port                                    | `6379`                         |

## Dependencies

This chart depends on:

- [PostgreSQL](https://github.com/bitnami/charts/tree/master/bitnami/postgresql) - For database storage
- [Redis](https://github.com/bitnami/charts/tree/master/bitnami/redis) - For caching (optional)

## Security

This chart follows security best practices:

- Runs as non-root user (UID 65532)
- Uses read-only root filesystem
- Drops all capabilities
- Configures proper security contexts

## Probes

The application includes:

- Liveness probe: Checks if the application is running properly
- Readiness probe: Ensures the application is ready to accept traffic

## Scaling

Horizontal Pod Autoscaling can be enabled by setting `autoscaling.enabled=true`.
