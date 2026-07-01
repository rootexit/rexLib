package rexMiddleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"github.com/rootexit/rexLib/rexCommon"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/zeromicro/go-zero/core/logc"
)

type RegionInterceptorMiddleware struct {
	Region *ip2region.Ip2Region
	debug  bool
}

func NewRegionInterceptorMiddleware(region *ip2region.Ip2Region, isDebug bool) *RegionInterceptorMiddleware {
	return &RegionInterceptorMiddleware{
		Region: region,
		debug:  isDebug,
	}
}

func (m *RegionInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		startTime := time.Now()

		clientIp := ""
		if ctx.Value(rexCtx.CtxClientIp{}) == nil {
			fullAddr := rexCommon.GetRemoteClientAddr(r)
			ips := strings.Split(fullAddr, ",")
			realAddr := ips[0]
			ip, port, ipType, err := rexCommon.ReturnIpAndPort(realAddr)
			if err != nil {
				logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
				http.Error(w, "Unknown IP format", http.StatusNotImplemented)
				return
			}
			if m.debug {
				logc.Infof(ctx, "RegionInterceptorMiddleware realAddr: %s, ip: %s, port: %s, ipType: %s", realAddr, ip, port, ipType)
			}

			ctx = context.WithValue(ctx, rexCtx.CtxClientIp{}, ip)
			ctx = context.WithValue(ctx, rexCtx.CtxClientPort{}, port)
			logc.Infof(ctx, "RegionInterceptorMiddleware Ip: %s, Port: %s", ip, port)
		} else {
			clientIp = ctx.Value(rexCtx.CtxClientIp{}).(string)
		}

		// note: 判断是否是国内的还是国外的IP地址

		info, _ := m.Region.MemorySearch(clientIp)
		ctx = context.WithValue(ctx, rexCtx.CtxCityId{}, info.CityId)
		ctx = context.WithValue(ctx, rexCtx.CtxCountry{}, info.Country)
		ctx = context.WithValue(ctx, rexCtx.CtxRegion{}, info.Region)
		ctx = context.WithValue(ctx, rexCtx.CtxProvince{}, info.Province)
		ctx = context.WithValue(ctx, rexCtx.CtxCity{}, info.City)
		ctx = context.WithValue(ctx, rexCtx.CtxISP{}, info.ISP)
		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "RegionInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}
