package integrationsdk

import "context"

// ProviderRunSnapshot contains installation-wide aggregate task observations.
// It contains no connection identifiers, task state, payloads, or credentials
// and is intended only for the administrator-scoped Monitoring surface.
type ProviderRunSnapshot struct {
	ObservedAt        string `json:"observed_at"`
	ReadyDue          int64  `json:"ready_due"`
	Failed            int64  `json:"failed"`
	FailedDue         int64  `json:"failed_due"`
	Processing        int64  `json:"processing"`
	ExpiredProcessing int64  `json:"expired_processing"`
	DeadLetter        int64  `json:"dead_letter"`
	// GoogleHTTP429Failures counts tasks whose latest persisted failure is 429;
	// it is not a historical event rate.
	GoogleHTTP429Failures int64 `json:"google_http_429_failures"`
	// These are durable failure transitions, not unique accounts. The 5-minute
	// count is the current aligned window; the hour count is the current and
	// preceding eleven aligned windows.
	GoogleHTTP429Events5m    int64  `json:"google_http_429_events_5m"`
	GoogleHTTP429EventsHour  int64  `json:"google_http_429_events_hour"`
	GoogleHTTP429WindowStart string `json:"google_http_429_window_start"`
	// GooglePush metrics use deduplicated routed-webhook inbox receipts from the
	// preceding hour. Timed counts exclude missing/invalid publish timestamps;
	// no-wake means no target background task was woken (with or without a
	// matching account route), not proof of a lost mail.
	GooglePushReceivedHour        int64  `json:"google_push_received_hour"`
	GooglePushNoWakeHour          int64  `json:"google_push_no_wake_hour"`
	GooglePushTimedHour           int64  `json:"google_push_timed_hour"`
	GooglePushDelayed5mHour       int64  `json:"google_push_delayed_5m_hour"`
	GooglePushMaxDelayMsHour      int64  `json:"google_push_max_delay_ms_hour"`
	GooglePushHotBucket           *int   `json:"google_push_hot_bucket,omitempty"`
	GooglePushHotBucketEventsHour int64  `json:"google_push_hot_bucket_events_hour"`
	OldestRunnableDueAt           string `json:"oldest_runnable_due_at,omitempty"`
	// GlobalRateUsed/Limit describe Integration's internal dispatch budget,
	// not Google's upstream API quota.
	GlobalRateUsed      int64  `json:"global_rate_used"`
	GlobalRateLimit     int64  `json:"global_rate_limit"`
	RateWindowStartedAt string `json:"rate_window_started_at,omitempty"`
}

// ProviderRunMonitoring is a source-owned system observation port. Callers
// must enforce unrestricted administrator permission before exposing it.
type ProviderRunMonitoring interface {
	ProviderRunSnapshot(context.Context) (ProviderRunSnapshot, error)
}
