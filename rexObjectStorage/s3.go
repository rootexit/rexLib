package rexObjectStorage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewCustomS3Client(ctx context.Context, osType string, AccessKey, AccessSecret, Region, AccessEndpoint string) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   osType,
			URL:           AccessEndpoint,
			SigningRegion: Region,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AccessKey, AccessSecret, "")), config.WithEndpointResolverWithOptions(customResolver),
	)
	if cfg.Region == "" {
		cfg.Region = Region
	}
	if err != nil {
		logx.Errorf("LoadDefaultConfig error: %v", err)
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if osType == ObjectStorageTypeCos {
			o.UsePathStyle = false // 使用virtual-host的方式访问
		}
	})
	return s3Client, nil
}
