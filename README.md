# 4-Layers Kubernetes Platform Architecture

This repository implements a comprehensive four-layer architecture for a multi-tenant Kubernetes-based platform, designed to provide separation of concerns, scalability, and maintainability for platform engineering teams. The architecture follows GitOps principles and leverages modern Kubernetes tools to create a robust foundation for multi-tenant application deployment.

## Architecture Overview

The four-layer approach provides clear boundaries between infrastructure provisioning, platform capabilities, tenant isolation, and application deployment. This separation enables different teams to own and manage their respective layers while maintaining consistency and security across the platform.

```
┌─────────────────────┐
│  4. Application     │  ← Built by Platform Teams, used by Development Teams
│     Layer           │
├─────────────────────┤
│  3. Tenant          │  ← Built by Platform Teams, used by Development Teams
│     Layer           │
├─────────────────────┤
│  2. Addons          │  ← Built by Platform Teams
│     Layer           │
├─────────────────────┤
│  1. Infrastructure  │  ← Built by Platform/Infrastructure Teams
│     Layer           │
└─────────────────────┘
```

## Layers Detailed

### 1. Infrastructure Layer 🏗️

**Purpose**: Provisions and manages the underlying Kubernetes cluster infrastructure.

**Responsibilities**:
- Kubernetes cluster creation and configuration
- Node management and scaling
- Network and storage infrastructure
- Basic cluster security and RBAC

**Key Components**:
- **Kind Configuration**: Multi-node cluster setup with control plane and worker nodes
- **Cluster Creation Script**: Automated cluster provisioning
- **Prerequisites**: Kind and kubectl tooling

**Usage**:
```bash
cd 1-infrastructure
./create-cluster.sh  # Creates a Kind cluster named '4-layers'
```

**Cleanup**:
```bash
kind delete cluster --name 4-layers
```

### 2. Addons Layer 🔧

**Purpose**: Installs and configures essential platform services and capabilities.

**Responsibilities**:
- Service mesh and networking (Cilium)
- Ingress controllers (NGINX Ingress)
- Policy engines (Kyverno)
- Custom resource management (Kro)
- Platform-wide security policies

**Key Components**:
- **Cilium**: eBPF-based networking and security
- **NGINX Ingress**: HTTP/HTTPS ingress management
- **Kyverno**: Policy-as-code for Kubernetes
- **Kro**: Custom resource graph definitions

**Installation**:
```bash
cd 2-addons
./install-addons.sh  # Installs all platform addons via Helm
```

**Resource Graph Definitions**:
- `application.yaml`: Defines the application deployment pattern
- `workspace.yaml`: Defines tenant workspace resources

### 3. Tenant Layer 🏢

**Purpose**: Creates and manages tenant-specific resources with proper isolation and governance.

**Responsibilities**:
- Namespace provisioning per tenant
- Network policies for tenant isolation
- Resource quotas and limits
- Tenant-specific RBAC policies
- Naming convention enforcement

**Key Features**:
- **Network Isolation**: Default-deny Cilium network policies
- **Naming Enforcement**: Kyverno policies requiring tenant prefixes
- **Resource Boundaries**: Namespace-based tenant separation

**Tenant Creation Process**:
1. Apply the workspace ResourceGraphDefinition (installed via addons layer)
2. Create tenant instances using the Workspace custom resource

**Example**:
```bash
# Tenant resources are created by applying workspace instances
kubectl apply -f 3-tenant/examples/team-a-workspace.yaml
```

**Generated Resources per Tenant**:
- Dedicated namespace
- Default-deny network policy
- Naming convention enforcement policy

### 4. Application Layer 🚀

**Purpose**: Defines standardized application deployment patterns and manages application instances.

**Responsibilities**:
- Application deployment templates
- Service exposure patterns
- Ingress configuration
- Application-specific configurations
- Development team self-service

**Key Features**:
- **Standardized Deployment**: Consistent patterns for Deployment, Service, and Ingress
- **Template-driven**: Parameterized application definitions
- **Multi-tenancy Ready**: Applications deployed within tenant boundaries

**Application Deployment Process**:
1. The application ResourceGraphDefinition is installed via the addons layer
2. Development teams create application instances using the MyApp custom resource

**Example**:
```bash
# Deploy an application instance
kubectl apply -f 4-application/examples/my-app.yaml
```

**Application Template Includes**:
- Kubernetes Deployment with configurable replicas and image
- ClusterIP Service for internal communication
- Ingress resource for external access

## Repository Structure

```
4-layers/
├── 1-infrastructure
│   ├── create-cluster.sh
│   └── kind-config.yaml
├── 2-addons
│   ├── flux
│   │   ├── git-repository.yaml
│   │   └── kustomization.yaml
│   ├── install-addons.sh
│   └── manifests
│       ├── helm-releases
│       ├── helm-repositories
│       ├── kustomization.yaml
│       └── ResourceGraphDefinitions
├── 3-tenant
│   └── examples
│       └── team-a-workspace.yaml
├── 4-application
│   └── examples
│       └── my-app.yaml
└── README.md
```

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) - Required for Kind
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) - Kubernetes in Docker
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) - Kubernetes CLI
- [Helm](https://helm.sh/docs/intro/install/) - Package manager for Kubernetes

### Quick Start

Follow these steps to set up the complete platform:

1. **Create the Infrastructure**:
   ```bash
   cd 1-infrastructure
   ./create-cluster.sh
   ```

2. **Install Platform Addons**:
   ```bash
   cd ../2-addons
   ./install-addons.sh
   ```

3. **Create a Tenant Workspace**:
   ```bash
   cd ../3-tenant
   kubectl apply -f examples/team-a-workspace.yaml
   ```

4. **Deploy an Application**:
   ```bash
   cd ../4-application
   kubectl apply -f examples/my-app.yaml
   ```

5. **Verify the Deployment**:
   ```bash
   kubectl get all -n team-a
   kubectl get ingress -n team-a
   ```

## Key Technologies

- **[Kind](https://kind.sigs.k8s.io/)**: Local Kubernetes clusters using Docker containers
- **[Cilium](https://cilium.io/)**: eBPF-based networking, observability, and security
- **[NGINX Ingress](https://kubernetes.github.io/ingress-nginx/)**: Kubernetes ingress controller
- **[Kyverno](https://kyverno.io/)**: Policy engine for Kubernetes
- **[Kro](https://kro.run/)**: Kubernetes resource orchestration and custom resource management

## Security Features

- **Network Segmentation**: Cilium network policies provide tenant isolation
- **Policy Enforcement**: Kyverno ensures compliance with naming conventions and security policies
- **Resource Boundaries**: Namespace-based separation with potential for resource quotas
- **Default-Deny**: All tenant workspaces start with restrictive network policies

## Multi-Tenancy

The platform provides soft multi-tenancy through:

- **Namespace Isolation**: Each tenant gets a dedicated namespace
- **Network Policies**: Default-deny traffic rules with selective allow rules
- **Naming Conventions**: Enforced prefixing to prevent resource conflicts
- **Policy Isolation**: Tenant-specific policies for governance

## Customization

### Adding New Addons

Add more `HelmReleases` under `2-addons/manifests/helm-releases`, then add the file name to `2-addons/manifests/helm-releases/kustomization.yaml`.

### Custom Application Patterns

Edit `2-addons/manifests/ResourceGraphDefinitions/application.yaml` to modify the application template or create new ResourceGraphDefinitions for different application types.

### Tenant Policies

Extend `2-addons/manifests/ResourceGraphDefinitions/workspace.yaml` to add additional tenant-specific resources like ResourceQuotas, LimitRanges, or custom RBAC policies.

## Production Considerations

Before using this architecture in production, consider at least:

1. **Replace Kind** with a production-grade Kubernetes distribution
2. **Implement proper RBAC** with fine-grained permissions
3. **Add resource quotas and limits** to prevent resource exhaustion
4. **Configure persistent storage** solutions
5. **Implement backup and disaster recovery** procedures
6. **Add monitoring and alerting** capabilities
7. **Secure ingress** with proper TLS certificates
