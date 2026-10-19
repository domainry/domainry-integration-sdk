package modulehost

// EventWorkNotifier is a non-blocking process-local hint after Integration
// commits provider events. It grants no authority and executes no work: the
// host event worker must still claim durable rows with normal owner controls.
type EventWorkNotifier interface {
	EventsCommitted()
}

// EventWorkNotifierHost is an optional host acceleration capability. Durable
// recovery polling remains necessary for restart and cross-process commits.
type EventWorkNotifierHost interface {
	EventWorkNotifier() EventWorkNotifier
}
