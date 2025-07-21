package rexMiddleware

import (
	"context"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
	"time"
)

type RegionInterceptorMiddleware struct {
	Region *ip2region.Ip2Region
}

func NewRegionInterceptorMiddleware(region *ip2region.Ip2Region) *RegionInterceptorMiddleware {
	return &RegionInterceptorMiddleware{
		Region: region,
	}
}

func (m *RegionInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		clientIp := ""
		if ctx.Value(rexCtx.CtxClientIp{}) == nil {
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
		} else {
			clientIp = ctx.Value(rexCtx.CtxClientIp{}).(string)
		}

		startTime := time.Now()

		info, _ := m.Region.MemorySearch(clientIp)
		ctx = context.WithValue(ctx, rexCtx.CtxCityId{}, info.CityId)
		ctx = context.WithValue(ctx, rexCtx.CtxCountry{}, info.Country)
		ctx = context.WithValue(ctx, rexCtx.CtxRegion{}, info.Region)
		ctx = context.WithValue(ctx, rexCtx.CtxProvince{}, info.Province)
		ctx = context.WithValue(ctx, rexCtx.CtxCity{}, info.City)
		ctx = context.WithValue(ctx, rexCtx.CtxISP{}, info.ISP)
		endTime := time.Now()
		logc.Infof(ctx, "地理位置中间件耗时: %v", endTime.Sub(startTime).Milliseconds())

		r = r.WithContext(ctx)
		next(w, r)
	}
}
