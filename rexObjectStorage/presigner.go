package rexObjectStorage

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.

type (
	Presigner interface {
		GetPresignClient() *s3.PresignClient
		GetObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
		PutObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
		DeleteObject(ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error)
		PresignPostObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error)
		PresignPutObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
		PresignHeadObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
	}
	defaultPresigner struct {
		presignClient *s3.PresignClient
	}
)

func NewPresigner(s3Client *s3.Client) Presigner {
	return &defaultPresigner{
		presignClient: s3.NewPresignClient(s3Client),
	}
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (p *defaultPresigner) GetPresignClient() *s3.PresignClient {
	return p.presignClient
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (p *defaultPresigner) GetObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return nil, err
	}
	return request, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (p *defaultPresigner) PutObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return nil, err
	}
	return request, err
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (p *defaultPresigner) DeleteObject(ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := p.presignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
		return nil, err
	}
	return request, err
}

func (p *defaultPresigner) PresignPostObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error) {
	request, err := p.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(lifetimeSecs) * time.Second
	})
	if err != nil {
		log.Printf("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}
	return request, nil
}

func (p *defaultPresigner) PresignPutObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(lifetimeSecs) * time.Second
	})
	if err != nil {
		log.Printf("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}
	return request, nil
}

func (p *defaultPresigner) PresignHeadObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := p.presignClient.PresignHeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(lifetimeSecs) * time.Second
	})
	if err != nil {
		log.Printf("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}
	return request, nil
}
