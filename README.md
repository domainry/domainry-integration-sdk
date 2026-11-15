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
original request. This reads owner evidence and, only while the exact request is
uncertain, may use the Provider's declared read-only reconciliation lookup.
It never invokes the mutation again, and `not_found` is not permission to replay.
`uncertain` includes in-flight and crash-interrupted writes; neither it nor
`failed` permits automatic re-execution. A success receipt acknowledges the provider effect;
mail acceptance and calendar notification requests do not prove delivery.
Hosts must recheck their current Identity/tool permission before executing or
showing receipts. Integration separately checks current account ownership,
revision, state and actual OAuth grants.

Embedded hosts use this Go port directly. SaaS uses three bounded authenticated
service POST endpoints (`write-access`, `write`, `write-receipt`), with redirects
and cookies disabled and no automatic client retry. These endpoints are not
part of the public product HTTP adapter or browser SDK. Product confirmation
and execution identities remain the host's responsibility.

## Small generated Provider content

`Operations.Call` accepts optional `ProviderCallContent` for at most 256 KiB of
host-generated content, alongside at most 256 KiB of operation JSON metadata.
Use `NewProviderCallContent` to freeze bytes/hash/MIME, and
`NewProviderCallTarget` with the current owner `Connection` plus the official
operation descriptor. The target binds configuration and secret references as
well as the existing revision token, including changes in the same millisecond.
Named connection, authenticated host actor/role/source and sensitive persistence
are mandatory. A guarded JSON-only mutation may omit content.

The embedded and SaaS paths use the same request. The JSON wire format bounds
base64 bytes; Integration turns them into the existing Connector content reader.
No generic upload route, new provider identity or file-storage service is added.
Large files must use their existing streaming/file or Data Exchange owner.

The host must first commit its immutable publication record and stable effect
identity. Integration claims that identity before I/O. Replays cannot change the
actor, source, target, metadata or content; failed, running and crash-interrupted
claims never execute again. Invocation evidence retains a guarded fingerprint,
not content, operation payload, response body or secret values. A succeeded
replay returns the original invocation without its sensitive response. A timeout
does not prove the effect was rejected; query the original invocation and use
read-only source inspection. Never generate a new request ID to bypass an
unresolved write. Provider acceptance is separate from indexing or delivery.

The SaaS client checks returned routing and guarded identity and never retries
the request automatically. Business approval, Case scope, index reconciliation
and the local publication lifecycle remain the host's responsibility.
