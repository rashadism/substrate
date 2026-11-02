# substrate-crds

Agent Substrate's CustomResourceDefinitions: `WorkerPool`, `SandboxConfig`, `CSIDriverConfig`.

Kept as a separate chart, as plain `templates/`, rather than under `substrate`'s special `crds/` directory — Helm only installs a `crds/` directory once and never updates it on `helm upgrade`. As ordinary templates here, CRD schema changes actually reach existing installs.

## Install

```bash
helm upgrade --install substrate-crds ./charts/substrate-crds
```

Install or upgrade this chart before installing or upgrading `substrate`.

## Uninstall

CRD removal deletes every CR of that kind cluster-wide. Uninstall `substrate` first, then this chart only when you're sure no `WorkerPool`/`SandboxConfig`/`CSIDriverConfig` anywhere in the cluster is still needed.
