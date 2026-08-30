package saashost

import (
	"context"

	integrationsdk "github.com/domainry/domainry-integration-sdk"
)

type Host interface{}

type Factory interface {
	integrationsdk.Factory
	OpenSaaS(context.Context, integrationsdk.ApplicationRef, Host) (integrationsdk.Binding, error)
}
