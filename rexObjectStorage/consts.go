package rexObjectStorage

import (
	"fmt"
	"log"
)

const (
	ObjectStorageTypeMinio = "minio"
	ObjectStorageTypeCos   = "cos"
	ObjectStorageTypeOss   = "oss"
	ObjectStorageTypeS3    = "s3"
)

var ObjectStorageTypeSupport = map[string]bool{
	ObjectStorageTypeMinio: true,
	ObjectStorageTypeCos:   true,
	ObjectStorageTypeOss:   true,
	ObjectStorageTypeS3:    true,
}

func CheckObjectStorageSupport(osType string) bool {
	if _, ok := ObjectStorageTypeSupport[osType]; ok {
		return true
	} else {
		return false
	}
}

func FormatBucketDomainByOsType(OsType string, name, region string, v ...any) (BucketInternetDomain, BucketInternalDomain, BucketAccelerateDomain string) {
	customDomain := ""
	if len(v) > 0 {
		customDomain = v[0].(string)
		if customDomain == "" || len(customDomain) <= 0 {
			customDomain = ""
		}
	}
	switch OsType {
	case ObjectStorageTypeCos:
		BucketInternetDomain = fmt.Sprintf("%s.cos.%s.myqcloud.com", name, region)
		BucketInternalDomain = fmt.Sprintf("%s.cos.%s.myqcloud.com", name, region)
		BucketAccelerateDomain = fmt.Sprintf("%s.cos.%s.myqcloud.com", name, region)
	case ObjectStorageTypeOss:
		BucketInternetDomain = fmt.Sprintf("%s.oss-%s.aliyuncs.com", name, region)
		BucketInternalDomain = fmt.Sprintf("%s.oss-%s-internal.com", name, region)
		BucketAccelerateDomain = fmt.Sprintf("%s.oss-accelerate.aliyuncs.com", name)
	case ObjectStorageTypeS3:
		BucketInternetDomain = fmt.Sprintf("%s.s3.%s.amazonaws.com", name, region)
		BucketInternalDomain = fmt.Sprintf("%s.s3.%s.amazonaws.com", name, region)
		BucketAccelerateDomain = fmt.Sprintf("%s.s3.%s.amazonaws.com", name, region)
	case ObjectStorageTypeMinio:
		BucketInternetDomain = fmt.Sprintf("%s.%s", name, customDomain)
		BucketInternalDomain = fmt.Sprintf("%s.%s", name, customDomain)
		BucketAccelerateDomain = fmt.Sprintf("%s.%s", name, customDomain)
	default:
		log.Printf("Unsupported bucket type %s.\n", OsType)
	}
	return BucketInternetDomain, BucketInternalDomain, BucketAccelerateDomain
}
