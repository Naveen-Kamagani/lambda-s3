package main

import (
	"aws-lambda-s3/handler"
	"aws-lambda-s3/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func generateExpectedBody(objectCount int, nextContinuationToken string) (string, error) {
	var objectList []string
	for i := 1; i <= objectCount; i++ {
		objectList = append(objectList, fmt.Sprintf("\"object%d.txt\"", i))
	}
	response := handler.OutputResponse{
		BucketName:   "test-bucket",
		Key:          "",
		ErrorMessage: "",
		Timestamp:    time.Now(),
		Body:         fmt.Sprintf(`{"objects":[%s],"nextContinuationToken":"%s"}`, strings.Join(objectList, ","), nextContinuationToken),
	}
	body, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func TestHandleRequest(t *testing.T) {
	tests := []struct {
		name               string
		request            events.APIGatewayProxyRequest
		handler            MyHandler
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "GET request for all objects",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     "GET",
				PathParameters: map[string]string{"objectKey": ""},
			},
			handler: MyHandler{
				Handler: handler.Handler{
					S3Client: &repositories.S3Client{
						Client: &repositories.S3Mock{
							MockTest: "less-than-50",
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				body, _ := generateExpectedBody(30, "")
				return body
			}(),
		},
		{
			name: "GET request for single object",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     "GET",
				PathParameters: map[string]string{"objectKey": "object1.txt"},
			},
			handler: MyHandler{
				Handler: handler.Handler{
					S3Client: &repositories.S3Client{
						Client: &repositories.S3Mock{
							MockTest: "valid",
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				response := handler.OutputResponse{
					BucketName:   "test-bucket",
					Key:          "object1.txt",
					ErrorMessage: "",
					Timestamp:    time.Now(),
					Body:         "{\"SomethingCool\":\"cool beans\"}", // Expected content of object1.txt
				}
				body, _ := json.Marshal(response)
				return string(body)
			}(),
		},
		{
			name: "DELETE request for single object",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     "DELETE",
				PathParameters: map[string]string{"objectKey": "object1.txt"},
			},
			handler: MyHandler{
				Handler: handler.Handler{
					S3Client: &repositories.S3Client{
						Client: &repositories.S3Mock{
							MockTest: "valid",
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				response := handler.OutputResponse{
					BucketName:   "test-bucket",
					Key:          "object1.txt",
					ErrorMessage: "",
					Timestamp:    time.Now(),
					Body:         "Deleted object1.txt successfully", // Expected delete confirmation
				}
				body, _ := json.Marshal(response)
				return string(body)
			}(),
		},
		{
			name: "DELETE request with nil object key",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     "DELETE",
				PathParameters: map[string]string{"objectKey": ""},
			},
			handler: MyHandler{
				Handler: handler.Handler{
					S3Client: &repositories.S3Client{
						Client: &repositories.S3Mock{
							MockTest: "nil key",
						},
					},
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: func() string {
				response := handler.OutputResponse{
					BucketName:   "test-bucket",
					Key:          "",
					ErrorMessage: "object key is passed as nil",
					Timestamp:    time.Now(),
					Body:         "",
				}
				body, _ := json.Marshal(response)
				return string(body)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("BUCKET_NAME", "test-bucket")
			response, err := tt.handler.handleRequest(tt.request)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)

			// Compare the body while ignoring the timestamp
			var expectedOutput, actualOutput handler.OutputResponse
			_ = json.Unmarshal([]byte(tt.expectedBody), &expectedOutput)
			_ = json.Unmarshal([]byte(response.Body), &actualOutput)

			expectedOutput.Timestamp = actualOutput.Timestamp // Ignore timestamp for comparison
			assert.Equal(t, expectedOutput, actualOutput)
		})
	}
}
