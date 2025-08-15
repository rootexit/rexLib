package rexMiddleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexHeaders"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net"
	"net/http"
	"strings"
	"time"
)

type PathHttpInterceptorMiddleware struct {
}

func NewPathHttpInterceptorMiddleware() *PathHttpInterceptorMiddleware {
	return &PathHttpInterceptorMiddleware{}
}

func (m *PathHttpInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := context.WithValue(r.Context(), rexCtx.CtxFullMethod{}, r.URL.Path)
		ctx = context.WithValue(ctx, rexCtx.CtxRequestURI{}, r.RequestURI)
		ctx = context.WithValue(ctx, rexCtx.CtxStartTime{}, startTime.UnixMilli())

		fullAddr := httpx.GetRemoteAddr(r)
		logc.Infof(ctx, "fullAddr:%v", fullAddr)
		ips := strings.Split(fullAddr, ",")
		realAddr := ips[0]
		ip, port, ipType, err := returnIpAndPort(realAddr)
		if err != nil {
			logc.Infof(ctx, "unknown ip format: %s", err)
			http.Error(w, "unknown ip format", http.StatusNotImplemented)
			return
		}
		logc.Infof(ctx, "realAddr: %s, ip: %s, port: %s, ipType: %s", realAddr, ip, port, ipType)

		ctx = context.WithValue(ctx, rexCtx.CtxClientIp{}, ip)
		ctx = context.WithValue(ctx, rexCtx.CtxClientPort{}, port)
		logc.Infof(ctx, "IP: %s, Port: %s", ip, port)
		if err != nil {
			logx.Infof("解析ip报错: %s", err)
			http.Error(w, "不支持的ip类型", http.StatusNotImplemented)
			return
		}

		requestID := uuid.NewString()
		ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, requestID)
		w.Header().Set(rexHeaders.HeaderXRequestIDFor, requestID)

		// 获取 User-Agent
		userAgent := r.Header.Get(rexHeaders.HeaderUserAgent)
		ctx = context.WithValue(ctx, rexCtx.CtxUserAgent{}, userAgent)

		endTime := time.Now()
		logc.Infof(ctx, "路径ip处理中间件耗时: %v", endTime.Sub(startTime).Milliseconds())

		r = r.WithContext(ctx)
		next(w, r)
	}
}

func returnIpAndPort(ipStr string) (ip, port, ipType string, err error) {
	ip = ""
	port = ""
	ipType = "IPv4"
	// note: 先判断是否是 IPv6
	if strings.Count(ipStr, ":") > 1 {
		if strings.Contains(ipStr, "[") {
			// note: ipv6带端口
			host, p, err := net.SplitHostPort(ipStr)
			if err != nil {
				logx.Errorf("❌ 无效的带端口 IPv6:", err)
				return ip, port, ipType, err
			}
			ip = host
			port = p
		} else {
			ip = ipStr
			port = ""
		}
	} else {
		// note: 判断不是ipv4
		if strings.Contains(ipStr, ":") {
			host, p, err := net.SplitHostPort(ipStr)
			if err != nil {
				logx.Errorf("❌ 无效的带端口 IPv6:", err)
				return ip, port, ipType, err
			}
			ip = host
			port = p
		} else {
			// 纯 IP，无端口
			ip = ipStr
			port = ""
		}
	}

	// 解析 IP
	netIp := net.ParseIP(strings.Trim(ip, "[]"))
	if netIp == nil {
		logx.Errorf("❌ 无效的 IP 地址")
		return ip, port, ipType, errors.New("unknown ip format")
	}
	// 类型判断
	if netIp.To4() != nil {
		ipType = "IPv4"
		return ip, port, ipType, nil
	} else {
		ipType = "IPv6"
		return ip, port, ipType, nil
	}
}
