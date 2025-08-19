package rexCloudConf

type CloudflareWafAndLoadBalancingCloudConf struct {
	RequestWafTraceIdHeader           string `json:",default=unknown"` // 请求
	RequestLoadBalancingTraceIdHeader string `json:",default=unknown"` // 请求追踪ID
	RealIpHeader                      string `json:",default=CF-Connecting-IP"`
}
