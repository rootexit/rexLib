package rexDatabase

type CommonStatus int8

const (
	AllStatusDisable CommonStatus = iota + 1
	AllStatusEnable
)

type CommonSessionStatus int8

const (
	SessionStatusAnonymous CommonSessionStatus = iota + 1
	SessionStatusAuthenticated
	SessionStatusRevoked
)

type CommonVerifiedStatus int8

const (
	VerifiedStatusUnSubmitted CommonVerifiedStatus = iota + 1
	VerifiedStatusDraft
	VerifiedStatusSubmitted
	VerifiedStatusRejected
	VerifiedStatusVerified
)

type CommonContentStatus int8

const (
	ContentStatusDraft CommonContentStatus = iota + 1
	ContentStatusSubmitted
	ContentStatusRejected
	ContentStatusUnPublished
	ContentStatusPublished
	ContentStatusPinnedOnly
	ContentStatusPinned
)
