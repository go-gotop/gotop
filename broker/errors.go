package broker

import "errors"

var (
	ErrPublisherNotConfigured  = errors.New("broker: publisher not configured")
	ErrSubscriberNotConfigured = errors.New("broker: subscriber not configured")
)
