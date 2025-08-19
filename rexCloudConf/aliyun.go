package rexCloudConf

type AliyunWafAndLoadBalancingCloudConf struct {
	RequestWafTraceIdHeader           string `json:",default=Eagleeye-Traceid"` // 请求
	RequestLoadBalancingTraceIdHeader string `json:",default=X-B3-Traceid"`     // 请求追踪ID
	RealIpHeader                      string `json:",default=X-Forwarded-For"`
}
