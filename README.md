# Domainry Integration SDK

Deployment-neutral contracts between Domainry Runtime and Integration running as an embedded module or remote SaaS service.

## Package layout

- The root package is the stable `Factory`, `Binding`, Catalog, Requirements, Delivery, and Web Push entrypoint.
- `modulehost` describes database, dialect, migration, provider registry, and secret-cipher capabilities borrowed by an embedded module.
- `remote` contains the SaaS client implementation.
- `saashost` describes SaaS composition.

The SDK intentionally has no public `persistence` package. Integration-owned tables and DML stay in the Integration implementation; Runtime consumes business capabilities through the root Binding.

Embedded modules use the host database, transaction boundary, SQL dialect, migration lock, and the host-owned `_schema_migrations` ledger. Run `go test ./...` before publishing an immutable SDK version.
