package rexOauth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rootexit/rexLib/rexCodes"
	"github.com/rootexit/rexLib/rexHeaders"
	"io"
	"net/http"
)

type (
	Oauth interface {
		GetConfig() OauthConfig
		Code2Token(code string) (resultCode int32, result *OauthTokenRespData, err error)
		RefreshAccessToken(rt string) (resultCode int32, result *OauthTokenRespData, err error)
	}
	defaultOauth struct {
		clientId      string
		clientSecret  string
		issuer        string
		redirectUrl   string
		authorizeUrl  string
		tokenUrl      string
		revokeUrl     string
		introspectUrl string
		client        *http.Client
	}
)

func NewOauth(c *OauthConfig, client *http.Client) (Oauth, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}
	if c.ClientId == "" || c.ClientSecret == "" || c.Issuer == "" || c.RedirectUrl == "" || c.AuthorizeUrl == "" || c.TokenUrl == "" {
		return nil, errors.New("config is invalid")
	}
	return &defaultOauth{
		clientId:      c.ClientId,
		clientSecret:  c.ClientSecret,
		issuer:        c.Issuer,
		redirectUrl:   c.RedirectUrl,
		authorizeUrl:  fmt.Sprintf("%s%s", c.Issuer, c.AuthorizeUrl),
		tokenUrl:      fmt.Sprintf("%s%s", c.Issuer, c.TokenUrl),
		revokeUrl:     fmt.Sprintf("%s%s", c.Issuer, c.RevokeUrl),
		introspectUrl: fmt.Sprintf("%s%s", c.Issuer, c.IntrospectUrl),
		client:        client,
	}, nil
}

func (o *defaultOauth) GetConfig() OauthConfig {
	return OauthConfig{
		ClientId:     o.clientId,
		ClientSecret: o.clientSecret,
		Issuer:       o.issuer,
		RedirectUrl:  o.redirectUrl,
		AuthorizeUrl: o.authorizeUrl,
		TokenUrl:     o.tokenUrl,
		RevokeUrl:    o.revokeUrl,
	}
}

func (o *defaultOauth) Code2Token(code string) (resultCode int32, result *OauthTokenRespData, err error) {

	tmpResult := OauthTokenResp{}
	// note: 撰写请求，并发送
	sendBody := OauthTokenReq{
		GrantType:    "authorization_code",
		ClientId:     o.clientId,
		ClientSecret: o.clientSecret,
		RedirectUri:  o.redirectUrl,
		Code:         code,
	}
	sendBodyBt, err := json.Marshal(sendBody)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl, bytes.NewBuffer(sendBodyBt))
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	tempAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", o.clientId, o.clientSecret))))
	req.Header.Set(rexHeaders.HeaderAuthorization, tempAuth)
	var res *http.Response
	res, err = o.client.Do(req)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	if err = json.Unmarshal(body, &tmpResult); err != nil {
		return rexCodes.FAIL, nil, err
	}
	if tmpResult.Code != rexCodes.OK {
		return tmpResult.Code, nil, fmt.Errorf("request err: %v", tmpResult.Msg)
	}
	return tmpResult.Code, &tmpResult.Data, nil
}

func (o *defaultOauth) RefreshAccessToken(rt string) (resultCode int32, result *OauthTokenRespData, err error) {

	tmpResult := OauthTokenResp{}
	// note: 撰写请求，并发送
	sendBody := OauthRefreshTokenReq{
		GrantType:    "refresh_token",
		RefreshToken: rt,
	}
	sendBodyBt, err := json.Marshal(sendBody)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl, bytes.NewBuffer(sendBodyBt))
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	tempAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", o.clientId, o.clientSecret))))
	req.Header.Set(rexHeaders.HeaderAuthorization, tempAuth)
	var res *http.Response
	res, err = o.client.Do(req)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rexCodes.FAIL, nil, err
	}
	if err = json.Unmarshal(body, &tmpResult); err != nil {
		return rexCodes.FAIL, nil, err
	}
	if tmpResult.Code != rexCodes.OK {
		return tmpResult.Code, nil, fmt.Errorf("request err: %v", tmpResult.Msg)
	}
	return tmpResult.Code, &tmpResult.Data, nil
}

func (o *defaultOauth) RevokeRefreshToken() (resultCode int32, result *OauthRevokeResp, err error) {
	return 0, nil, nil
}
