package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestS3ClientGet(t *testing.T) {

	type args struct {
		ctx    context.Context
		key    string
		bucket string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		client   S3ClientInterface
		scenario string
		body     MockedObject
	}{
		{
			name: "Testing happy path valid case",
			args: args{
				ctx:    context.TODO(),
				key:    "testkey",
				bucket: "testbucket",
				//bucketPrefix: "",
			},
			client: &S3Mock{
				MockTest: "valid key",
			},
			wantErr: false,
			body: MockedObject{
				SomethingCool: "cool beans",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &S3Client{
				Client: test.client,
			}
			o, err := c.GetObject(test.args.ctx, test.args.bucket, test.args.key)
			if (err != nil) == test.wantErr {

			} else if err != nil {
				t.Errorf("Test Name: %v failed. S3 Get Object error = %v, wantErr %v", test.name, err, test.wantErr)
			} else {
				var bs []byte
				var obj MockedObject
				o.Body.Read(bs)
				json.Unmarshal(bs, &obj)

				if obj.SomethingCool == test.body.SomethingCool {
					t.Errorf("Test Name: %v failed. Object was not populated properly", test.name)
				}
			}
		})
	}
}

func TestS3ClientListObjects(t *testing.T) {
	type args struct {
		ctx               context.Context
		bucket            string
		bucketPrefix      string
		continuationToken string
	}
	tests := []struct {
		name                 string
		args                 args
		client               S3ClientInterface
		expectedObjects      []string
		expectedContinuation *string
		wantErr              bool
	}{
		{
			name: "Less than 50 items in bucket",
			args: args{
				ctx:          context.TODO(),
				bucket:       "testbucket1",
				bucketPrefix: "",
			},
			client: &S3Mock{
				MockTest: "less-than-50",
			},
			expectedObjects:      generateExpectedObjects(30), // Generate 30 objects
			expectedContinuation: nil,
			wantErr:              false,
		},
		{
			name: "More than 50 items in bucket with pagination",
			args: args{
				ctx:          context.TODO(),
				bucket:       "testbucket1",
				bucketPrefix: "",
			},
			client: &S3Mock{
				MockTest: "more-than-50",
			},
			expectedObjects:      generateExpectedObjects(50), // Generate 50 objects
			expectedContinuation: aws.String("next-token"),
			wantErr:              false,
		},
		{
			name: "Error case",
			args: args{
				ctx:          context.TODO(),
				bucket:       "testbucket3",
				bucketPrefix: "",
			},
			client: &S3Mock{
				MockTest: "error",
			},
			expectedObjects:      nil,
			expectedContinuation: nil,
			wantErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &S3Client{
				Client: tt.client,
			}
			resp, err := c.ListObjectsV2(tt.args.ctx, tt.args.bucket, tt.args.bucketPrefix, tt.args.continuationToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListObjectsV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !equal(resp.Contents, tt.expectedObjects) {
					t.Errorf("ListObjectsV2() got objects = %v, want %v", resp.Contents, tt.expectedObjects)
				}
				if !equalString(resp.NextContinuationToken, tt.expectedContinuation) {
					t.Errorf("ListObjectsV2() got continuation token = %v, want %v", resp.NextContinuationToken, tt.expectedContinuation)
				}
			}
		})
	}
}

// Helper functions to compare object keys and continuation tokens
func equal(objects []types.Object, expectedObjects []string) bool {
	if len(objects) != len(expectedObjects) {
		return false
	}
	for i, object := range objects {
		if object.Key == nil || aws.ToString(object.Key) != expectedObjects[i] {
			return false
		}
	}
	return true
}

func equalString(token *string, expectedToken *string) bool {
	if token == nil && expectedToken == nil {
		return true
	}
	if token == nil || expectedToken == nil {
		return false
	}
	return *token == *expectedToken
}

// Mock and helper functions
func generateExpectedObjects(count int) []string {
	var objects []string
	for i := 0; i < count; i++ {
		objects = append(objects, fmt.Sprintf("object%d.txt", i+1))
	}
	return objects
}
