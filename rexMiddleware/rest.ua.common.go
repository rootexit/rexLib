package rexMiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexHeaders"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/zeromicro/go-zero/core/logc"
)

type UaParserInterceptorMiddleware struct {
	Uaparser *uaparser.Parser
	debug    bool
}

func NewUaParserInterceptorMiddleware(uaparser *uaparser.Parser, isDebug bool) *UaParserInterceptorMiddleware {
	return &UaParserInterceptorMiddleware{
		Uaparser: uaparser,
		debug:    isDebug,
	}
}

func (m *UaParserInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		ctx := r.Context()

		userAgent := ""
		if ctx.Value(rexCtx.CtxUserAgent{}) == nil {
			userAgent = r.Header.Get(rexHeaders.HeaderUserAgent)
			ctx = context.WithValue(ctx, rexCtx.CtxUserAgent{}, userAgent)
		} else {
			userAgent = ctx.Value(rexCtx.CtxUserAgent{}).(string)
		}
		if m.debug {
			logc.Infof(ctx, "UaParserInterceptorMiddleware userAgent: %s", userAgent)
		}

		client := m.Uaparser.Parse(userAgent)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgentFamily{}, client.UserAgent.Family)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgentMajor{}, client.UserAgent.Major)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgentMinor{}, client.UserAgent.Minor)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgentPatch{}, client.UserAgent.Patch)

		ctx = context.WithValue(ctx, rexCtx.CtxOsFamily{}, client.Os.Family)
		ctx = context.WithValue(ctx, rexCtx.CtxOsMajor{}, client.Os.Major)
		ctx = context.WithValue(ctx, rexCtx.CtxOsMinor{}, client.Os.Minor)
		ctx = context.WithValue(ctx, rexCtx.CtxOsPatch{}, client.Os.Patch)
		ctx = context.WithValue(ctx, rexCtx.CtxOsPatchMinor{}, client.Os.PatchMinor)

		ctx = context.WithValue(ctx, rexCtx.CtxDeviceFamily{}, client.Device.Family)
		ctx = context.WithValue(ctx, rexCtx.CtxDeviceBrand{}, client.Device.Brand)
		ctx = context.WithValue(ctx, rexCtx.CtxDeviceModel{}, client.Device.Model)
		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "UaParserInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}
