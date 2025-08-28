package rexUserAgent

import (
	"context"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"github.com/rootexit/rexLib/rexCommon"
	"github.com/rootexit/rexLib/rexCrypto"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/ua-parser/uap-go/uaparser"
)

func UserAgentUtilsFromIpAndUa(ip, ua string, region *ip2region.Ip2Region, uaparser *uaparser.Parser) (Client, error) {
	// note: 做一次hash用于快速查询
	ipHash := rexCrypto.NewHash().Sha256([]byte(ip))
	tmpC := Client{
		IP:     ip,
		KName:  "",
		IPHash: ipHash,
	}
	ip, port, _, err := rexCommon.ReturnIpAndPort(ip)
	if err != nil {
		return tmpC, nil
	}
	tmpC.Port = port
	tmpC.UserAgent = ua
	info, _ := region.MemorySearch(ip)
	tmpC.CityId = info.CityId
	tmpC.Country = info.Country
	tmpC.Region = info.Region
	tmpC.Province = info.Province
	tmpC.City = info.City
	tmpC.ISP = info.ISP

	uaClient := uaparser.Parse(ua)

	tmpC.UserAgentFamily = uaClient.UserAgent.Family
	tmpC.UserAgentMajor = uaClient.UserAgent.Major
	tmpC.UserAgentMinor = uaClient.UserAgent.Minor
	tmpC.UserAgentPatch = uaClient.UserAgent.Patch
	tmpC.OsFamily = uaClient.Os.Family
	tmpC.OsMajor = uaClient.Os.Major
	tmpC.OsMinor = uaClient.Os.Minor
	tmpC.OsPatch = uaClient.Os.Patch
	tmpC.OsPatchMinor = uaClient.Os.PatchMinor
	tmpC.DeviceFamily = uaClient.Device.Family
	tmpC.DeviceBrand = uaClient.Device.Brand
	tmpC.DeviceModel = uaClient.Device.Model
	return tmpC, nil
}

func UserAgentUtils(ctx context.Context) Client {
	// note: 做一次hash用于快速查询
	ip := ""
	if ctx.Value(rexCtx.CtxClientIp{}) != nil {
		ip = ctx.Value(rexCtx.CtxClientIp{}).(string)
	}
	ipHash := rexCrypto.NewHash().Sha256([]byte(ip))
	tmpC := Client{
		IP:     ip,
		KName:  "",
		IPHash: ipHash,
	}
	if ctx.Value(rexCtx.CtxClientPort{}) != nil {
		tmpC.Port = ctx.Value(rexCtx.CtxClientPort{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgent{}) != nil {
		tmpC.UserAgent = ctx.Value(rexCtx.CtxUserAgent{}).(string)
	}
	if ctx.Value(rexCtx.CtxCityId{}) != nil {
		tmpC.CityId = ctx.Value(rexCtx.CtxCityId{}).(int64)
	}
	if ctx.Value(rexCtx.CtxCountry{}) != nil {
		tmpC.Country = ctx.Value(rexCtx.CtxCountry{}).(string)
	}
	if ctx.Value(rexCtx.CtxRegion{}) != nil {
		tmpC.Region = ctx.Value(rexCtx.CtxRegion{}).(string)
	}
	if ctx.Value(rexCtx.CtxProvince{}) != nil {
		tmpC.Province = ctx.Value(rexCtx.CtxProvince{}).(string)
	}
	if ctx.Value(rexCtx.CtxCity{}) != nil {
		tmpC.City = ctx.Value(rexCtx.CtxCity{}).(string)
	}
	if ctx.Value(rexCtx.CtxISP{}) != nil {
		tmpC.ISP = ctx.Value(rexCtx.CtxISP{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentFamily{}) != nil {
		tmpC.UserAgentFamily = ctx.Value(rexCtx.CtxUserAgentFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentMajor{}) != nil {
		tmpC.UserAgentMajor = ctx.Value(rexCtx.CtxUserAgentMajor{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentMinor{}) != nil {
		tmpC.UserAgentMinor = ctx.Value(rexCtx.CtxUserAgentMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentPatch{}) != nil {
		tmpC.UserAgentPatch = ctx.Value(rexCtx.CtxUserAgentPatch{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsFamily{}) != nil {
		tmpC.OsFamily = ctx.Value(rexCtx.CtxOsFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsMajor{}) != nil {
		tmpC.OsMajor = ctx.Value(rexCtx.CtxOsMajor{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsMinor{}) != nil {
		tmpC.OsMinor = ctx.Value(rexCtx.CtxOsMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsPatch{}) != nil {
		tmpC.OsPatch = ctx.Value(rexCtx.CtxOsPatch{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsPatchMinor{}) != nil {
		tmpC.OsPatchMinor = ctx.Value(rexCtx.CtxOsPatchMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceFamily{}) != nil {
		tmpC.DeviceFamily = ctx.Value(rexCtx.CtxDeviceFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceBrand{}) != nil {
		tmpC.DeviceBrand = ctx.Value(rexCtx.CtxDeviceBrand{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceModel{}) != nil {
		tmpC.DeviceModel = ctx.Value(rexCtx.CtxDeviceModel{}).(string)
	}
	return tmpC
}

func UserAgentUtilsWithFunc(ctx context.Context, callback func(client Client)) Client {
	// note: 做一次hash用于快速查询
	ip := ctx.Value(rexCtx.CtxClientIp{}).(string)
	ipHash := rexCrypto.NewHash().Sha256([]byte(ip))
	tmpC := Client{
		IP:     ip,
		KName:  "",
		IPHash: ipHash,
	}
	if ctx.Value(rexCtx.CtxClientPort{}) != nil {
		tmpC.Port = ctx.Value(rexCtx.CtxClientPort{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgent{}) != nil {
		tmpC.UserAgent = ctx.Value(rexCtx.CtxUserAgent{}).(string)
	}
	if ctx.Value(rexCtx.CtxCityId{}) != nil {
		tmpC.CityId = ctx.Value(rexCtx.CtxCityId{}).(int64)
	}
	if ctx.Value(rexCtx.CtxCountry{}) != nil {
		tmpC.Country = ctx.Value(rexCtx.CtxCountry{}).(string)
	}
	if ctx.Value(rexCtx.CtxRegion{}) != nil {
		tmpC.Region = ctx.Value(rexCtx.CtxRegion{}).(string)
	}
	if ctx.Value(rexCtx.CtxProvince{}) != nil {
		tmpC.Province = ctx.Value(rexCtx.CtxProvince{}).(string)
	}
	if ctx.Value(rexCtx.CtxCity{}) != nil {
		tmpC.City = ctx.Value(rexCtx.CtxCity{}).(string)
	}
	if ctx.Value(rexCtx.CtxISP{}) != nil {
		tmpC.ISP = ctx.Value(rexCtx.CtxISP{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentFamily{}) != nil {
		tmpC.UserAgentFamily = ctx.Value(rexCtx.CtxUserAgentFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentMajor{}) != nil {
		tmpC.UserAgentMajor = ctx.Value(rexCtx.CtxUserAgentMajor{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentMinor{}) != nil {
		tmpC.UserAgentMinor = ctx.Value(rexCtx.CtxUserAgentMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxUserAgentPatch{}) != nil {
		tmpC.UserAgentPatch = ctx.Value(rexCtx.CtxUserAgentPatch{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsFamily{}) != nil {
		tmpC.OsFamily = ctx.Value(rexCtx.CtxOsFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsMajor{}) != nil {
		tmpC.OsMajor = ctx.Value(rexCtx.CtxOsMajor{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsMinor{}) != nil {
		tmpC.OsMinor = ctx.Value(rexCtx.CtxOsMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsPatch{}) != nil {
		tmpC.OsPatch = ctx.Value(rexCtx.CtxOsPatch{}).(string)
	}
	if ctx.Value(rexCtx.CtxOsPatchMinor{}) != nil {
		tmpC.OsPatchMinor = ctx.Value(rexCtx.CtxOsPatchMinor{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceFamily{}) != nil {
		tmpC.DeviceFamily = ctx.Value(rexCtx.CtxDeviceFamily{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceBrand{}) != nil {
		tmpC.DeviceBrand = ctx.Value(rexCtx.CtxDeviceBrand{}).(string)
	}
	if ctx.Value(rexCtx.CtxDeviceModel{}) != nil {
		tmpC.DeviceModel = ctx.Value(rexCtx.CtxDeviceModel{}).(string)
	}
	if callback != nil {
		callback(tmpC)
	}
	return tmpC
}
