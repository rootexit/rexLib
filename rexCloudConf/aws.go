package rexCloudConf

type AwsWafAndLoadBalancingCloudConf struct {
	RequestWafTraceIdHeader           string `json:",default=unknown"`         // 请求
	RequestLoadBalancingTraceIdHeader string `json:",default=X-Amzn-Trace-Id"` // 请求追踪ID
	RealIpHeader                      string `json:",default=X-Forwarded-For"`
}
