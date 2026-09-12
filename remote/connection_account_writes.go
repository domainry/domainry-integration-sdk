package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	sdk "github.com/domainry/domainry-integration-sdk"
)

func (c *remoteClient) AuthorizeConnectionAccountWrite(ctx context.Context, subject sdk.ConnectionAccountSubject, key string, input sdk.ConnectionAccountWriteOperation) (sdk.ConnectionAccountWriteAccess, error) {
	if err := input.Validate(); err != nil {
		return sdk.ConnectionAccountWriteAccess{}, err
	}
	var out sdk.ConnectionAccountWriteAccess
	err := c.accountWriteCall(ctx, subject, key, "write-access", input, &out)
	return out, err
}

func (c *remoteClient) WriteConnectionAccount(ctx context.Context, subject sdk.ConnectionAccountSubject, key string, input sdk.ConnectionAccountWriteRequest) (sdk.ConnectionAccountWriteResult, error) {
	return c.accountWriteResult(ctx, subject, key, "write", input)
}

func (c *remoteClient) ReadConnectionAccountWriteReceipt(ctx context.Context, subject sdk.ConnectionAccountSubject, key string, input sdk.ConnectionAccountWriteRequest) (sdk.ConnectionAccountWriteResult, error) {
	return c.accountWriteResult(ctx, subject, key, "write-receipt", input)
}

func (c *remoteClient) accountWriteResult(ctx context.Context, subject sdk.ConnectionAccountSubject, key, action string, input sdk.ConnectionAccountWriteRequest) (sdk.ConnectionAccountWriteResult, error) {
	if err := input.Validate(); err != nil {
		return sdk.ConnectionAccountWriteResult{}, err
	}
	var out sdk.ConnectionAccountWriteResult
	err := c.accountWriteCall(ctx, subject, key, action, input, &out)
	if err != nil {
		return sdk.ConnectionAccountWriteResult{}, err
	}
	if out.Source != input.ExpectedSource || out.Status != sdk.AccountWriteNotFound && out.Status != sdk.AccountWriteSucceeded && out.Status != sdk.AccountWriteFailed && out.Status != sdk.AccountWriteUncertain || out.Status == sdk.AccountWriteSucceeded && len(out.Receipt) == 0 || out.Status != sdk.AccountWriteSucceeded && len(out.Receipt) != 0 {
		return sdk.ConnectionAccountWriteResult{}, fmt.Errorf("Integration account write response is invalid; query the original request")
	}
	return out, nil
}

// One bounded POST, no redirects, cookies or automatic application retries.
// Any transport/response error may follow a completed write; the only recovery
// is ReadConnectionAccountWriteReceipt with the same host-frozen request.
func (c *remoteClient) accountWriteCall(ctx context.Context, subject sdk.ConnectionAccountSubject, key, action string, input, out any) error {
	if err := subject.Validate(); err != nil {
		return err
	}
	if key == "" || strings.TrimSpace(key) != key {
		return fmt.Errorf("Integration account write key is invalid")
	}
	reference, err := url.Parse(connectionAccountPath("/"+url.PathEscape(key)+"/"+action, subject))
	if err != nil {
		return fmt.Errorf("Integration account write path is invalid")
	}
	payload, err := json.Marshal(input)
	if err != nil || len(payload) > 2<<20 {
		return fmt.Errorf("Integration account write request is invalid")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base.ResolveReference(reference).String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("Integration account write request is invalid")
	}
	request.GetBody = nil
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("X-Domainry-Runtime-ID", c.runtimeID)
	client := *c.client
	client.Jar = nil
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Integration account write response unavailable; query the original request")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Integration account write returned HTTP %d; query the original request", response.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(response.Body, (64<<10)+1))
	if err != nil || len(raw) > 64<<10 {
		return fmt.Errorf("Integration account write response unavailable; query the original request")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil || d.Decode(new(any)) != io.EOF {
		return fmt.Errorf("Integration account write response invalid; query the original request")
	}
	return nil
}

var _ sdk.ConnectionAccountWrites = (*remoteClient)(nil)
var _ sdk.ConnectionAccountWritesBinding = (*binding)(nil)
