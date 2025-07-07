# Mini-Atlas

A scaled-down version of Atlas, the Internal Developer Platform used at Banking Circle. Mini-Atlas provides developers with a self-service platform to create and manage their cloud-native applications and infrastructure on Kubernetes.

## What is Mini-Atlas?

Mini-Atlas allows developers to declaratively provision and manage:

- **Workspaces** - Isolated tenant environments with network policies and naming conventions
- **Web Applications** - Containerized applications with deployments, services, and ingress
- **Infrastructure** - PostgreSQL databases, Redis caches, and other backing services
- **Messaging** - Kafka topics on a shared message bus

## Architecture

Mini-Atlas follows a multi-layered approach:

### 1. Infrastructure Layer (`1-infrastructure/`)
Creates the foundational Kubernetes cluster using **kind** (Kubernetes in Docker) for local development. In production Atlas, this layer uses Terraform to provision AKS clusters.

### 2. Addons Layer (`2-addons/`)
Installs platform components using **FluxCD** for GitOps-based deployment:

- **Helm Repositories & Releases** - Core platform software
- **Shared Resources** - Multi-tenant infrastructure (Kafka cluster)
- **Abstractions** - Custom resource definitions using **KRO** (Kubernetes Resource Orchestrator)

### 3. User Space (`user-space/`)
Developer-facing resources where teams can provision their applications and infrastructure using high-level abstractions.

## Key Technologies

- **[KRO](https://kro.run)** - Kubernetes Resource Orchestrator for creating custom abstractions
- **FluxCD** - GitOps continuous delivery for Kubernetes
- **Cilium** - Cloud-native networking and security
- **Strimzi** - Kubernetes-native Apache Kafka
- **CloudNativePG** - PostgreSQL operator for Kubernetes
- **Kyverno** - Policy-as-code for Kubernetes security and governance

## Quick Start

### Prerequisites

- Docker
- kubectl
- Helm (v3+)
- kind

### 1. Create the Cluster

```bash
cd 1-infrastructure
./create-cluster.sh
```

### 2. Install Platform Addons

```bash
cd ../2-addons
./install-addons.sh
```

This will:
- Install Cilium CNI with kube-proxy replacement
- Bootstrap FluxCD
- Deploy all platform components via GitOps

### 3. Verify Installation

```bash
# Check cluster status
kubectl get nodes

# Check FluxCD status
flux get kustomizations

# Check platform components
kubectl get pods -A
```

## Using Mini-Atlas

### Creating a Workspace

Workspaces provide isolated environments for teams:

```yaml
apiVersion: kro.run/v1alpha1
kind: Workspace
metadata:
  name: team-a
spec:
  name: team-a
```

This creates:
- A dedicated namespace
- Network policies for isolation
- Naming convention enforcement via Kyverno policies

### Deploying a Web Application

```yaml
apiVersion: kro.run/v1alpha1
kind: WebApplication
metadata:
  name: my-app
  namespace: team-a
spec:
  name: my-app
  namespace: team-a
  image: nginx
  tag: latest
  replicas: 2
  host: my-app.example.com
```

This provisions:
- Kubernetes Deployment
- Service for internal communication
- Ingress for external access

### Provisioning Infrastructure

```yaml
apiVersion: kro.run/v1alpha1
kind: Infrastructure
metadata:
  name: team-a-db
  namespace: team-a
spec:
  name: team-a-db
  namespace: team-a
  database: team-a-01
```

This creates:
- PostgreSQL cluster with CloudNativePG
- Redis deployment and service
- Database credentials as Kubernetes secrets

### Creating Kafka Topics

```yaml
apiVersion: kro.run/v1alpha1
kind: Topic
metadata:
  name: team-a-events
  namespace: team-a
spec:
  name: team-a-events
  namespace: team-a
```

This provisions a Kafka topic on the shared cluster with appropriate partitioning and retention policies.

## Project Structure

```
mini-atlas/
├── 1-infrastructure/          # Cluster creation
│   ├── create-cluster.sh      # Kind cluster setup script
│   └── kind-config.yaml       # Kind configuration
├── 2-addons/                  # Platform components
│   ├── flux/                  # FluxCD bootstrap config
│   ├── install-addons.sh      # Addon installation script
│   └── manifests/             # Kubernetes manifests
│       ├── abstractions/      # KRO resource definitions
│       ├── helm-releases/     # Helm chart deployments
│       ├── helm-repositories/ # Helm repository configs
│       ├── namespaces/        # Namespace definitions
│       └── shared-resources/  # Multi-tenant resources
└── user-space/                # Developer workspace
    └── team-a/                # Example team workspace
        ├── workspace.yaml     # Workspace definition
        ├── webapp.yaml        # Web application
        ├── infra.yaml         # Infrastructure resources
        └── topic.yaml         # Messaging topic
```

## Available Abstractions

### Workspace
- **Purpose**: Isolated tenant environment
- **Creates**: Namespace, network policies, naming conventions
- **Schema**: `name: string`

### WebApplication
- **Purpose**: Containerized web application
- **Creates**: Deployment, Service, Ingress
- **Schema**: `name`, `namespace`, `image`, `tag`, `replicas`, `host`

### Infrastructure
- **Purpose**: Database and caching services
- **Creates**: PostgreSQL cluster, Redis deployment, secrets
- **Schema**: `name`, `namespace`, `database`

### Topic
- **Purpose**: Kafka messaging topic
- **Creates**: KafkaTopic with partitioning and retention
- **Schema**: `name`, `namespace`

## Differences from Production Atlas

| Feature | Mini-Atlas | Production Atlas |
|---------|------------|------------------|
| Cluster | Kind (local) | AKS (Terraform) |
| Abstractions | KRO | Complex custom operators |
| Scale | Single cluster | Multi-cluster |
| Networking | Cilium (basic) | Advanced mesh networking |
| Security | Basic policies | Enterprise security controls |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Test your changes with a fresh cluster
4. Submit a pull request

## Cleanup

To remove the cluster:

```bash
kind delete cluster --name mini-atlas
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
