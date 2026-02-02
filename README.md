# Shoulders

Shoulders is a reference implementation of an Internal Developer Platform (IDP) that demonstrates how to use Crossplane to provide a self-service platform for developers to create and manage their cloud-native applications and infrastructure on Kubernetes.

The name originates from the quote _"If I have seen further it is by standing on the shoulders of Giants"_ by Isaac Newton. 
The Shoulders platform is composed by a set of open source tools that work together to provide a platform that can be used to build and deploy applications on Kubernetes. Those applications will, then, run on the shoulders of the maintainers and contributors of all those open source tools.

## What is Shoulders?

Shoulders allows developers to declaratively provision and manage:

- **Workspaces** - Isolated tenant environments with network policies and naming conventions.
- **Web Applications** - Containerized applications with deployments, services, and Cilium-powered Ingress.
- **State Stores** - PostgreSQL databases (via CloudNativePG) and Redis caches.
- **Event Streams** - Full Kafka clusters and multiple topics (via Strimzi).
- **Observability** - Built-in LGTM stack (Loki, Grafana, Tempo, Mimir/Prometheus) for comprehensive monitoring.

## Architecture

Shoulders follows a multi-layered approach:

### 1. Cluster Layer (`1-cluster/`)
Creates the foundational Kubernetes cluster using **kind** (Kubernetes in Docker) for local development. 

### 2. Addons Layer (`2-addons/`)
Installs platform components using **FluxCD** for GitOps-based deployment:

- **Helm Repositories & Releases** - Core platform software.
- **Abstractions** - Custom resource definitions using **Crossplane** for composable infrastructure.

### 3. User Space (`3-user-space/`)
Developer-facing resources where teams can provision their applications and infrastructure using high-level abstractions.

## Key Technologies

- **[Crossplane](https://crossplane.io)** - Composable infrastructure for creating custom abstractions.
- **[FluxCD](https://fluxcd.io)** - GitOps continuous delivery for Kubernetes.
- **[Cilium](https://cilium.io)** - Cloud-native networking, security, and Gateway API implementation.
- **[Strimzi](https://strimzi.io)** - Kubernetes-native Apache Kafka operator.
- **[CloudNativePG](https://cloudnative-pg.io)** - PostgreSQL operator for Kubernetes.
- **[Kyverno](https://kyverno.io)** - Policy-as-code for Kubernetes security and governance.

## Observability

Shoulders comes with a pre-configured observability stack to help developers monitor their applications and infrastructure.

- **[Loki](https://grafana.com/oss/loki/)** - High-availability log aggregation system.
- **[Grafana](https://grafana.com/oss/grafana/)** - The open observability platform for visualization and dashboards.
- **[Tempo](https://grafana.com/oss/tempo/)** - Easy-to-use, high-scale distributed tracing backend.
- **[Prometheus](https://prometheus.io)** - Systems monitoring and alerting toolkit.
- **[Grafana Alloy](https://grafana.com/oss/alloy/)** - Our preferred collector for logs, metrics, and traces.

### Accessing Grafana

You can access the Grafana dashboard to view logs, metrics, and traces by port-forwarding to the Grafana service:

```bash
kubectl port-forward svc/kube-prometheus-stack-grafana -n observability 3000:80
```

Then, open [http://localhost:3000](http://localhost:3000) in your browser. The default credentials are `admin` / `prom-operator`.

## Quick Start


#### Option A (scripts)

### Prerequisites

- Docker
- kubectl
- Helm (v3+)
- kind

### 1. Create the Cluster

```bash
./1-cluster/create-cluster.sh
```

### 2. Install Platform Addons

```bash
./2-addons/install-addons.sh
```

This will:
- Install Cilium CNI with kube-proxy replacement.
- Bootstrap FluxCD.
- Deploy all platform components via GitOps.


#### Option B (CLI, single command):

```bash
cd shoulders-cli
go build -o shoulders
./shoulders up
```

### 3. Verify Installation

```bash
# Check cluster status
kubectl get nodes

# Check FluxCD status
flux get kustomizations

# Check platform components
kubectl get pods -A
```

## Using Shoulders

### Creating a Workspace

Workspaces provide isolated environments for teams:

```yaml
apiVersion: shoulders.io/v1alpha1
kind: Workspace
metadata:
  name: team-a
spec: {}
```

This creates:
- A dedicated namespace
- Network policies for isolation
- Naming convention enforcement via Kyverno policies

### Deploying a Web Application

```yaml
apiVersion: shoulders.io/v1alpha1
kind: WebApplication
metadata:
  name: my-app
  namespace: team-a
spec:
  image: nginx
  tag: latest
  replicas: 2
  host: my-app.example.com
```

This provisions:
- Kubernetes Deployment
- Service for internal communication
- Ingress for external access

### Provisioning State Stores

```yaml
apiVersion: shoulders.io/v1alpha1
kind: StateStore
metadata:
  name: team-a-db
  namespace: team-a
spec:
  database: team-a-01
```

This creates:
- PostgreSQL cluster with CloudNativePG
- Redis deployment and service
- Database credentials as Kubernetes secrets

### Provisioning Event Streams

```yaml
apiVersion: shoulders.io/v1alpha1
kind: EventStream
metadata:
  name: team-a-01
  namespace: team-a
spec:
  topics:
    - name: logs
    - name: events
      partitions: 5
      config:
        retention.ms: "604800000"
```

This provisions a dedicated Kafka cluster and multiple topics with appropriate partitioning and retention policies.

## Project Structure

```
shoulders/
├── 1-cluster/                 # Cluster creation
│   ├── create-cluster.sh      # Kind cluster setup script
│   └── kind-config.yaml       # Kind configuration
├── 2-addons/                  # Platform components
│   ├── flux/                  # FluxCD bootstrap config
│   ├── install-addons.sh      # Addon installation script
│   └── manifests/             # Kubernetes manifests
│       ├── crossplane/        # Crossplane XRDs and Compositions
│       ├── gateway/           # Gateway API resources
│       ├── helm-releases/     # Helm chart deployments
│       ├── helm-repositories/ # Helm repository configs
│       └── namespaces/        # Namespace definitions
└── 3-user-space/              # Developer workspace
    └── team-a/                # Example team workspace
        ├── workspace.yaml     # Workspace definition
        ├── webapp.yaml        # Web application
        ├── state-store.yaml   # DB and redis resources
        └── event-stream.yaml # Messaging topic
```

## Available Abstractions

### Workspace
- **Purpose**: Isolated tenant environment
- **Creates**: Namespace, network policies, naming conventions
- **Schema**: No required fields

### WebApplication
- **Purpose**: Containerized web application
- **Creates**: Deployment, Service, Ingress
- **Schema**: `image`, `tag`, `replicas`, `host`

### StateStore
- **Purpose**: Database and caching services
- **Creates**: PostgreSQL cluster, Redis deployment, secrets
- **Schema**: `database`

### EventStream
- **Purpose**: Kafka cluster and messaging topics
- **Creates**: Kafka cluster, NodePool, and multiple KafkaTopics
- **Schema**: `topics` array with `name`, `partitions`, `replicas`, and `config`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Test your changes with a fresh cluster
4. Submit a pull request

## Cleanup

To remove the cluster:

```bash
kind delete cluster --name shoulders
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
