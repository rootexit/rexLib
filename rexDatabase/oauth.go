package rexDatabase

//EXISTING_USER       // 老用户正常登录
//AUTO_REGISTER       // 自动注册
//FIRST_LOGIN         // 首次登录
//ACCOUNT_LINKED      // 第三方账号绑定
//INVITED_REGISTER    // 邀请注册

const (
	OauthMetaLoginTypeExistingUser    = "EXISTING_USER"
	OauthMetaLoginTypeAutoRegister    = "AUTO_REGISTER"
	OauthMetaLoginTypeFirstLogin      = "FIRST_LOGIN"
	OauthMetaLoginTypeAccountLinked   = "ACCOUNT_LINKED"
	OauthMetaLoginTypeInvitedRegister = "INVITED_REGISTER"
)
