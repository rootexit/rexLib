package rexDatabase

type CommonStatus int8

const (
	AllStatusDisable CommonStatus = iota + 1
	AllStatusEnable
)

const (
	SessionStatusAnonymous CommonStatus = iota + 1
	SessionStatusAuthenticated
	SessionStatusRevoked
)

const (
	VerifiedStatusDraft CommonStatus = iota + 1
	VerifiedStatusSubmitted
	VerifiedStatusRejected
	VerifiedStatusVerified
)

const (
	ContentStatusDraft CommonStatus = iota + 1
	ContentStatusSubmitted
	ContentStatusRejected
	ContentStatusUnPublished
	ContentStatusPublished
	ContentStatusPinnedOnly
	ContentStatusPinned
)
