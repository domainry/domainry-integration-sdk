package modulehost

import (
	"context"

	integrationsdk "github.com/domainry/domainry-integration-sdk"
)

type Factory interface {
	integrationsdk.Factory
	OpenModule(context.Context, integrationsdk.ApplicationRef, Host) (integrationsdk.Binding, error)
}
