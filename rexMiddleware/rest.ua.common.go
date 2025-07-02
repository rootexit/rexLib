package rexMiddleware

import (
	"context"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"time"
)

type UaParserInterceptorMiddleware struct {
	Uaparser *uaparser.Parser
}

func NewUaParserInterceptorMiddleware(uaparser *uaparser.Parser) *UaParserInterceptorMiddleware {
	return &UaParserInterceptorMiddleware{
		Uaparser: uaparser,
	}
}

func (m *UaParserInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		userAgent := ""
		if ctx.Value(CtxUserAgent) == nil {
			userAgent = r.Header.Get(CtxUserAgent)
			ctx = context.WithValue(ctx, CtxUserAgent, userAgent)
		} else {
			userAgent = ctx.Value(CtxUserAgent).(string)
		}

		startTime := time.Now()
		client := m.Uaparser.Parse(userAgent)
		ctx = context.WithValue(ctx, CtxUserAgentFamily, client.UserAgent.Family)
		ctx = context.WithValue(ctx, CtxUserAgentMajor, client.UserAgent.Major)
		ctx = context.WithValue(ctx, CtxUserAgentMinor, client.UserAgent.Minor)
		ctx = context.WithValue(ctx, CtxUserAgentPatch, client.UserAgent.Patch)

		ctx = context.WithValue(ctx, CtxOsFamily, client.Os.Family)
		ctx = context.WithValue(ctx, CtxOsMajor, client.Os.Major)
		ctx = context.WithValue(ctx, CtxOsMinor, client.Os.Minor)
		ctx = context.WithValue(ctx, CtxOsPatch, client.Os.Patch)
		ctx = context.WithValue(ctx, CtxOsPatchMinor, client.Os.PatchMinor)

		ctx = context.WithValue(ctx, CtxDeviceFamily, client.Device.Family)
		ctx = context.WithValue(ctx, CtxDeviceBrand, client.Device.Brand)
		ctx = context.WithValue(ctx, CtxDeviceModel, client.Device.Model)
		endTime := time.Now()
		logc.Infof(ctx, "设备解析中间件耗时: %v", endTime.Sub(startTime).Milliseconds())

		r = r.WithContext(ctx)
		next(w, r)
	}
}
