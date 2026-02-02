# shoulders-cli

Developer CLI for the Shoulders Internal Developer Platform. The bootstrap flow uses Go-native APIs (Kind + Helm SDK) rather than shelling out to scripts.

## Requirements
- Go 1.25+
- Access to a Kubernetes cluster (local Kind cluster recommended)
- `kubectl`, `kind`, `helm`, and `docker` available in PATH for lifecycle commands

## Install (local)
```bash
cd shoulders-cli
go mod tidy
go build -o shoulders
```

## Usage
```bash
./shoulders up
./shoulders status
./shoulders workspace create team-a
./shoulders workspace use team-a
./shoulders app init hello --image nginx:1.26 --replicas 1
./shoulders app list
./shoulders infra add-db app-db --type postgres --tier dev
./shoulders infra add-stream events --topics "logs,events" --partitions 3 --replicas 3 \
	--config cleanup.policy=compact
./shoulders dashboard
./shoulders logs hello
```

## Configuration
The current workspace context is stored at `~/.shoulders/config.yaml`.

## Output formats
Use `-o table|json|yaml` for supported list and status commands.

## Notes
- `shoulders app init` supports `--dry-run` to emit YAML instead of applying it.
- `shoulders logs` attempts a Loki query first and falls back to direct pod log streaming (no `kubectl`).
- `shoulders up` provisions the cluster via the Kind Go API and installs Cilium + Flux without running shell scripts. It pulls the Cilium chart and Flux install manifest from their upstream URLs.
- `shoulders infra add-stream` supports `--partitions`, `--replicas`, and repeatable `--config key=value` entries.
