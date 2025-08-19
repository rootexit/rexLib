package rexCustomAwsSign

import (
	"github.com/rootexit/rexLib/rexHeaders"
	"time"
)

const (
	timeFormat = "20060102T150405Z"
	//authHeaderPrefix = "AWS4-HMAC-SHA256"
	authHeaderPrefix = "LIL4-HMAC-SHA256"
	shortTimeFormat  = "20060102"
	awsV4Request     = "lil4_request"

	// emptyStringSHA256 is a SHA256 of an empty string
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
	maxSkew           = 5 * time.Minute
)

var ignoredHeaders = map[string]string{
	rexHeaders.HeaderAuthorization: "",
	rexHeaders.HeaderUserAgent:     "",
	rexHeaders.HeaderUserAgent:     "",
	"X-Amzn-Trace-Id":              "",
}
