package rexDatabase

type MinisiteType string

const (
	MinisiteTypeInternal MinisiteType = "internal"
	MinisiteTypeExternal MinisiteType = "external"
)

type MinisiteVisibilityType string

const (
	MinisiteVisibilityTypePrivate MinisiteVisibilityType = "private"
	MinisiteVisibilityTypePublic  MinisiteVisibilityType = "public"
)
