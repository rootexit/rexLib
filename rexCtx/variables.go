package rexCtx

type (
	CtxSiteDomain struct{}
	CtxStartTime  struct{}
	CtxClientIp   struct{}
	CtxClientPort struct{}
	CtxFullMethod struct{}

	CtxRequestURI      struct{}
	CtxRequestID       struct{}
	CtxTenantId        struct{}
	CtxTenants         struct{}
	CtxDomainId        struct{}
	CtxUserId          struct{}
	CtxAdminId         struct{}
	CtxRole            struct{}
	CtxClaimsAudience  struct{}
	CtxClaimsExpiresAt struct{}
	CtxClaimsId        struct{}
	CtxClaimsIssuedAt  struct{}
	CtxClaimsIssuer    struct{}
	CtxClaimsNotBefore struct{}
	CtxClaimsSubject   struct{}
	CtxUserAgent       struct{}
	CtxXAccessKeyFor   struct{}
	CtxSessionIDFor    struct{}
	CtxXAuthMethodFor  struct{}
	CtxXAccountFor     struct{}
	CtxCityId          struct{}
	CtxCountry         struct{}
	CtxRegion          struct{}
	CtxProvince        struct{}
	CtxCity            struct{}
	CtxISP             struct{}
	CtxUserAgentFamily struct{}
	CtxUserAgentMajor  struct{}
	CtxUserAgentMinor  struct{}
	CtxUserAgentPatch  struct{}
	CtxOsFamily        struct{}
	CtxOsMajor         struct{}
	CtxOsMinor         struct{}
	CtxOsPatch         struct{}
	CtxOsPatchMinor    struct{}
	CtxDeviceFamily    struct{}
	CtxDeviceBrand     struct{}
	CtxDeviceModel     struct{}
)
