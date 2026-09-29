package integrationsdk

import "context"

const (
	ProviderQuotaStatusAvailable     = "available"
	ProviderQuotaStatusNotConfigured = "not_configured"
	ProviderQuotaStatusUnavailable   = "unavailable"
)

// ProviderQuotaUsage is one provider-reported quota observation. Usage and
// limit come from the same provider monitoring surface and are never inferred
// from Integration's local admission counters.
type ProviderQuotaUsage struct {
	Service     string `json:"service"`
	QuotaMetric string `json:"quota_metric"`
	LimitName   string `json:"limit_name"`
	Location    string `json:"location,omitempty"`
	Window      string `json:"window"`
	Usage       int64  `json:"usage"`
	Limit       int64  `json:"limit"`
	UsedPercent int64  `json:"used_percent"`
}

// ProviderQuotaSnapshot is installation-level official provider evidence.
// Status must be available before any usage value is interpreted. ErrorCode is
// deliberately bounded and contains no upstream response body or credential
// material.
type ProviderQuotaSnapshot struct {
	Source         string               `json:"source"`
	Status         string               `json:"status"`
	ObservedAt     string               `json:"observed_at"`
	DataThroughAt  string               `json:"data_through_at,omitempty"`
	Usage          []ProviderQuotaUsage `json:"usage,omitempty"`
	MaxUsedPercent *int64               `json:"max_used_percent,omitempty"`
	Unmatched      int64                `json:"unmatched,omitempty"`
	Truncated      int64                `json:"truncated,omitempty"`
	ErrorCode      string               `json:"error_code,omitempty"`
}

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
	// Foreground counts cover direct Provider operations and deliveries. They
	// are kept separate from background task transitions so operators can tell
	// interactive/write-path pressure from scheduled synchronization pressure.
	GoogleHTTP429ForegroundEvents5m   int64 `json:"google_http_429_foreground_events_5m"`
	GoogleHTTP429ForegroundEventsHour int64 `json:"google_http_429_foreground_events_hour"`
	// Gmail reports quota pressure as HTTP 403 with either rateLimitExceeded or
	// userRateLimitExceeded as well as through generic HTTP 429 responses. These
	// fields keep those verified Gmail rejection reasons distinct from 429.
	GoogleGmailProjectRateLimitFailures             int64  `json:"google_gmail_project_rate_limit_failures"`
	GoogleGmailProjectRateLimitEvents5m             int64  `json:"google_gmail_project_rate_limit_events_5m"`
	GoogleGmailProjectRateLimitEventsHour           int64  `json:"google_gmail_project_rate_limit_events_hour"`
	GoogleGmailProjectRateLimitForegroundEvents5m   int64  `json:"google_gmail_project_rate_limit_foreground_events_5m"`
	GoogleGmailProjectRateLimitForegroundEventsHour int64  `json:"google_gmail_project_rate_limit_foreground_events_hour"`
	GoogleGmailUserRateLimitFailures                int64  `json:"google_gmail_user_rate_limit_failures"`
	GoogleGmailUserRateLimitEvents5m                int64  `json:"google_gmail_user_rate_limit_events_5m"`
	GoogleGmailUserRateLimitEventsHour              int64  `json:"google_gmail_user_rate_limit_events_hour"`
	GoogleGmailUserRateLimitForegroundEvents5m      int64  `json:"google_gmail_user_rate_limit_foreground_events_5m"`
	GoogleGmailUserRateLimitForegroundEventsHour    int64  `json:"google_gmail_user_rate_limit_foreground_events_hour"`
	GoogleGmailRateLimitWindowStart                 string `json:"google_gmail_rate_limit_window_start"`
	// Feishu exposes endpoint/app/tenant frequency limits only through official
	// rejection responses, not through a proactive usage API. These fields count
	// the canonical 429 or legacy code 99991400 signal without inventing a
	// remaining-quota percentage.
	FeishuRateLimitFailures             int64  `json:"feishu_rate_limit_failures"`
	FeishuRateLimitEvents5m             int64  `json:"feishu_rate_limit_events_5m"`
	FeishuRateLimitEventsHour           int64  `json:"feishu_rate_limit_events_hour"`
	FeishuRateLimitForegroundEvents5m   int64  `json:"feishu_rate_limit_foreground_events_5m"`
	FeishuRateLimitForegroundEventsHour int64  `json:"feishu_rate_limit_foreground_events_hour"`
	FeishuRateLimitWindowStart          string `json:"feishu_rate_limit_window_start"`
	// GooglePush metrics use deduplicated routed-webhook inbox receipts from the
	// preceding hour. Timed counts exclude missing/invalid publish timestamps;
	// no-wake means no target background task was woken (with or without a
	// matching account route), not proof of a lost mail.
	GooglePushReceivedHour        int64 `json:"google_push_received_hour"`
	GooglePushNoWakeHour          int64 `json:"google_push_no_wake_hour"`
	GooglePushTimedHour           int64 `json:"google_push_timed_hour"`
	GooglePushDelayed5mHour       int64 `json:"google_push_delayed_5m_hour"`
	GooglePushMaxDelayMsHour      int64 `json:"google_push_max_delay_ms_hour"`
	GooglePushHotBucket           *int  `json:"google_push_hot_bucket,omitempty"`
	GooglePushHotBucketEventsHour int64 `json:"google_push_hot_bucket_events_hour"`
	// Sync completion metrics correlate a verified routed Push wake with the
	// successful commit of the exact background task wake sequence. Only
	// receipts with a trusted publish timestamp participate in delay samples.
	GoogleSyncCompletedHour          int64 `json:"google_sync_completed_hour"`
	GoogleSyncTimedHour              int64 `json:"google_sync_timed_hour"`
	GoogleSyncDelayed5mHour          int64 `json:"google_sync_delayed_5m_hour"`
	GoogleSyncMaxEndToEndDelayMsHour int64 `json:"google_sync_max_end_to_end_delay_ms_hour"`
	GoogleSyncPending                int64 `json:"google_sync_pending"`
	GoogleSyncPendingOver5m          int64 `json:"google_sync_pending_over_5m"`
	GoogleSyncOldestPendingAgeMs     int64 `json:"google_sync_oldest_pending_age_ms"`
	// Gmail history evidence is emitted only by successful incremental history
	// transitions. MessagesWithoutPush counts messages recovered when no
	// verified Push receipt covered the claimed wake sequence; CursorExpired
	// counts history cursors rejected by Gmail and recovered by reconciliation.
	GoogleGmailHistoryMessagesWithoutPushHour int64  `json:"google_gmail_history_messages_without_push_hour"`
	GoogleGmailHistoryCursorExpiredHour       int64  `json:"google_gmail_history_cursor_expired_hour"`
	OldestRunnableDueAt                       string `json:"oldest_runnable_due_at,omitempty"`
	// GoogleOfficialQuota is populated only from Google Cloud Monitoring's
	// consumer_quota metrics using installation-owned Cloud credentials. It is
	// intentionally unrelated to per-user Gmail OAuth tokens.
	GoogleOfficialQuota ProviderQuotaSnapshot `json:"google_official_quota"`
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
