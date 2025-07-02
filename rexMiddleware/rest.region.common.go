package rexMiddleware

import (
	"context"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
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
		if ctx.Value(CtxClientIp) == nil {
			fullAddr := httpx.GetRemoteAddr(r)
			fullAddrAndPort := strings.Split(fullAddr, ":")
			ctx = context.WithValue(ctx, CtxClientIp, fullAddrAndPort[0])
			logx.Infof("client ip : %s", fullAddrAndPort[0])
			clientIp = fullAddrAndPort[0]
		} else {
			clientIp = ctx.Value(CtxClientIp).(string)
		}

		startTime := time.Now()

		info, _ := m.Region.MemorySearch(clientIp)
		ctx = context.WithValue(ctx, CtxCityId, info.CityId)
		ctx = context.WithValue(ctx, CtxCountry, info.Country)
		ctx = context.WithValue(ctx, CtxRegion, info.Region)
		ctx = context.WithValue(ctx, CtxProvince, info.Province)
		ctx = context.WithValue(ctx, CtxCity, info.City)
		ctx = context.WithValue(ctx, CtxISP, info.ISP)
		endTime := time.Now()
		logc.Infof(ctx, "地理位置中间件耗时: %v", endTime.Sub(startTime).Milliseconds())

		r = r.WithContext(ctx)
		next(w, r)
	}
}
