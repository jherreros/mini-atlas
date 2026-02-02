# Copilot instructions for Shoulders

## Big picture architecture
- The repo is a 3-layer platform: `1-cluster/` creates a local kind cluster, `2-addons/` installs platform services via FluxCD, and `3-user-space/` contains team-facing examples.
- Flux kustomizations in `2-addons/flux/kustomizations.yaml` enforce install order: helm repositories → namespaces → helm releases → crossplane → gateway.
- Crossplane defines the IDP abstractions in `2-addons/manifests/crossplane/definitions/` (XRDs) and implements them with compositions in `2-addons/manifests/crossplane/compositions/`.

## Crossplane composition patterns
- Compositions use **Pipeline** mode and functions. Patch/transform is used for static resource wiring (`function-patch-and-transform`), while dynamic resource generation uses `function-go-templating` (see `state-store-composition.yaml`).
- `application-composition.yaml` maps `WebApplication` specs into a Deployment, Service, and Gateway API `HTTPRoute` with `spec.host` → `spec.hostnames[0]` and a `cilium-gateway` parent ref.
- `workspace-composition.yaml` creates a Namespace, a default-deny CiliumNetworkPolicy, and a Kyverno ClusterPolicy that prefixes workload names with the workspace name.
- `state-store-composition.yaml` conditionally emits CloudNativePG and Redis resources with Go templates; note the fixed `app-secret` secret and base64-encoded creds.

## Developer workflows
- Cluster bootstrap: `1-cluster/create-cluster.sh` creates a kind cluster named `shoulders`.
- Platform install: `2-addons/install-addons.sh` installs Cilium via Helm, ensures Flux CLI is present, then `flux install` + `kubectl apply -f 2-addons/flux/`.
- Observability access: README documents Grafana port-forward (`svc/kube-prometheus-stack-grafana` in `observability`).

## Conventions and integration points
- All developer-facing resources live in `3-user-space/`; use the examples in `3-user-space/team-a/` as canonical templates (e.g., `webapp.yaml`).
- Most XRDs are **Namespaced** (see `application-xrd.yaml`), so example resources should always set `metadata.namespace` to a workspace namespace. `Workspace` is cluster-scoped per `workspace-xrd.yaml`.
- Platform components are declared as Flux-managed manifests under `2-addons/manifests/` (helm repos/releases, namespaces, gateway, crossplane).
- Crossplane RBAC for composed resources is defined in `2-addons/manifests/crossplane/rbac/crossplane-composed-resources.yaml`—add permissions here when compositions introduce new APIs.
