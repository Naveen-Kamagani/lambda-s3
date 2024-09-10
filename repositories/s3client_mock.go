package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Mock struct {
	MockTest string
}

type MockedObject struct {
	SomethingCool string
}

func (s *S3Mock) PutObject(_ context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	o := &s3.PutObjectOutput{}
	switch strings.ToLower(s.MockTest) {
	case "valid":
		if params.Bucket == nil || params.Key == nil {
			return nil, errors.New("bucket or key was empty")
		}
		if params.Body == nil {
			return nil, errors.New("passed in body was nil")
		}
		if params.ACL != types.ObjectCannedACLPrivate {
			return nil, errors.New("defaulted ACL incorrect")
		}
		return o, nil
	case "withacl":
		if params.Bucket == nil || params.Key == nil {
			return nil, errors.New("bucket or key was empty")
		}
		if params.Body == nil {
			return nil, errors.New("passed in body was nil")
		}
		if params.ACL != types.ObjectCannedACLBucketOwnerFullControl {
			return nil, errors.New("ACL incorrect, should be bucket owner full control")
		}
		return o, nil
	case "error":
		return nil, errors.New("some error")
	case "generic":
		return nil, nil
	default:
		return nil, &smithy.GenericAPIError{
			Message: "test s3 api error",
		}
	}
}

func (s *S3Mock) GetObject(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {

	b := MockedObject{
		SomethingCool: "cool beans",
	}
	bs, _ := json.Marshal(b)

	o := &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewBuffer(bs)),
	}

	switch strings.ToLower(s.MockTest) {
	case "valid":
		if params.Bucket == nil || params.Key == nil {
			return nil, errors.New("bucket or key was empty")
		}
		return o, nil
	case "error":
		return nil, errors.New("some error")
	case "invalid key":
		return nil, errors.New("invalid key or object key does not exist in the bucket")
	default:
		return nil, &smithy.GenericAPIError{
			Message: "test s3 api error",
		}
	}
}

func (s *S3Mock) DeleteObject(_ context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	o := &s3.DeleteObjectOutput{}

	switch strings.ToLower(s.MockTest) {
	case "valid":
		if params.Bucket == nil || params.Key == nil {
			return nil, errors.New("bucket or key was empty")
		}
		return o, nil
	case "error":
		return nil, errors.New("error deleting object")
	case "invalid":
		return nil, errors.New("object key is passed as nil - ")
	default:
		return nil, &smithy.GenericAPIError{
			Message: "test s3 api error",
		}
	}
}

// func (s *S3Mock) ListObjectsV2(_ context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
// 	o := &s3.ListObjectsV2Output{}

// 	switch strings.ToLower(s.MockTest) {
// 	case "less-than-50":
// 		if params.Bucket == nil {
// 			return nil, errors.New("bucket was empty")
// 		}
// 		// Simulate less than 50 objects
// 		o.Contents = make([]types.Object, 30)
// 		for i := 0; i < 30; i++ {
// 			o.Contents[i] = types.Object{Key: aws.String(fmt.Sprintf("object%d.txt", i+1))}
// 		}
// 		//o.IsTruncated = false // No more objects to retrieve
// 		o.NextContinuationToken = nil
// 		return o, nil

// 	case "more-than-50":
// 		if params.Bucket == nil {
// 			return nil, errors.New("bucket was empty")
// 		}
// 		// Simulate more than 50 objects
// 		o.Contents = make([]types.Object, 50)
// 		for i := 0; i < 50; i++ {
// 			o.Contents[i] = types.Object{Key: aws.String(fmt.Sprintf("object%d.txt", i+1))}
// 		}
// 		//o.IsTruncated = true // More objects to retrieve
// 		o.NextContinuationToken = aws.String("next-token")
// 		return o, nil

// 	case "next-page":
// 		if params.Bucket == nil {
// 			return nil, errors.New("bucket was empty")
// 		}
// 		if params.ContinuationToken == nil || *params.ContinuationToken != "next-token" {
// 			return nil, errors.New("invalid continuation token")
// 		}
// 		// Simulate the next page of objects
// 		o.Contents = make([]types.Object, 25)
// 		for i := 0; i < 25; i++ {
// 			o.Contents[i] = types.Object{Key: aws.String(fmt.Sprintf("object%d.txt", 51+i))}
// 		}
// 		//o.IsTruncated = false // No more objects to retrieve
// 		o.NextContinuationToken = nil
// 		return o, nil

// 	case "error":
// 		return nil, errors.New("some error")
// 	default:
// 		return nil, &smithy.GenericAPIError{
// 			Message: "test s3 api error",
// 		}
// 	}
// }

func (s *S3Mock) ListObjectsV2(_ context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if params == nil {
		params = &s3.ListObjectsV2Input{}
	}
	switch strings.ToLower(s.MockTest) {
	case "less-than-50":
		return &s3.ListObjectsV2Output{
			Contents: GenerateMockS3Objects(30), // Generate 30 objects
		}, nil
	case "more-than-50":
		if params.ContinuationToken == nil || *params.ContinuationToken == "" {
			return &s3.ListObjectsV2Output{
				Contents:              GenerateMockS3Objects(50), // Generate first 50 objects
				NextContinuationToken: aws.String("next-token"),
			}, nil
		} else {
			return &s3.ListObjectsV2Output{
				Contents: GenerateMockS3ObjectsWithOffset(5, 50), // Generate the next 5 objects starting from index 51
			}, nil
		}
	case "error":
		return nil, fmt.Errorf("error listing objects")
	default:
		return nil, &smithy.GenericAPIError{
			Message: "test s3 api error",
		}
	}
}

// GenerateS3Objects creates a slice of S3 objects with keys "object1.txt", "object2.txt", ..., "objectN.txt"
func GenerateMockS3Objects(count int) []types.Object {
	objects := make([]types.Object, count)
	for i := 0; i < count; i++ {
		objects[i] = types.Object{
			Key: aws.String(fmt.Sprintf("object%d.txt", i+1)),
		}
	}
	return objects
}

// GenerateS3ObjectsWithOffset creates a slice of S3 objects with an offset in their numbering
func GenerateMockS3ObjectsWithOffset(count, offset int) []types.Object {
	objects := make([]types.Object, count)
	for i := 0; i < count; i++ {
		objects[i] = types.Object{
			Key: aws.String(fmt.Sprintf("object%d.txt", i+offset+1)),
		}
	}
	return objects
}
