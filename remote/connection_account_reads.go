package remote

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	sdk "github.com/domainry/domainry-integration-sdk"
)

func (c *remoteClient) AuthorizeConnectionAccountRead(ctx context.Context, subject sdk.ConnectionAccountSubject, key string, input sdk.ConnectionAccountReadOperation) (sdk.ConnectionAccountReadAccess, error) {
	if err := subject.Validate(); err != nil {
		return sdk.ConnectionAccountReadAccess{}, err
	}
	if err := input.Validate(); err != nil {
		return sdk.ConnectionAccountReadAccess{}, err
	}
	var out sdk.ConnectionAccountReadAccess
	err := c.call(ctx, http.MethodPost, connectionAccountPath("/"+url.PathEscape(strings.TrimSpace(key))+"/read-access", subject), input, &out)
	return out, err
}

func (c *remoteClient) ReadConnectionAccount(ctx context.Context, subject sdk.ConnectionAccountSubject, key string, input sdk.ConnectionAccountReadRequest) (sdk.ConnectionAccountReadResult, error) {
	if err := subject.Validate(); err != nil {
		return sdk.ConnectionAccountReadResult{}, err
	}
	if err := input.Validate(); err != nil {
		return sdk.ConnectionAccountReadResult{}, err
	}
	var out sdk.ConnectionAccountReadResult
	err := c.call(ctx, http.MethodPost, connectionAccountPath("/"+url.PathEscape(strings.TrimSpace(key))+"/read", subject), input, &out)
	return out, err
}

var _ sdk.ConnectionAccountReads = (*remoteClient)(nil)
var _ sdk.ConnectionAccountReadsBinding = (*binding)(nil)
