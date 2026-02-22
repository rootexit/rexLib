package rexDatabase

type CommonBootstrapStatus int8

const (
	BootstrapStatusPending = iota + 1
	BootstrapStatusInProgress
	BootstrapStatusCompleted
)
