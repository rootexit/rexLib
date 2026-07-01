package rexMiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"github.com/oschwald/geoip2-golang/v2"
	"github.com/rootexit/rexLib/rexCommon"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexDatabase"
	"github.com/zeromicro/go-zero/core/logc"
)

type GlobalRegionInterceptorMiddleware struct {
	cityDB *geoip2.Reader
	asnDB  *geoip2.Reader
	region *ip2region.Ip2Region
	debug  bool
}

func NewGlobalRegionInterceptorMiddleware(cityDB, asnDB *geoip2.Reader, region *ip2region.Ip2Region, isDebug bool) *GlobalRegionInterceptorMiddleware {
	return &GlobalRegionInterceptorMiddleware{
		cityDB: cityDB,
		asnDB:  asnDB,
		region: region,
		debug:  isDebug,
	}
}

func (m *GlobalRegionInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		startTime := time.Now()

		clientIp := ""
		clientPort := ""
		clientInfo := rexDatabase.Client{}
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
			clientPort = ctx.Value(rexCtx.CtxClientPort{}).(string)
		}
		if m.debug {
			logc.Infof(ctx, "RegionInterceptorMiddleware clientIp: %s", clientIp)
		}

		internal, reason, err := rexCommon.IsInternalIP(clientIp)
		if err != nil {
			logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
			http.Error(w, "Unknown IP format", http.StatusNotImplemented)
			return
		}
		if internal {
			clientInfo = rexDatabase.Client{
				ClientNetwork: rexDatabase.ClientNetwork{
					IpAddress: clientIp,
					Port:      clientPort,
					Network:   reason,
					Isp:       reason,
				},
				ClientLocation: rexDatabase.ClientLocation{
					Continent:      "内网",
					Country:        "内网",
					Province:       "内网",
					City:           "内网",
					Longitude:      0,
					Latitude:       0,
					TimeZone:       "localhost",
					AccuracyRadius: 0,
				},
			}
			ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)
			endTime := time.Now()
			if m.debug {
				logc.Infof(ctx, "RegionInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
			}
			r = r.WithContext(ctx)
			next(w, r)
			return
		}

		ip, err := netip.ParseAddr(clientIp)
		if err != nil {
			logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
			http.Error(w, "Unknown IP format", http.StatusNotImplemented)
			return
		}
		// note: 判断是否是国内的还是国外的IP地址
		city, err := m.cityDB.City(ip)
		if err != nil {
			logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
			http.Error(w, "Unknown IP format", http.StatusNotImplemented)
			return
		}
		asn, err := m.asnDB.ASN(ip)
		if err != nil {
			logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
			http.Error(w, "Unknown IP format", http.StatusNotImplemented)
			return
		}
		if m.debug {
			asnJson, _ := json.Marshal(asn)
			logc.Infof(ctx, "RegionInterceptorMiddleware asn: %v", string(asnJson))
		}
		clientInfo = rexDatabase.Client{
			ClientNetwork: rexDatabase.ClientNetwork{
				IpAddress:                    clientIp,
				Port:                         clientPort,
				Network:                      asn.Network.String(),
				Isp:                          "unknown",
				AutonomousSystemOrganization: asn.AutonomousSystemOrganization,
				AutonomousSystemNumber:       asn.AutonomousSystemNumber,
			},
			ClientLocation: rexDatabase.ClientLocation{
				Continent:      city.Continent.Names.SimplifiedChinese,
				Country:        city.Country.Names.SimplifiedChinese,
				Province:       "",
				City:           city.City.Names.SimplifiedChinese,
				Longitude:      0,
				Latitude:       0,
				TimeZone:       city.Location.TimeZone,
				AccuracyRadius: city.Location.AccuracyRadius,
			},
			ClientUa:  rexDatabase.ClientUa{},
			ClientBot: rexDatabase.ClientBot{},
		}
		if len(city.Subdivisions) > 0 {
			tmpSubdivisions := make([]string, len(city.Subdivisions))
			for i, subdivision := range city.Subdivisions {
				tmpSubdivisions[i] = subdivision.Names.SimplifiedChinese
			}
			tmpProvince := strings.Join(tmpSubdivisions, " ")
			clientInfo.ClientLocation.Province = tmpProvince
		} else if len(city.Subdivisions) == 1 {
			clientInfo.ClientLocation.Province = city.Subdivisions[0].Names.SimplifiedChinese
		} else {
			clientInfo.ClientLocation.Province = ""
		}
		if city.Location.HasData() {
			clientInfo.ClientLocation.Longitude = *city.Location.Longitude
			clientInfo.ClientLocation.Latitude = *city.Location.Latitude
		}
		if m.debug {
			logc.Infof(ctx, "RegionInterceptorMiddleware Country ISOCode: %s", city.Country.ISOCode)
		}
		if city.Country.ISOCode == "CN" {
			// 中国
			info, err := m.region.MemorySearch(clientIp)
			if err != nil {
				logc.Errorf(ctx, "RegionInterceptorMiddleware unknown ip format: %s", err)
				http.Error(w, "Unknown IP format", http.StatusNotImplemented)
				return
			}
			clientInfo.ClientNetwork.Isp = info.ISP
			clientInfo.ClientLocation.Country = info.Country
			clientInfo.ClientLocation.Province = info.Province
			clientInfo.ClientLocation.City = info.City
		}
		ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)
		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "RegionInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}
