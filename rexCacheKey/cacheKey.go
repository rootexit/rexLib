package rexCacheKey

const (
	CsrfCacheKey = "%s:csrf-%s"

	OnlineKey = "%s:online-%s"

	ACCESS_TOKEN_KEY = "%s:accessToken-%s"

	REFRESH_TOKEN_KEY = "%s:refreshToken-%s"

	ENCRYPTION_DATA_KEY = "%s:encryption-%s"

	HTTP_ACCESS_TOKEN_KEY = "%s:httpAccessToken-%d"

	ACCOUNT_AUTH_KEY = "%s:accountAuth-%s"

	EMS_AUTH_KEY = "%s:emsAuth-%s"

	SMS_AUTH_KEY = "%s:smsAuth-%s"

	WECHAT_KEY_CRON = "%s:wechatCron-TID%d-key-%s"

	// note: 服务名+公众号的Appid
	WECHAT_APPID_CRON = "%s:wx-access-%s"

	WECHAT_USER_TOKEN = "%s:wechat-usertokenn-%s"

	WECHAT_JSAPI_TICKET_KEY_CRON = "%s:jsApiTicket-TID%d-key-%s"
	WECHAT_TICKET_APPID_CRON     = "%s:wx-ticket-%s"

	// note: 服务名+抖音的Appid
	DOUYIN_APPID_CRON = "%s:dy-access-%s"

	DOUYIN_USER_TOKEN = "%s:dy-usertokenn-%s"

	DOUYIN_JSAPI_TICKET_KEY_CRON = "%s:jsApiTicket-TID%d-key-%s"
	DOUYIN_TICKET_APPID_CRON     = "%s:dy-ticket-%s"

	KeyChangeCacheKey = "%s-cache-key-%s"
)
