package integrationsdk

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// AuthoringDomain is the deployment-neutral Integration contribution to the
// Runtime authoring catalog. Runtime may aggregate this value, but the
// capability routes, schemas, examples, and source evidence are owned here.
type AuthoringDomain struct {
	Key          string                `json:"key"`
	Capabilities []AuthoringCapability `json:"capabilities"`
}

type AuthoringCapability struct {
	Key                      string                       `json:"key"`
	Status                   string                       `json:"status"`
	Lifecycle                string                       `json:"lifecycle"`
	AllowedContexts          []string                     `json:"allowed_contexts,omitempty"`
	Parameters               []AuthoringParameter         `json:"parameters,omitempty"`
	Requires                 []string                     `json:"requires,omitempty"`
	Conflicts                []string                     `json:"conflicts,omitempty"`
	Permissions              []string                     `json:"permissions,omitempty"`
	AuditEvents              []string                     `json:"audit_events,omitempty"`
	ValidationEndpoint       string                       `json:"validation_endpoint,omitempty"`
	PreviewEndpoint          string                       `json:"preview_endpoint,omitempty"`
	SimulationEndpoint       string                       `json:"simulation_endpoint,omitempty"`
	ConfigurationRoutes      []string                     `json:"configuration_routes,omitempty"`
	ResourceOperations       *AuthoringResourceOperations `json:"resource_operations,omitempty"`
	ResourceKeyPathParameter string                       `json:"resource_key_path_parameter,omitempty"`
	Errors                   []AuthoringError             `json:"errors,omitempty"`
	Examples                 []AuthoringExample           `json:"examples,omitempty"`
	InputSchema              *AuthoringSchema             `json:"input_schema,omitempty"`
	OutputSchema             *AuthoringSchema             `json:"output_schema,omitempty"`
	OutputVariables          []AuthoringOutput            `json:"output_variables,omitempty"`
	ReferenceContracts       []AuthoringReference         `json:"reference_contracts,omitempty"`
	Execution                *AuthoringExecution          `json:"execution,omitempty"`
	Sources                  []AuthoringSource            `json:"sources"`
}

type AuthoringResourceOperations struct {
	PersistenceMode string                   `json:"persistence_mode"`
	Validate        string                   `json:"validate"`
	Upsert          string                   `json:"upsert"`
	UpsertHeaders   []AuthoringRequestHeader `json:"upsert_headers"`
	SuccessSchema   *AuthoringSchema         `json:"success_schema"`
	Get             string                   `json:"get"`
	Versions        string                   `json:"versions"`
	Simulate        string                   `json:"simulate,omitempty"`
	Rollback        string                   `json:"rollback,omitempty"`
	Delete          string                   `json:"delete,omitempty"`
}

type AuthoringRequestHeader struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	ValueSource string `json:"value_source"`
	Description string `json:"description"`
}

type AuthoringSchema struct {
	Schema               string                     `json:"$schema,omitempty"`
	Ref                  string                     `json:"$ref,omitempty"`
	Type                 string                     `json:"type,omitempty"`
	Properties           map[string]AuthoringSchema `json:"properties,omitempty"`
	Definitions          map[string]AuthoringSchema `json:"$defs,omitempty"`
	Required             []string                   `json:"required,omitempty"`
	Items                *AuthoringSchema           `json:"items,omitempty"`
	OneOf                []AuthoringSchema          `json:"oneOf,omitempty"`
	Enum                 []any                      `json:"enum,omitempty"`
	Const                any                        `json:"const,omitempty"`
	Default              any                        `json:"default,omitempty"`
	Format               string                     `json:"format,omitempty"`
	Minimum              *float64                   `json:"minimum,omitempty"`
	Maximum              *float64                   `json:"maximum,omitempty"`
	MinLength            *int                       `json:"minLength,omitempty"`
	MaxLength            *int                       `json:"maxLength,omitempty"`
	MinItems             *int                       `json:"minItems,omitempty"`
	MaxItems             *int                       `json:"maxItems,omitempty"`
	AdditionalProperties *bool                      `json:"additionalProperties,omitempty"`
	Description          string                     `json:"description,omitempty"`
}

type AuthoringOutput struct {
	Name        string `json:"name"`
	JSONPointer string `json:"json_pointer"`
	Type        string `json:"type"`
	VisibleTo   string `json:"visible_to"`
}

type AuthoringReference struct {
	Kind             string `json:"kind"`
	InputJSONPointer string `json:"input_json_pointer"`
	ScopeFrom        string `json:"scope_from,omitempty"`
	ResolverEndpoint string `json:"resolver_endpoint"`
}

type AuthoringExecution struct {
	ReadSet         []string `json:"read_set,omitempty"`
	WriteSet        []string `json:"write_set,omitempty"`
	BoundaryClass   string   `json:"boundary_class"`
	Transaction     string   `json:"transaction"`
	Idempotency     string   `json:"idempotency"`
	SideEffects     []string `json:"side_effects,omitempty"`
	SideEffectLevel string   `json:"side_effect_level"`
	Compensation    string   `json:"compensation,omitempty"`
	PermissionModel string   `json:"permission_model"`
	ChangeControl   string   `json:"change_control,omitempty"`
}

type AuthoringExample struct {
	Name               string         `json:"name"`
	Value              map[string]any `json:"value"`
	ExpectedErrorCodes []string       `json:"expected_error_codes,omitempty"`
}

type AuthoringParameter struct {
	Key           string         `json:"key"`
	Type          string         `json:"type"`
	Required      bool           `json:"required,omitempty"`
	Default       any            `json:"default,omitempty"`
	Enum          []string       `json:"enum,omitempty"`
	Minimum       *float64       `json:"minimum,omitempty"`
	Maximum       *float64       `json:"maximum,omitempty"`
	MinLength     *int           `json:"min_length,omitempty"`
	MaxLength     *int           `json:"max_length,omitempty"`
	RequiredWhen  map[string]any `json:"required_when,omitempty"`
	ConflictsWith []string       `json:"conflicts_with,omitempty"`
	ItemSchema    string         `json:"item_schema,omitempty"`
	Format        string         `json:"format,omitempty"`
	ReadOnly      bool           `json:"read_only,omitempty"`
}

type AuthoringError struct {
	Code          string   `json:"code"`
	FieldPath     string   `json:"field_path,omitempty"`
	ParameterKeys []string `json:"parameter_keys,omitempty"`
	MessageKey    string   `json:"message_key"`
}

type AuthoringSource struct {
	Kind   string `json:"kind"`
	Path   string `json:"path"`
	Symbol string `json:"symbol,omitempty"`
}

type AuthoringSelection struct {
	ProviderKey  string
	OperationKey string
}

type authoringConnector struct {
	Key        string                       `json:"key"`
	Providers  []authoringConnectorProvider `json:"providers,omitempty"`
	Operations []authoringOperation         `json:"operations,omitempty"`
}

type authoringConnectorProvider struct {
	Key          string           `json:"key"`
	ConfigFields []authoringField `json:"config_fields,omitempty"`
	SecretFields []authoringField `json:"secret_fields,omitempty"`
}

type authoringOperation struct {
	Key    string           `json:"key"`
	Input  []authoringField `json:"input,omitempty"`
	Output []authoringField `json:"output,omitempty"`
}

type authoringField struct {
	Key         string                   `json:"key"`
	Description string                   `json:"description,omitempty"`
	Type        string                   `json:"type"`
	Validation  authoringFieldValidation `json:"validation,omitempty"`
	Required    bool                     `json:"required"`
	Default     any                      `json:"default,omitempty"`
}

type authoringFieldValidation struct {
	MinLength int      `json:"min_length,omitempty"`
	MaxLength int      `json:"max_length,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Options   []string `json:"options,omitempty"`
}

var integrationConnectionStatuses = []string{"active", "configured", "degraded", "disabled", "draft", "verified"}
var integrationInvocationStatuses = []string{"cancelled", "failed", "queued", "running", "succeeded"}
var integrationEventStatuses = []string{"dead_letter", "failed", "ignored", "processed", "processing", "quarantined", "received"}

const (
	integrationAuthoringSourcePath    = "github.com/domainry/domainry-integration-sdk/authoring.go"
	exactActionSameKeyPermissionModel = "exact_action_same_key_permission"
)

// IntegrationAuthoringDomain returns the canonical source-owned Integration
// authoring contribution. It publishes only routes actually owned by the
// Integration HTTP surface; Runtime publication handoff is intentionally not
// represented as Integration outbox authoring.
func IntegrationAuthoringDomain() AuthoringDomain {
	capabilities := []AuthoringCapability{
		integrationCatalogCapability(),
		integrationBindingValidationCapability(nil, nil, nil),
		integrationConnectionCapability("integration.connection"),
		integrationConnectionCapability("integration.connection.rotate"),
		integrationConnectionCommandCapability("integration.connection.disable", "POST /tenant-admin/integrations/connections/{connectionKey}/disable", "integration_connection_disabled", "Management.SetConnectionStatus"),
		integrationConnectionCommandCapability("integration.connection.delete", "DELETE /tenant-admin/integrations/connections/{connectionKey}", "integration_connection_deleted", "Management.DeleteConnection"),
		integrationOperationTestCapability(nil),
		integrationInvocationListCapability(),
		integrationEventListCapability(),
		integrationEventReplayCapability(),
	}
	sort.Slice(capabilities, func(i, j int) bool { return capabilities[i].Key < capabilities[j].Key })
	return AuthoringDomain{Key: "integration", Capabilities: capabilities}
}

// SpecializeIntegrationAuthoringCapability closes an Integration capability
// schema and its examples against one Connector-owned catalog definition.
// connectorDefinition accepts the complete JSON definition so this SDK does
// not transfer Connector ownership into Runtime or Integration.
func SpecializeIntegrationAuthoringCapability(key string, connectorDefinition json.RawMessage, selection AuthoringSelection) (AuthoringCapability, bool, error) {
	var connector authoringConnector
	if err := json.Unmarshal(connectorDefinition, &connector); err != nil {
		return AuthoringCapability{}, false, fmt.Errorf("decode Connector authoring projection: %w", err)
	}
	key, selection.ProviderKey, selection.OperationKey = strings.TrimSpace(key), strings.TrimSpace(selection.ProviderKey), strings.TrimSpace(selection.OperationKey)
	switch key {
	case "integration.connection", "integration.connection.rotate":
		for index := range connector.Providers {
			provider := &connector.Providers[index]
			if provider.Key == selection.ProviderKey {
				capability := integrationConnectionCapability(key)
				capability.InputSchema = integrationConnectionInputSchema(provider)
				capability.Examples = integrationConnectionExamples(connector.Key, provider.Key, provider)
				return capability, true, nil
			}
		}
	case "integration.operation_test":
		for index := range connector.Operations {
			operation := &connector.Operations[index]
			if operation.Key == selection.OperationKey {
				return integrationOperationTestCapability(operation), true, nil
			}
		}
	case "integration.binding_validation":
		var provider *authoringConnectorProvider
		for index := range connector.Providers {
			if connector.Providers[index].Key == selection.ProviderKey {
				provider = &connector.Providers[index]
				break
			}
		}
		var operation *authoringOperation
		for index := range connector.Operations {
			if connector.Operations[index].Key == selection.OperationKey {
				operation = &connector.Operations[index]
				break
			}
		}
		return integrationBindingValidationCapability(&connector, provider, operation), true, nil
	}
	return AuthoringCapability{}, false, nil
}

func integrationCatalogCapability() AuthoringCapability {
	closed, open := false, true
	connection := integrationConnectionOutputItemSchema()
	return AuthoringCapability{
		Key: "integration.catalog", Status: "supported", Lifecycle: "owner_catalog_read_only", Permissions: []string{ActionIntegrationConnectorsList},
		ConfigurationRoutes: []string{"GET /tenant-admin/integrations/connectors"},
		InputSchema:         &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Properties: map[string]AuthoringSchema{}},
		OutputSchema: &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"connectors", "connections", "connections_available", "count"}, Properties: map[string]AuthoringSchema{
			"connectors":  {Type: "array", Items: &AuthoringSchema{Type: "object", AdditionalProperties: &open}},
			"connections": {Type: "array", Items: &connection}, "connections_available": {Type: "boolean"}, "count": {Type: "integer"},
		}},
		OutputVariables: []AuthoringOutput{{Name: "connectors", JSONPointer: "/connectors", Type: "connector_list", VisibleTo: "subsequent_capability_calls"}, {Name: "count", JSONPointer: "/count", Type: "integer", VisibleTo: "subsequent_capability_calls"}},
		Execution:       readExecution([]string{"integration.connector_catalog_projection"}),
		Examples:        readExamples(),
		Sources:         []AuthoringSource{{Kind: "owner_sdk", Path: integrationAuthoringSourcePath, Symbol: "IntegrationAuthoringDomain"}, {Kind: "owner_sdk", Path: "github.com/domainry/domainry-integration-sdk/sdk.go", Symbol: "Catalog.ListConnectorDefinitions"}, {Kind: "owner_catalog", Path: "github.com/domainry/domainry-connectors/catalog/catalog.go", Symbol: "Catalog"}},
	}
}

func integrationBindingValidationCapability(connector *authoringConnector, provider *authoringConnectorProvider, operation *authoringOperation) AuthoringCapability {
	closed, open := false, true
	connectionDraft := AuthoringSchema{Type: "object", AdditionalProperties: &open}
	input, output := AuthoringSchema{Type: "object", AdditionalProperties: &open}, AuthoringSchema{Type: "object", AdditionalProperties: &open}
	connectorKey, operationKey := AuthoringSchema{Type: "string"}, AuthoringSchema{Type: "string"}
	if connector != nil {
		connectorKey.Const, connectorKey.Enum = connector.Key, []any{connector.Key}
	}
	if provider != nil {
		connectionDraft = *integrationConnectionInputSchema(provider)
	}
	if operation != nil {
		operationKey.Const, operationKey.Enum = operation.Key, []any{operation.Key}
		input, output = integrationProtocolObjectSchema(operation.Input), integrationProtocolObjectSchema(operation.Output)
	}
	return AuthoringCapability{
		Key: "integration.binding_validation", Status: "supported", Lifecycle: "side_effect_free_validation", Requires: []string{"integration.catalog", "integration.connection"}, Permissions: []string{ActionIntegrationConnectionsValidate},
		Parameters:         []AuthoringParameter{{Key: "connector_key", Type: "connector_key", Required: true}, {Key: "operation_key", Type: "operation_key"}, {Key: "connection_key", Type: "connection_key"}, {Key: "connection_draft", Type: "integration_connection_draft"}, {Key: "input", Type: "operation_input"}, {Key: "output", Type: "operation_output"}},
		ValidationEndpoint: "POST /tenant-admin/integrations/connections/{connectionKey}/validate", ConfigurationRoutes: []string{"POST /tenant-admin/integrations/connections/{connectionKey}/validate"}, ResourceKeyPathParameter: "connectionKey",
		InputSchema: &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"connector_key"}, Properties: map[string]AuthoringSchema{
			"connector_key": connectorKey, "operation_key": operationKey, "connection_key": {Type: "string"}, "connection_draft": connectionDraft, "input": input, "output": output,
		}},
		OutputSchema: &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"valid", "connection_ready", "errors"}, Properties: map[string]AuthoringSchema{
			"valid": {Type: "boolean"}, "connection_ready": {Type: "boolean"}, "errors": {Type: "array", Items: &AuthoringSchema{Type: "object", AdditionalProperties: &open}},
		}},
		OutputVariables:    []AuthoringOutput{{Name: "valid", JSONPointer: "/valid", Type: "boolean", VisibleTo: "subsequent_capability_calls"}, {Name: "connection_ready", JSONPointer: "/connection_ready", Type: "boolean", VisibleTo: "subsequent_capability_calls"}},
		ReferenceContracts: []AuthoringReference{connectorReference(), operationKeyReference(), connectionReference()},
		Execution:          readExecution([]string{"integration.connector_catalog_projection", "integration.connection", "integration.secret"}),
		Errors:             []AuthoringError{authoringError("backend.integration.binding.operation_not_found", "operation_key"), authoringError("backend.integration.binding.connection_connector_mismatch", "connection_key"), authoringError("backend.integration.binding.protocol_field_unknown", "input"), authoringError("backend.integration.binding.protocol_type_mismatch", "input")},
		Examples:           integrationBindingExamples(connector, operation),
		Sources:            authoringSources("SpecializeIntegrationAuthoringCapability", "github.com/domainry/domainry-integration-sdk/management.go", "ConnectionInput"),
	}
}

func integrationConnectionCapability(key string) AuthoringCapability {
	rotate := key == "integration.connection.rotate"
	lifecycle, audit := "immediate_audited_configuration", "integration_connection_upserted"
	if rotate {
		lifecycle, audit = "explicit_connection_rotation", "integration_connection_rotated"
	}
	capability := AuthoringCapability{
		Key: key, Status: "supported", Lifecycle: lifecycle, Requires: []string{"integration.catalog"}, Permissions: []string{ActionIntegrationConnectionsUpsert}, AuditEvents: []string{audit},
		Parameters:         []AuthoringParameter{{Key: "key", Type: "connection_key"}, {Key: "connector_key", Type: "connector_key", Required: true}, {Key: "provider_key", Type: "provider_key", Required: true}, {Key: "name", Type: "string"}, {Key: "status", Type: "string", Default: "configured", Enum: cloneStrings(integrationConnectionStatuses)}, {Key: "config", Type: "provider_config"}, {Key: "secret_refs", Type: "provider_secret_refs"}},
		ValidationEndpoint: "POST /tenant-admin/integrations/connections/{connectionKey}/validate", ConfigurationRoutes: []string{"PUT /tenant-admin/integrations/connections/{connectionKey}"}, ResourceKeyPathParameter: "connectionKey",
		InputSchema: integrationConnectionInputSchema(nil), OutputSchema: integrationConnectionOutputSchema(),
		OutputVariables:    []AuthoringOutput{{Name: "connection_key", JSONPointer: "/key", Type: "connection_key", VisibleTo: "subsequent_capability_calls"}},
		ReferenceContracts: []AuthoringReference{connectorReference(), providerReference()},
		Execution:          &AuthoringExecution{ReadSet: []string{"integration.connector_catalog_projection", "integration.secret"}, WriteSet: []string{"integration.connection"}, BoundaryClass: "integration_owner", Transaction: "integration_connection_transaction", Idempotency: "connection_key", SideEffects: []string{audit}, SideEffectLevel: "internal", PermissionModel: exactActionSameKeyPermissionModel, ChangeControl: "direct_on_configuring_runtime_change_plan_on_existing_runtime"},
		Errors:             integrationConnectionErrors(), Examples: integrationConnectionExamples("example_connector", "default", nil),
		Sources: authoringSources("SpecializeIntegrationAuthoringCapability", "github.com/domainry/domainry-integration-sdk/management.go", "Management.UpsertConnection"),
	}
	if rotate {
		capability.Requires = []string{"integration.connection"}
	} else {
		capability.ConfigurationRoutes = append(capability.ConfigurationRoutes, "GET /tenant-admin/integrations/connections/{connectionKey}", "DELETE /tenant-admin/integrations/connections/{connectionKey}")
		capability.ResourceOperations = &AuthoringResourceOperations{PersistenceMode: "audited_resource", Validate: capability.ValidationEndpoint, Upsert: "PUT /tenant-admin/integrations/connections/{connectionKey}", UpsertHeaders: directAuthoringHeaders(), SuccessSchema: directAuthoringSuccessSchema(), Get: "GET /tenant-admin/integrations/connections/{connectionKey}", Delete: "DELETE /tenant-admin/integrations/connections/{connectionKey}"}
	}
	return capability
}

func integrationConnectionCommandCapability(key, route, audit, symbol string) AuthoringCapability {
	parameters := []AuthoringParameter{{Key: "connection_key", Type: "connection_key", Required: true}}
	output := integrationConnectionOutputSchema()
	if key == "integration.connection.delete" {
		closed := false
		output = &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed}
	}
	actionKey := ActionIntegrationConnectionsDisable
	if key == "integration.connection.delete" {
		actionKey = ActionIntegrationConnectionsDelete
	}
	return AuthoringCapability{
		Key: key, Status: "supported", Lifecycle: "audited_connection_command", Requires: []string{"integration.connection"}, Permissions: []string{actionKey}, AuditEvents: []string{audit}, Parameters: parameters,
		ConfigurationRoutes: []string{route}, ResourceKeyPathParameter: "connectionKey", InputSchema: parameterObjectSchema(parameters), OutputSchema: output,
		OutputVariables: []AuthoringOutput{{Name: "connection_key", JSONPointer: "/key", Type: "connection_key", VisibleTo: "subsequent_capability_calls"}}, ReferenceContracts: []AuthoringReference{connectionReference()},
		Execution: &AuthoringExecution{ReadSet: []string{"integration.connection"}, WriteSet: []string{"integration.connection"}, BoundaryClass: "integration_owner", Transaction: "integration_connection_transaction", Idempotency: "connection_state_transition", SideEffects: []string{audit}, SideEffectLevel: "internal", PermissionModel: exactActionSameKeyPermissionModel},
		Errors:    integrationConnectionErrors(), Examples: commandExamples("connection_key", "erp_primary", "backend.integration.connection.missing_key"), Sources: authoringSources("IntegrationAuthoringDomain", "github.com/domainry/domainry-integration-sdk/management.go", symbol),
	}
}

func integrationOperationTestCapability(operation *authoringOperation) AuthoringCapability {
	open := true
	input, output, operationSchema := AuthoringSchema{Type: "object", AdditionalProperties: &open}, AuthoringSchema{Type: "object", AdditionalProperties: &open}, AuthoringSchema{Type: "string"}
	operationKey := "test_connection"
	if operation != nil {
		operationKey = operation.Key
		operationSchema.Const, operationSchema.Enum = operation.Key, []any{operation.Key}
		input, output = integrationProtocolObjectSchema(operation.Input), integrationProtocolObjectSchema(operation.Output)
	}
	closed := false
	return AuthoringCapability{
		Key: "integration.operation_test", Status: "supported", Lifecycle: "explicit_confirmed_test", Requires: []string{"integration.connection", "integration.catalog"}, Permissions: []string{ActionIntegrationConnectionsTestOperation}, AuditEvents: []string{"integration_operation_tested"},
		Parameters:         []AuthoringParameter{{Key: "operation", Type: "operation_key", Required: true}, {Key: "input", Type: "operation_input"}, {Key: "confirm", Type: "boolean", Required: true}},
		ValidationEndpoint: "POST /tenant-admin/integrations/connections/{connectionKey}/test-operation", ConfigurationRoutes: []string{"POST /tenant-admin/integrations/connections/{connectionKey}/test-operation"}, ResourceKeyPathParameter: "connectionKey",
		InputSchema:     &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"operation", "confirm"}, Properties: map[string]AuthoringSchema{"operation": operationSchema, "input": input, "confirm": {Type: "boolean", Const: true}}},
		OutputSchema:    &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"connection", "operation", "response", "receipt"}, Properties: map[string]AuthoringSchema{"connection": {Type: "object", AdditionalProperties: &open}, "operation": {Type: "string"}, "response": output, "receipt": {Type: "object", AdditionalProperties: &open}}},
		OutputVariables: []AuthoringOutput{{Name: "response", JSONPointer: "/response", Type: "operation_response", VisibleTo: "subsequent_capability_calls"}, {Name: "receipt", JSONPointer: "/receipt", Type: "integration_delivery_receipt", VisibleTo: "subsequent_capability_calls"}}, ReferenceContracts: []AuthoringReference{operationReference()},
		Execution: &AuthoringExecution{ReadSet: []string{"integration.connection", "integration.connector_catalog_projection", "integration.secret"}, WriteSet: []string{"integration.invocation"}, BoundaryClass: "integration_owner", Transaction: "integration_test_invocation", Idempotency: "explicit_test_invocation", SideEffects: []string{"provider_test_call", "integration_operation_tested"}, SideEffectLevel: "external_confirmed", Compensation: "explicit_test_invocations_are_audited_and_never_automatically_retried", PermissionModel: exactActionSameKeyPermissionModel, ChangeControl: "explicit_confirm_required"},
		Errors:    []AuthoringError{authoringError("backend.integration.operation_test_confirmation_required", "confirm"), authoringError("backend.automation.connector_operation_not_found", "operation"), authoringError("backend.integration.binding.protocol_field_unknown", "input"), authoringError("backend.integration.binding.protocol_type_mismatch", "input")},
		Examples:  integrationOperationTestExamples(operationKey, operation), Sources: authoringSources("SpecializeIntegrationAuthoringCapability", "github.com/domainry/domainry-integration-sdk/management.go", "Management.TestConnection"),
	}
}

func integrationInvocationListCapability() AuthoringCapability {
	parameters := []AuthoringParameter{{Key: "connector_key", Type: "connector_key"}, {Key: "connection_key", Type: "connection_key"}, {Key: "operation", Type: "operation_key"}, {Key: "status", Type: "string", Enum: cloneStrings(integrationInvocationStatuses)}, {Key: "created_from", Type: "datetime"}, {Key: "limit", Type: "integer", Default: 100, Minimum: float64Pointer(1), Maximum: float64Pointer(200)}}
	return integrationReadListCapability("integration.invocation.list", "GET /tenant-admin/integrations/invocations", "invocations", integrationInvocationSchema(), parameters, []AuthoringReference{connectorReference(), connectionReference(), operationReference()}, "Operations.ListInvocations")
}

func integrationEventListCapability() AuthoringCapability {
	parameters := []AuthoringParameter{{Key: "provider", Type: "string"}, {Key: "event_type", Type: "string"}, {Key: "status", Type: "string", Enum: cloneStrings(integrationEventStatuses)}, {Key: "limit", Type: "integer", Default: 100, Minimum: float64Pointer(1), Maximum: float64Pointer(100)}}
	return integrationReadListCapability("integration.event.list", "GET /tenant-admin/integrations/events", "events", integrationEventSchema(), parameters, nil, "Operations.ListEvents")
}

func integrationReadListCapability(key, route, collection string, item AuthoringSchema, parameters []AuthoringParameter, references []AuthoringReference, symbol string) AuthoringCapability {
	closed := false
	representative := map[string]any{"limit": 25}
	if key == "integration.invocation.list" {
		representative["connector_key"] = "webhook"
		representative["created_from"] = "2026-01-01T00:00:00Z"
	} else {
		representative["provider"] = "slack"
	}
	actionKey := ActionIntegrationEventsList
	if key == "integration.invocation.list" {
		actionKey = ActionIntegrationInvocationsList
	}
	return AuthoringCapability{
		Key: key, Status: "supported", Lifecycle: "owner_activity_query", Permissions: []string{actionKey}, Parameters: parameters, ConfigurationRoutes: []string{route},
		InputSchema: parameterObjectSchema(parameters), OutputSchema: &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{collection, "count"}, Properties: map[string]AuthoringSchema{collection: {Type: "array", Items: &item}, "count": {Type: "integer", Minimum: float64Pointer(0)}}},
		OutputVariables: []AuthoringOutput{{Name: collection, JSONPointer: "/" + collection, Type: "integration_activity_list", VisibleTo: "subsequent_capability_calls"}, {Name: "count", JSONPointer: "/count", Type: "integer", VisibleTo: "subsequent_capability_calls"}}, ReferenceContracts: references,
		Execution: readExecution([]string{"integration." + collection}), Errors: []AuthoringError{authoringErrorWithParameters("backend.integration.query_limit_invalid", "limit", []string{"actual", "maximum", "minimum"})},
		Examples: []AuthoringExample{{Name: "minimal_valid", Value: map[string]any{}}, {Name: "representative", Value: representative}, {Name: "invalid_with_repair", Value: map[string]any{"limit": 201}, ExpectedErrorCodes: []string{"backend.integration.query_limit_invalid"}}},
		Sources:  authoringSources("IntegrationAuthoringDomain", "github.com/domainry/domainry-integration-sdk/operations.go", symbol),
	}
}

func integrationEventReplayCapability() AuthoringCapability {
	parameters := []AuthoringParameter{{Key: "event_id", Type: "integration_event_id", Required: true}}
	return AuthoringCapability{
		Key: "integration.event.replay", Status: "supported", Lifecycle: "audited_event_replay", Requires: []string{"integration.event.list"}, Permissions: []string{ActionIntegrationEventsReplay}, AuditEvents: []string{"integration_event_replayed"}, Parameters: parameters,
		ConfigurationRoutes: []string{"POST /tenant-admin/integrations/events/{eventID}/replay"}, ResourceKeyPathParameter: "eventID", InputSchema: parameterObjectSchema(parameters), OutputSchema: integrationEventOutputSchema(),
		OutputVariables: []AuthoringOutput{{Name: "event_id", JSONPointer: "/id", Type: "integration_event_id", VisibleTo: "subsequent_capability_calls"}, {Name: "status", JSONPointer: "/status", Type: "string", VisibleTo: "subsequent_capability_calls"}},
		Execution:       &AuthoringExecution{ReadSet: []string{"integration.event", "integration.event_mapping"}, WriteSet: []string{"integration.event"}, BoundaryClass: "integration_owner", Transaction: "integration_event_transaction", Idempotency: "event_id_and_current_state", SideEffects: []string{"integration_event_replayed"}, SideEffectLevel: "internal", PermissionModel: exactActionSameKeyPermissionModel},
		Errors:          []AuthoringError{authoringError("backend.integration.event_not_found", "event_id"), authoringError("backend.integration.event_replay_failed", "event_id")},
		Examples:        commandExamples("event_id", "evt_1001", "backend.integration.event_not_found"), Sources: authoringSources("IntegrationAuthoringDomain", "github.com/domainry/domainry-integration-sdk/operations.go", "Operations.ReplayEvent"),
	}
}

const jsonSchemaDraft = "https://json-schema.org/draft/2020-12/schema"

func integrationConnectionInputSchema(provider *authoringConnectorProvider) *AuthoringSchema {
	closed := false
	config, secrets := AuthoringSchema{Type: "object", AdditionalProperties: &closed, Properties: map[string]AuthoringSchema{}}, AuthoringSchema{Type: "object", AdditionalProperties: &closed, Properties: map[string]AuthoringSchema{}}
	if provider != nil {
		config, secrets = integrationFieldObjectSchema(provider.ConfigFields, false), integrationFieldObjectSchema(provider.SecretFields, true)
	}
	return &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Required: []string{"connector_key", "provider_key"}, Properties: map[string]AuthoringSchema{
		"key": {Type: "string"}, "connector_key": {Type: "string"}, "provider_key": {Type: "string"}, "name": {Type: "string"}, "status": {Type: "string", Enum: stringEnums(integrationConnectionStatuses), Default: "configured"}, "config": config, "secret_refs": secrets,
	}}
}

func integrationFieldObjectSchema(fields []authoringField, secretRefs bool) AuthoringSchema {
	closed := false
	result := AuthoringSchema{Type: "object", AdditionalProperties: &closed, Properties: map[string]AuthoringSchema{}, Required: []string{}}
	for _, field := range fields {
		property := integrationFieldValueSchema(field)
		if secretRefs {
			property = AuthoringSchema{Type: "string", Description: "Reference to an Integration-owned secret."}
		}
		result.Properties[field.Key] = property
		if field.Required {
			result.Required = append(result.Required, field.Key)
		}
	}
	return result
}

func integrationProtocolObjectSchema(fields []authoringField) AuthoringSchema {
	return integrationFieldObjectSchema(fields, false)
}

func integrationFieldValueSchema(field authoringField) AuthoringSchema {
	property := AuthoringSchema{Description: field.Description, Default: field.Default, MinLength: positiveIntPointer(field.Validation.MinLength), MaxLength: positiveIntPointer(field.Validation.MaxLength), Minimum: field.Validation.Min, Maximum: field.Validation.Max, Enum: stringEnums(field.Validation.Options)}
	switch field.Type {
	case "integer":
		property.Type = "integer"
	case "decimal", "number":
		property.Type = "number"
	case "boolean", "bool":
		property.Type = "boolean"
	case "json", "object", "map":
		property.Type = "object"
		open := true
		property.AdditionalProperties = &open
	case "array", "list":
		property.Type = "array"
		property.Items = &AuthoringSchema{}
	default:
		property.Type = "string"
	}
	return property
}

func integrationConnectionOutputSchema() *AuthoringSchema {
	schema := integrationConnectionOutputItemSchema()
	schema.Schema = jsonSchemaDraft
	return &schema
}

func integrationConnectionOutputItemSchema() AuthoringSchema {
	closed, open := false, true
	return AuthoringSchema{Type: "object", AdditionalProperties: &closed, Required: []string{"key", "connector_key", "provider_key", "status"}, Properties: map[string]AuthoringSchema{
		"key": {Type: "string"}, "workspace_id": {Type: "string"}, "connector_key": {Type: "string"}, "provider_key": {Type: "string"}, "name": {Type: "string"}, "status": {Type: "string", Enum: stringEnums(integrationConnectionStatuses)}, "config": {Type: "object", AdditionalProperties: &open}, "secret_refs": {Type: "object", AdditionalProperties: &open}, "created_by": {Type: "string"}, "created_at": {Type: "string", Format: "date-time"}, "updated_at": {Type: "string", Format: "date-time"},
	}}
}

func integrationInvocationSchema() AuthoringSchema {
	closed, open := false, true
	return AuthoringSchema{Type: "object", AdditionalProperties: &closed, Required: []string{"id", "connector_key", "operation", "status"}, Properties: map[string]AuthoringSchema{
		"id": {Type: "string"}, "workspace_id": {Type: "string"}, "connector_key": {Type: "string"}, "provider_key": {Type: "string"}, "connection_key": {Type: "string"}, "operation": {Type: "string"}, "status": {Type: "string", Enum: stringEnums(integrationInvocationStatuses)}, "duration_ms": {Type: "integer"}, "request_ref": {Type: "string"}, "response_ref": {Type: "string"}, "error": {Type: "string"}, "metadata": {Type: "object", AdditionalProperties: &open}, "created_at": {Type: "string", Format: "date-time"}, "updated_at": {Type: "string", Format: "date-time"},
	}}
}

func integrationEventSchema() AuthoringSchema {
	closed, open := false, true
	return AuthoringSchema{Type: "object", AdditionalProperties: &closed, Required: []string{"id", "provider", "event_type", "external_id", "status", "attempt_count"}, Properties: map[string]AuthoringSchema{
		"id": {Type: "string"}, "workspace_id": {Type: "string"}, "connector_key": {Type: "string"}, "connection_key": {Type: "string"}, "provider": {Type: "string"}, "event_type": {Type: "string"}, "external_id": {Type: "string"}, "status": {Type: "string", Enum: stringEnums(integrationEventStatuses)}, "payload": {Type: "object", AdditionalProperties: &open}, "error": {Type: "string"}, "attempt_count": {Type: "integer"}, "next_retry_at": {Type: "string", Format: "date-time"}, "last_attempt_at": {Type: "string", Format: "date-time"}, "received_at": {Type: "string", Format: "date-time"}, "updated_at": {Type: "string", Format: "date-time"},
	}}
}

func integrationEventOutputSchema() *AuthoringSchema {
	schema := integrationEventSchema()
	schema.Schema = jsonSchemaDraft
	return &schema
}

func parameterObjectSchema(parameters []AuthoringParameter) *AuthoringSchema {
	closed := false
	result := &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &closed, Properties: map[string]AuthoringSchema{}}
	for _, parameter := range parameters {
		property := AuthoringSchema{Type: "string", Default: parameter.Default, Minimum: parameter.Minimum, Maximum: parameter.Maximum, Enum: stringEnums(parameter.Enum)}
		switch parameter.Type {
		case "integer":
			property.Type = "integer"
		case "boolean":
			property.Type = "boolean"
		case "object", "operation_input", "provider_config", "provider_secret_refs":
			property.Type = "object"
			open := true
			property.AdditionalProperties = &open
		}
		result.Properties[parameter.Key] = property
		if parameter.Required {
			result.Required = append(result.Required, parameter.Key)
		}
	}
	return result
}

func integrationConnectionExamples(connectorKey, providerKey string, provider *authoringConnectorProvider) []AuthoringExample {
	config, secrets := map[string]any{}, map[string]any{}
	if provider != nil {
		for _, field := range provider.ConfigFields {
			if field.Required {
				config[field.Key] = integrationExampleFieldValue(field)
			}
		}
		for _, field := range provider.SecretFields {
			if field.Required {
				secrets[field.Key] = "$secret." + field.Key
			}
		}
	}
	return []AuthoringExample{
		{Name: "minimal_valid", Value: map[string]any{"connector_key": connectorKey, "provider_key": providerKey, "status": "draft"}},
		{Name: "representative", Value: map[string]any{"key": connectorKey + "_primary", "connector_key": connectorKey, "provider_key": providerKey, "name": "Primary " + connectorKey, "status": "configured", "config": config, "secret_refs": secrets}},
		{Name: "invalid_with_repair", Value: map[string]any{"connector_key": connectorKey, "provider_key": providerKey, "status": "unknown"}, ExpectedErrorCodes: []string{"backend.integration.connection.invalid_status"}},
	}
}

func integrationBindingExamples(connector *authoringConnector, operation *authoringOperation) []AuthoringExample {
	connectorKey, operationKey := "example_connector", ""
	input, output := map[string]any{}, map[string]any{}
	if connector != nil {
		connectorKey = connector.Key
	}
	if operation != nil {
		operationKey = operation.Key
		for _, field := range operation.Input {
			if field.Required {
				input[field.Key] = integrationExampleFieldValue(field)
			}
		}
		for _, field := range operation.Output {
			if field.Required {
				output[field.Key] = integrationExampleFieldValue(field)
			}
		}
	}
	representative := map[string]any{"connector_key": connectorKey, "connection_key": connectorKey + "_primary"}
	if operationKey != "" {
		representative["operation_key"], representative["input"], representative["output"] = operationKey, input, output
	}
	return []AuthoringExample{
		{Name: "minimal_valid", Value: map[string]any{"connector_key": connectorKey, "connection_key": connectorKey + "_primary"}},
		{Name: "representative", Value: representative},
		{Name: "invalid_with_repair", Value: map[string]any{"connector_key": connectorKey, "operation_key": operationKey, "input": map[string]any{"unknown": true}}, ExpectedErrorCodes: []string{"backend.integration.binding.protocol_field_unknown"}},
	}
}

func integrationOperationTestExamples(operationKey string, operation *authoringOperation) []AuthoringExample {
	minimalInput, representativeInput := map[string]any{}, map[string]any{}
	if operation != nil {
		for _, field := range operation.Input {
			value := integrationExampleFieldValue(field)
			representativeInput[field.Key] = value
			if field.Required {
				minimalInput[field.Key] = value
			}
		}
	}
	minimal := map[string]any{"operation": operationKey, "confirm": true}
	if len(minimalInput) > 0 {
		minimal["input"] = minimalInput
	}
	return []AuthoringExample{
		{Name: "minimal_valid", Value: minimal},
		{Name: "representative", Value: map[string]any{"operation": operationKey, "input": representativeInput, "confirm": true}},
		{Name: "invalid_with_repair", Value: map[string]any{"operation": operationKey, "input": minimalInput, "confirm": false}, ExpectedErrorCodes: []string{"backend.integration.operation_test_confirmation_required"}},
	}
}

func integrationExampleFieldValue(field authoringField) any {
	if field.Default != nil {
		return field.Default
	}
	switch field.Type {
	case "integer":
		return 1
	case "decimal", "number":
		return 1.0
	case "boolean", "bool":
		return true
	case "json", "object", "map":
		return map[string]any{}
	case "array", "list":
		return []any{}
	default:
		return "example_" + field.Key
	}
}

func integrationConnectionErrors() []AuthoringError {
	return []AuthoringError{
		authoringError("backend.integration.connection.missing_key", "key"), authoringError("backend.integration.connection.key_invalid", "key"), authoringError("backend.integration.connection.missing_connector", "connector_key"), authoringError("backend.integration.connector.not_found", "connector_key"), authoringError("backend.integration.connection.connector_immutable", "connector_key"), authoringError("backend.integration.connection.provider_immutable", "provider_key"), authoringError("backend.integration.connection.provider_secret_required", "secret_refs"), authoringError("backend.integration.connection.provider_secret_unknown", "secret_refs"), authoringError("backend.integration.connection.invalid_status", "status"),
	}
}

func readExecution(readSet []string) *AuthoringExecution {
	return &AuthoringExecution{ReadSet: readSet, BoundaryClass: "integration_owner", Transaction: "read_only", Idempotency: "naturally_idempotent", SideEffectLevel: "none", PermissionModel: exactActionSameKeyPermissionModel}
}

func readExamples() []AuthoringExample {
	return []AuthoringExample{{Name: "minimal_valid", Value: map[string]any{}}, {Name: "representative", Value: map[string]any{}}}
}

func commandExamples(key, value, errorCode string) []AuthoringExample {
	return []AuthoringExample{{Name: "minimal_valid", Value: map[string]any{key: value}}, {Name: "representative", Value: map[string]any{key: value}}, {Name: "invalid_with_repair", Value: map[string]any{key: ""}, ExpectedErrorCodes: []string{errorCode}}}
}

func connectorReference() AuthoringReference {
	return AuthoringReference{Kind: "connector_key", InputJSONPointer: "/connector_key", ResolverEndpoint: "/tenant-admin/platform-capabilities/references/connector_key"}
}
func providerReference() AuthoringReference {
	return AuthoringReference{Kind: "provider_key", InputJSONPointer: "/provider_key", ScopeFrom: "/connector_key", ResolverEndpoint: "/tenant-admin/platform-capabilities/references/provider_key"}
}
func operationReference() AuthoringReference {
	return AuthoringReference{Kind: "operation_key", InputJSONPointer: "/operation", ScopeFrom: "/connector_key", ResolverEndpoint: "/tenant-admin/platform-capabilities/references/operation_key"}
}
func operationKeyReference() AuthoringReference {
	return AuthoringReference{Kind: "operation_key", InputJSONPointer: "/operation_key", ScopeFrom: "/connector_key", ResolverEndpoint: "/tenant-admin/platform-capabilities/references/operation_key"}
}
func connectionReference() AuthoringReference {
	return AuthoringReference{Kind: "connection_key", InputJSONPointer: "/connection_key", ResolverEndpoint: "/tenant-admin/platform-capabilities/references/connection_key"}
}

func authoringSources(authoringSymbol, contractPath, contractSymbol string) []AuthoringSource {
	return []AuthoringSource{{Kind: "owner_sdk", Path: integrationAuthoringSourcePath, Symbol: authoringSymbol}, {Kind: "owner_sdk", Path: contractPath, Symbol: contractSymbol}}
}

func authoringError(code, path string) AuthoringError {
	return AuthoringError{Code: code, FieldPath: path, MessageKey: code}
}
func authoringErrorWithParameters(code, path string, parameters []string) AuthoringError {
	return AuthoringError{Code: code, FieldPath: path, ParameterKeys: parameters, MessageKey: code}
}

func directAuthoringHeaders() []AuthoringRequestHeader {
	return []AuthoringRequestHeader{{Name: "Builder-Task-ID", Required: true, ValueSource: "builder_task_id", Description: "Stable identity of the project-owned builder task."}, {Name: "Idempotency-Key", Required: true, ValueSource: "request_fingerprint", Description: "Stable key for this capability, resource, and canonical payload."}, {Name: "Expected-Schema-Hash", Required: true, ValueSource: "expected_resource_hash", Description: "Last observed resource hash, or empty when the resource does not exist."}}
}

func directAuthoringSuccessSchema() *AuthoringSchema {
	open, closed := true, false
	return &AuthoringSchema{Schema: jsonSchemaDraft, Type: "object", AdditionalProperties: &open, Required: []string{"resource", "resource_hash", "snapshot_hash", "available_successors"}, Properties: map[string]AuthoringSchema{
		"resource": {OneOf: []AuthoringSchema{{Type: "object", AdditionalProperties: &open}, {Type: "array"}}}, "resource_hash": {Type: "string", MinLength: positiveIntPointer(1)}, "snapshot_hash": {Type: "string", MinLength: positiveIntPointer(1)},
		"available_successors": {Type: "array", Items: &AuthoringSchema{Type: "object", AdditionalProperties: &closed, Required: []string{"key", "domain", "status", "detail_endpoint"}, Properties: map[string]AuthoringSchema{"key": {Type: "string"}, "domain": {Type: "string"}, "status": {Type: "string"}, "detail_endpoint": {Type: "string"}, "validation_endpoint": {Type: "string"}}}},
	}}
}

func stringEnums(values []string) []any {
	result := make([]any, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}
func cloneStrings(values []string) []string { return append([]string(nil), values...) }
func float64Pointer(value float64) *float64 { return &value }
func positiveIntPointer(value int) *int {
	if value <= 0 {
		return nil
	}
	return &value
}
