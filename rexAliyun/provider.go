package rexAliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type defaultCredentials struct {
	config *oss.Config
}

func (defCre *defaultCredentials) GetAccessKeyID() string {
	return defCre.config.AccessKeyID
}

func (defCre *defaultCredentials) GetAccessKeySecret() string {
	return defCre.config.AccessKeySecret
}

func (defCre *defaultCredentials) GetSecurityToken() string {
	return defCre.config.SecurityToken
}

type defaultCredentialsProvider struct {
	config *oss.Config
}

func (defBuild *defaultCredentialsProvider) GetCredentials() oss.Credentials {
	return &defaultCredentials{config: defBuild.config}
}

func NewDefaultCredentialsProvider(accessID, accessKey, token string) (defaultCredentialsProvider, error) {
	var provider defaultCredentialsProvider
	if accessID == "" {
		return provider, fmt.Errorf("access key id is empty!")
	}
	if accessKey == "" {
		return provider, fmt.Errorf("access key secret is empty!")
	}
	config := &oss.Config{
		AccessKeyID:     accessID,
		AccessKeySecret: accessKey,
		SecurityToken:   token,
	}
	return defaultCredentialsProvider{
		config,
	}, nil
}
