# substrate-postgres

Bundled, certificate-authenticated PostgreSQL for Agent Substrate's control-plane store. Deploys the StatefulSet and renders the two connection-string Secrets the `substrate` chart reads (`substrate.postgres.readWriteConnectionStringSecretRef`/`ownerConnectionStringSecretRef`).

## Prerequisites

Install before this chart:
- The `substrate` chart's CRDs and identity-bootstrap hook (creates `postgres-ca-pool` among the other CA pools).
- The pod-certificate controller, running and healthy — this StatefulSet's pod stays pending without it.

`substrate`'s own `ate-api-server` pods sit in `CreateContainerConfigError` until this chart's connection Secrets exist — install `substrate` first.

## One manual step after install

This chart does not create the database roles, schema, or grants — only the server and the connection Secrets pointing at certificates that don't authorize anything until the roles exist. `helm install`/`helm upgrade` prints the exact command to run once the StatefulSet is ready (see `templates/NOTES.txt`; re-print it any time with `helm get notes <release>`). It runs the role/schema SQL straight from `pkg/postgressetup/setup.sql` in this same checkout — run it from the repo root. The step is idempotent, safe to re-run after an upgrade or a `postgres-0` recreation.

It creates the two roles, two passwordless login users bound to them, and the schema named under `postgres` in `values.yaml` — the `substrate` chart's own `postgres.schema`/`ownerRole`/`readWriteRole` values must match.

## Values

| Value | Default | Description |
|---|---|---|
| `image` | pinned digest | PostgreSQL image |
| `storage` | `20Gi` | Data volume size |
| `resources` | see `values.yaml` | Requests/limits for the `postgres` container |
| `connectionSecrets.readWrite` | `substrate-postgres-readwrite` | Name of the rendered read/write connection Secret |
| `connectionSecrets.owner` | `substrate-postgres-owner` | Name of the rendered owner connection Secret |
| `postgres.schema`, `postgres.ownerRole`, `postgres.ownerUser`, `postgres.readWriteRole`, `postgres.readWriteUser` | see `values.yaml` | Identities the role/schema SQL creates; must match the `substrate` chart's corresponding `postgres.*` values |
