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

type CommonContentStatus int32

// note: 文章状态, 1->预草稿,2->草稿,3->提审,4->驳回,5->下架,6->上架（仅列表）,7->仅置顶,8->推荐(置顶+列表)

const (
	ContentStatusPreDraft CommonContentStatus = iota + 1
	ContentStatusDraft
	ContentStatusSubmitted
	ContentStatusRejected
	ContentStatusUnPublished
	ContentStatusPublished
	ContentStatusTopOnly
	ContentStatusRecommend
)
