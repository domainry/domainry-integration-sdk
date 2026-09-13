package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/domainry/domainry-foundation/requestcontext"
	sdk "github.com/domainry/domainry-integration-sdk"
	"net/http"
)

func (b *binding) SubjectLifecycle() sdk.SubjectLifecycle { return b.client }
func (c *remoteClient) subjectCall(ctx context.Context, operation string, r sdk.SubjectErasureRequest, plan json.RawMessage) (json.RawMessage, error) {
	if err := r.Validate(operation == "prepare" || operation == "erase"); err != nil {
		return nil, err
	}
	if requestcontext.WorkspaceID(ctx) != r.WorkspaceID {
		return nil, fmt.Errorf("Integration remote subject scope mismatch")
	}
	var result json.RawMessage
	err := c.call(ctx, http.MethodPost, "/integration/v1/subjects/"+operation, struct {
		Request sdk.SubjectErasureRequest `json:"request"`
		Plan    json.RawMessage           `json:"plan,omitempty"`
	}{r, plan}, &result)
	return result, err
}
func (c *remoteClient) PreviewSubject(ctx context.Context, r sdk.SubjectErasureRequest) (json.RawMessage, error) {
	return c.subjectCall(ctx, "preview", r, nil)
}
func (c *remoteClient) ExportSubject(ctx context.Context, r sdk.SubjectErasureRequest) (json.RawMessage, error) {
	return c.subjectCall(ctx, "export", r, nil)
}
func (c *remoteClient) PrepareSubjectErasure(ctx context.Context, r sdk.SubjectErasureRequest) (json.RawMessage, error) {
	return c.subjectCall(ctx, "prepare", r, nil)
}
func (c *remoteClient) ErasePreparedSubject(ctx context.Context, r sdk.SubjectErasureRequest, p json.RawMessage) (json.RawMessage, error) {
	return c.subjectCall(ctx, "erase", r, p)
}

var _ sdk.SubjectLifecycleBinding = (*binding)(nil)
var _ sdk.SubjectLifecycle = (*remoteClient)(nil)
