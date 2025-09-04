package rexOauth

type OauthConfig struct {
	ClientId      string `json:",default=client_id"`
	ClientSecret  string `json:",default=client_secret"`
	Issuer        string `json:",default=https://idp.xx.com"`
	RedirectUrl   string `json:",default=https://app.xx.com/oauth/callback"`
	AuthorizeUrl  string `json:",default=/oauth/authorize"`
	TokenUrl      string `json:",default=/oauth/token"`
	RevokeUrl     string `json:",default=/oauth/revoke"`
	IntrospectUrl string `json:",default=/oauth/introspect"`
}

type OauthRefreshTokenReq struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type OauthTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type OauthTokenResp struct {
	Code      int32              `json:"code"`
	Msg       string             `json:"msg"`
	Path      string             `json:"path"`
	RequestId string             `json:"request_id"`
	Data      OauthTokenRespData `json:"data"`
}

type OauthTokenRespData struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type OauthRevokeResp struct {
	Code      int32  `json:"code"`
	Msg       string `json:"msg"`
	Path      string `json:"path"`
	RequestId string `json:"request_id"`
	Data      string `json:"data"`
}
