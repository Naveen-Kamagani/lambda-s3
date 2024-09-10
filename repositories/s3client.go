package repositories

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ClientInterface interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type S3Client struct {
	Client S3ClientInterface
}

func InitializeS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func (c *S3Client) PutObject(ctx context.Context, body io.Reader, bucket string, key string, acl types.ObjectCannedACL) (*s3.PutObjectOutput, error) {

	if body == nil {
		return nil, errors.New("no body was passed in")
	}
	if bucket == "" {
		return nil, errors.New("no bucket was passed in")
	}
	if key == "" {
		return nil, errors.New("no key was passed in")
	}
	if acl == "" {
		acl = types.ObjectCannedACLPrivate
	}
	i := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
		ACL:    acl,
	}
	return c.Client.PutObject(ctx, i)
}

func (c *S3Client) GetObject(ctx context.Context, bucket string, key string) (*s3.GetObjectOutput, error) {

	if bucket == "" {
		return nil, errors.New("no bucket was passed in")
	}
	if key == "" {
		return nil, errors.New("no key was passed in")
	}
	i := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	return c.Client.GetObject(ctx, i)
}

func (c *S3Client) DeleteObject(ctx context.Context, bucket string, key string) (*s3.DeleteObjectOutput, error) {

	if bucket == "" {
		return nil, errors.New("no bucket was passed in")
	}
	if key == "" {
		return nil, errors.New("no key was passed in")
	}
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	return c.Client.DeleteObject(ctx, input)
}

func (c *S3Client) ListObjectsV2(ctx context.Context, bucket, bucketPrefix string) (*s3.ListObjectsV2Output, error) {
	if bucket == "" {
		return nil, errors.New("no bucket was passed in")
	}

	// Prepare the input for the ListObjectsV2 API call
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int32(50),
		Prefix:  aws.String(bucketPrefix),
		//ContinuationToken: continuationToken,
	}

	// Make the API call
	return c.Client.ListObjectsV2(context.TODO(), input)
}

// func (c *S3Client) ListObjects(ctx context.Context, bucket string) ([]string, error) {
// 	if bucket == "" {
// 		return nil, errors.New("no bucket was passed in")
// 	}

// 	var allObjects []string
// 	var continuationToken *string

// 	for {
// 		// Prepare the input for the ListObjectsV2 API call
// 		input := &s3.ListObjectsV2Input{
// 			Bucket:            aws.String(bucket),
// 			MaxKeys:           aws.Int32(50),
// 			ContinuationToken: continuationToken,
// 		}

// 		// Make the API call
// 		resp, err := c.Client.ListObjectsV2(context.TODO(), input)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to list objects, %v", err)
// 		}

// 		// Collect object keys
// 		for _, object := range resp.Contents {
// 			allObjects = append(allObjects, aws.ToString(object.Key))
// 		}

// 		// Check if there are more objects to retrieve
// 		if resp.NextContinuationToken == nil {
// 			break
// 		}

// 		// Update the continuation token for the next request
// 		continuationToken = resp.NextContinuationToken
// 	}

// 	return allObjects, nil
// }
