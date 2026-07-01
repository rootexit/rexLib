package rexMiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexDatabase"
	"github.com/rootexit/rexLib/rexHeaders"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/zeromicro/go-zero/core/logc"
)

type GlobalUaParserInterceptorMiddleware struct {
	Uaparser *uaparser.Parser
	debug    bool
}

func NewGlobalUaParserInterceptorMiddleware(uaparser *uaparser.Parser, isDebug bool) *GlobalUaParserInterceptorMiddleware {
	return &GlobalUaParserInterceptorMiddleware{
		Uaparser: uaparser,
		debug:    isDebug,
	}
}

func (m *GlobalUaParserInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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

		clientInfo := rexDatabase.Client{}

		if ctx.Value(rexCtx.CtxClientInfo{}) == nil {
			clientInfo = rexDatabase.Client{}
			ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)
		} else {
			clientInfo = ctx.Value(rexCtx.CtxClientInfo{}).(rexDatabase.Client)
		}
		if m.debug {
			logc.Infof(ctx, "UaParserInterceptorMiddleware clientInfo: %s", clientInfo)
		}

		uaParseClient := m.Uaparser.Parse(userAgent)

		clientInfo.ClientUa = rexDatabase.ClientUa{
			UserAgent:       userAgent,
			UserAgentFamily: uaParseClient.UserAgent.Family,
			UserAgentMajor:  uaParseClient.UserAgent.Major,
			UserAgentMinor:  uaParseClient.UserAgent.Minor,
			UserAgentPatch:  uaParseClient.UserAgent.Patch,
			OsFamily:        uaParseClient.Os.Family,
			OsMajor:         uaParseClient.Os.Major,
			OsMinor:         uaParseClient.Os.Minor,
			OsPatch:         uaParseClient.Os.Patch,
			OsPatchMinor:    uaParseClient.Os.PatchMinor,
			DeviceFamily:    uaParseClient.Device.Family,
			DeviceBrand:     uaParseClient.Device.Brand,
			DeviceModel:     uaParseClient.Device.Model,
		}

		ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)

		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "UaParserInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}
