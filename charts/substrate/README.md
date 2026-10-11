# substrate

Agent Substrate's control plane: the CRDs, the identity/CA bootstrap, the pod-certificate controller, `atelet`, `ate-controller`, `ate-api-server`, `atenet-router`, `atenet-egress`, a bundled S3-compatible object store (rustfs) for snapshots, and a default gVisor `SandboxConfig`.

## Prerequisites

- A Kubernetes cluster with the `ClusterTrustBundle` and `PodCertificateRequest` beta APIs enabled.
- Every image this chart references built from source into a registry you control — no pre-built image is published anywhere for any component here. `ateSetup.image`, `podCertificateController.image`, `atelet.image`, `atelet.snapshotPluginImage`, `ateController.image`, `ateApiServer.image`, `ateApiServer.snapshotPluginImage`, `atenetRouter.image`, `atenetEgress.image` all need a value.

Install this chart before `substrate-postgres` — its pod-certificate controller issues the certs that chart's PostgreSQL pod needs to start. `ate-api-server` sits in `CreateContainerConfigError` until `substrate-postgres`'s Secrets exist too; don't `--wait` on either chart alone.

## Identity and CA bootstrap

Unlike PostgreSQL's role/schema setup, this chart's CA pools, JWT authority pool, and authentication config are created automatically by a `pre-install,pre-upgrade` hook that runs the equivalent `ate-setup create` steps for you — no manual step needed here.

## Values

| Value | Default | Description |
|---|---|---|
| `otel.endpoint` | `http://localhost:4317` | Shared OTLP endpoint every component exports through; a harmless default since logs/metrics collection is out of scope for this chart |
| `atelet.imageCredentialProviderBinDir`, `atelet.imageCredentialProviderConfig` | GKE paths | Kubelet image-credential-provider paths; set these for your node distribution |
| `ateApiServer.replicas` | `2` | `ate-api-server` replica count |
| `ateApiServer.actorJWTIssuer` | this chart's own `api` Service | Issuer URL placed in actor JWTs; must be HTTPS with no query or fragment |
| `atenetRouter.envoyImage`, `atenetEgress.envoyImage` | pinned digest | Envoy builds for the router and egress gateway |
| `gcpManagedPrometheus.enabled` | `false` | Renders a GKE-only `PodMonitoring` for the router's Envoy sidecar |
| `rustfs.accessKey`, `rustfs.secretKey` | generated | Left unset; the chart generates and persists a random pair on first install rather than shipping a fixed credential |
| `postgres.readWriteConnectionStringSecretRef.name`, `postgres.ownerConnectionStringSecretRef.name` | `substrate-postgres-readwrite`, `substrate-postgres-owner` | Secrets this chart reads the PostgreSQL connection strings from, each under the key `connectionString` |
| `postgres.schema`, `postgres.ownerRole`, `postgres.readWriteRole` | `substrate`, `substrate_owner`, `substrate_readwrite` | Must match the `substrate-postgres` chart's corresponding `postgres.*` values |
| `sandboxConfig.gvisor` | see `values.yaml` | The default gVisor `SandboxConfig`'s pinned pause image and per-architecture asset URLs |
