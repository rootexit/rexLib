package rexDatabase

// note: 1 -> 一次性短链, 2 -> 有效期短链, 3 -> 永久短链
const (
	ShortLinkTypeOneTime   int32 = iota + 1 // 一次性短链
	ShortLinkTypeValidTime                  // 有效期短链
	ShortLinkTypePermanent                  // 永久短链
)
