package modulehost

import (
	"context"

	connector "github.com/domainry/domainry-connector-sdk"
	ormmigration "github.com/domainry/domainry-orm/migration"
	"github.com/domainry/domainry-orm/sqlhost"
)

type Database = sqlhost.Database

type Dialect interface {
	Identifier(string) string
	Table(string) string
	Placeholder(int) string
	Insert(string, []string) string
}

type SchemaMigration = ormmigration.Migration

type MigrationRegistrar interface {
	Driver() string
	Schema() string
	ApplyOwnedMigrations(context.Context, string, []SchemaMigration) error
}

type ProviderRegistry interface {
	Provider(connectorKey, providerKey string) (connector.Adapter, bool)
	Descriptors() []connector.ProviderDescriptor
}

type SecretMaterialCipher interface {
	DecryptSecretMaterial(ctx context.Context, workspaceID, secretKey, ciphertext string) (string, error)
}

type Host interface {
	Database() Database
	Dialect() Dialect
	Migrations() MigrationRegistrar
	Providers() ProviderRegistry
	SecretCipher() SecretMaterialCipher
}
