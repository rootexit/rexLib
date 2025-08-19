package rexCors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	allowMethods  = "Access-Control-Allow-Methods"
	allowHeaders  = "Access-Control-Allow-Headers"
	exposeHeaders = "Access-Control-Expose-Headers"
	maxAgeHeader  = "Access-Control-Max-Age"
)

type Cors struct {
	AllowedOrigins  []string `json:",default=[http://localhost:8888]"`
	AllowMethods    []string `json:",default=[GET,POST]"`
	AllowHeaders    []string `json:",default=[Content-Type,Origin,X-CSRF-Token,Authorization,X-Request-ID,Host,X-Lil-Content-Sha256,X-Lil-Date,X-Amz-Content-Sha256,X-Amz-Date,Range]"`
	ExposeHeaders   []string `json:",default=[Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers]"`
	MaxAgeHeaderVal int64    `json:",default=86400"`
}

func CustomNotAllowedFn() func(http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		// note: 请求了没注册的方法
		// 可选：写入你希望暴露的 CORS 响应头
		w.Header().Set("Content-Type", "application/json")

		// 可选：自定义错误响应体
		resp := map[string]string{
			"error":   "Method Not Allowed",
			"message": "The requested method is not allowed on this endpoint.",
		}

		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func CustomMiddlewareFn(allowMethodsInput []string, allowHeadersInput []string, exposeHeadersInput []string, maxAgeHeaderVal int64) func(header http.Header) {
	return func(header http.Header) {
		tempMethods := strings.Join(allowMethodsInput, ",")
		header.Set(allowMethods, tempMethods)
		tempAllowHeaders := strings.Join(allowHeadersInput, ",")
		header.Set(allowHeaders, tempAllowHeaders)
		tempExposeHeaders := strings.Join(exposeHeadersInput, ",")
		header.Set(exposeHeaders, tempExposeHeaders)
		header.Set(maxAgeHeader, fmt.Sprintf("%d", maxAgeHeaderVal))
	}
}
