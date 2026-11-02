# substrate-crds

Agent Substrate's CustomResourceDefinitions: `WorkerPool`, `SandboxConfig`, `CSIDriverConfig`.

Kept as a separate chart, as plain `templates/`, rather than under `substrate`'s special `crds/` directory. Helm only installs a chart's `crds/` directory once and never updates it on `helm upgrade` — a permanent Helm behavior, not a bug. A chart's CRDs tend to gain fields over time, so a chart that vendors them under `crds/` means every existing install silently stops receiving schema updates after its first install. Templates in an ordinary chart are reconciled on every upgrade, so this chart is the one and only place these CRDs are managed: install or upgrade it on its own, independent of the `substrate` chart's own release cadence.

## Install

```bash
helm upgrade --install substrate-crds ./charts/substrate-crds
```

Install or upgrade this chart before installing or upgrading `substrate`.

## Uninstall

CRD removal deletes every CR of that kind cluster-wide. Uninstall `substrate` first, then this chart only when you're sure no `WorkerPool`/`SandboxConfig`/`CSIDriverConfig` anywhere in the cluster is still needed.
