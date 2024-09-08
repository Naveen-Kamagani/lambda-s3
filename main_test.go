package main

import (
	"aws-lambda-s3/handler"
	"aws-lambda-s3/repositories"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func generateExpectedBody(objectCount int, nextContinuationToken string) string {
	var objectList []string
	for i := 1; i <= objectCount; i++ {
		objectList = append(objectList, fmt.Sprintf("\"object%d.txt\"", i))
	}
	return fmt.Sprintf(`{"objects":[%s],"nextContinuationToken":"%s"}`, strings.Join(objectList, ","), nextContinuationToken)
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
			expectedBody:       generateExpectedBody(30, ""),
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
							MockTest: "valid key",
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "{\"SomethingCool\":\"cool beans\"}", // Expected content of object1.txt
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
			expectedBody:       "Deleted object1.txt successfully", // Expected delete confirmation
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
			expectedBody:       "object key is passed as nil - ", // Expected error message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("BUCKET_NAME", "test-bucket")
			response, err := tt.handler.handleRequest(tt.request)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tt.expectedBody, response.Body)
		})
	}
}
