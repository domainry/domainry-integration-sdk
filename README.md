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
- `IntegrationHTTPAdapterContract` is the canonical Module HTTP route and
  governance contract. Current Go source and contract tests are authoritative;
  Runtime and Plane must not re-author Integration paths or payload types.
- `modulehost` describes database, dialect, migration, provider registry, secret-cipher, and Runtime Trigger capabilities borrowed by an embedded module.
- `remote` contains the SaaS client implementation.
- `saashost` describes SaaS composition.
- `browser` publishes `@domainry/integration-client` for Integration-owned
  admin activity and Web Push clients; these APIs do not belong to the
  Runtime browser client.

The SDK intentionally has no public `persistence` package. Integration-owned tables and DML stay in the Integration implementation; Runtime consumes business capabilities through the root Binding.

`ConnectionAccountReadsBinding` exposes a separate current-user read port.
`AuthorizeConnectionAccountRead` checks one exact operation contract against
current account ownership, readiness and actual provider-declared grants.
`ReadConnectionAccount` resolves credentials inside Integration and returns a
source reference plus bounded payload. A host supplies a freshly authorized
`ConnectionAccountSubject`; browser/model requests cannot provide identity,
scopes, credentials or configuration. Hosts must also recheck their current
Identity/tool permission before exposing live or retained results.

Read payloads are sensitive and are not persisted in Integration invocation
evidence. Replaying a succeeded request returns `payload_available: false` with
the original invocation reference. A fresh read uses a new request ID. Retained
results must be authorized again and match the current source/account revision.

Connector definitions and Provider schemas remain source-owned by
`domainry-connectors`. The Integration Catalog port exposes their materialized
projection without transferring ownership to Integration or Runtime.

Embedded modules use the host database, transaction boundary, SQL dialect,
migration lock, and the host-owned `_schema_migrations` ledger.

The contract tests reject Runtime source paths, non-`/integration` product
routes and Integration outbox capability keys.
Run `go test ./...` and `npm test --prefix browser` before publishing an
immutable SDK version.

`ConnectionAccountWritesBinding` is an optional host-only mutation port, separate
from account reads. `AuthorizeConnectionAccountWrite` returns the current source
revision for the host's authorization/confirmation. `WriteConnectionAccount`
accepts that frozen source, exact typed payload and stable host execution ID;
the owner rechecks account/grants, claims once and persists only a fingerprint
and a validated calendar/mail receipt. Changed input under the same actor and
request ID conflicts. Invocation IDs are independent of caller payload.

After a timeout/disconnect, use `ReadConnectionAccountWriteReceipt` with the
original request. This reads owner evidence without vendor I/O. `uncertain`
includes in-flight and crash-interrupted writes; neither it nor `failed` permits
automatic re-execution. A success receipt acknowledges the provider effect;
mail acceptance and calendar notification requests do not prove delivery.
Hosts must recheck their current Identity/tool permission before executing or
showing receipts. Integration separately checks current account ownership,
revision, state and actual OAuth grants.

Embedded hosts use this Go port directly. SaaS uses three bounded authenticated
service POST endpoints (`write-access`, `write`, `write-receipt`), with redirects
and cookies disabled and no automatic client retry. These endpoints are not
part of the public product HTTP adapter or browser SDK. Product confirmation
and execution identities remain the host's responsibility.
