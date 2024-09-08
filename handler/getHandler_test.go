package handler

import (
	"aws-lambda-s3/repositories"
	"fmt"
	"net/http"
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

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name               string
		request            events.APIGatewayProxyRequest
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "less-than-50",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": ""},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       generateExpectedBody(30, ""),
		},
		{
			name: "valid key",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "object1.txt"},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "{\"SomethingCool\":\"cool beans\"}",
		},
		{
			name: "invalid key",
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"objectKey": "non-existent.txt"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "invalid key or object key does not exist in the bucket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				S3Client: &repositories.S3Client{
					Client: &repositories.S3Mock{
						MockTest: tt.name,
					},
				},
			}
			response, err := h.GetHandler(tt.request, "test-bucket")
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tt.expectedBody, response.Body)
		})
	}
}
