# Domainry Integration SDK

Deployment-neutral contracts between Domainry Runtime and Integration running as an embedded module or remote SaaS service.

## Package layout

- The root package is the stable `Factory`, `Binding`, Catalog projection,
  Requirements, Delivery, Management, Operations, inbound Event, Runtime
  Trigger receipt, local-worker, and Web Push entrypoint.
- `IntegrationAuthoringDomain` is the canonical Integration capability
  disclosure, including schemas, references, examples, repair errors and
  source symbols. `SpecializeIntegrationAuthoringCapability` derives
  Connector-specific examples from the source-owned Connector definition.
- `IntegrationHTTPAdapterContract` is the canonical Module HTTP route,
  governance and OpenAPI contract. Runtime and Plane may aggregate it, but
  must not re-author Integration paths or schemas.
- `modulehost` describes database, dialect, migration, provider registry, secret-cipher, and Runtime Trigger capabilities borrowed by an embedded module.
- `remote` contains the SaaS client implementation.
- `saashost` describes SaaS composition.
- `browser` publishes `@domainry/integration-client` for Integration-owned
  admin activity and Web Push clients; these APIs do not belong to the
  Runtime browser client.

The SDK intentionally has no public `persistence` package. Integration-owned tables and DML stay in the Integration implementation; Runtime consumes business capabilities through the root Binding.

Connector definitions and Provider schemas remain source-owned by
`domainry-connectors`. The Integration Catalog port exposes their materialized
projection without transferring ownership to Integration or Runtime.

Embedded modules use the host database, transaction boundary, SQL dialect,
migration lock, and the host-owned `_schema_migrations` ledger.

The contract tests reject Runtime source paths, non-`/integration` product
routes and Integration outbox capability keys.
Run `go test ./...` and `npm test --prefix browser` before publishing an
immutable SDK version.
