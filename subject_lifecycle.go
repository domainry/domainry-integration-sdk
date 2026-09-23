package integrationsdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// SubjectLifecycle is a privileged owner port. Browser management adapters do
// not expose it. Provenance comes from the authenticated Runtime coordinator.
type SubjectLifecycleBinding interface{ SubjectLifecycle() SubjectLifecycle }

// SubjectLifecyclePersistenceBinding is implemented by an embedded owner whose
// database also contains Lifecycle's shared subject fence and step journal.
// Hosts bind it only after Lifecycle has installed those tables.
type SubjectLifecyclePersistenceBinding interface {
	BindSubjectLifecyclePersistence(context.Context) error
}

type SubjectRecordReference struct {
	ObjectKey string `json:"object_key"`
	RecordID  string `json:"record_id"`
}
type SubjectErasureRequest struct {
	WorkspaceID           string                   `json:"workspace_id"`
	SubjectID             string                   `json:"subject_id"`
	RequestID             string                   `json:"request_id,omitempty"`
	Resources             []SubjectRecordReference `json:"resources,omitempty"`
	PublicationMessageIDs []string                 `json:"publication_message_ids,omitempty"`
	EventIDs              []string                 `json:"event_ids,omitempty"`
	LegalHolds            json.RawMessage          `json:"legal_holds,omitempty"`
}

func (r SubjectErasureRequest) Validate(erasure bool) error {
	if !boundedAccountWriteValue(r.WorkspaceID, 191) || !boundedAccountWriteValue(r.SubjectID, 255) || erasure && !boundedAccountWriteValue(r.RequestID, 255) {
		return fmt.Errorf("Integration subject scope is invalid")
	}
	if len(r.Resources) > 10001 || len(r.PublicationMessageIDs) > 10000 || len(r.EventIDs) > 10000 {
		return fmt.Errorf("Integration subject provenance exceeds limit")
	}
	for _, ref := range r.Resources {
		if !boundedAccountWriteValue(ref.ObjectKey, 255) || !boundedAccountWriteValue(ref.RecordID, 255) {
			return fmt.Errorf("Integration subject resource is invalid")
		}
	}
	for _, ids := range [][]string{r.PublicationMessageIDs, r.EventIDs} {
		for _, id := range ids {
			if !boundedAccountWriteValue(id, 255) {
				return fmt.Errorf("Integration subject provenance is invalid")
			}
		}
	}
	if erasure && len(r.LegalHolds) > 0 {
		var holds []json.RawMessage
		if json.Unmarshal(r.LegalHolds, &holds) != nil || len(holds) > 0 {
			return fmt.Errorf("Integration erasure blocked by legal hold")
		}
	}
	return nil
}

type SubjectLifecycle interface {
	PreviewSubject(context.Context, SubjectErasureRequest) (json.RawMessage, error)
	ExportSubject(context.Context, SubjectErasureRequest) (json.RawMessage, error)
	PrepareSubjectErasure(context.Context, SubjectErasureRequest) (json.RawMessage, error)
	ErasePreparedSubject(context.Context, SubjectErasureRequest, json.RawMessage) (json.RawMessage, error)
}
