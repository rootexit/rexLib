package rexMiddleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rootexit/rexLib/rexCommon"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexHeaders"
	"github.com/rootexit/rexLib/rexUlid"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type PathHttpInterceptorMiddleware struct {
	isAllowedInheritRequestId bool
	debug                     bool
}

func NewPathHttpInterceptorMiddleware(isAllowedInheritRequestId, isDebug bool) *PathHttpInterceptorMiddleware {
	return &PathHttpInterceptorMiddleware{
		isAllowedInheritRequestId: isAllowedInheritRequestId,
		debug:                     isDebug,
	}
}

func (m *PathHttpInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := context.WithValue(r.Context(), rexCtx.CtxFullMethod{}, r.URL.Path)
		ctx = context.WithValue(ctx, rexCtx.CtxRequestURI{}, r.RequestURI)
		ctx = context.WithValue(ctx, rexCtx.CtxStartTime{}, startTime.UnixMilli())

		fullAddr := httpx.GetRemoteAddr(r)
		if m.debug {
			logc.Infof(ctx, "PathHttpInterceptorMiddleware fullAddr :%v", fullAddr)
		}
		ips := strings.Split(fullAddr, ",")
		realAddr := ips[0]
		ip, port, ipType, err := rexCommon.ReturnIpAndPort(realAddr)
		if err != nil {
			logc.Errorf(ctx, "PathHttpInterceptorMiddleware unknown ip format: %s", err)
			http.Error(w, "Unknown IP format", http.StatusNotImplemented)
			return
		}
		if m.debug {
			logc.Infof(ctx, "PathHttpInterceptorMiddleware realAddr: %s, ip: %s, port: %s, ipType: %s", realAddr, ip, port, ipType)
		}
		ctx = context.WithValue(ctx, rexCtx.CtxClientIp{}, ip)
		ctx = context.WithValue(ctx, rexCtx.CtxClientPort{}, port)
		ctx = context.WithValue(ctx, rexCtx.CtxClientType{}, ipType)
		if err != nil {
			logc.Errorf(ctx, "PathHttpInterceptorMiddleware parse ip err: %s", err)
			http.Error(w, "Unsupported IP types", http.StatusNotImplemented)
			return
		}

		requestID := rexUlid.NewString()
		if m.isAllowedInheritRequestId {
			if r.Header.Get(rexHeaders.HeaderXRequestIdFor) != "" {
				// 如果请求头中有 HeaderXRequestIDFor，则使用它
				if m.debug {
					logc.Infof(ctx, "PathHttpInterceptorMiddleware inheritance request ID: %s", requestID)
				}
				requestID = r.Header.Get(rexHeaders.HeaderXRequestIdFor)
			}
		}

		ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, requestID)
		w.Header().Set(rexHeaders.HeaderXRequestIdFor, requestID)

		// 获取 User-Agent
		userAgent := r.Header.Get(rexHeaders.HeaderUserAgent)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgent{}, userAgent)

		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "PathHttpInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}
