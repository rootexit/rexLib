package rexHeaders

const (
	HeaderAuthorization  = "Authorization"
	HeaderHost           = "Host"
	HeaderXAuthMethodFor = "X-AuthMethod-For"
	HeaderXAccountFor    = "X-Account-For"
	HeaderXAccessKeyFor  = "X-AccessKey-For"
	HeaderXSessionIdFor  = "X-SessionId-For"

	HeaderXAdminIDFor  = "X-AdminID-For"
	HeaderXUserIDFor   = "X-UserId-For"
	HeaderXTenantIDFor = "X-TenantId-For"
	HeaderXDomainIDFor = "X-DomainId-For"

	HeaderXRequestIDFor  = "X-RequestId-For"
	HeaderXForwardedFor  = "X-Forwarded-For"
	HeaderUserAgent      = "User-Agent"
	HeaderAcceptEncoding = "Accept-Encoding"
	HeaderConnection     = "Connection"
	HeaderContentLength  = "Content-Length"
	HeaderAccept         = "Accept"
	HeaderContentType    = "Content-Type"

	// note: 模仿aws签名算法实现的头
	HeaderXRExDate          = "X-REx-Date"
	HeaderXRExContentSha256 = "X-REx-Content-Sha256"
)
